package main

import (
	"kpi-apz-lab-3/scripts/client"
	"net/http"
	"time"
)

func main() {
	Client := &http.Client{}
	method := "POST"
	url := "http://localhost:17000/"
	delay := 2 * time.Second

	client.SendMessage(Client, method, url, "figure 0.0 0.0\n")
	time.Sleep(delay)

	client.SendMessage(Client, method, url, "figure 0.0 1.0\n")
	time.Sleep(delay)

	client.SendMessage(Client, method, url, "figure 1.0 0.0\n")
	time.Sleep(delay)

	client.SendMessage(Client, method, url, "figure 1.0 1.0\n")

}
