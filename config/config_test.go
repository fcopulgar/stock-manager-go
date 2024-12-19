package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	// Establecer una variable de entorno para la prueba
	os.Setenv("TEST_KEY", "test_value")
	defer os.Unsetenv("TEST_KEY") // Limpiar al finalizar

	// Llamar a LoadConfig (aunque no tengamos .env, no debería fallar)
	LoadConfig()

	value := GetEnv("TEST_KEY")
	if value != "test_value" {
		t.Errorf("Expected GetEnv('TEST_KEY') to return 'test_value', got '%s'", value)
	}
}

func TestLoadConfig_NoEnvFile(t *testing.T) {
	// No creamos un archivo .env para esta prueba
	// Solo esperamos que no falle y no modifique variables preexistentes

	// Establecer una variable de entorno antes de LoadConfig
	os.Setenv("PRE_EXISTING_KEY", "pre_value")
	defer os.Unsetenv("PRE_EXISTING_KEY")

	LoadConfig()

	// Verificar que la variable no se haya alterado
	value := GetEnv("PRE_EXISTING_KEY")
	if value != "pre_value" {
		t.Errorf("Expected 'pre_value', got '%s'", value)
	}
}

func TestLoadConfig_WithEnvFile(t *testing.T) {
	// Creamos un .env temporal
	content := "ENV_TEST_KEY=env_test_value\n"
	err := os.WriteFile(".env", []byte(content), 0644)
	if err != nil {
		t.Fatalf("Error creating .env file: %v", err)
	}
	defer os.Remove(".env") // Limpiar el archivo al final

	// Asegurarnos de que la variable no esté en el entorno
	os.Unsetenv("ENV_TEST_KEY")

	LoadConfig()

	value := GetEnv("ENV_TEST_KEY")
	if value != "env_test_value" {
		t.Errorf("Expected 'env_test_value' from .env, got '%s'", value)
	}
}
