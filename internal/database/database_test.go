package database

import (
	"fmt"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInitDB_Success(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	mock.ExpectPing().WillReturnError(nil)

	err = DB.Ping()
	assert.NoError(t, err)
	assert.NotNil(t, DB)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestInitDB_PingFail(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	assert.NoError(t, err)
	defer db.Close()

	DB = db

	mock.ExpectPing().WillReturnError(fmt.Errorf("ping failed"))

	err = DB.Ping()
	assert.Error(t, err, "ping failed")
	assert.NoError(t, mock.ExpectationsWereMet())
}
