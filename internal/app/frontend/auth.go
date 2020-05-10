package frontend

import (
	"bufio"
	"fmt"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

// Login helps the user log in to his/her library account
func Login(jar *cookiejar.Jar) error {
	// Uses cookiejar to manage the cookies
	client := &http.Client{
		Jar: jar,
	}

	// Sends a POST request
	username, password := getCredentials()
	loginURL := URL + "/login"
	res, err := client.PostForm(loginURL, url.Values{
		"username": {username},
		"password": {password},
	})
	if err != nil {
		fmt.Println(ErrRequestFailed.Error())
		return err
	}
	defer res.Body.Close()

	// Outputs the response
	data, err := readResponse(res)
	if err != nil {
		return err
	}
	if _, ok := data["data"]; ok {
		fmt.Printf("Welcome %v!\n", username)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// Logout helps the user log out of his/her library account
func Logout(jar *cookiejar.Jar) error {
	// Uses cookiejar to manage the cookies
	client := &http.Client{
		Jar: jar,
	}

	// Sends a GET request
	logoutURL := URL + "/logout"
	res, err := client.Get(logoutURL)
	if err != nil {
		fmt.Println(ErrRequestFailed.Error())
		return err
	}
	defer res.Body.Close()

	// Outputs the response
	data, err := readResponse(res)
	if err != nil {
		return err
	}
	if _, ok := data["data"]; ok {
		fmt.Println("Successfully logged out!")
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// Me shows the currently logged-in user
func Me(jar *cookiejar.Jar) error {
	// Uses cookiejar to manage the cookies
	client := &http.Client{
		Jar: jar,
	}

	// Sends a GET request
	meURL := URL + "/user/me"
	res, err := client.Get(meURL)
	if err != nil {
		fmt.Println(ErrRequestFailed.Error())
		return err
	}
	defer res.Body.Close()

	// Outputs the response
	data, err := readResponse(res)
	if err != nil {
		return err
	}
	if dataBody, ok := data["data"]; ok {
		fmt.Printf("Current user ID: %v\n", dataBody)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// Status shows current login status
func Status(jar *cookiejar.Jar) error {
	// Uses cookiejar to manage the cookies
	client := &http.Client{
		Jar: jar,
	}

	// Sends a GET request
	statURL := URL + "/status"
	res, err := client.Get(statURL)
	if err != nil {
		fmt.Println(ErrRequestFailed.Error())
		return err
	}
	defer res.Body.Close()

	// Outputs the response
	data, err := readResponse(res)
	if err != nil {
		return err
	}
	if dataBody, ok := data["data"]; ok {
		stat, _ := dataBody.(bool)
		if stat {
			fmt.Println("Online")
		} else {
			fmt.Println("Offline")
		}
	}
	return nil
}

func getCredentials() (string, string) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter Username: ")
	scanner.Scan()
	username := scanner.Text()

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println()

	return strings.TrimSpace(username), password
}
