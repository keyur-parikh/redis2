package parser

import (
	"testing"
)

func TestValidCheckParsing(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		wantCode int
		wantErr  string
	}{
		{
			name:     "Valid input with 2 elements",
			input:    []byte("*2\r\n$4\r\nECHO\r\n$3\r\nhey\r\n"),
			wantCode: 0,
			wantErr:  "",
		},
		{
			name:     "Missing asterisk",
			input:    []byte("2\r\n$3\r\nGET\r\n$5\r\nmykey\r\n"),
			wantCode: 2,
			wantErr:  "expecting * at the beginning",
		},
		{
			name:     "Incomplete data",
			input:    []byte("*2\r\n$3\r\nGET\r\n$5\r\nmyk"),
			wantCode: 1,
			wantErr:  "couldnt find the element",
		},
		{
			name:     "Missing CRLF",
			input:    []byte("*2\r\n$3\r\nGET$5\r\nmykey\r\n"),
			wantCode: 2,
			wantErr:  "Incorrect Sized Message or missing CRLF",
		},
		{
			name:     "Wrong Size Data",
			input:    []byte("*2\r\n$3\r\nGETS\r\n$5\r\nmykey\r\n"),
			wantCode: 2,
			wantErr:  "Incorrect Sized Message or missing CRLF",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, err := ValidCheckParsing(&tt.input, nil)
			if code != tt.wantCode {
				t.Errorf("got code %d, want %d", code, tt.wantCode)
			}
			if err != nil {
				if err.Error() != tt.wantErr {
					t.Errorf("got error %v, want error %v", err, tt.wantErr)
				}
			} else if tt.wantErr != "" {
				t.Errorf("expected error %v but got nil", tt.wantErr)
			}
		})
	}
}
