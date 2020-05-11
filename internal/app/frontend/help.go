package frontend

import (
	"fmt"
	"strings"
)

// ShowHelp shows a list of commands
func ShowHelp() error {
	fmt.Println("COMMANDS:")

	printRequiredPrivilege("public")
	printCommand("help", "Shows a list of commands")
	printCommand("exit", "Quit")
	fmt.Println()
	printCommand("login", "Log in to your library account")
	printCommand("logout", "Log out of your library account")
	printCommand("status", "Shows the current login status")
	fmt.Println()
	printCommand("show books", "Shows all books in the library")
	printCommand("show book", "Shows the book of given ID")
	printCommand("find books", "Finds books by title / author / ISBN")
	fmt.Println()

	printRequiredPrivilege("admin")
	printCommand("add book", "Adds a new book to the library")
	printCommand("update book", "Updates data of a book")
	printCommand("remove book", "Removes a book from the library")
	fmt.Println()
	printCommand("add user", "Adds a new user to the database")
	printCommand("update user", "Updates data of a user")
	printCommand("remove user", "Removes a user from the database")
	printCommand("show users", "Shows all users in the library")
	printCommand("show user", "Shows the user of given ID")
	fmt.Println()

	printRequiredPrivilege("user")
	printCommand("me", "Shows the current logged-in user")
	fmt.Println()
	printCommand("borrow book", "Borrows a book from the library")
	printCommand("return book", "Returns a book to the library")
	printCommand("check ddl", "Checks the deadline to return a book")
	printCommand("extend ddl", "Extends the deadline to return a book")
	printCommand("show list", "Shows all books that you've borrowed")
	printCommand("show overdue", "Shows all overdue books that you've borrowed")
	printCommand("show history", "Shows all records")

	return nil
}

func printRequiredPrivilege(level string) {
	indent := 3
	fmt.Print(strings.Repeat(" ", indent))
	switch level {
	case "user":
		fmt.Println("User privilege required:")
	case "admin":
		fmt.Println("Admin privilege required:")
	case "root":
		fmt.Println("Super Admin privilege required:")
	default:
		fmt.Println("Public:")
	}
}

func printCommand(cmd, usage string) {
	indent := 6
	fmt.Print(strings.Repeat(" ", indent))
	fmt.Printf("%-15s%s\n", cmd, usage)
}
