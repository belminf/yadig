package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/belminf/yadig/aws"
	"github.com/belminf/yadig/config"
)

type matchedEnisType []aws.MatchedEni
type regionResultsType map[string]*matchedEnisType
type profileResultsType map[string]regionResultsType

var ip string

// Collect flags
func init() {

}

// Collect ENIs
func addMatchedEnis(sess *aws.ProfileSessionType, ip string, matchedEnis *matchedEnisType, wg *sync.WaitGroup) {
	defer wg.Done()
	*matchedEnis = aws.InterfacesWithIP(sess, ip)
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
	results := make(profileResultsType)
	wg := sync.WaitGroup{}
	for _, pr := range yadigConfig.ProfileRegions {
		results[pr.Profile] = make(regionResultsType)
		sess := aws.ProfileSession(pr.Profile, pr.Region)
		pr.Region = sess.Region
		me := make(matchedEnisType, 0, 0)
		results[pr.Profile][pr.Region] = &me
		wg.Add(1)
		go addMatchedEnis(sess, ip, &me, &wg)
	}
	wg.Wait()

	// Print results
	for p, rmap := range results {
		for r, emap := range rmap {
			for _, e := range *emap {
				fmt.Printf("[%s/%s] %s", p, r, e.Display)
			}
		}
	}
}
