package utils

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"sync"
)

func ExecCmd(command string, args ...string) error {
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
