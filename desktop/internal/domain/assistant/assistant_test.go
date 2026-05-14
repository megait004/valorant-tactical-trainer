package assistant

import "testing"

func TestRecommendEconomy(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		query    Query
		expected string
	}{
		{name: "eco when low credits", query: Query{Credits: 1200}, expected: "Eco"},
		{name: "half buy after loss", query: Query{Credits: 2500, PreviousOutcome: "loss"}, expected: "Light / Half Buy"},
		{name: "full buy above threshold", query: Query{Credits: 4400}, expected: "Full Buy"},
		{name: "force buy middle credits", query: Query{Credits: 3200, PreviousOutcome: "win"}, expected: "Force Buy"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			advice := RecommendEconomy(test.query)
			if advice.Plan != test.expected {
				t.Fatalf("expected %s, got %s", test.expected, advice.Plan)
			}
		})
	}
}

func TestSeedCardsHaveSafetyNotes(t *testing.T) {
	t.Parallel()

	cards := SeedCards()
	if len(cards) == 0 {
		t.Fatal("expected seed cards")
	}
	for _, card := range cards {
		if card.ID == "" || card.Title == "" || card.SafetyNotes == "" {
			t.Fatalf("expected complete safe card: %+v", card)
		}
	}
}
