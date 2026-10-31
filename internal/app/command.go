package app

import (
	"fmt"
	"io"
	"strings"

	"example.com/go-cli-template/internal/service"
	"github.com/samber/oops"
	"github.com/spf13/cobra"
)

func execute(args []string, stdout, stderr io.Writer, version string, greetings *service.Greeting) error {
	if len(args) > 1 && (args[0] == "-h" || args[0] == "--help" || args[0] == "--version") {
		return oops.Errorf("unexpected arguments: %s", strings.Join(args[1:], " "))
	}
	root := newCommand(version, greetings)
	root.SetOut(stdout)
	root.SetErr(stderr)
	root.SetArgs(args)
	return root.Execute()
}

func newCommand(version string, greetings *service.Greeting) *cobra.Command {
	root := &cobra.Command{
		Use:           "go-cli-template",
		Short:         "A Go CLI template",
		Version:       version,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	root.CompletionOptions.DisableDefaultCmd = true
	root.SetVersionTemplate("{{.Name}} {{.Version}}\n")
	root.SetHelpFunc(func(cmd *cobra.Command, _ []string) {
		fmt.Fprintf(cmd.OutOrStdout(), "Usage: %s\n\n%s\n", cmd.UseLine(), cmd.Short)
		if cmd.HasAvailableSubCommands() {
			fmt.Fprintln(cmd.OutOrStdout(), "\nCommands:")
			for _, child := range cmd.Commands() {
				fmt.Fprintf(cmd.OutOrStdout(), "  %-12s %s\n", child.Name(), child.Short)
			}
		}
		fmt.Fprintln(cmd.OutOrStdout(), "\nOptions:")
		fmt.Fprint(cmd.OutOrStdout(), cmd.Flags().FlagUsages())
	})
	var name string
	greet := &cobra.Command{
		Use:   "greet [--name NAME]",
		Short: "Print a greeting",
		Args: func(_ *cobra.Command, args []string) error {
			if len(args) != 0 {
				return oops.Errorf("unexpected arguments: %s", strings.Join(args, " "))
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, _ []string) error {
			message, err := greetings.Greet(name)
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(cmd.OutOrStdout(), message)
			return oops.Wrapf(err, "write greeting")
		},
	}
	greet.Flags().StringVar(&name, "name", "world", "Name to greet")
	root.AddCommand(greet)
	return root
}
