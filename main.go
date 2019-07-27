package main

import (
	"fmt"
	"os/user"

	ini "gopkg.in/go-ini/ini.v1"
)

func main() {
	for _, s := range profiles() {
		fmt.Println(s)
	}
}

func profiles() (ret []string) {

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
