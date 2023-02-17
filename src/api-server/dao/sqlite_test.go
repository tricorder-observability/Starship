package dao

import "testing"

func TestInitSqlitGorm(t *testing.T) {
	testDbFilePath := "test_un_deploy_code"
	_, err := InitSqlite(testDbFilePath)
	if err != nil {
		t.Errorf("init sqlite gorm error")
	}
}
