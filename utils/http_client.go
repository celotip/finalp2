package utils

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

func RequestGET(url string, headers map[string]string) ([]byte, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		log.Printf("Error: %s", body)
		return nil, fmt.Errorf("Failed to get data, status code: %d", resp.StatusCode)
	}

	return body, nil

}
