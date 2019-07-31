package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os/user"
	"path/filepath"
	"strings"

	ini "gopkg.in/go-ini/ini.v1"
	yaml "gopkg.in/yaml.v2"
)

//AWSProfileRegion is a profile and region pair
type AWSProfileRegion struct {
	Profile  string `yaml:"profile"`
	Region   string `yaml:"region"`
	AliasVal string `yaml:"alias,omitempty"`
}

//Alias for profile/region pair
func (pr *AWSProfileRegion) Alias() string {
	if pr.AliasVal != "" {
		return pr.AliasVal
	}
	if pr.Region != "" {
		return fmt.Sprintf("%s/%s", pr.Profile, pr.Region)
	}

	return pr.Profile
}

//Config struct for yadig
type Config struct {
	ProfileRegions []AWSProfileRegion `yaml:"search"`
}

//LoadConfig loads a configuration file
func LoadConfig(awsConfigPath string) *Config {
	configFile, err := ioutil.ReadFile(*configPath())

	// If cannot open file, continue
	if err != nil {
		return loadAWSProfiles(awsConfigPath)
	}

	config := Config{}
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Printf("YAML error: %v\n", err)
	}
	return &config
}

func loadAWSProfiles(configIniPath string) *Config {

	// Load INI from typical path
	cfg, err := ini.Load(configIniPath)
	if err != nil {
		fmt.Printf("Fail to read: %v", err)
	}

	// Filter profile results
	config := Config{}
	config.ProfileRegions = make([]AWSProfileRegion, 1, 1)
	for _, s := range cfg.SectionStrings() {

		splitSect := strings.SplitN(s, " ", 2)

		// Ignore DEFAULT
		if splitSect[0] != "profile" {
			continue
		}

		config.ProfileRegions = append(
			config.ProfileRegions,
			AWSProfileRegion{
				Profile: splitSect[1],
				Region:  "",
			},
		)
	}

	return &config
}

func configPath() *string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	path := filepath.Join(usr.HomeDir, ".config", "yadig", "config.yaml")
	return &path
}
