package main

import (
	"fmt"
	"log"

	"github.com/matejb/talks/2026/gops-2026/layouts-examples/08-monorepo/libs/library"
)

func main() {
	store := library.NewStore()

	// Simulate some borrows
	_ = store.BorrowBook(1, 1)
	_ = store.BorrowBook(2, 2)

	// Worker job: list all borrowed books
	borrowed, err := store.BorrowedBooks()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("=== Borrowed Books Report ===")
	if len(borrowed) == 0 {
		fmt.Println("No books are currently borrowed.")
		return
	}
	for _, b := range borrowed {
		fmt.Printf("- %q by %s (borrowed by %s)\n", b.Title, b.Author, b.Borrowed.Name)
	}
	fmt.Printf("\nTotal: %d book(s) borrowed.\n", len(borrowed))
}
