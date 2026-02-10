package config

import (
	"bytes"
	"net/http"
	"os"
)

func NewSupabaseRequest(method, path string, body []byte) (*http.Request, error) {
	url := os.Getenv("SUPABASE_URL") + "/rest/v1" + path

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	key := os.Getenv("SUPABASE_SERVICE_KEY")

	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
