package query_handling

import (
	"github.com/keyur/redis2/internal/definitions"
	"testing"
)

func QueryHandlingTest(t *testing.T) {
	tests := []struct {
		function       string
		name           string
		key            string
		value          string
		requestContext definitions.RequestContext
		wantErr        string
	}{
		{
			function:       "get",
			name:           "Successful GET",
			key:            "hello",
			value:          "",
			requestContext: definitions.RequestContext{Connection: nil, KVStore: map[string]string{"hello": "world"}},
			wantErr:        "",
		},
		{
			function:       "get",
			name:           "UnSuccessful GET",
			key:            "hell",
			value:          "",
			requestContext: definitions.RequestContext{Connection: nil, KVStore: map[string]string{"hello": "world"}},
			wantErr:        "not a valid value",
		},
		{
			function:       "set",
			name:           "Successful SET",
			key:            "hello",
			value:          "world",
			requestContext: definitions.RequestContext{Connection: nil, KVStore: map[string]string{}},
			wantErr:        "",
		},
		{
			function:       "set",
			name:           "Successful SET replacement",
			key:            "hello",
			value:          "keyur",
			requestContext: definitions.RequestContext{Connection: nil, KVStore: map[string]string{"hello": "world"}},
			wantErr:        "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.function == "get" {
				err := Get(tt.key, &tt.requestContext)
				if err != nil {
					if err.Error() != tt.wantErr {
						t.Errorf("got error %v, want error %v", err, tt.wantErr)
					}
				} else if tt.wantErr != "" {
					t.Errorf("expected error %v but got nil", err)
				}
			}
			if tt.function == "set" {
				err := Set(tt.key, tt.value, false, 0, &tt.requestContext)
				if err != nil {
					t.Errorf("Got this error %v", err)
				}
				if tt.requestContext.KVStore[tt.key] != tt.value {
					t.Errorf("The value for %v was supposed to be %v but got %v", tt.key, tt.value, tt.requestContext.KVStore[tt.key])
				}
			}

		})
	}
}
