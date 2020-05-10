package frontend

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	dataMsg, errMsg, err := readResponse(res)
	if err != nil {
		return err
	}
	if dataMsg != nil {
		fmt.Printf("Welcome %v!\n", username)
	} else if errMsg != nil {
		fmt.Println(errMsg)
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
	dataMsg, errMsg, err := readResponse(res)
	if err != nil {
		return err
	}
	if dataMsg != nil {
		fmt.Println("Successfully logged out!")
	} else if errMsg != nil {
		fmt.Println(errMsg)
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
	dataMsg, errMsg, err := readResponse(res)
	if err != nil {
		return err
	}
	if dataMsg != nil {
		fmt.Printf("Current user ID: %v\n", dataMsg)
	} else if errMsg != nil {
		fmt.Println(errMsg)
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
	dataMsg, _, err := readResponse(res)
	if err != nil {
		return err
	}
	if dataMsg != nil {
		stat, _ := dataMsg.(bool)
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

func readResponse(res *http.Response) (dataMsg, errMsg interface{}, err error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(ErrReadResponseFailed.Error())
		return nil, nil, err
	}
	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(ErrInvalidResponse.Error())
		return nil, nil, err
	}
	if dataMsg, ok := data["data"]; ok {
		return dataMsg, nil, nil
	} else if errMsg, ok := data["error"]; ok {
		return nil, errMsg, nil
	}
	return nil, nil, nil
}
