package testutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

type ErrResponse struct {
	Message string   `json:"message"`
	Details []string `json:"details,omitempty"`
}

func RespondJSON(ctx context.Context, w http.ResponseWriter, body any, status int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	bodyBytes, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		rsp := ErrResponse{
			Message: http.StatusText(http.StatusInternalServerError),
		}

		if err := json.NewEncoder(w).Encode(rsp); err != nil {
			fmt.Printf("write error response: error: %v", err)
		}
		return
	}
	w.WriteHeader(status)
	if _, err := fmt.Fprintf(w, "%s", bodyBytes); err != nil {
		fmt.Printf("write response error: %v", err)
	}

}

func AssertJSON(t *testing.T, want, got []byte) {
	t.Helper()

	var jw, jg any
	if err := json.Unmarshal(want, &jw); err != nil {
		t.Fatalf("cannot unmarshal want %v: %v", want, err)

		if err := json.Unmarshal(got, &jg); err != nil {
			t.Fatalf("cannnot unmarshal got %q: %v", got, err)
		}
		if diff := cmp.Diff(jg, jw); diff != "" {
			t.Errorf("got differs: (-got +want)\n%s", diff)
		}
	}
}

func AssertResponse(t *testing.T, got *http.Response, status int, body []byte) {
	t.Helper()
	t.Cleanup(func() { _ = got.Body.Close() })

	gb, err := io.ReadAll(got.Body)
	if err != nil {
		t.Fatal(err)
	}

	if got.StatusCode != status {
		t.Fatalf("want status %d, but got %d, body: %q", status, got.StatusCode, gb)
	}
	if len(gb) == 0 && len(body) == 0 {
		return
	}

	AssertJSON(t, body, gb)
}

func LoadFile(t *testing.T, path string) []byte {
	t.Helper()
	bt, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("cannot reand from %q: %v", path, err)
	}
	return bt
}
