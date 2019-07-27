package profiles

import (
	"fmt"
	"os/user"

	ini "gopkg.in/go-ini/ini.v1"
)

//FromCredentials returns all profiles in the ~/.aws/credentials file
func FromCredentials() (ret []string) {

	// Get user directory
	usr, err := user.Current()

	// Load INI from typical path
	cfg, err := ini.Load(usr.HomeDir + "/.aws/credentials")
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
