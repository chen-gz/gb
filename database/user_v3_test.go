package database

//import (
//	"testing"
//	_ "testing"
//)
//
//func TestConnectToMariaDB(t *testing.T) {
//	connect_to_mariadb()
//}
//
//func TestGetMariaDBUsers(t *testing.T) {
//	GetMariaDBUsers()
//}
//func TestGetMariaDBDatabases(t *testing.T) {
//	GetMariaDBDatabases()
//}

//func TestUserDbInit(t *testing.T) {
//	db_user, _ := UserDbInit()
//	assert.NotNil(t, db_user, "should not be nil")
//}
//func TestUserAdd(t *testing.T) {
//	db_user, _ := UserDbInit()
//	user := User{
//		Email: "chen-gz@outlook.com",
//		Name:  "Guangzong Chen",
//	}
//	err := UserAdd(db_user, user, "Connie")
//	assert.Nil(t, err, "should be nil")
//}
//func TestV3Login(t *testing.T) {
//	db_user, _ := UserDbInit()
//	result := V3Login(db_user, "chen-gz@outlook.com", "Connie")
//	assert.True(t, result, "should be true")
//	fmt.Println(result)
//
//	result = V3Login(db_user, "chen-gz@outlook.com", "haha")
//	assert.False(t, result, "should be false")
//	fmt.Println(result)
//}
