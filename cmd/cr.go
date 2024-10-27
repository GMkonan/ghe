package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

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

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		scannerForStderr := bufio.NewScanner(stderr)
		for scannerForStderr.Scan() {
			fmt.Println(scannerForStderr.Text())
		}

		if err := scannerForStderr.Err(); err != nil {
			log.Printf("error reading stderr: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		scannerForStdout := bufio.NewScanner(stdout)
		for scannerForStdout.Scan() {
			fmt.Println(scannerForStdout.Text())
		}

		if err := scannerForStdout.Err(); err != nil {
			log.Fatal(err)
			log.Printf("error reading stdout: %v", err)
		}
	}()

	wg.Wait()

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
		return fmt.Errorf("command failed: %v", err)
	}

	return nil
}

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

func getDirName() string {
	currentDir, err := os.Getwd()

	if err != nil {
		fmt.Printf("Error grabing the path to dir")
	}
	dirName := filepath.Base(currentDir)

	// make sure dirName is valid for github
	if strings.Contains(dirName, " ") {
		dirName = strings.ReplaceAll(dirName, " ", "-")
	}
	return dirName
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
			dirName := getDirName()
			orgAndProject := org + "/" + dirName
			ghFullArgs = append(ghFullArgs, orgAndProject)
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

		// fmt.Println(ghFullArgs)
		errCmd := ExecGhCmd(ghFullArgs[0], ghFullArgs[1:]...)

		if errCmd != nil {
			fmt.Println("Error creating repo")
		}
	},
}
