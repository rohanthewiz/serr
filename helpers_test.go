package serr

import (
	"errors"
	"testing"
)

func TestLastNTokens(t *testing.T) {
	type args struct {
		str       string
		separator string
		numTokens int
	}
	tests := []struct {
		name           string
		args           args
		wantLastTokens string
	}{
		{
			name:           "Test Last N Tokens 1",
			args:           args{str: "abc/def/ghi", separator: "/", numTokens: 2},
			wantLastTokens: "def/ghi",
		},
		{
			name:           "Test Last N Tokens 2",
			args:           args{str: "abcdefg", separator: "/", numTokens: 2},
			wantLastTokens: "abcdefg",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotLastTokens := LastNTokens(tt.args.str, tt.args.separator, tt.args.numTokens); gotLastTokens != tt.wantLastTokens {
				t.Errorf("Last2Tokens() = %v, want %v", gotLastTokens, tt.wantLastTokens)
			}
		})
	}
}

func TestStringFromErr(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		want    string
		wantAlt string
	}{
		{
			name: "Nil error",
			err:  nil,
			want: "",
		},
		{
			name: "Standard error",
			err:  errors.New("standard error"),
			want: "standard error",
		},
		{
			name:    "SErr with message",
			err:     NewSErr("serr message"),
			want:    "serr message [error_attrs] => function->rohanthewiz/serr.TestStringFromErr; location->serr/helpers_test.go:58",
			wantAlt: "serr message [error_attrs] => location->serr/helpers_test.go:58; function->rohanthewiz/serr.TestStringFromErr",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := StringFromErr(tt.err); !(got == tt.want || (tt.wantAlt != "" && got == tt.wantAlt)) {
				t.Errorf("StringFromErr() = %v, want %v or %v", got, tt.want, tt.wantAlt)
			}
		})
	}
}
