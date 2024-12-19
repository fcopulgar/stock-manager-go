package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Set an environment variable for the test
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY")

	// Call LoadConfig (even if we don't have .env, it should not fail)
	LoadConfig()

	value := GetEnv("TEST_KEY")
	if value != "test_value" {
		t.Errorf("Expected GetEnv('TEST_KEY') to return 'test_value', got '%s'", value)
	}
}

func TestLoadConfig_NoEnvFile(t *testing.T) {
	// We do not create an .env file for this test.
	// We just hope that it does not fail and does not modify pre-existing variables.

	// Set an environment variable before LoadConfig
	os.Setenv("PRE_EXISTING_KEY", "pre_value")
	defer os.Unsetenv("PRE_EXISTING_KEY")

	LoadConfig()

	// Verify that the variable has not been altered
	value := GetEnv("PRE_EXISTING_KEY")
	if value != "pre_value" {
		t.Errorf("Expected 'pre_value', got '%s'", value)
	}
}

func TestLoadConfig_WithEnvFile(t *testing.T) {
	// We create a temporary .env
	content := "ENV_TEST_KEY=env_test_value\n"
	err := os.WriteFile(".env", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Error creating .env file: %v", err)
	}
	defer os.Remove(".env")

	// Make sure that the variable is not in the environment
	os.Unsetenv("ENV_TEST_KEY")

	LoadConfig()

	value := GetEnv("ENV_TEST_KEY")
	if value != "env_test_value" {
		t.Errorf("Expected 'env_test_value' from .env, got '%s'", value)
	}
}
