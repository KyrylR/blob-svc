package handlers

import (
	"blob-svc/internal/service/requests"
	"fmt"
	"log"
	"net/http"

	"gitlab.com/distributed_lab/ape"
	"gitlab.com/distributed_lab/ape/problems"
)

func CreateAccount(w http.ResponseWriter, r *http.Request) {
	_, err := requests.NewCreateAccountRequest(r)
	if err != nil {
		ape.RenderErr(w, problems.BadRequest(err)...)
		return
	}

	resp, err := http.Get("http://127.0.0.1:80")
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(resp)
}
