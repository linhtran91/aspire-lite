package handlers

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLimitOffset(t *testing.T) {
	cases := []struct {
		name       string
		input      url.Values
		wantLimit  int
		wantOffset int
	}{
		{
			name: "success",
			input: url.Values{
				"page": []string{"2"},
				"size": []string{"10"},
			},
			wantLimit:  10,
			wantOffset: 11,
		},
		{
			name:       "empty",
			input:      url.Values{},
			wantLimit:  10,
			wantOffset: 1,
		},
		{
			name: "high offset",
			input: url.Values{
				"page": []string{"1"},
				"size": []string{"10000000"},
			},
			wantLimit:  10,
			wantOffset: 1,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			size, page := getLimitOffset(tc.input)
			assert.Equal(t, size, tc.wantLimit)
			assert.Equal(t, page, tc.wantOffset)
		})
	}
}
