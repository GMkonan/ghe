package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

func ConfirmAction(cmd *cobra.Command) bool {
	reader := bufio.NewReader(os.Stdin)
	listOfConfirmArgs := make([]string, 0)

	visibility, _ := cmd.Flags().GetBool("private")
	if visibility {
		listOfConfirmArgs = append(listOfConfirmArgs, "Private")
	} else {
		listOfConfirmArgs = append(listOfConfirmArgs, "Public")
	}
	push, _ := cmd.Flags().GetBool("push")
	if push {
		listOfConfirmArgs = append(listOfConfirmArgs, "Yes")
	} else {
		listOfConfirmArgs = append(listOfConfirmArgs, "No")
	}

	for {
		// fmt.Print("Do you want to proceed? [y/N]: ")

		confirmationText := fmt.Sprintf(`
Confirm the options:
Visibility: %s
Push local commits: %s
Do you want to proceed? [y/N]
`, listOfConfirmArgs[0], listOfConfirmArgs[1])
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
