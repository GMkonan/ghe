package cmd

import (
	"fmt"
	"os"

	"ghe/utils"

	"github.com/spf13/cobra"
)

func init() {

	var Org string
	crCmd.Flags().StringVar(&Org, "org", "", "Define a organization to create your repo into")

	var IsPrivate bool
	crCmd.Flags().BoolVar(&IsPrivate, "private", false, "Create repo as private")

	var Push bool
	crCmd.Flags().BoolVar(&Push, "push", false, "Push to repo")

	var Confirm bool
	crCmd.Flags().BoolVar(&Confirm, "confirm", false, "Confirm create repo action without prompting")
}

var crCmd = &cobra.Command{
	Use:   "cr",
	Short: "Create a repo on github",
	Long: `Create a repo with (repo create) from gh but with sane defaults
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// upstream?
		ghBaseArgs := []string{"gh", "repo", "create"}
		ghFullArgs := make([]string, len(ghBaseArgs))
		ghPositionalArgs := make([]string, 0)

		copy(ghFullArgs, ghBaseArgs)

		org, _ := cmd.Flags().GetString("org")
		if org != "" {
			dirName := utils.GetDirName()
			orgAndProject := org + "/" + dirName
			ghPositionalArgs = append(ghPositionalArgs, orgAndProject)
		}

		isRepo, err := os.Stat(".git")

		if err != nil && !os.IsNotExist(err) {
			fmt.Println("error on os stat")
		}

		if os.IsNotExist(err) {
			utils.ExecCmd("git", "init")
		}

		if os.IsExist(err) && isRepo.IsDir() {
			utils.ExecCmd("git", "init")
		}

		isPrivate, _ := cmd.Flags().GetBool("private")
		if isPrivate == true {
			ghPositionalArgs = append(ghPositionalArgs, "--private")
		} else {
			ghPositionalArgs = append(ghPositionalArgs, "--public")
		}

		push, _ := cmd.Flags().GetBool("push")
		if push == true {
			ghPositionalArgs = append(ghPositionalArgs, "--push")
		}

		ghFullArgs = append(ghFullArgs, "--source=.")

		ghFullArgs = append(ghFullArgs, "--remote=origin")

		ghFullArgs = append(ghFullArgs, ghPositionalArgs...)

		confirm, _ := cmd.Flags().GetBool("confirm")

		if !confirm && !utils.ConfirmAction(cmd) {
			fmt.Println("Operation cancelled")
			return
		}

		errCmd := utils.ExecCmd(ghFullArgs[0], ghFullArgs[1:]...)

		if errCmd != nil {
			fmt.Println("Error creating repo")
		}
	},
}
