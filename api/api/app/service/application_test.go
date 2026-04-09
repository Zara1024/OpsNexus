package service

import (
	"dodevops-api/api/app/model"
	ccmodel "dodevops-api/api/configcenter/model"
	"testing"
	"time"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

func TestBuildJenkinsBaseURL(t *testing.T) {
	svc := &ApplicationService{}

	tests := []struct {
		name    string
		account *ccmodel.AccountAuth
		want    string
	}{
		{
			name: "plain host with port",
			account: &ccmodel.AccountAuth{
				Host: "180.76.231.65",
				Port: 8080,
			},
			want: "http://180.76.231.65:8080",
		},
		{
			name: "host already includes scheme and port",
			account: &ccmodel.AccountAuth{
				Host: "http://180.76.231.65:8080/",
				Port: 0,
			},
			want: "http://180.76.231.65:8080",
		},
		{
			name: "https default port omitted",
			account: &ccmodel.AccountAuth{
				Host: "jenkins.example.com",
				Port: 443,
			},
			want: "https://jenkins.example.com",
		},
		{
			name: "host already includes https",
			account: &ccmodel.AccountAuth{
				Host: "https://jenkins.example.com/",
				Port: 443,
			},
			want: "https://jenkins.example.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.buildJenkinsBaseURL(tt.account); got != tt.want {
				t.Fatalf("buildJenkinsBaseURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestBuildJenkinsJobURL(t *testing.T) {
	svc := &ApplicationService{}
	account := &ccmodel.AccountAuth{
		Host: "http://180.76.231.65:8080/",
		Port: 0,
	}

	tests := []struct {
		name    string
		jobName string
		want    string
	}{
		{
			name:    "simple job",
			jobName: "s3-api",
			want:    "http://180.76.231.65:8080/job/s3-api",
		},
		{
			name:    "folder job",
			jobName: "team-a/release-api",
			want:    "http://180.76.231.65:8080/job/team-a/job/release-api",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := svc.buildJenkinsJobURL(account, tt.jobName); got != tt.want {
				t.Fatalf("buildJenkinsJobURL() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestUpdateTaskFailureUpdatesDeploymentStatusWhenAllTasksFinished(t *testing.T) {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		t.Fatalf("open sqlite db: %v", err)
	}

	if err := db.AutoMigrate(&model.QuickDeployment{}, &model.QuickDeploymentTask{}); err != nil {
		t.Fatalf("auto migrate quick deployment tables: %v", err)
	}

	svc := &ApplicationService{db: db}

	deployment := model.QuickDeployment{
		Title:           "mock-release",
		BusinessGroupID: 56,
		BusinessDeptID:  2,
		Status:          2,
		TaskCount:       1,
		ExecutionMode:   1,
		CreatorID:       89,
		CreatorName:     "管理员",
	}
	if err := db.Create(&deployment).Error; err != nil {
		t.Fatalf("create deployment: %v", err)
	}

	task := model.QuickDeploymentTask{
		DeploymentID: deployment.ID,
		AppID:        22,
		AppName:      "opsnexus-test-app",
		AppCode:      "opsnexus-test-app",
		Environment:  "test",
		JenkinsEnvID: 61,
		Status:       2,
		ExecuteOrder: 1,
	}
	if err := db.Create(&task).Error; err != nil {
		t.Fatalf("create task: %v", err)
	}

	svc.updateTaskFailure(&task, time.Now().Add(-2*time.Second), "任务启动失败")

	var updated model.QuickDeployment
	if err := db.First(&updated, deployment.ID).Error; err != nil {
		t.Fatalf("reload deployment: %v", err)
	}

	if updated.Status != 4 {
		t.Fatalf("deployment status = %d, want %d after task failure", updated.Status, 4)
	}
}
