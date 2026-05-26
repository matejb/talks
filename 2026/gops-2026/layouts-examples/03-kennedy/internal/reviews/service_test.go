package reviews

import "testing"

func TestService_CreateAndList(t *testing.T) {
	svc := NewService()

	// No reviews initially.
	reviews, err := svc.Reviews(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 0 {
		t.Fatalf("expected 0 reviews, got %d", len(reviews))
	}

	// Create a review.
	err = svc.CreateReview(Review{BookID: 1, UserID: 1, Text: "Great book!", Rating: 5})
	if err != nil {
		t.Fatal(err)
	}

	// List reviews for book 1.
	reviews, err = svc.Reviews(1)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 1 {
		t.Fatalf("expected 1 review, got %d", len(reviews))
	}
	if reviews[0].Text != "Great book!" {
		t.Fatalf("expected 'Great book!', got '%s'", reviews[0].Text)
	}
	if reviews[0].Rating != 5 {
		t.Fatalf("expected rating 5, got %d", reviews[0].Rating)
	}
	if reviews[0].ID != 1 {
		t.Fatalf("expected ID 1, got %d", reviews[0].ID)
	}

	// Book 2 has no reviews.
	reviews, err = svc.Reviews(2)
	if err != nil {
		t.Fatal(err)
	}
	if len(reviews) != 0 {
		t.Fatalf("expected 0 reviews for book 2, got %d", len(reviews))
	}

	// Create another review for the same book.
	err = svc.CreateReview(Review{BookID: 1, UserID: 2, Text: "Decent", Rating: 3})
	if err != nil {
		t.Fatal(err)
	}

	reviews, _ = svc.Reviews(1)
	if len(reviews) != 2 {
		t.Fatalf("expected 2 reviews, got %d", len(reviews))
	}
}
