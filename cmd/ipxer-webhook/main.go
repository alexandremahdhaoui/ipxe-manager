package main

import (
	"fmt"
	"os"
)

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
	_, _ = fmt.Fprintf(os.Stdout, "Starting %s version %s (%s) %s\n", Name, Version, CommitSHA, BuildTimestamp)

	// TODO: implement me
	panic("implement me")
}
