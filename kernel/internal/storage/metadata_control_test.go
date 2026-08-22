package storage

import (
	"bytes"
	"testing"
)

func TestValidateUploadRejectsControlCharacters(t *testing.T) {
	payload := []byte("synthetic")
	checksum := SHA256Hex(payload)

	tests := []struct {
		name     string
		fileName string
		metadata map[string]string
	}{
		{name: "tab in file name", fileName: "report\tfinal.txt"},
		{name: "delete in file name", fileName: "report\x7ffinal.txt"},
		{name: "tab in metadata", fileName: "report.txt", metadata: map[string]string{"fixture": "bad\tvalue"}},
		{name: "delete in metadata", fileName: "report.txt", metadata: map[string]string{"fixture": "bad\x7fvalue"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := validateUpload(Upload{
				Body:          bytes.NewReader(payload),
				ContentLength: int64(len(payload)),
				ContentType:   "text/plain",
				FileName:      test.fileName,
				SHA256:        checksum,
				Metadata:      test.metadata,
			}, 1024)
			if err == nil {
				t.Fatal("validateUpload() error = nil, want control-character rejection")
			}
			assertStorageFailureCode(t, err, codeMetadataInvalid)
		})
	}
}
