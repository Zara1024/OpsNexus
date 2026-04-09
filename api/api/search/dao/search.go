package dao

import (
	"fmt"
	"strconv"
	"strings"

	searchModel "dodevops-api/api/search/model"
	. "dodevops-api/pkg/db"
)

func SearchClusters(keyword string, limit int) ([]searchModel.GlobalSearchResult, error) {
	type row struct {
		ID          uint
		Name        string
		Version     string
		Status      int
		Description string
	}

	var rows []row
	like := buildSearchLike(keyword)
	err := Db.Table("k8s_cluster").
		Select(`
			id,
			COALESCE(name, '') AS name,
			COALESCE(version, '') AS version,
			COALESCE(status, 0) AS status,
			COALESCE(description, '') AS description
		`).
		Where("name LIKE ? OR version LIKE ? OR description LIKE ?", like, like, like).
		Order("updated_at DESC, id DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	results := make([]searchModel.GlobalSearchResult, 0, len(rows))
	for _, item := range rows {
		id := strconv.FormatUint(uint64(item.ID), 10)
		results = append(results, searchModel.GlobalSearchResult{
			Type:        "cluster",
			TypeLabel:   "K8s 集群",
			ID:          id,
			Title:       item.Name,
			Subtitle:    compactText("版本 "+item.Version, clusterStatusText(item.Status)),
			Description: defaultText(item.Description, "Kubernetes 集群"),
			Status:      clusterStatusText(item.Status),
			Route:       "/k8s/cluster/" + id,
			Tags:        compactTags(item.Version),
			AccessPath:  "/k8s/list",
		})
	}
	return results, nil
}

func SearchHosts(keyword string, limit int) ([]searchModel.GlobalSearchResult, error) {
	type row struct {
		ID        uint
		HostName  string
		Name      string
		GroupName string
		PrivateIP string
		PublicIP  string
		SSHIP     string
		OS        string
		Status    int
	}

	var rows []row
	like := buildSearchLike(keyword)
	err := Db.Table("cmdb_host ch").
		Select(`
			ch.id,
			COALESCE(ch.host_name, '') AS host_name,
			COALESCE(ch.name, '') AS name,
			COALESCE(cg.name, '') AS group_name,
			COALESCE(ch.private_ip, '') AS private_ip,
			COALESCE(ch.public_ip, '') AS public_ip,
			COALESCE(ch.ssh_ip, '') AS ssh_ip,
			COALESCE(ch.os, '') AS os,
			COALESCE(ch.status, 0) AS status
		`).
		Joins("LEFT JOIN cmdb_group cg ON cg.id = ch.group_id").
		Where(`
			COALESCE(ch.host_name, '') LIKE ? OR
			COALESCE(ch.name, '') LIKE ? OR
			COALESCE(ch.private_ip, '') LIKE ? OR
			COALESCE(ch.public_ip, '') LIKE ? OR
			COALESCE(ch.ssh_ip, '') LIKE ? OR
			COALESCE(cg.name, '') LIKE ?
		`, like, like, like, like, like, like).
		Order("ch.update_time DESC, ch.id DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	results := make([]searchModel.GlobalSearchResult, 0, len(rows))
	for _, item := range rows {
		title := defaultText(item.HostName, item.Name)
		description := compactText(
			defaultText(item.GroupName, ""),
			defaultText(item.PrivateIP, defaultText(item.PublicIP, item.SSHIP)),
		)
		results = append(results, searchModel.GlobalSearchResult{
			Type:        "host",
			TypeLabel:   "主机",
			ID:          strconv.FormatUint(uint64(item.ID), 10),
			Title:       title,
			Subtitle:    compactText(defaultText(item.Name, ""), hostStatusText(item.Status)),
			Description: defaultText(description, "CMDB 主机"),
			Status:      hostStatusText(item.Status),
			Route:       "/cmdb/ecs",
			Tags:        compactTags(item.OS, item.PrivateIP, item.PublicIP),
			AccessPath:  "/cmdb/ecs",
		})
	}
	return results, nil
}

func SearchApplications(keyword string, limit int) ([]searchModel.GlobalSearchResult, error) {
	type row struct {
		ID              uint
		Name            string
		Code            string
		Description     string
		RepoURL         string
		ProgrammingLang string
		Status          int
	}

	var rows []row
	like := buildSearchLike(keyword)
	err := Db.Table("app_application").
		Select(`
			id,
			COALESCE(name, '') AS name,
			COALESCE(code, '') AS code,
			COALESCE(description, '') AS description,
			COALESCE(repo_url, '') AS repo_url,
			COALESCE(programming_lang, '') AS programming_lang,
			COALESCE(status, 0) AS status
		`).
		Where(`
			COALESCE(name, '') LIKE ? OR
			COALESCE(code, '') LIKE ? OR
			COALESCE(description, '') LIKE ? OR
			COALESCE(repo_url, '') LIKE ? OR
			COALESCE(programming_lang, '') LIKE ?
		`, like, like, like, like, like).
		Order("updated_at DESC, id DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	results := make([]searchModel.GlobalSearchResult, 0, len(rows))
	for _, item := range rows {
		results = append(results, searchModel.GlobalSearchResult{
			Type:        "application",
			TypeLabel:   "应用",
			ID:          strconv.FormatUint(uint64(item.ID), 10),
			Title:       item.Name,
			Subtitle:    compactText(defaultText(item.Code, ""), appStatusText(item.Status)),
			Description: defaultText(item.Description, defaultText(item.RepoURL, "应用管理记录")),
			Status:      appStatusText(item.Status),
			Route:       "/app/application",
			Tags:        compactTags(item.ProgrammingLang, item.Code),
			AccessPath:  "/app/application",
		})
	}
	return results, nil
}

func SearchAlerts(keyword string, limit int) ([]searchModel.GlobalSearchResult, error) {
	type row struct {
		ID           uint
		AlertDesc    string
		BusinessLine string
		AlertLevel   string
		Status       int
		Handler      string
		DetailURL    string
	}

	var rows []row
	like := buildSearchLike(keyword)
	err := Db.Table("monitor_incident").
		Select(`
			id,
			COALESCE(alert_desc, '') AS alert_desc,
			COALESCE(business_line, '') AS business_line,
			COALESCE(alert_level, '') AS alert_level,
			COALESCE(status, 0) AS status,
			COALESCE(handler, '') AS handler,
			COALESCE(detail_url, '') AS detail_url
		`).
		Where(`
			COALESCE(alert_desc, '') LIKE ? OR
			COALESCE(business_line, '') LIKE ? OR
			COALESCE(alert_level, '') LIKE ? OR
			COALESCE(handler, '') LIKE ?
		`, like, like, like, like).
		Order("update_time DESC, id DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	results := make([]searchModel.GlobalSearchResult, 0, len(rows))
	for _, item := range rows {
		results = append(results, searchModel.GlobalSearchResult{
			Type:        "alert",
			TypeLabel:   "告警事件",
			ID:          strconv.FormatUint(uint64(item.ID), 10),
			Title:       defaultText(item.AlertDesc, fmt.Sprintf("告警事件 #%d", item.ID)),
			Subtitle:    compactText(defaultText(item.BusinessLine, ""), incidentStatusText(item.Status)),
			Description: defaultText(item.DetailURL, defaultText(item.Handler, "监控告警事件")),
			Status:      incidentStatusText(item.Status),
			Route:       "/monitor/alert-center",
			Tags:        compactTags(item.AlertLevel, item.Handler),
			AccessPath:  "/monitor/alert-center",
		})
	}
	return results, nil
}

func SearchAdmins(keyword string, limit int) ([]searchModel.GlobalSearchResult, error) {
	type row struct {
		ID       uint
		Username string
		Nickname string
		Email    string
		Phone    string
		Status   int
	}

	var rows []row
	like := buildSearchLike(keyword)
	err := Db.Table("sys_admin").
		Select(`
			id,
			COALESCE(username, '') AS username,
			COALESCE(nickname, '') AS nickname,
			COALESCE(email, '') AS email,
			COALESCE(phone, '') AS phone,
			COALESCE(status, 0) AS status
		`).
		Where(`
			COALESCE(username, '') LIKE ? OR
			COALESCE(nickname, '') LIKE ? OR
			COALESCE(email, '') LIKE ? OR
			COALESCE(phone, '') LIKE ?
		`, like, like, like, like).
		Order("create_time DESC, id DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	results := make([]searchModel.GlobalSearchResult, 0, len(rows))
	for _, item := range rows {
		results = append(results, searchModel.GlobalSearchResult{
			Type:        "user",
			TypeLabel:   "用户",
			ID:          strconv.FormatUint(uint64(item.ID), 10),
			Title:       defaultText(item.Username, item.Nickname),
			Subtitle:    compactText(defaultText(item.Nickname, ""), adminStatusText(item.Status)),
			Description: defaultText(item.Email, defaultText(item.Phone, "系统用户")),
			Status:      adminStatusText(item.Status),
			Route:       "/system/admin",
			Tags:        compactTags(item.Email, item.Phone),
			AccessPath:  "/system/admin",
		})
	}
	return results, nil
}

func SearchMenus(keyword string, limit int) ([]searchModel.GlobalSearchResult, error) {
	type row struct {
		ID         uint
		MenuName   string
		URL        string
		Value      string
		MenuType   uint
		MenuStatus uint
	}

	var rows []row
	like := buildSearchLike(keyword)
	err := Db.Table("sys_menu").
		Select(`
			id,
			COALESCE(menu_name, '') AS menu_name,
			COALESCE(url, '') AS url,
			COALESCE(value, '') AS value,
			COALESCE(menu_type, 0) AS menu_type,
			COALESCE(menu_status, 0) AS menu_status
		`).
		Where(`
			COALESCE(menu_name, '') LIKE ? OR
			COALESCE(url, '') LIKE ? OR
			COALESCE(value, '') LIKE ?
		`, like, like, like).
		Order("sort ASC, id DESC").
		Limit(limit).
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	results := make([]searchModel.GlobalSearchResult, 0, len(rows))
	for _, item := range rows {
		route := "/system/menu"
		accessPath := "/system/menu"
		if strings.TrimSpace(item.URL) != "" {
			route = ensurePath(item.URL)
			accessPath = route
		}
		results = append(results, searchModel.GlobalSearchResult{
			Type:        "menu",
			TypeLabel:   "菜单",
			ID:          strconv.FormatUint(uint64(item.ID), 10),
			Title:       item.MenuName,
			Subtitle:    compactText(defaultText(item.URL, item.Value), menuTypeText(item.MenuType)),
			Description: compactText(menuStatusText(item.MenuStatus), defaultText(item.Value, "系统菜单")),
			Status:      menuStatusText(item.MenuStatus),
			Route:       route,
			Tags:        compactTags(item.Value),
			AccessPath:  accessPath,
		})
	}
	return results, nil
}

func buildSearchLike(keyword string) string {
	return "%" + strings.TrimSpace(keyword) + "%"
}

func compactText(parts ...string) string {
	items := compactTags(parts...)
	return strings.Join(items, " | ")
}

func compactTags(parts ...string) []string {
	items := make([]string, 0, len(parts))
	seen := make(map[string]struct{})
	for _, part := range parts {
		text := strings.TrimSpace(part)
		if text == "" {
			continue
		}
		if _, ok := seen[text]; ok {
			continue
		}
		seen[text] = struct{}{}
		items = append(items, text)
	}
	return items
}

func defaultText(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func ensurePath(url string) string {
	url = strings.TrimSpace(url)
	if url == "" {
		return ""
	}
	if strings.HasPrefix(url, "/") {
		return url
	}
	return "/" + url
}

func clusterStatusText(status int) string {
	switch status {
	case 1:
		return "创建中"
	case 2:
		return "运行中"
	case 3:
		return "离线"
	default:
		return "未知"
	}
}

func hostStatusText(status int) string {
	switch status {
	case 1:
		return "认证成功"
	case 2:
		return "未认证"
	case 3:
		return "认证失败"
	default:
		return "未知"
	}
}

func appStatusText(status int) string {
	switch status {
	case 1:
		return "启用"
	case 2:
		return "禁用"
	default:
		return "未知"
	}
}

func incidentStatusText(status int) string {
	switch status {
	case 1:
		return "待处理"
	case 2:
		return "处理中"
	case 3:
		return "已归档"
	default:
		return "未知"
	}
}

func adminStatusText(status int) string {
	switch status {
	case 1:
		return "启用"
	case 2:
		return "禁用"
	default:
		return "未知"
	}
}

func menuTypeText(menuType uint) string {
	switch menuType {
	case 1:
		return "目录"
	case 2:
		return "菜单"
	case 3:
		return "按钮"
	default:
		return "未知"
	}
}

func menuStatusText(status uint) string {
	switch status {
	case 1:
		return "禁用"
	case 2:
		return "启用"
	default:
		return "未知"
	}
}
