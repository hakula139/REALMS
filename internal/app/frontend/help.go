package frontend

import (
	"fmt"
)

// ShowHelp shows a list of commands
func ShowHelp() error {
	fmt.Println("COMMANDS:")
	printCommand("help", "Shows a list of commands")
	printCommand("login", "Log in to your library account")
	printCommand("logout", "Log out of your library account")
	printCommand("me", "Shows the current logged-in user")
	printCommand("status", "Shows the current login status")
	printCommand("exit", "Quit")
	return nil
}

func printCommand(cmd, usage string) {
	fmt.Printf("   %-15s%s\n", cmd, usage)
}
