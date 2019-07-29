package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/defaults"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//MatchedEni holds info on ENIs that matched an IP
type MatchedEni struct {
	OwnerID     string
	InstanceID  string
	Description string
	VpcID       string
	Attached    bool
}

var (
	//ConfigIniPath is the path AWS credentials file
	ConfigIniPath = defaults.SharedConfigFilename()
	credsIniPath  = defaults.SharedCredentialsFilename()
)

//InterfacesWithIP gets all ENIs that match the IP in the AWS profile and region
func InterfacesWithIP(profile string, region string, ip string) []MatchedEni {

	ret := []MatchedEni{}

	sess := session.Must(session.NewSessionWithOptions(session.Options{
		Config:  aws.Config{Region: aws.String(region)},
		Profile: profile,
	}))

	svc := ec2.New(sess)

	params := &ec2.DescribeNetworkInterfacesInput{
		Filters: []*ec2.Filter{
			{
				Name:   aws.String("private-ip-address"),
				Values: aws.StringSlice([]string{ip}),
			},
		},
	}

	err := svc.DescribeNetworkInterfacesPages(
		params,
		func(page *ec2.DescribeNetworkInterfacesOutput, lastPage bool) bool {
			for _, eni := range page.NetworkInterfaces {
				thisEni := MatchedEni{
					Description: *eni.Description,
					Attached:    *eni.Status != "available",
					VpcID:       *eni.VpcId,
				}

				if eni.Attachment != nil {
					thisEni.OwnerID = *eni.Attachment.InstanceOwnerId
					thisEni.InstanceID = *eni.Attachment.InstanceId
				}

				ret = append(ret, thisEni)
			}
			return true
		},
	)
	if err != nil {
		fmt.Println("Error", err)
	}
	return ret
}
