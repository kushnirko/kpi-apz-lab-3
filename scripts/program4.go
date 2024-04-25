package main

import (
	"fmt"
	"github.com/roman-mazur/architecture-lab-3/scripts/client"
	"net/http"
	"time"
)

func main() {
	Client := &http.Client{}
	method := "POST"
	url := "http://localhost:17000/"
	cmds := "green\nfigure 0.0 0.0\n"
	delay := 1 * time.Second
	delta := 0.2
	x := 0.0
	y := 0.0

	client.SendMessage(Client, method, url, cmds)

	for x < 1 {
		cmds = fmt.Sprintf("move %f %f\nupdate\n", x, y)
		client.SendMessage(Client, method, url, cmds)
		x += delta
		time.Sleep(delay)
	}

	for y < 1 {
		cmds = fmt.Sprintf("move %f %f\nupdate\n", x, y)
		client.SendMessage(Client, method, url, cmds)
		y += delta
		time.Sleep(delay)
	}

	for x > 0.1 {
		cmds = fmt.Sprintf("move %f %f\nupdate\n", x, y)
		client.SendMessage(Client, method, url, cmds)
		x -= delta
		time.Sleep(delay)
	}

	for y > 0.1 {
		cmds = fmt.Sprintf("move %f %f\nupdate\n", x, y)
		client.SendMessage(Client, method, url, cmds)
		y -= delta
		time.Sleep(delay)
	}

}
