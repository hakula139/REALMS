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
	client := &http.Client{
		Jar: jar,
	}
	jsonStr, err := json.Marshal(input)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, booksMgrURL, bytes.NewBuffer(jsonStr))
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
