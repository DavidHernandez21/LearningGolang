package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"
)

// generateLoad sends count requests to the server.
func generateLoad(count int, source string, addr string) error {
	if source == "" {
		return fmt.Errorf("-source must be set to a markdown source file")
	}
	if addr == "" {
		return fmt.Errorf("-addr must be set to the address of the server (e.g., http://localhost:8080)")
	}

	src, err := os.ReadFile(source)
	if err != nil {
		return fmt.Errorf("error reading source: %v", err)
	}
	reader := bytes.NewReader(src)

	url := addr + "/render"

	for i := 0; i < count; i++ {
		reader.Seek(0, io.SeekStart)

		resp, err := http.Post(url, "text/markdown", reader)
		if err != nil {
			return fmt.Errorf("error writing request: %v", err)
		}
		if _, err := io.Copy(io.Discard, resp.Body); err != nil {
			return fmt.Errorf("error reading response body: %v", err)
		}
		resp.Body.Close()
	}

	return nil
}

func BenchmarkLoad(b *testing.B) {
	source := "README.md"
	addr := "http://127.0.0.1:8080"
	if err := generateLoad(b.N, source, addr); err != nil {
		b.Errorf("generateLoad got err %v want nil", err)
	}
}
