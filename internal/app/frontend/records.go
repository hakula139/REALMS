package frontend

import (
	"bufio"
	"fmt"
	"net/http/cookiejar"
	"os"
	"strings"
	"time"
)

const maxExtendTimes = 3

// recordInput should only be used in debug mode
type recordInput struct {
	BorrowDate time.Time `json:"borrow_date,omitempty"`
}

// BorrowBook borrows a book from the library
func BorrowBook(jar *cookiejar.Jar) error {
	bookID := getBookID()
	var input recordInput
	if debugMode {
		if err := getRecordInput(&input); err != nil {
			return err
		}
	}

	// Sends a POST request
	res, err := sendRecordRequest("POST", jar, &input, bookID, addMode)
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
		record, ok := dataBody.(map[string]interface{})
		if !ok {
			fmt.Println(ErrInvalidResponse.Error())
			return nil
		}
		fmt.Printf("Successfully borrowed book %v\n", record["book_id"])
		fmt.Printf("Your return date is: %v\n", record["return_date"])
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// ExtendDeadline extends the deadline to return a book
func ExtendDeadline(jar *cookiejar.Jar) error {
	bookID := getBookID()

	// Sends a PATCH request
	res, err := sendRecordRequest("PATCH", jar, nil, bookID, updateMode)
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
		record, ok := dataBody.(map[string]interface{})
		if !ok {
			fmt.Println(ErrInvalidResponse.Error())
			return nil
		}
		fmt.Printf("Successfully extended the return date of book %v\n", record["book_id"])
		fmt.Printf("Your return date is: %v\n", record["return_date"])
		fmt.Printf("You have extended %v/%v times\n", record["extend_times"], maxExtendTimes)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// ReturnBook returns a book to the library
func ReturnBook(jar *cookiejar.Jar) error {
	bookID := getBookID()

	// Sends a DELETE request
	res, err := sendRecordRequest("DELETE", jar, nil, bookID, removeMode)
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
		fmt.Printf("Successfully returned book %v\n", bookID)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

// ShowBookList shows all books that the user has borrowed
func ShowBookList(jar *cookiejar.Jar) error {
	return showBookList(jar, showMode)
}

// ShowOverdueList shows all overdue books that the user has borrowed
func ShowOverdueList(jar *cookiejar.Jar) error {
	return showBookList(jar, showOverdueMode)
}

// ShowHistory shows all records of the user
func ShowHistory(jar *cookiejar.Jar) error {
	return showBookList(jar, showHistoryMode)
}

// ShowBorrowed shows a book that the user has borrowed
func ShowBorrowed(jar *cookiejar.Jar) error {
	bookID := getBookID()

	// Sends a GET request
	res, err := sendRecordRequest("GET", jar, nil, bookID, showMode)
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
		switch record := dataBody.(type) {
		case map[string]interface{}:
			printRecord(record)
		case []interface{}:
			printRecords(record, showMode)
		}
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

func showBookList(jar *cookiejar.Jar, mode int) error {
	// Sends a GET request
	res, err := sendRecordRequest("GET", jar, nil, 0, mode)
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
		records := dataBody.([]interface{})
		printRecords(records, mode)
	} else if errBody, ok := data["error"]; ok {
		fmt.Println(errBody)
	}
	return nil
}

func getRecordInput(input *recordInput) error {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("(Format: yyyy-mm-dd)")
	fmt.Print("Borrow date: ")
	scanner.Scan()
	dateStr := scanner.Text()

	fmt.Println("(Format: hh:mm:ss)")
	fmt.Print("Borrow time: ")
	scanner.Scan()
	timeStr := scanner.Text()

	borrowDateStr := dateStr + "T" + timeStr + "Z"
	borrowDate, err := time.Parse(time.RFC3339, borrowDateStr)
	if err != nil {
		fmt.Println("Wrong format")
		return ErrInvalidInput
	}
	input.BorrowDate = borrowDate
	return nil
}

func printRecords(records []interface{}, mode int) {
	if len(records) == 0 {
		fmt.Println("No records found")
		return
	}
	fmt.Printf("%s\t%s\t  %s\t\t%s",
		"ID",
		"Book ID",
		"Return Date",
		"Extended",
	)
	if mode == showHistoryMode {
		fmt.Printf("  %s\n",
			"Returned Date",
		)
		fmt.Println(strings.Repeat("-", 70))
	} else {
		fmt.Println()
		fmt.Println(strings.Repeat("-", 48))
	}
	for _, elem := range records {
		record := elem.(map[string]interface{})
		fmt.Printf("%v\t", record["id"])
		fmt.Printf("%v\t  ", record["book_id"])
		fmt.Printf("%v\t", record["return_date"])
		fmt.Printf("%v/%v\t  ", record["extend_times"], maxExtendTimes)
		if mode == showHistoryMode {
			if returnedDate := record["real_return_date"]; returnedDate != nil {
				fmt.Println(returnedDate)
			} else {
				fmt.Println("N/A")
			}
		} else {
			fmt.Println()
		}
	}
}

func printRecord(record map[string]interface{}) {
	fmt.Printf("Record %v\n", record["id"])
	fmt.Printf("   Book ID:     %v\n", record["book_id"])
	fmt.Printf("   Return Date: %v\n", record["return_date"])
	fmt.Printf("   Extended:    %v/%v\n", record["extend_times"], maxExtendTimes)
}
