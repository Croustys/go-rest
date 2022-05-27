package main

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestMain(m *testing.M) {
	exitVal := m.Run()
	os.Exit(exitVal)
}
func TestEnvVars(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		t.Fatal("Error initializing godotenv")
	}

	MONGO_URI := os.Getenv("MONGO_URI")
	BCRYPT_SECRET := os.Getenv("BCRYPT_SECRET")
	JWT_TOKEN := os.Getenv("JWT_TOKEN")

	if MONGO_URI == "" || BCRYPT_SECRET == "" || JWT_TOKEN == "" {
		t.Fatal("One of the Env vars are missing")
	}
}
func TestSetup_store(t *testing.T) {
	originalPath := "/"
	store := setup_store()

	if store.Options.Path != originalPath {
		t.Error("Original Path does not Equal to '/'")
	}

	if !store.Options.HttpOnly {
		t.Error("Store should be HTTP only")
	}

	if store.Options.MaxAge == -1 || store.Options.MaxAge == 0 {
		t.Error("Cookie shouldn't have a 0 or -1 MaxAge")
	}
}
