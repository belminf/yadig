package profiles

import (
	"fmt"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	testCredsIniPath     = filepath.Join("testdata", "credentials")
	testProfilesExpected = []string{"personal", "work"}
)

func TestCredsIni(t *testing.T) {
	profiles := FromCredentials(testCredsIniPath)

	if !reflect.DeepEqual(profiles, testProfilesExpected) {
		t.Fatal(fmt.Sprintf("Profiles (%s) didn't match expected (%s)", profiles, testProfilesExpected))
	}
}
