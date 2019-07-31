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
	Display string
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

func newMatchedEni(svc *ec2.EC2, eni *ec2.NetworkInterface) *MatchedEni {
	if *eni.Status == "available" {
		return &MatchedEni{
			Display: fmt.Sprintf("Unattached in %s", *eni.VpcId),
		}
	} else if *eni.Attachment.InstanceOwnerId == "amazon-elb" {
		return &MatchedEni{
			Display: fmt.Sprintf("ELB (%s)", *eni.Description),
		}
	}

	// Default
	return &MatchedEni{
		Display: fmt.Sprintf(
			"%s (%s)",
			*eni.Attachment.InstanceId,
			ec2NameTag(svc, *eni.Attachment.InstanceId),
		),
	}
}

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
				ret = append(ret, *newMatchedEni(svc, eni))
			}
			return true
		},
	)
	if err != nil {
		fmt.Printf(
			"[ERROR] %s/%s for %s: %s\n",
			profileSession.Profile,
			profileSession.Region,
			ip,
			err,
		)
	}
	return ret
}

func ec2NameTag(svc *ec2.EC2, instanceID string) string {
	params := &ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: aws.String("instance-id"),
				Values: []*string{
					aws.String(instanceID),
				},
			},
		},
	}

	resp, err := svc.DescribeInstances(params)

	if err != nil || len(resp.Reservations) != 1 {
		return "COULD NOT FETCH INSTANCE"
	}

	for _, tag := range resp.Reservations[0].Instances[0].Tags {
		if *tag.Key == "Name" {
			return *tag.Value
		}
	}

	return "COULD NOT FIND TAG"
}
