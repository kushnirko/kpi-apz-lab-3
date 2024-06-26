package main

import (
	"fmt"
	"kpi-apz-lab-3/scripts/client"
	"net/http"
	"time"
)

func main() {
	Client := &http.Client{}
	method := "POST"
	url := "http://localhost:17000/"
	cmds := "reset\nwhite\nfigure 0.0 0.0\n"
	delay := 1 * time.Second
	delta := 0.1
	x := 0.0

	client.SendMessage(Client, method, url, cmds)

	for x <= 1 {
		cmds = fmt.Sprintf("move %f %f\nupdate\n", x, x)
		client.SendMessage(Client, method, url, cmds)
		x += delta
		time.Sleep(delay)
	}

}
