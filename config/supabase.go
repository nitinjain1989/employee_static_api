package config

import (
	"os"

	"github.com/supabase-community/supabase-go"
)

func NewSupabaseClient() *supabase.Client {
	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_KEY")

	client, err := supabase.NewClient(url, key, nil)
	if err != nil {
		panic(err)
	}
	return client
}
