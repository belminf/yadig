package main

import (
	"fmt"
	"yadig/profiles"
)

func main() {
	for _, s := range profiles.FromCredentials() {
		fmt.Println(s)
	}
}
