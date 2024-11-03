package cmd

import (
	"fmt"
	"os"

	"ghe/utils"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(rpCmd)

	var Org string
	rpCmd.Flags().StringVar(&Org, "org", "", "Define a organization to create your repo into")

	var IsPrivate bool
	rpCmd.Flags().BoolVar(&IsPrivate, "private", false, "Create repo as private")

	var Push bool
	rpCmd.Flags().BoolVar(&Push, "push", false, "Push to repo")

	var Confirm bool
	rpCmd.Flags().BoolVar(&Confirm, "confirm", false, "Confirm create repo action without prompting")
}

var rpCmd = &cobra.Command{
	Use:   "cr",
	Short: "Create a repo on github",
	Long: `Create a repo with (repo create) from gh but with sane defaults
		- public
		- no push
		- name from current folder
		- 
	`,
	Run: func(cmd *cobra.Command, args []string) {

		// upstream?
		ghBaseArgs := []string{"gh", "repo", "create"}
		ghFullArgs := make([]string, len(ghBaseArgs))

		copy(ghFullArgs, ghBaseArgs)

		org, _ := cmd.Flags().GetString("org")
		if org != "" {
			dirName := utils.GetDirName()
			orgAndProject := org + "/" + dirName
			ghFullArgs = append(ghFullArgs, orgAndProject)
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
		// source
		ghFullArgs = append(ghFullArgs, "--source=.")

		// remote
		ghFullArgs = append(ghFullArgs, "--remote=origin")

		isPrivate, _ := cmd.Flags().GetBool("private")
		if isPrivate == true {
			ghFullArgs = append(ghFullArgs, "--private")
		} else {
			ghFullArgs = append(ghFullArgs, "--public")
		}

		push, _ := cmd.Flags().GetBool("push")
		if push == true {
			ghFullArgs = append(ghFullArgs, "--push")
		}

		confirm, _ := cmd.Flags().GetBool("confirm")
		if !confirm && !utils.ConfirmAction() {
			fmt.Println("Operation cancelled")
			return
		}

		// fmt.Println(ghFullArgs)
		errCmd := utils.ExecCmd(ghFullArgs[0], ghFullArgs[1:]...)

		if errCmd != nil {
			fmt.Println("Error creating repo")
		}
	},
}
