package main

import (
	"flag"
	"fmt"
	"strings"
	"sync"

	"github.com/belminf/yadig/aws"
	"github.com/belminf/yadig/profiles"
)

type matchedEnisType []aws.MatchedEni
type regionResultsType map[string]*matchedEnisType
type profileResultsType map[string]regionResultsType
type arrayFlags []string

var ip string
var regions arrayFlags

func (r *arrayFlags) String() string {
	return strings.Join(*r, ",")
}

func (r *arrayFlags) Set(value string) error {
	*r = append(*r, value)
	return nil
}

// Collect flags
func init() {
	flag.StringVar(&ip, "ip", "", "IP address")
	flag.Var(&regions, "region", "Regions")
}

// Collect ENIs
func addMatchedEnis(profile, region, ip string, matchedEnis *matchedEnisType, wg *sync.WaitGroup) {
	defer wg.Done()
	*matchedEnis = aws.InterfacesWithIP(profile, region, ip)
}

func main() {
	flag.Parse()

	// Collect results
	results := make(profileResultsType)
	wg := sync.WaitGroup{}
	for _, p := range profiles.FromConfig(aws.ConfigIniPath) {
		fmt.Println(p)
		results[p] = make(regionResultsType)
		for _, r := range regions {
			fmt.Println(r)
			me := make(matchedEnisType, 0, 0)
			results[p][r] = &me
			wg.Add(1)
			go addMatchedEnis(p, r, ip, &me, &wg)
		}
	}
	wg.Wait()

	// Print results
	for p, rmap := range results {
		for r, emap := range rmap {
			for _, e := range *emap {
				fmt.Printf("Profile: %s, Region: %s", p, r)
				fmt.Println(e)
			}
		}
	}
}
