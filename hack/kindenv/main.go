package main

import (
	"fmt"
	"os"
	"reflect"
	"strings"
)

const (

	// Available commands
	setupCommand    = "setup"
	teardownCommand = "teardown"
	helpCommand     = "usage"
)

// ----------------------------------------------------- USAGE ------------------------------------------------------ //

const (
	banner        = "# KINDENV\n\n"
	usageTemplate = `## Usage

%s [command]

Available commands:
  - %q
  - %q
`
)

func usage() error {
	arg0 := fmt.Sprintf("go run \"%s/hack/kindenv\"", os.Getenv("PWD"))
	_, _ = fmt.Fprintf(os.Stdout, usageTemplate, arg0, setupCommand, teardownCommand)

	return nil
}

// ----------------------------------------------------- MAIN ------------------------------------------------------- //

func main() {
	_, _ = fmt.Fprint(os.Stdout, banner)

	// 1. Print usageTemplate or

	if len(os.Args) < 2 { //nolint:gomnd // if no specified subcommand then print usageTemplate and exit.
		_ = usage()

		os.Exit(1)
	}

	// 2. Switch command.

	var command func() error

	switch os.Args[1] {
	case setupCommand:
		command = setup
	case teardownCommand:
		command = teardown
	case helpCommand:
		command = usage
	}

	// 3. Execute command

	if err := command(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	os.Exit(0)
}

// ----------------------------------------------------- HELPERS ---------------------------------------------------- //

func formatExpectedEnvList[T any]() string {
	optionalEnvs := make([]string, 0)
	requiredEnvs := make([]string, 0)

	observedMaxStrLen := 0

	rt := reflect.TypeFor[T]()
	for i := range rt.NumField() {
		field := rt.Field(i)
		val, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}

		substr := strings.Split(val, ",")
		switch len(substr) {
		case 0:
			continue
		case 1:
			optionalEnvs = append(optionalEnvs, substr[0])
		default:
			requiredEnvs = append(requiredEnvs, substr[0])
		}

		if envStrLen := len(substr[0]); envStrLen > observedMaxStrLen {
			observedMaxStrLen = envStrLen
		}
	}

	envs := ""
	for _, s := range requiredEnvs {
		envs = fmt.Sprintf("%s- %s %s[Required]\n", envs, s, fmtSpaces(s, observedMaxStrLen))
	}

	for _, s := range optionalEnvs {
		envs = fmt.Sprintf("%s- %s %s[Optional]\n", envs, s, fmtSpaces(s, observedMaxStrLen))
	}

	return envs
}

func fmtSpaces(s string, maxLen int) string {
	return strings.Repeat(" ", maxLen-len(s))
}
