package database

import (
	"testing"
	_ "testing"
)

func TestConnectToMariaDB(t *testing.T) {
	connect_to_mariadb()
}

func TestGetMariaDBUsers(t *testing.T) {
	GetMariaDBUsers()
}
func TestGetMariaDBDatabases(t *testing.T) {
	GetMariaDBDatabases()
}
