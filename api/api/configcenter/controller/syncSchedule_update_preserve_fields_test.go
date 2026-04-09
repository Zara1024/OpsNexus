package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ccmodel "dodevops-api/api/configcenter/model"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

type syncScheduleUpdateTestResponse struct {
	Code int                  `json:"code"`
	Data ccmodel.SyncSchedule `json:"data"`
}

func TestUpdateSyncSchedulePreservesCreatedAtAndRuntimeFields(t *testing.T) {
	db := setupSyncScheduleTestDB(t)
	gin.SetMode(gin.TestMode)

	createdAt := util.HTime{Time: time.Date(2026, 4, 2, 12, 30, 25, 0, time.Local)}
	lastRunAt := util.HTime{Time: time.Date(2026, 4, 2, 12, 45, 0, 0, time.Local)}
	nextRunAt := util.HTime{Time: time.Date(2026, 4, 5, 4, 5, 0, 0, time.Local)}

	seed := ccmodel.SyncSchedule{
		Name:        "opsnexus-update-preserve-fields",
		CronExpr:    "5 4 * * 0",
		KeyTypes:    "[1]",
		Status:      0,
		LastRunTime: &lastRunAt,
		NextRunTime: &nextRunAt,
		SyncLog:     "manual sync finished",
		Remark:      "before edit",
		CreatedAt:   createdAt,
		UpdatedAt:   createdAt,
	}
	if err := db.Create(&seed).Error; err != nil {
		t.Fatalf("seed schedule: %v", err)
	}

	router := gin.New()
	router.PUT("/api/v1/config/sync-schedule", NewSyncScheduleController().Update)

	body := []byte(`{"id":1,"name":"opsnexus-update-preserve-fields-edit","cronExpr":"25 4 * * 0","keyTypes":"[1]","status":1,"remark":"after edit"}`)
	req := httptest.NewRequest(http.MethodPut, "/api/v1/config/sync-schedule", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var response syncScheduleUpdateTestResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if got := response.Data.CreatedAt.Format("2006-01-02 15:04:05"); got != "2026-04-02 12:30:25" {
		t.Fatalf("response createdAt = %s, want 2026-04-02 12:30:25; body = %s", got, recorder.Body.String())
	}
	if response.Data.LastRunTime == nil {
		t.Fatalf("response lastRunTime = nil, want preserved value; body = %s", recorder.Body.String())
	}
	if got := response.Data.LastRunTime.Format("2006-01-02 15:04:05"); got != "2026-04-02 12:45:00" {
		t.Fatalf("response lastRunTime = %s, want 2026-04-02 12:45:00; body = %s", got, recorder.Body.String())
	}
	if response.Data.SyncLog != "manual sync finished" {
		t.Fatalf("response syncLog = %q, want manual sync finished; body = %s", response.Data.SyncLog, recorder.Body.String())
	}

	var persisted ccmodel.SyncSchedule
	if err := db.First(&persisted, seed.ID).Error; err != nil {
		t.Fatalf("query persisted schedule: %v", err)
	}

	if got := persisted.CreatedAt.Format("2006-01-02 15:04:05"); got != "2026-04-02 12:30:25" {
		t.Fatalf("persisted createdAt = %s, want 2026-04-02 12:30:25", got)
	}
	if persisted.LastRunTime == nil {
		t.Fatalf("persisted lastRunTime = nil, want preserved value")
	}
	if got := persisted.LastRunTime.Format("2006-01-02 15:04:05"); got != "2026-04-02 12:45:00" {
		t.Fatalf("persisted lastRunTime = %s, want 2026-04-02 12:45:00", got)
	}
	if persisted.SyncLog != "manual sync finished" {
		t.Fatalf("persisted syncLog = %q, want manual sync finished", persisted.SyncLog)
	}
	if persisted.Name != "opsnexus-update-preserve-fields-edit" {
		t.Fatalf("persisted name = %q, want opsnexus-update-preserve-fields-edit", persisted.Name)
	}
	if persisted.Status != 1 {
		t.Fatalf("persisted status = %d, want 1", persisted.Status)
	}
	if persisted.Remark != "after edit" {
		t.Fatalf("persisted remark = %q, want after edit", persisted.Remark)
	}
}
