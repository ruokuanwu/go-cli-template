package cli

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

const usage = `Usage: go-cli-template <command> [options]

Commands:
  greet       Print a greeting

Options:
  -h, --help  Show help
  --version   Show version
`

// Run executes the CLI and returns a process exit code without exiting the process.
func Run(args []string, stdout, stderr io.Writer, version string) int {
	if len(args) == 0 {
		fmt.Fprint(stdout, usage)
		return 0
	}
	switch args[0] {
	case "-h", "--help", "--version":
		if len(args) != 1 {
			fmt.Fprintln(stderr, "unexpected arguments:", strings.Join(args[1:], " "))
			return 2
		}
		if args[0] == "--version" {
			fmt.Fprintln(stdout, "go-cli-template", version)
		} else {
			fmt.Fprint(stdout, usage)
		}
		return 0
	case "greet":
		return greet(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "unknown command: %s\nRun 'go-cli-template --help' for usage.\n", args[0])
		return 2
	}
}

func greet(args []string, stdout, stderr io.Writer) int {
	flags := flag.NewFlagSet("greet", flag.ContinueOnError)
	flags.SetOutput(stderr)
	name := flags.String("name", "world", "Name to greet")
	flags.Usage = func() {}
	if err := flags.Parse(args); err != nil {
		if err == flag.ErrHelp {
			fmt.Fprintln(stdout, "Usage: go-cli-template greet [--name NAME]")
			fmt.Fprintln(stdout, "  --name string  Name to greet (default: world)")
			return 0
		}
		return 2
	}
	if flags.NArg() != 0 {
		fmt.Fprintln(stderr, "unexpected arguments:", strings.Join(flags.Args(), " "))
		return 2
	}
	if strings.TrimSpace(*name) == "" {
		fmt.Fprintln(stderr, "name must not be empty")
		return 2
	}
	fmt.Fprintf(stdout, "Hello, %s!\n", *name)
	return 0
}
