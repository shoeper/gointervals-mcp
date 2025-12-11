package intervalsclient

import (
	"context"
	"encoding/base64"
	"net/http"
)

// BasicAuth creates a RequestEditorFn that adds the Authorization header.
func BasicAuth(apiKey string) func(ctx context.Context, req *http.Request) error {
	return func(ctx context.Context, req *http.Request) error {
		// Create the standard "Basic <encoded>" string
		auth := "API_KEY:" + apiKey
		basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))

		// Inject the header
		req.Header.Set("Authorization", basicAuth)
		//req.Header.Set("Accept", "application/json")
		//req.Header.Set("Content-Type", "application/json")

		//log.Printf("Request: %s %s", req.Method, req.URL.String())
		//for k, v := range req.Header {
		//	log.Printf("%s: %v", k, v)
		//}

		//if req.Body != nil {
		//	body, _ := io.ReadAll(req.Body)
		//	log.Printf("Body: %s", body)
		//	req.Body = io.NopCloser(bytes.NewBuffer(body))
		//}

		return nil
	}
}
