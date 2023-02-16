package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDBinstance(t *testing.T) {
	db := DBinstance()
	if db == nil {
		t.Error("Database connection failed")
	}

	require.NotNil(t, db)

}
