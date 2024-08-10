package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

type response struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

type Config struct {
	Next     string
	Previous string
}

func GetLocations(isNew bool, config *Config) error {
	var url string
	if isNew {
		url = (*config).Next
	} else {
		url = (*config).Previous
	}
	response := response{}

	if url == "" {
		return nil
	}

	res, err := http.Get(url)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()

	err1 := json.Unmarshal(body, &response)
	if err1 != nil {
		return err1
	}

	if res.StatusCode > 299 {
		return errors.New(fmt.Sprintln("Failed with the status code: %d", res.StatusCode))
	}

	for _, result := range response.Results {
		fmt.Println(result.Name)
	}

	(*config).Next = response.Next
	(*config).Previous = response.Previous
	return err
}
