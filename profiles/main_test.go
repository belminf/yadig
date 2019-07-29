package profiles

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	testConfigIniPath    = filepath.Join("testdata", "config")
	testProfilesExpected = []string{"personal", "work"}
)

func TestCredsIni(t *testing.T) {
	profiles := FromConfig(testConfigIniPath)

	if !reflect.DeepEqual(profiles, testProfilesExpected) {
		t.Fatal(fmt.Sprintf("Profiles (%s) didn't match expected (%s)", profiles, testProfilesExpected))
	}
}
