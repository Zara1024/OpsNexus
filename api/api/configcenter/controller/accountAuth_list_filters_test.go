package controller

import (
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

type accountAuthListTestResponse struct {
	Code int `json:"code"`
	Data struct {
		Total int64                 `json:"total"`
		List  []ccmodel.AccountAuth `json:"list"`
	} `json:"data"`
}

func setupAccountAuthListTestDB(t *testing.T) *gorm.DB {
	t.Helper()

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&ccmodel.AccountAuth{}); err != nil {
		t.Fatalf("auto migrate account auth: %v", err)
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

	seed := []ccmodel.AccountAuth{
		{Alias: "saas3-mysql-root", Host: "127.0.0.1", Port: 3306, Name: "root", Password: "plain-1", Type: 1, Remark: "mysql"},
		{Alias: "opsnexus-mock-jenkins", Host: "10.0.0.1", Port: 18080, Name: "mock", Password: "plain-2", Type: 4, Remark: "jenkins"},
		{Alias: "opsnexus-generic-account", Host: "10.0.0.2", Port: 22, Name: "ops", Password: "plain-3", Type: 6, Remark: "generic"},
	}

	for i := range seed {
		if err := seed[i].EncryptPassword(); err != nil {
			t.Fatalf("encrypt seed password %d: %v", i, err)
		}
		if err := db.Create(&seed[i]).Error; err != nil {
			t.Fatalf("create seed %d: %v", i, err)
		}
	}

	return db
}

func TestAccountAuthListFiltersByAlias(t *testing.T) {
	setupAccountAuthListTestDB(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/config/accountauth/list", NewAccountAuthController().List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/config/accountauth/list?page=1&pageSize=10&alias=opsnexus-mock-jenkins", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var response accountAuthListTestResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if response.Data.Total != 1 {
		t.Fatalf("total = %d, want 1; body = %s", response.Data.Total, recorder.Body.String())
	}
	if len(response.Data.List) != 1 {
		t.Fatalf("list len = %d, want 1; body = %s", len(response.Data.List), recorder.Body.String())
	}
	if response.Data.List[0].Alias != "opsnexus-mock-jenkins" {
		t.Fatalf("alias = %q, want %q", response.Data.List[0].Alias, "opsnexus-mock-jenkins")
	}
}

func TestAccountAuthListFiltersByType(t *testing.T) {
	setupAccountAuthListTestDB(t)
	gin.SetMode(gin.TestMode)

	router := gin.New()
	router.GET("/api/v1/config/accountauth/list", NewAccountAuthController().List)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/config/accountauth/list?page=1&pageSize=10&type=6", nil)
	recorder := httptest.NewRecorder()

	router.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d, body = %s", recorder.Code, http.StatusOK, recorder.Body.String())
	}

	var response accountAuthListTestResponse
	if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}

	if response.Data.Total != 1 {
		t.Fatalf("total = %d, want 1; body = %s", response.Data.Total, recorder.Body.String())
	}
	if len(response.Data.List) != 1 {
		t.Fatalf("list len = %d, want 1; body = %s", len(response.Data.List), recorder.Body.String())
	}
	if response.Data.List[0].Type != 6 {
		t.Fatalf("type = %d, want 6", response.Data.List[0].Type)
	}
}
