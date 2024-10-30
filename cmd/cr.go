package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"ghe/utils"

	"github.com/spf13/cobra"
)

func confirmAction() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		// fmt.Print("Do you want to proceed? [y/N]: ")
		confirmationText := fmt.Sprintf(`Confirm the options:
			Visibility: Public
			Push local commits: Yes
			Do you want to proceed? [y/N]
		`)
		fmt.Print(confirmationText)

		response, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			return false
		}

		response = strings.ToLower(strings.TrimSpace(response))

		switch response {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			fmt.Println("Please answer 'y' or 'n'")
		}
	}
}

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

		// if true no need to do anything since it's already a repo
		if err != nil && !isRepo.IsDir() {
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
		if !confirm && !confirmAction() {
			fmt.Println("Operation cancelled")
			return
		}

		fmt.Println(ghFullArgs)
		// errCmd := utils.ExecCmd(ghFullArgs[0], ghFullArgs[1:]...)
		//
		// if errCmd != nil {
		// 	fmt.Println("Error creating repo")
		// }
	},
}
