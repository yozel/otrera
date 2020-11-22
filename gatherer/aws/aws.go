package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/yozel/otrera/gatherer"
	"gopkg.in/ini.v1"
)

func ListProfiles(filename string) ([]string, error) {
	var r []string
	cfg, err := ini.Load(filename)
	err = errors.Wrapf(err, "Failed to read file")
	if err != nil {
		return nil, err
	}
	sections := cfg.Sections()
	for _, section := range sections {
		var pn string
		sn := section.Name()
		if sn == "DEFAULT" {
			continue
		} else if sn == "default" {
			pn = sn
		} else {
			if sn[0:8] != "profile " {
				log.Printf("Can't parse section: %s\n", sn)
				continue
			}
			pn = sn[8:]
		}
		r = append(r, pn)
	}
	return r, nil
}

func getInstanceDetails(profile, region string) ([]*ec2.Instance, error) {
	sess, err := session.NewSessionWithOptions(
		session.Options{
			Config:            aws.Config{Region: aws.String(region)},
			Profile:           profile,
			SharedConfigState: session.SharedConfigEnable,
			// AssumeRoleTokenProvider: stscreds.StdinTokenProvider, // For MFA
		})
	err = errors.Wrapf(err, "Can't create aws session")
	if err != nil {
		log.Fatal(err)
	}

	client := ec2.New(sess)

	var r []*ec2.Instance
	_, err = client.DescribeInstances(&ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatal(err)
	}

	err = client.DescribeInstancesPages(&ec2.DescribeInstancesInput{},
		func(page *ec2.DescribeInstancesOutput, lastPage bool) bool {
			for _, reservation := range page.Reservations {
				for _, instance := range reservation.Instances {
					r = append(r, instance)
				}
			}
			return !lastPage
		})
	err = errors.Wrapf(err, "Can't retrieve instance list from aws")
	if err != nil {
		return nil, err
	}

	return r, nil
}

type EC2InstanceObject struct {
	ec2instance *ec2.Instance
}

func (e *EC2InstanceObject) Name() string {
	return *e.ec2instance.InstanceId
}

func (e *EC2InstanceObject) Content() interface{} {
	return *e.ec2instance
}

func (e *EC2InstanceObject) Copy() gatherer.RawObject {
	return gatherer.RawObject{IName: e.Name(), IContent: e.Content()}
}

func DescribeEC2Instances(options map[string]string) ([]gatherer.RawObjectInterface, error) {
	profile := options["profile"]
	region := options["region"]
	instances, err := getInstanceDetails(profile, region)
	err = errors.Wrapf(err, fmt.Sprintf("Can't get []*ec2.Instance for %s %s", profile, region))
	if err != nil {
		return nil, err
	}
	result := []gatherer.RawObjectInterface{}
	for _, instance := range instances {
		result = append(result, &EC2InstanceObject{instance})
	}
	return result, nil
}
