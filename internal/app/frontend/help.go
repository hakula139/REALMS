package frontend

import (
	"fmt"
)

// ShowHelp shows a list of commands
func ShowHelp() error {
	fmt.Println("COMMANDS:")
	printCommand("help", "Shows a list of commands\n")
	printCommand("login", "Log in to your library account")
	printCommand("logout", "Log out of your library account")
	printCommand("me", "Shows the current logged-in user")
	printCommand("status", "Shows the current login status\n")
	printCommand("add book", "Adds a new book to the library")
	printCommand("update book", "Updates data of a book")
	printCommand("remove book", "Removes a book from the library")
	printCommand("show books", "Shows all books in the library")
	printCommand("show book", "Shows the book of given ID")
	printCommand("find books", "Finds books by title / author / ISBN\n")
	printCommand("add user", "Adds a new user to the database")
	printCommand("update user", "Updates data of a user")
	printCommand("remove user", "Removes a user from the database")
	printCommand("show users", "Shows all users in the library")
	printCommand("show user", "Shows the user of given ID\n")
	printCommand("exit", "Quit")
	return nil
}

func printCommand(cmd, usage string) {
	fmt.Printf("   %-15s%s\n", cmd, usage)
}
