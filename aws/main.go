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

//ProfileSessionType holds info about an AWS profile session
type ProfileSessionType struct {
	Profile string
	Region  string
	Session *session.Session
}

var (
	//ConfigIniPath is the path AWS credentials file
	ConfigIniPath = defaults.SharedConfigFilename()
	credsIniPath  = defaults.SharedCredentialsFilename()
)

//ProfileSession returns an AWS session for a profile
func ProfileSession(profile, region string) *ProfileSessionType {
	profileSession := ProfileSessionType{
		Profile: profile,
	}

	sessOptions := session.Options{
		Profile:           profile,
		SharedConfigState: session.SharedConfigEnable,
	}

	//Force region if provided
	if region != "" {
		sessOptions.Config = aws.Config{Region: aws.String(region)}
	}

	profileSession.Session = session.Must(session.NewSessionWithOptions(sessOptions))
	profileSession.Region = *profileSession.Session.Config.Region

	return &profileSession
}

//InterfacesWithIP gets all ENIs that match the IP in the AWS profile and region ("" = default region)
func InterfacesWithIP(profileSession *ProfileSessionType, ip string) []MatchedEni {

	ret := []MatchedEni{}
	svc := ec2.New(profileSession.Session)

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
		fmt.Printf(
			"Error for %s/%s: %s\n",
			profileSession.Profile,
			profileSession.Region,
			err,
		)
	}
	return ret
}
