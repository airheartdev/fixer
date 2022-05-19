package fixer

import (
	"testing"
)

func TestCurrenciesString(t *testing.T) {
	for _, tt := range []struct {
		cs   Currencies
		want string
	}{
		{nil, ""},
		{Currencies{}, ""},
		{Currencies{SEK, DKK}, "DKK,SEK"},
		{Currencies{USD, AUD, EUR}, "AUD,EUR,USD"},
	} {
		if got := tt.cs.String(); got != tt.want {
			t.Fatalf("(Currencies{%s}).String() = %q, want %q", tt.cs, got, tt.want)
		}
	}
}
