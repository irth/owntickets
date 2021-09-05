package owntickets_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/irth/owntickets/internal/owntickets"
	"github.com/stretchr/testify/assert"
)

func TestStrToBool(t *testing.T) {
	falsey := []string{
		"", "asdasfasdfasdf",
		"no", "No", "NO", "n", "N",
		"false", "False", "FALSE", "FaLsE",
		"0",
	}
	truthy := []string{
		"yes", "YES", "Yes", "y", "Y",
		"True", "true", "TRUE", "tRuE",
		"1",
	}

	for _, v := range falsey {
		assert.False(t, owntickets.StrToBool(v))
	}
	for _, v := range truthy {
		assert.True(t, owntickets.StrToBool(v))
	}
}

func TestLoadFromEnv(t *testing.T) {
	c := owntickets.Config{}
	c.LoadFromEnv()
	assert.Empty(t, c.PasswordHash)
	assert.False(t, c.RequirePasswordForTicketCreation)
	assert.Empty(t, c.TicketCreationPasswordHash)

	os.Setenv("OWNTICKETS_PASSWORD_HASH", "2137")
	os.Setenv("OWNTICKETS_TICKET_PASSWORD_HASH", "2138")
	os.Setenv("OWNTICKETS_REQUIRE_PASSWORD", "yes")
	os.Setenv("OWNTICKETS_DATABASE", "db")
	c.LoadFromEnv()
	os.Unsetenv("OWNTICKETS_PASSWORD_HASH")
	os.Unsetenv("OWNTICKETS_TICKET_PASSWORD_HASH")
	os.Unsetenv("OWNTICKETS_REQUIRE_PASSWORD")
	os.Unsetenv("OWNTICKETS_DATABASE")
	assert.Equal(t, "2137", c.PasswordHash)
	assert.Equal(t, "2138", c.TicketCreationPasswordHash)
	assert.True(t, c.RequirePasswordForTicketCreation)
	assert.Equal(t, "db", c.Database)
}

func TestLoadFromFile(t *testing.T) {
	tempDir := t.TempDir()
	configPath := path.Join(tempDir, "config.json")
	f, err := os.Create(configPath)
	assert.NoError(t, err)
	defer f.Close()
	json.NewEncoder(f).Encode(map[string]interface{}{
		"passwordHash":                     "2137",
		"requirePasswordForTicketCreation": true,
		"ticketCreationPasswordHash":       "lol",
		"database":                         "db",
	})
	f.Close()

	var c owntickets.Config
	c.LoadFromFile(configPath)
	assert.Equal(t, "2137", c.PasswordHash)
	assert.Equal(t, "lol", c.TicketCreationPasswordHash)
	assert.True(t, c.RequirePasswordForTicketCreation)
	assert.Equal(t, "db", c.Database)
}

func TestValidate(t *testing.T) {
	var c owntickets.Config
	assert.Error(t, c.Validate())
	c.Database = "a"
	c.PasswordHash = "lol"
	assert.NoError(t, c.Validate())
	c.RequirePasswordForTicketCreation = true
	assert.Error(t, c.Validate())
	c.TicketCreationPasswordHash = "a"
	assert.NoError(t, c.Validate())
	c.Database = ""
	assert.Error(t, c.Validate())
}
