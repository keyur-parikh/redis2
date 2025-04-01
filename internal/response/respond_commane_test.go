package response

import (
	"testing"
)

func RespondCommandTesting(t *testing.T) {
	tests := []struct {
		name    string
		input   []string
		wantErr string
	}{
		{
			name:    "Valid ECHO",
			input:   []string{"ECHO", "hello"},
			wantErr: "",
		},
		{
			name:    "Another Valid ECHO",
			input:   []string{"ECHO", "hello"},
			wantErr: "",
		},
		{
			name:    "Invalid ECHO",
			input:   []string{"ECH", "hello"},
			wantErr: "can't understand the command",
		},
		{
			name:    "Valid GET",
			input:   []string{"GeT", "hello"},
			wantErr: "",
		},
		{
			name:    "Invalid GET",
			input:   []string{"GeT", "hello", "sir"},
			wantErr: "GET must be followed by one value",
		},
		{
			name:    "Valid SET",
			input:   []string{"seT", "keyur", "parikh"},
			wantErr: "",
		},
		{
			name:    "InValid SET",
			input:   []string{"seT", "keyur"},
			wantErr: "SET must be followed by two values",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := RespondCommand(tt.input, nil)
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
