package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/belminf/yadig/aws"
	"github.com/belminf/yadig/config"
)

// Collect flags
func init() {

}

// Collect ENIs
func getMatchedEnis(sess *aws.ProfileSessionType, ip string, prAlias string, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, i := range aws.InterfacesWithIP(sess, ip) {
		fmt.Printf("[%s] %s\n", prAlias, i.Display)
	}
}

func main() {
	flag.Parse()
	ip := flag.Arg(0)
	yadigConfig := config.LoadConfig(aws.ConfigIniPath)

	if ip == "" {
		fmt.Println("[ERROR] Provided no IP")
		os.Exit(1)
	}

	// Collect results
	wg := sync.WaitGroup{}
	fmt.Println("")
	for _, pr := range yadigConfig.ProfileRegions {
		sess := aws.ProfileSession(pr.Profile, pr.Region)
		wg.Add(1)
		go getMatchedEnis(sess, ip, pr.Alias(), &wg)
	}
	wg.Wait()
	fmt.Println("")
}
