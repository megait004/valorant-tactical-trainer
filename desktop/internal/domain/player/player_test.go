package player

import "testing"

func TestNormalizeTag(t *testing.T) {
	t.Parallel()

	got := NormalizeTag(" #VN2 ")
	if got != "VN2" {
		t.Fatalf("expected VN2, got %q", got)
	}
}

func TestNormalizeRegionDefaultsToAP(t *testing.T) {
	t.Parallel()

	got := NormalizeRegion(" ")
	if got != "ap" {
		t.Fatalf("expected ap, got %q", got)
	}
}

func TestIsValidLookup(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		tag  string
		want bool
	}{
		{name: "player", tag: "VN2", want: true},
		{name: " player ", tag: " #VN2 ", want: true},
		{name: "", tag: "VN2", want: false},
		{name: "player", tag: "", want: false},
	}

	for _, test := range tests {
		got := IsValidLookup(test.name, test.tag)
		if got != test.want {
			t.Fatalf("IsValidLookup(%q, %q) = %v, want %v", test.name, test.tag, got, test.want)
		}
	}
}
