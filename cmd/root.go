package cmd

import (
	"fmt"
	"ghe/utils"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ghe",
	Short: "github CLI enhanced",
	Long: `ghe or GH enhanced is a wrapper of gh cli to make the experience of using it better
	`,
}

func Execute() {
	rootCmd.AddCommand(crCmd)
	cmd, args, _ := rootCmd.Find(os.Args[1:])

	// Handle err when it's actually an err and not a gh command case like below
	if cmd.Use == rootCmd.Use && args[0] != "cr" { // && cmd.Flags().Parse(os.Args[1:]) != pflag.ErrHelp {
		fmt.Println("Running GH command instead...")
		utils.ExecCmd("gh", args...)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
