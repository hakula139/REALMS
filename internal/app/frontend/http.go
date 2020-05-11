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
	if mode == showMode {
		booksMgrURL += "/books"
		if bookID != 0 {
			booksMgrURL += "/" + strconv.Itoa(bookID)
		}
		return http.Get(booksMgrURL)
	}

	if mode == findMode {
		booksMgrURL += "/books/find"
	} else {
		booksMgrURL += "/admin/books"
		if mode != addMode {
			booksMgrURL += "/" + strconv.Itoa(bookID)
		}
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
	if (mode == showMode && userID != 0) || mode == updateMode || mode == removeMode {
		usersMgrURL += "/" + strconv.Itoa(userID)
	}
	return sendRequest(method, jar, input, usersMgrURL)
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
