package api

import (
	"errors"
	"fmt"
	"io"
	"net/http"
)

func GetLocations(url string) ([]byte, error) {
	res, err := http.Get(url)
	if err != nil {
		return make([]byte, 0), err
	}

	if res.StatusCode > 299 {
		return make([]byte, 0), errors.New(fmt.Sprintln("Failed with the status code: %d", res.StatusCode))
	}

	body, err1 := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err1 != nil {
		return body, err1
	}
	return body, nil
}
