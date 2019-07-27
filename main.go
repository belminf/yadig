package main

import (
	"fmt"

	"github.com/belminf/yadig/profiles"
)

func main() {
	for _, s := range profiles.FromCredentials() {
		fmt.Println(s)
	}
}
