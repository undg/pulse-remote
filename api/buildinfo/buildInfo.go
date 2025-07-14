package buildinfo

import (
	"fmt"
	"runtime"
)

var (
	// These variables are replaced by ldflags at build time
	GitVersion = "v0.0.0-main"
	GitCommit  = "0000000"
	BuildTime  = "1970-01-01T00:00:00Z" // build date in ISO8601 format
)

type BuildInfo struct {
	GitVersion string `json:"gitVersion" doc:"The version of the Git repository"`
	GitCommit  string `json:"gitCommit" doc:"The Git commit hash"`
	BuildDate  string `json:"buildDate" doc:"The date and time of the build in ISO8601 format"`
	GoVersion  string `json:"goVersion" doc:"The version of Go used for compilation"`
	Compiler   string `json:"compiler" doc:"The compiler used for building"`
	Platform   string `json:"platform" doc:"The operating system and architecture"`
}

// Get returns the overall codebase version. It's for detecting
// what code a binary was built from.
func Get() *BuildInfo {
	// These variables come from -ldflags settings and in
	// their absence fallback to the constants above
	return &BuildInfo{
		GitVersion: GitVersion,
		GitCommit:  GitCommit,
		BuildDate:  BuildTime,
		GoVersion:  runtime.Version(),
		Compiler:   runtime.Compiler,
		Platform:   fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}
