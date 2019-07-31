package main

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/belminf/yadig/aws"
	"github.com/belminf/yadig/config"
)

type regionResultType struct {
	Alias string
	Enis  []aws.MatchedEni
}
type regionResultsType map[string]*regionResultType
type profileResultsType map[string]regionResultsType

var ip string

// Collect flags
func init() {

}

// Collect ENIs
func addMatchedEnis(sess *aws.ProfileSessionType, ip string, regionResult *regionResultType, wg *sync.WaitGroup) {
	defer wg.Done()
	regionResult.Enis = aws.InterfacesWithIP(sess, ip)
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
		rr := regionResultType{Alias: pr.Alias()}
		results[pr.Profile][pr.Region] = &rr
		wg.Add(1)
		go addMatchedEnis(sess, ip, &rr, &wg)
	}
	wg.Wait()

	// Print results
	fmt.Println("")
	for _, rmap := range results {
		for _, rrmap := range rmap {
			for _, e := range rrmap.Enis {
				fmt.Printf("[%s] %s", rrmap.Alias, e.Display)
			}
		}
	}
	fmt.Println("")
}
