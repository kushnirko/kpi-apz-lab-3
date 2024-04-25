package main

import (
	"github.com/roman-mazur/architecture-lab-3/scripts/client"
	"net/http"
)

func main() {
	Client := &http.Client{}
	method := "POST"
	url := "http://localhost:17000/"
	cmds := "white\nbgrect 0.25 0.25 0.75 0.75\nfigure 0.5 0.5\ngreen\nfigure 0.6 0.6\nupdate\n"

	client.SendMessage(Client, method, url, cmds)
}
