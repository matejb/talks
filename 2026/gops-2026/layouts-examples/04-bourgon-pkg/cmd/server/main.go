package main

import (
	"log"
	"net/http"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/04-bourgon-pkg/pkg/httphandler"
	"github.com/matejb/talks/2026/gops-2026/layouts-examples/04-bourgon-pkg/pkg/memory"
)

func main() {
	us := memory.NewUserService()
	bs := memory.NewBookService(us)
	rs := memory.NewReviewService()

	h := &httphandler.Handler{
		Users:   us,
		Books:   bs,
		Reviews: rs,
	}

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Bourgon pkg/ layout example running on :8084")
	log.Fatal(http.ListenAndServe(":8084", mux))
}
