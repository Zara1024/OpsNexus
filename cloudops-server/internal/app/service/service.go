package service

import (
	"fmt"

	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/app/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/app/repository"
)

// Services 聚合 app 模块服务。
type Services struct {
	Dashboard *DashboardService
}

func New(db *gorm.DB) *Services {
	repo := repository.New(db)
	return &Services{Dashboard: NewDashboardService(repo)}
}

type DashboardService struct {
	repo *repository.Repository
}

func NewDashboardService(repo *repository.Repository) *DashboardService {
	return &DashboardService{repo: repo}
}

func (s *DashboardService) Overview() (*model.Overview, error) {
	res := &model.Overview{}
	var err error

	if res.Asset.HostTotal, err = s.repo.CountHosts(); err != nil {
		return nil, err
	}
	if res.Asset.HostOnline, err = s.repo.CountHostByStatus(1); err != nil {
		return nil, err
	}
	if res.Asset.HostOffline, err = s.repo.CountHostByStatus(2); err != nil {
		return nil, err
	}
	if res.Asset.HostAlert, err = s.repo.CountHostAlert(); err != nil {
		return nil, err
	}
	if res.Asset.K8sClusterTotal, err = s.repo.CountK8sClusters(); err != nil {
		return nil, err
	}
	// 当前阶段未接入数据库实例管理，先返回 0。
	res.Asset.DBInstanceTotal = 0

	if res.Alert.ActiveTotal, err = s.repo.CountFiringAlerts(); err != nil {
		return nil, err
	}
	if res.Alert.TodayHandled, err = s.repo.CountTodayHandledAlerts(); err != nil {
		return nil, err
	}

	severityRows, err := s.repo.AlertSeverityStats()
	if err != nil {
		return nil, err
	}
	res.Alert.SeverityStats = []model.SeverityStat{
		{Severity: "P1", Count: 0},
		{Severity: "P2", Count: 0},
		{Severity: "P3", Count: 0},
		{Severity: "P4", Count: 0},
	}
	idx := map[string]int{"P1": 0, "P2": 1, "P3": 2, "P4": 3}
	for _, row := range severityRows {
		sev := toString(row["severity"])
		if i, ok := idx[sev]; ok {
			res.Alert.SeverityStats[i].Count = toInt64(row["count"])
		}
	}
	res.RecentAlertTop5 = topNSeverity(res.Alert.SeverityStats, 5)

	cpu, memory, disk, err := s.repo.ResourceAvg()
	if err != nil {
		return nil, err
	}
	res.Resource.CPUAvg = cpu
	res.Resource.MemoryAvg = memory
	res.Resource.DiskAvg = disk

	opRows, err := s.repo.RecentOps(10)
	if err != nil {
		return nil, err
	}
	res.RecentOps = make([]model.RecentOp, 0, len(opRows))
	for _, row := range opRows {
		res.RecentOps = append(res.RecentOps, model.RecentOp{
			ID:        toInt64(row["id"]),
			Username:  toString(row["username"]),
			Module:    toString(row["module"]),
			Action:    toString(row["action"]),
			Path:      toString(row["path"]),
			CreatedAt: toString(row["created_at"]),
		})
	}

	return res, nil
}

func (s *DashboardService) Trends(days int) (*model.Trends, error) {
	if days != 30 {
		days = 7
	}
	res := &model.Trends{Days: days}

	alertRows, err := s.repo.AlertDailyTrend(days)
	if err != nil {
		return nil, err
	}
	res.AlertTrend = make([]model.DailyTrend, 0, len(alertRows))
	for _, row := range alertRows {
		res.AlertTrend = append(res.AlertTrend, model.DailyTrend{Date: toString(row["date"]), Count: toInt64(row["count"])})
	}

	resourceRows, err := s.repo.ResourceDailyTrend(days)
	if err != nil {
		return nil, err
	}
	res.ResourceTrend = make([]model.ResourceTrend, 0, len(resourceRows))
	for _, row := range resourceRows {
		res.ResourceTrend = append(res.ResourceTrend, model.ResourceTrend{
			Date:   toString(row["date"]),
			CPU:    toFloat64(row["cpu"]),
			Memory: toFloat64(row["memory"]),
			Disk:   toFloat64(row["disk"]),
		})
	}

	return res, nil
}

func topNSeverity(items []model.SeverityStat, n int) []model.SeverityStat {
	if len(items) <= n {
		return items
	}
	return items[:n]
}

func toString(v any) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

func toInt64(v any) int64 {
	switch val := v.(type) {
	case int64:
		return val
	case int32:
		return int64(val)
	case int:
		return int64(val)
	case float64:
		return int64(val)
	case float32:
		return int64(val)
	default:
		return 0
	}
}

func toFloat64(v any) float64 {
	switch val := v.(type) {
	case float64:
		return val
	case float32:
		return float64(val)
	case int64:
		return float64(val)
	case int:
		return float64(val)
	default:
		return 0
	}
}
