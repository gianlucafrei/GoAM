package postgres_adapter

import (
	"testing"

	"goiam/internal/db"

	"github.com/stretchr/testify/require"
)

func TestUserDb(t *testing.T) {
	conn, err := setupTestDB(t)
	require.NoError(t, err)
	defer conn.Close()

	userDB, err := NewPostgresUserDB(conn)
	require.NoError(t, err)
	db.UserDBTests(t, userDB)
}
