/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

package imp

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestGeneratedClientDebugRedactsSensitiveData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Set-Cookie", "session=raw-cookie")
		w.Header().Set("Content-Type", "application/json")
		_, _ = io.Copy(io.Discard, r.Body)
		_, _ = w.Write([]byte(`{"access_token":"raw-access-token","name":"visible"}`))
	}))
	defer server.Close()

	var logs bytes.Buffer
	originalWriter := log.Writer()
	log.SetOutput(&logs)
	defer log.SetOutput(originalWriter)

	cfg := NewConfiguration()
	cfg.Debug = true
	cfg.HTTPClient = server.Client()
	client := NewAPIClient(cfg)

	req, err := http.NewRequest(http.MethodPost, server.URL+"/import?X-Amz-Signature=raw-signature&safe=visible", strings.NewReader(`{"secretAccessKey":"raw-secret","name":"visible"}`))
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("Authorization", "Bearer raw-bearer")
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.callAPI(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	got := logs.String()
	for _, secret := range []string{"raw-bearer", "raw-signature", "raw-secret", "raw-cookie", "raw-access-token"} {
		if strings.Contains(got, secret) {
			t.Fatalf("generated client debug log leaked %q: %s", secret, got)
		}
	}
	if !strings.Contains(got, "safe=visible") || !strings.Contains(got, `"name":"visible"`) || !strings.Contains(got, "******") {
		t.Fatalf("generated client debug log lost expected context or masks: %s", got)
	}
}
