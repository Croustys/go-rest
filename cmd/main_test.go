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
