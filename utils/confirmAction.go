package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ConfirmAction() bool {
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
