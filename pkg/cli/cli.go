package cli

import (
	"fmt"
	"runtime"

	"github.com/alecthomas/kong"
)

// VersionFlag is a flag type that can be used to display a version number, stored in the "version" variable.
type VersionFlag bool

// BeforeApply writes the version variable and terminates with a 0 exit status.
func (v VersionFlag) BeforeApply(app *kong.Kong, vars kong.Vars) error {
	versionstring := fmt.Sprintf("%s commit:%s release:%s build:%s date:%s goVersion:%s platform:%s/%s", vars["binName"], vars["commit"], vars["version"], vars["builtBy"], vars["date"], runtime.Version(), runtime.GOOS, runtime.GOARCH)
	fmt.Fprintln(app.Stdout, versionstring)
	app.Exit(0)
	return nil
}
