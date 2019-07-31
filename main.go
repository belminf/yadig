package main

import (
	"flag"
	"fmt"
	"sync"

	"github.com/belminf/yadig/aws"
	"github.com/belminf/yadig/profiles"
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

	// Collect results
	results := make(profileResultsType)
	wg := sync.WaitGroup{}
	for _, p := range profiles.FromConfig(aws.ConfigIniPath) {
		results[p] = make(regionResultsType)
		sess := aws.ProfileSession(p, "")
		r := sess.Region
		me := make(matchedEnisType, 0, 0)
		results[p][r] = &me
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
