package profiles

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	ini "gopkg.in/go-ini/ini.v1"
)

//FromCredentials returns all profiles in the ~/.aws/credentials file
func FromCredentials() (ret []string) {

	// Get profrile list
	credsIni := external.DefaultSharedCredentialsFilename()
	// Load INI from typical path
	cfg, err := ini.Load(credsIni)
	if err != nil {
		fmt.Printf("Fail to read: %v", err)
	}

	// Filter profile results
	for _, s := range cfg.SectionStrings() {

		// Ignore DEFAULT
		if s == "DEFAULT" {
			continue
		}

		ret = append(ret, s)
	}

	return
}
