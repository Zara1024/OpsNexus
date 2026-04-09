package controller

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ccmodel "dodevops-api/api/configcenter/model"
	"dodevops-api/common/util"
	pkgdb "dodevops-api/pkg/db"

	"github.com/gin-gonic/gin"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

type ecsAuthListTestResponse struct {
	Code int `json:"code"`
	Data struct {
		Total int64             `json:"total"`
		List  []ccmodel.EcsAuth `json:"list"`
	} `json:"data"`
}

func setupEcsAuthListTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&ccmodel.EcsAuth{}); err != nil {
		t.Fatalf("auto migrate ecs auth: %v", err)
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

	now := util.HTime{Time: time.Now()}
	seed := []ccmodel.EcsAuth{
		{Name: "opsnexus-seed-password", Type: 1, Username: "root", Password: "plain-1", Port: 22, CreateTime: now, Remark: "password"},
		{Name: "opsnexus-seed-publickey", Type: 3, Username: "root", Port: 22, CreateTime: now, Remark: "public-key"},
		{Name: "opsnexus-seed-privatekey", Type: 2, Username: "ubuntu", PublicKey: "PRIVATE KEY", Port: 2222, CreateTime: now, Remark: "private-key"},
	}

	for i := range seed {
		if err := db.Create(&seed[i]).Error; err != nil {
			t.Fatalf("create seed %d: %v", i, err)
		}
	}

	return db
}

func TestEcsAuthListFiltersByName(t *testing.T) {
	setupEcsAuthListTestDB(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/config/ecsauthlist", NewEcsAuthController().GetEcsAuthList)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/config/ecsauthlist?page=1&pageSize=10&name=opsnexus-seed-publickey", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var response ecsAuthListTestResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if response.Data.Total != 1 {
		t.Fatalf("total = %d, want 1; body = %s", response.Data.Total, recorder.Body.String())
	}
	if len(response.Data.List) != 1 {
		t.Fatalf("list len = %d, want 1; body = %s", len(response.Data.List), recorder.Body.String())
	}
	if response.Data.List[0].Name != "opsnexus-seed-publickey" {
		t.Fatalf("name = %q, want %q", response.Data.List[0].Name, "opsnexus-seed-publickey")
	}
}
