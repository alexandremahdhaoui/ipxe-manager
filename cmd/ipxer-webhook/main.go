package main

import "fmt"

const (
	Name = "ipxer-webhook"
)

var (
	Version        = "dev" //nolint:gochecknoglobals // set by ldflags
	CommitSHA      = "n/a" //nolint:gochecknoglobals // set by ldflags
	BuildTimestamp = "n/a" //nolint:gochecknoglobals // set by ldflags
)

// ------------------------------------------------- Main ----------------------------------------------------------- //

func main() {
	fmt.Printf("Starting %s version %s (%s) %s\n", Name, Version, CommitSHA, BuildTimestamp)
}
