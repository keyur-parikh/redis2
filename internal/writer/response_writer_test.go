package writer

import (
	"github.com/keyur/redis2/internal/definitions"
	"testing"
)

func ResponseWritingTest(t *testing.T) {
	tests := []struct {
		name           string
		content        []string
		requestContext definitions.RequestContext
		wantErr        string
	}{
		{
			name:           "Hello Value",
			content:        []string{"hello", "value"},
			requestContext: definitions.RequestContext{Connection: nil, KVStore: map[string]string{"hello": "world"}},
			wantErr:        "*2\\r\\n$5\\r\\nhello\\r\\n$5\\r\\nvalue\\r\\n",
		},
		{
			name:           "Dir tmp/redis/files",
			content:        []string{"dir, tmp/redis/files"},
			requestContext: definitions.RequestContext{Connection: nil, KVStore: map[string]string{"hello": "world"}},
			wantErr:        "*2\\r\\n$3\\r\\ndir\\r\\n$16\\r\\n/tmp/redis-files\\r\\n",
		},
		{
			name:           "dbfilename dumb.rdb",
			content:        []string{"dbfilename", "dump.rdb"},
			requestContext: definitions.RequestContext{Connection: nil, KVStore: map[string]string{"hello": "world"}},
			wantErr:        "*2\\r\\n$10\\r\\ndbfilename\\r\\n$8\\r\\ndump.rdb\\r\\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ArrayResponseWriter(tt.content, &tt.requestContext)
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
