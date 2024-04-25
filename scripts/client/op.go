package client

import (
	"fmt"
	"net/http"
	"strings"
)

func SendMessage(client *http.Client, method string, url string, message string) {
	in := strings.NewReader(message)
	req, err := http.NewRequest(method, url, in)
	_, err = client.Do(req)

	if err != nil {
		fmt.Println("Bad request", err)
		return
	}
}
