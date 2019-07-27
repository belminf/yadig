package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/belminf/yadig/profiles"
)

func main() {
	credsIniPath := external.DefaultSharedCredentialsFilename()
	for _, s := range profiles.FromCredentials(credsIniPath) {
		fmt.Println(s)
	}
}
