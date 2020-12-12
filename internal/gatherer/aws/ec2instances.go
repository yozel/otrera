package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pkg/errors"
	"github.com/yozel/otrera/internal/types"
)

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
					if *instance.State.Name == "terminated" {
						continue
					}
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

func (e *EC2InstanceObject) Copy() types.RawObject {
	return types.RawObject{IName: e.Name(), IContent: e.Content()}
}

func DescribeEC2Instances(options map[string]string) ([]types.RawObjectInterface, error) {
	profile := options["profile"]
	region := options["region"]
	instances, err := getInstanceDetails(profile, region)
	err = errors.Wrapf(err, fmt.Sprintf("Can't get []*ec2.Instance for %s %s", profile, region))
	if err != nil {
		return nil, err
	}
	result := []types.RawObjectInterface{}
	for _, instance := range instances {
		result = append(result, &EC2InstanceObject{instance})
	}
	return result, nil
}
