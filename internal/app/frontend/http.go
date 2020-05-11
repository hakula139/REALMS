package frontend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/cookiejar"
	"strconv"
)

func sendBookRequest(
	method string,
	jar *cookiejar.Jar,
	input interface{},
	bookID int,
	mode int,
) (res *http.Response, err error) {
	booksMgrURL := URL
	switch mode {
	case addMode:
		booksMgrURL += "/admin/books"
	case updateMode:
		fallthrough
	case removeMode:
		booksMgrURL += "/admin/books/" + strconv.Itoa(bookID)
	case showMode:
		booksMgrURL += "/books"
		if bookID != 0 {
			booksMgrURL += "/" + strconv.Itoa(bookID)
		}
		return http.Get(booksMgrURL)
	case findMode:
		booksMgrURL += "/books/find"
	}
	return sendRequest(method, jar, input, booksMgrURL)
}

func sendUserRequest(
	method string,
	jar *cookiejar.Jar,
	input interface{},
	userID int,
	mode int,
) (res *http.Response, err error) {
	usersMgrURL := URL + "/admin/users"
	switch mode {
	case addMode:
		// Does nothing
	case updateMode:
		fallthrough
	case removeMode:
		usersMgrURL += "/" + strconv.Itoa(userID)
	case showMode:
		if userID != 0 {
			usersMgrURL += "/" + strconv.Itoa(userID)
		}
	}
	return sendRequest(method, jar, input, usersMgrURL)
}

func sendRecordRequest(
	method string,
	jar *cookiejar.Jar,
	input interface{},
	bookID int,
	mode int,
) (res *http.Response, err error) {
	recordsMgrURL := URL + "/user"
	switch mode {
	case addMode:
		fallthrough
	case updateMode:
		fallthrough
	case removeMode:
		recordsMgrURL += "/books/" + strconv.Itoa(bookID)
	case showMode:
		recordsMgrURL += "/books"
		if bookID != 0 {
			recordsMgrURL += "/" + strconv.Itoa(bookID)
		}
	case showOverdueMode:
		recordsMgrURL += "/overdue"
	case showHistoryMode:
		recordsMgrURL += "/history"
	}
	return sendRequest(method, jar, input, recordsMgrURL)
}

func sendRequest(
	method string,
	jar *cookiejar.Jar,
	input interface{},
	url string,
) (res *http.Response, err error) {
	client := &http.Client{
		Jar: jar,
	}
	jsonStr, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return client.Do(req)
}

func readResponse(res *http.Response) (data map[string]interface{}, err error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(ErrReadResponseFailed.Error())
		return nil, err
	}
	if err := json.Unmarshal(body, &data); err != nil {
		fmt.Println(ErrInvalidResponse.Error())
		return nil, err
	}
	return data, nil
}
