package main

import "testing"

func TestFilterProfanity(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no profanity",
			input: "I had something interesting for breakfast",
			want:  "I had something interesting for breakfast",
		},
		{
			name:  "single profane word",
			input: "I hear Mastodon is better than Chirpy. sharbert I need to migrate",
			want:  "I hear Mastodon is better than Chirpy. **** I need to migrate",
		},
		{
			name:  "profanity is case insensitive",
			input: "This is a KERFUFFLE opinion I need to share with the world",
			want:  "This is a **** opinion I need to share with the world",
		},
		{
			name:  "multiple profane words",
			input: "kerfuffle sharbert fornax",
			want:  "**** **** ****",
		},
		{
			name:  "profanity with punctuation is not filtered",
			input: "Sharbert!",
			want:  "Sharbert!",
		},
		{
			name:  "empty string",
			input: "",
			want:  "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := filterProfanity(tt.input); got != tt.want {
				t.Errorf("filterProfanity(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}
