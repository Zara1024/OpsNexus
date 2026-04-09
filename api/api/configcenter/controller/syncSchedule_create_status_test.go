package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	ccmodel "dodevops-api/api/configcenter/model"
	pkgdb "dodevops-api/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type syncScheduleCreateTestResponse struct {
	Code int                  `json:"code"`
	Data ccmodel.SyncSchedule `json:"data"`
}

func setupSyncScheduleTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&ccmodel.SyncSchedule{}); err != nil {
		t.Fatalf("auto migrate sync schedule: %v", err)
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

	return db
}

func TestCreateSyncSchedulePreservesDisabledStatus(t *testing.T) {
	db := setupSyncScheduleTestDB(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.POST("/api/v1/config/sync-schedule", NewSyncScheduleController().Create)

	body := []byte(`{"name":"opsnexus-disabled-schedule","cronExpr":"0 3 * * 0","keyTypes":"[1]","status":0,"remark":"disabled on create"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/config/sync-schedule", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var response syncScheduleCreateTestResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if response.Data.Status != 0 {
		t.Fatalf("response status = %d, want 0; body = %s", response.Data.Status, recorder.Body.String())
	}

	var persisted ccmodel.SyncSchedule
	if err := db.Where("name = ?", "opsnexus-disabled-schedule").First(&persisted).Error; err != nil {
		t.Fatalf("query persisted schedule: %v", err)
	}

	if persisted.Status != 0 {
		t.Fatalf("persisted status = %d, want 0", persisted.Status)
	}
}
