package version

import (
	"fmt"
	"runtime"

	"go.uber.org/zap"
)

// These fields are set during an official build
// Global vars set from command-line arguments
var (
	buildVersion = "--"
	buildHash    = "--"
	buildTime    = "--"
	gitBranch    = "--"
)

// BuildDetails returns a string containing details about the SigNoz query-service binary.
func BuildDetails() string {
	licenseInfo := `Check SigNoz Github repo for license details`

	return fmt.Sprintf(`
SigNoz version   : %v
Commit SHA-1     : %v
Commit timestamp : %v
Branch           : %v
Go version       : %v

For SigNoz Official Documentation,  visit https://signoz.io/docs
For SigNoz Community Slack,         visit http://signoz.io/slack
For discussions about SigNoz,       visit https://community.signoz.io

%s.
Copyright 2022 SigNoz
`,
		buildVersion, buildHash, buildTime, gitBranch,
		runtime.Version(), licenseInfo)
}

// PrintVersion prints version and other helpful information.
func PrintVersion() {
	zap.S().Infof("\n%s\n", BuildDetails())
}

func GetVersion() string {
	return buildVersion
}
