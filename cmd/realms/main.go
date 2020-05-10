package main

import (
	"bufio"
	"fmt"
	"net/http/cookiejar"
	"os"

	"github.com/hakula139/REALMS/internal/app/frontend"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "REALMS",
		Version: "v0.1.0",
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Hakula Chen",
				Email: "i@hakula.xyz",
			},
		},
		Usage:  "REALMS Establishes A Library Management System",
		Action: router,
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("realms: failed to start: " + err.Error())
	}
}

func router(c *cli.Context) error {
	fmt.Println("Welcome to REALMS! Check the manual using the command 'help'.")
	jar, _ := cookiejar.New(nil)
	for {
		fmt.Print("> ")
		scanner := bufio.NewScanner(os.Stdin)
		scanner.Scan()
		op := scanner.Text()
		switch op {
		case "":
			// Does nothing
		case "help":
			frontend.ShowHelp()
		case "login":
			if err := frontend.Login(jar); err != nil {
				fmt.Println(err.Error())
			}
		case "logout":
			if err := frontend.Logout(jar); err != nil {
				fmt.Println(err.Error())
			}
		case "me":
			if err := frontend.Me(jar); err != nil {
				fmt.Println(err.Error())
			}
		case "status":
			if err := frontend.Status(jar); err != nil {
				fmt.Println(err.Error())
			}
		case "exit":
			fmt.Println("Bye!")
			return nil
		default:
			fmt.Println("Invalid operation! Check the manual using the command 'help'.")
		}
	}
}
