package service

import (
	"context"
	"encoding/csv"
	"strconv"

	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/audit/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/audit/repository"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/middleware"
)

// Services 聚合 audit 模块服务。
type Services struct {
	Log *LogService
}

func New(db *gorm.DB) *Services {
	repo := repository.New(db)
	return &Services{Log: NewLogService(repo)}
}

// ListQuery 前端查询参数。
type ListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Keyword   string `form:"keyword"`
	Module    string `form:"module"`
	Action    string `form:"action"`
	StartTime string `form:"start_time"`
	EndTime   string `form:"end_time"`
}

type LogService struct {
	repo *repository.Repository
}

func NewLogService(repo *repository.Repository) *LogService { return &LogService{repo: repo} }

func (s *LogService) CreateFromRecord(ctx context.Context, record *middleware.AuditRecord) error {
	if record == nil {
		return nil
	}
	return s.repo.Create(ctx, &model.AuditLog{
		UserID:       record.UserID,
		Username:     record.Username,
		IP:           record.IP,
		Method:       record.Method,
		Path:         record.Path,
		RequestBody:  record.RequestBody,
		ResponseCode: record.ResponseCode,
		Module:       record.Module,
		Action:       record.Action,
		ResourceType: record.ResourceType,
		ResourceID:   record.ResourceID,
		Description:  record.Description,
		DurationMS:   record.DurationMS,
	})
}

func (s *LogService) List(query ListQuery) ([]model.AuditLog, int64, error) {
	return s.repo.List(repository.LogQuery{
		Page:      query.Page,
		PageSize:  query.PageSize,
		Keyword:   query.Keyword,
		Module:    query.Module,
		Action:    query.Action,
		StartTime: query.StartTime,
		EndTime:   query.EndTime,
	})
}

func (s *LogService) ListLoginLogs(query ListQuery) ([]model.AuditLog, int64, error) {
	return s.repo.ListLoginLogs(repository.LogQuery{
		Page:      query.Page,
		PageSize:  query.PageSize,
		Keyword:   query.Keyword,
		Module:    query.Module,
		Action:    query.Action,
		StartTime: query.StartTime,
		EndTime:   query.EndTime,
	})
}

func (s *LogService) ExportCSV(query ListQuery) (string, error) {
	list, err := s.repo.ListForExport(repository.LogQuery{
		Keyword:   query.Keyword,
		Module:    query.Module,
		Action:    query.Action,
		StartTime: query.StartTime,
		EndTime:   query.EndTime,
	})
	if err != nil {
		return "", err
	}

	buf := make([]byte, 0, 4096)
	w := csv.NewWriter(byteBuffer{&buf})
	_ = w.Write([]string{"ID", "用户名", "IP", "方法", "路径", "模块", "动作", "状态码", "耗时(ms)", "时间"})
	for _, item := range list {
		_ = w.Write([]string{
			strconv.FormatInt(item.ID, 10),
			item.Username,
			item.IP,
			item.Method,
			item.Path,
			item.Module,
			item.Action,
			strconv.Itoa(item.ResponseCode),
			strconv.FormatInt(item.DurationMS, 10),
			item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	w.Flush()
	return string(buf), nil
}

type byteBuffer struct{ b *[]byte }

func (bb byteBuffer) Write(p []byte) (n int, err error) {
	*bb.b = append(*bb.b, p...)
	return len(p), nil
}
