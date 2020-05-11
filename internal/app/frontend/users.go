package frontend

import (
	"bufio"
	"fmt"
	"net/http/cookiejar"
	"os"
	"strings"
	"syscall"

	"golang.org/x/crypto/ssh/terminal"
)

type userModel struct {
	ID       uint   `json:"id,omitempty"`
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
	Level    uint   `json:"level,omitempty"`
}

// AddUser adds a new user to the database
func AddUser(jar *cookiejar.Jar) error {
	var input userModel
	if err := getUserInput(&input, addMode); err != nil {
		return err
	}

	// Sends a POST request
	res, err := sendUserRequest("POST", jar, &input, 0, addMode)
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
		user, ok := dataBody.(map[string]interface{})
		if !ok {
			fmt.Println(ErrInvalidResponse.Error())
			return nil
		}
		fmt.Printf("Successfully added user %v\n", user["id"])
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// UpdateUser updates data of a user
func UpdateUser(jar *cookiejar.Jar) error {
	userID := getUserID()
	var input userModel
	if err := getUserInput(&input, updateMode); err != nil {
		return err
	}

	// Sends a PATCH request
	res, err := sendUserRequest("PATCH", jar, &input, userID, updateMode)
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
		fmt.Printf("Successfully updated user %v\n", userID)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// RemoveUser removes a user from the database
func RemoveUser(jar *cookiejar.Jar) error {
	userID := getUserID()

	// Sends a DELETE request
	res, err := sendUserRequest("DELETE", jar, nil, userID, removeMode)
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
		fmt.Printf("Successfully removed user %v\n", userID)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// ShowUsers shows all users in the library
func ShowUsers(jar *cookiejar.Jar) error {
	// Sends a GET request
	res, err := sendUserRequest("GET", jar, nil, 0, showMode)
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
		users := dataBody.([]interface{})
		printUsers(users)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// ShowUser shows the user of given ID
func ShowUser(jar *cookiejar.Jar) error {
	userID := getUserID()

	// Sends a GET request
	res, err := sendUserRequest("GET", jar, nil, userID, showMode)
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
		switch user := dataBody.(type) {
		case map[string]interface{}:
			printUser(user)
		case []interface{}:
			printUsers(user)
		}
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

func getUserInput(input *userModel, mode int) error {
	scanner := bufio.NewScanner(os.Stdin)

	username := ""
	if mode == addMode {
		fmt.Print("Enter Username: ")
		scanner.Scan()
		username = scanner.Text()
		if username == "" {
			fmt.Println("Username shouldn't be empty")
			return ErrInvalidInput
		}
	}

	fmt.Print("Enter Password: ")
	bytePassword, _ := terminal.ReadPassword(int(syscall.Stdin))
	password := string(bytePassword)
	fmt.Println()
	if password == "" {
		fmt.Println("Password shouldn't be empty")
		return ErrInvalidInput
	}

	fmt.Print("Enter Password again: ")
	bytePasswordConfirm, _ := terminal.ReadPassword(int(syscall.Stdin))
	passwordConfirm := string(bytePasswordConfirm)
	fmt.Println()
	if password != passwordConfirm {
		fmt.Println("Password doesn't match")
		return ErrInvalidInput
	}

	var level uint
	fmt.Println("(1: User, 2: Admin, 3: Super Admin)")
	fmt.Print("Enter Privilege Level: ")
	fmt.Scanln(&level)
	if level < 1 || level > 3 {
		fmt.Println("Level should be an integer between 1 and 3")
		return ErrInvalidInput
	}

	input.Username = strings.TrimSpace(username)
	input.Password = password
	input.Level = level
	return nil
}

func getUserID() int {
	var userID int
	fmt.Print("User ID: ")
	fmt.Scanln(&userID)
	return userID
}

func printUsers(users []interface{}) {
	if len(users) == 0 {
		fmt.Println("No users found")
		return
	}
	width := 25
	fmt.Printf("ID\t%-*s%-s\n",
		width, "Username",
		"Level",
	)
	fmt.Println(strings.Repeat("-", 13+width))
	for _, elem := range users {
		user := elem.(map[string]interface{})
		fmt.Printf("%v\t", user["id"])
		fmt.Printf("%-*s", width, slice(user["username"].(string), width-2))
		fmt.Printf("%v\n", user["level"])
	}
}

func printUser(user map[string]interface{}) {
	fmt.Printf("User %v\n", user["id"])
	fmt.Printf("   Username: %v\n", user["username"])
	fmt.Printf("   Level:    %v\n", user["level"])
}
