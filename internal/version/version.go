package version

// These are set by the linker
var gitHash = "Unknown"
var buildTime = "Unknown"
var versionNumber = "Unknown"

type Version struct {
	GitHash       string `json:"gitHash"`
	BuildTime     string `json:"buildTime"`
	VersionNumber string `json:"versionNumber"`
}

// Info returns the current version info
func Info() *Version {
	return &Version{
		GitHash:       gitHash,
		BuildTime:     buildTime,
		VersionNumber: versionNumber,
	}
}
