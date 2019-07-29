package profiles

import (
	"fmt"
	"strings"

	ini "gopkg.in/go-ini/ini.v1"
)

//FromConfig returns all profiles in the ~/.aws/config file
func FromConfig(configIniPath string) (ret []string) {

	// Load INI from typical path
	cfg, err := ini.Load(configIniPath)
	if err != nil {
		fmt.Printf("Fail to read: %v", err)
	}

	// Filter profile results
	for _, s := range cfg.SectionStrings() {

		splitSect := strings.SplitN(s, " ", 2)

		// Ignore DEFAULT
		if splitSect[0] != "profile" {
			continue
		}

		ret = append(ret, splitSect[1])
	}

	return
}
