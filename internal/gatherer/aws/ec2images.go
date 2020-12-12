package aws

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"github.com/yozel/otrera/internal/types"
)

type EC2ImageObject struct {
	ec2instance *ec2.Image
}

func (e *EC2ImageObject) Name() string {
	return *e.ec2instance.ImageId
}

func (e *EC2ImageObject) Content() interface{} {
	return *e.ec2instance
}

func (e *EC2ImageObject) Copy() types.RawObject {
	return types.RawObject{IName: e.Name(), IContent: e.Content()}
}

func getImageDetails(profile, region string) ([]*ec2.Image, error) {
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

	stsClient := sts.New(sess)
	cid, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, err
	}

	client := ec2.New(sess)

	output, err := client.DescribeImages(&ec2.DescribeImagesInput{Owners: []*string{cid.Account}})
	err = errors.Wrapf(err, "Can't retrieve instance list from aws")
	if err != nil {
		return nil, err
	}

	return output.Images, nil
}

func DescribeEC2Images(options map[string]string) ([]types.RawObjectInterface, error) {
	profile := options["profile"]
	region := options["region"]
	images, err := getImageDetails(profile, region)
	err = errors.Wrapf(err, fmt.Sprintf("Can't get []*ec2.Images for %s %s", profile, region))
	if err != nil {
		return nil, err
	}
	result := []types.RawObjectInterface{}
	for _, image := range images {
		result = append(result, &EC2ImageObject{image})
	}
	return result, nil
}
