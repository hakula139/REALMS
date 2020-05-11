package frontend

import (
	"bufio"
	"fmt"
	"net/http/cookiejar"
	"os"
	"strings"
)

type bookModel struct {
	ID        uint   `json:"id,omitempty"`
	Title     string `json:"title,omitempty"`
	Author    string `json:"author,omitempty"`
	Publisher string `json:"publisher,omitempty"`
	ISBN      string `json:"isbn,omitempty"`
}

type messageInput struct {
	Message string `json:"message,omitempty"`
}

// AddBook adds a new book to the library
func AddBook(jar *cookiejar.Jar) error {
	var input bookModel
	if err := getBookInput(&input, addMode); err != nil {
		return err
	}

	// Sends a POST request
	res, err := sendBookRequest("POST", jar, &input, 0, addMode)
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
		book, ok := dataBody.(map[string]interface{})
		if !ok {
			fmt.Println(ErrInvalidResponse.Error())
			return nil
		}
		fmt.Printf("Successfully added book %v\n", book["id"])
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// UpdateBook updates data of a book
func UpdateBook(jar *cookiejar.Jar) error {
	bookID := getBookID()
	var input bookModel
	if err := getBookInput(&input, updateMode); err != nil {
		return err
	}

	// Sends a PATCH request
	res, err := sendBookRequest("PATCH", jar, &input, bookID, updateMode)
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
		fmt.Printf("Successfully updated book %v\n", bookID)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// RemoveBook removes a book from the library
func RemoveBook(jar *cookiejar.Jar) error {
	bookID := getBookID()
	var input messageInput
	if err := getMessageInput(&input); err != nil {
		return err
	}

	// Sends a DELETE request
	res, err := sendBookRequest("DELETE", jar, &input, bookID, removeMode)
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
		fmt.Printf("Successfully removed book %v\n", bookID)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// ShowBooks shows all books in the library
func ShowBooks() error {
	// Sends a GET request
	res, err := sendBookRequest("GET", nil, nil, 0, showMode)
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
		books := dataBody.([]interface{})
		printBooks(books)
	}
	return nil
}

// ShowBook shows the book of given ID
func ShowBook() error {
	bookID := getBookID()

	// Sends a GET request
	res, err := sendBookRequest("GET", nil, nil, bookID, showMode)
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
		switch book := dataBody.(type) {
		case map[string]interface{}:
			printBook(book)
		case []interface{}:
			printBooks(book)
		}
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// FindBooks finds books by title / author / ISBN
func FindBooks(jar *cookiejar.Jar) error {
	var input bookModel
	if err := getBookInput(&input, findMode); err != nil {
		return err
	}

	// Sends a POST request
	res, err := sendBookRequest("POST", jar, &input, 0, findMode)
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
		books := dataBody.([]interface{})
		printBooks(books)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

func getBookInput(input *bookModel, mode int) error {
	scanner := bufio.NewScanner(os.Stdin)

	if mode == addMode {
		fmt.Print("Title (required): ")
	} else {
		fmt.Print("Title (optional): ")
	}
	scanner.Scan()
	title := scanner.Text()
	if mode == addMode && title == "" {
		fmt.Println("The book must have a title!")
		return ErrInvalidInput
	}
	input.Title = title

	fmt.Print("Author (optional): ")
	scanner.Scan()
	input.Author = scanner.Text()

	if mode != findMode {
		fmt.Print("Publisher (optional): ")
		scanner.Scan()
		input.Publisher = scanner.Text()
	}

	fmt.Print("ISBN (optional): ")
	scanner.Scan()
	input.ISBN = scanner.Text()

	return nil
}

func getMessageInput(input *messageInput) error {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Explanation (optional): ")
	scanner.Scan()
	input.Message = scanner.Text()
	return nil
}

func getBookID() int {
	var bookID int
	fmt.Print("Book ID: ")
	fmt.Scanln(&bookID)
	return bookID
}

func printBooks(books []interface{}) {
	if len(books) == 0 {
		fmt.Println("No books found")
		return
	}
	width := 25
	fmt.Printf("%s\t%-*s%-*s%-*s%-*s\n",
		"ID",
		width, "Title",
		width, "Author",
		width, "Publisher",
		width, "ISBN",
	)
	fmt.Println(strings.Repeat("-", 8+width*4))
	for _, elem := range books {
		book := elem.(map[string]interface{})
		fmt.Printf("%v\t", book["id"])
		fmt.Printf("%-*s", width, slice(book["title"].(string), width-2))
		fmt.Printf("%-*s", width, slice(book["author"].(string), width-2))
		fmt.Printf("%-*s", width, slice(book["publisher"].(string), width-2))
		fmt.Printf("%-*s\n", width, slice(book["isbn"].(string), width-2))
	}
}

func printBook(book map[string]interface{}) {
	fmt.Printf("Book %v\n", book["id"])
	fmt.Printf("   Title:     %v\n", book["title"])
	fmt.Printf("   Author:    %v\n", book["author"])
	fmt.Printf("   Publisher: %v\n", book["publisher"])
	fmt.Printf("   ISBN:      %v\n", book["isbn"])
}
