package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew проверяет создание и чтение конфигурации из файла.
func TestNew(t *testing.T) {
	// Создание временного файла конфигурации.
	content := []byte(`
source_url: "https://xkcd.com"
DbFile: "pkg/database/database.json"
parallel: 100
port: ":8080"
postgresql: "host=localhost dbname=postgres user=postgres port=5432 password=local sslmode=disable"
token_max_time: 30
concurrencyLimit: 5
rateLimit: 2
`)
	tmpfile, err := os.CreateTemp("", "config.*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // Очистка после завершения теста.

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Чтение конфигурации из файла.
	cfg, err := New(tmpfile.Name())
	assert.Equal(t, err, nil)
	assert.Equal(t, "https://xkcd.com", cfg.Url)
	assert.Equal(t, "pkg/database/database.json", cfg.DbFile)
	assert.Equal(t, 100, cfg.Parallel)
	assert.Equal(t, ":8080", cfg.Port)
	assert.Equal(t, "host=localhost dbname=postgres user=postgres port=5432 password=local sslmode=disable", cfg.Postgresql)
	assert.Equal(t, 30, cfg.Token_max_time)
	assert.Equal(t, 5, cfg.ConcurrencyLimit)
	assert.Equal(t, 2, cfg.RateLimit)
}

func TestNew_InvalidConfig(t *testing.T) {
	_, err := New("hahaha.yaml")
	assert.Error(t, err)
	assert.Equal(t, "can't read config: open hahaha.yaml: The system cannot find the file specified.", err.Error())
}
