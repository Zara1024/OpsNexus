package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	cmdbmodel "dodevops-api/api/cmdb/model"
	"dodevops-api/common/util"
	pkgdb "dodevops-api/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func setupCmdbSqlLogTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&cmdbmodel.CmdbSQLRecord{}); err != nil {
		t.Fatalf("auto migrate cmdb sql log: %v", err)
	}

	oldDB := pkgdb.Db
	pkgdb.Db = db

	t.Cleanup(func() {
		sqlDB, _ := db.DB()
		if sqlDB != nil {
			_ = sqlDB.Close()
		}
		pkgdb.Db = oldDB
	})

	seed := []cmdbmodel.CmdbSQLRecord{
		{InstanceID: "127.0.0.1", Database: "mysql", OperationType: "SELECT", SQLContent: "SELECT 1;", ExecUser: "admin", IP: "10.0.0.1", Result: "SUCCESS", QueryTime: util.HTime{Time: time.Now()}},
		{InstanceID: "127.0.0.1", Database: "mysql", OperationType: "SELECT", SQLContent: "SELECT 2;", ExecUser: "admin", IP: "10.0.0.1", Result: "SUCCESS", QueryTime: util.HTime{Time: time.Now().Add(time.Second)}},
		{InstanceID: "127.0.0.1", Database: "mysql", OperationType: "SELECT", SQLContent: "SELECT 3;", ExecUser: "admin", IP: "10.0.0.1", Result: "SUCCESS", QueryTime: util.HTime{Time: time.Now().Add(2 * time.Second)}},
	}

	for i := range seed {
		if err := db.Create(&seed[i]).Error; err != nil {
			t.Fatalf("create seed %d: %v", i, err)
		}
	}

	return db
}

func TestBatchDeleteCmdbSqlLogDeletesSelectedRows(t *testing.T) {
	db := setupCmdbSqlLogTestDB(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.DELETE("/api/v1/cmdb/sqlLog/batch/delete", BatchDeleteCmdbSqlLog)

	body, err := json.Marshal(cmdbmodel.BatchDeleteCmdbSqlLogDto{Ids: []uint{1, 2}})
	if err != nil {
		t.Fatalf("marshal body: %v", err)
	}

	req := httptest.NewRequest(http.MethodDelete, "/api/v1/cmdb/sqlLog/batch/delete", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var count int64
	if err := db.Model(&cmdbmodel.CmdbSQLRecord{}).Count(&count).Error; err != nil {
		t.Fatalf("count logs: %v", err)
	}
	if count != 1 {
		t.Fatalf("count = %d, want 1", count)
	}

	var remaining []cmdbmodel.CmdbSQLRecord
	if err := db.Order("id ASC").Find(&remaining).Error; err != nil {
		t.Fatalf("load remaining logs: %v", err)
	}
	if len(remaining) != 1 || remaining[0].ID != 3 {
		t.Fatalf("remaining logs = %+v, want only id=3", remaining)
	}
}
