package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// return error but also return success output
func ExecGhCmd(command string, args ...string) error {
	cmd := exec.Command(command, args...) //.Output()

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
		return fmt.Errorf("Failed on stderr pipe: %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Command Failed to start: %v", err)
	}
	scanner := bufio.NewScanner(stderr)
	for scanner.Scan() {
		fmt.Println(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
		return fmt.Errorf("error reading stderr: %v", err)
	}

	scanner2 := bufio.NewScanner(stdout)
	for scanner2.Scan() {
		fmt.Println(scanner2.Text())
	}

	if err := scanner2.Err(); err != nil {
		log.Fatal(err)
		return fmt.Errorf("error reading stdout: %v", err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
		return fmt.Errorf("command failed: %v", err)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(rpCmd)

	var Push bool
	rpCmd.Flags().BoolVar(&Push, "push", false, "Push to repo")
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
		fmt.Println("Start create repo command")

		currentDir, err := os.Getwd()

		if err != nil {
			fmt.Printf("Error grabing the path to dir")
		}
		dirName := filepath.Base(currentDir)

		// make sure dirName is valid for github
		if strings.Contains(dirName, " ") {
			dirName = strings.ReplaceAll(dirName, " ", "-")
		}

		ghBaseArgs := []string{"gh", "repo", "create", dirName, "--private", "--source=.", "--remote=upstream"}
		ghFullArgs := make([]string, len(ghBaseArgs))

		copy(ghFullArgs, ghBaseArgs)
		push, _ := cmd.Flags().GetBool("push")
		if push == true {
			ghFullArgs = append(ghFullArgs, "--push")
		}
		fmt.Println(ghFullArgs)
		errCmd := ExecGhCmd(ghFullArgs[0], ghFullArgs[1:]...)

		if errCmd != nil {
			fmt.Println("Error creating repo")
		}
	},
}
