package main

import (
	"log"
	"net/http"
)

func main() {
	us := newUserService()
	bs := newBookService(us)
	rs := newReviewService()

	h := &handler{Users: us, Books: bs, Reviews: rs}

	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	log.Println("Flat layout example running on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
