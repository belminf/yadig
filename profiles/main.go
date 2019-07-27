package profiles

import (
	"fmt"

	ini "gopkg.in/go-ini/ini.v1"
)

//FromCredentials returns all profiles in the ~/.aws/credentials file
func FromCredentials(credsIniPath string) (ret []string) {

	// Load INI from typical path
	cfg, err := ini.Load(credsIniPath)
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
