package domain

import "testing"

func TestHash(t *testing.T) {
	var testdata = []struct {
		in  string
		out string
	}{
		{"http://example.com/link1", "bd0904098d30fb9916bfb8a7cca263a2f8fe5f4a528ee31739d79fa264b2ca54"},
		{"http://example.com/link2", "a811307d3e5ceaf5cea08cd9ced1fbb93e40a4d94dc894547b89b8476ea07667"},
		{"http://example.com/link3", "77a8ca7e27edc81a31bbbb0f2afe832946626defd398952e2ed84daaceb1ce42"},
		{"http://kenanbek.me", "3996a77a38ab7f5d1709ac1dc8a05025bb4cd9d070826e0d489ede51e2ea68cc"},
		{"https://kenanbek.github.io/", "49f43babeb15041df63635e8161c316a6c516e4a8dca1210f9844aaef3c9e04e"},
	}

	for _, tt := range testdata {
		t.Run(tt.in, func(t *testing.T) {
			s := Hash(tt.in)
			if s != tt.out {
				t.Errorf("got %s, expected %s", s, tt.out)
			}
		})
	}
}
