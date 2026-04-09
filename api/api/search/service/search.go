package service

import (
	"errors"
	"sort"
	"strings"

	searchDao "dodevops-api/api/search/dao"
	searchModel "dodevops-api/api/search/model"
	systemDao "dodevops-api/api/system/dao"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
)

var globalSearchTypeOrder = []string{"cluster", "host", "application", "alert", "user", "menu"}

type GlobalSearchService interface {
	GlobalSearch(c *gin.Context, query searchModel.GlobalSearchQuery)
}

type GlobalSearchServiceImpl struct{}

func NewGlobalSearchService() GlobalSearchService {
	return &GlobalSearchServiceImpl{}
}

func (s *GlobalSearchServiceImpl) GlobalSearch(c *gin.Context, query searchModel.GlobalSearchQuery) {
	query.Keyword = strings.TrimSpace(query.Keyword)
	if query.Keyword == "" {
		result.Failed(c, 400, "搜索关键字不能为空")
		return
	}

	if query.Limit <= 0 {
		query.Limit = 8
	}
	if query.Limit > 20 {
		query.Limit = 20
	}

	searchTypes := normalizeSearchTypes(query.Types)
	adminID, _ := jwt.GetAdminId(c)
	accessiblePaths := buildAccessiblePathSet(adminID)

	typeResults := make(map[string][]searchModel.GlobalSearchResult)
	for _, resourceType := range searchTypes {
		items, err := searchByType(resourceType, query.Keyword, query.Limit)
		if err != nil {
			result.Failed(c, 500, "全局搜索失败: "+err.Error())
			return
		}

		filtered := make([]searchModel.GlobalSearchResult, 0, len(items))
		for _, item := range items {
			accessPath := item.AccessPath
			if accessPath == "" {
				accessPath = item.Route
			}
			if accessPath != "" && !isPathAccessible(accessPath, accessiblePaths) {
				continue
			}
			filtered = append(filtered, item)
		}

		if len(filtered) > 0 {
			typeResults[resourceType] = filtered
		}
	}

	groups := make([]searchModel.GlobalSearchGroup, 0, len(typeResults))
	results := make([]searchModel.GlobalSearchResult, 0)
	for _, resourceType := range globalSearchTypeOrder {
		items := typeResults[resourceType]
		if len(items) == 0 {
			continue
		}
		groups = append(groups, searchModel.GlobalSearchGroup{
			Type:      resourceType,
			TypeLabel: items[0].TypeLabel,
			Count:     len(items),
			Items:     items,
		})
		results = append(results, items...)
	}

	sort.SliceStable(results, func(i, j int) bool {
		leftIndex := searchTypeIndex(results[i].Type)
		rightIndex := searchTypeIndex(results[j].Type)
		if leftIndex != rightIndex {
			return leftIndex < rightIndex
		}
		return strings.ToLower(results[i].Title) < strings.ToLower(results[j].Title)
	})

	result.Success(c, searchModel.GlobalSearchResponse{
		Keyword: query.Keyword,
		Total:   len(results),
		Groups:  groups,
		Results: results,
	})
}

func searchByType(resourceType, keyword string, limit int) ([]searchModel.GlobalSearchResult, error) {
	switch resourceType {
	case "cluster":
		return searchDao.SearchClusters(keyword, limit)
	case "host":
		return searchDao.SearchHosts(keyword, limit)
	case "application":
		return searchDao.SearchApplications(keyword, limit)
	case "alert":
		return searchDao.SearchAlerts(keyword, limit)
	case "user":
		return searchDao.SearchAdmins(keyword, limit)
	case "menu":
		return searchDao.SearchMenus(keyword, limit)
	default:
		return nil, errors.New("不支持的搜索类型: " + resourceType)
	}
}

func normalizeSearchTypes(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return append([]string(nil), globalSearchTypeOrder...)
	}

	allowed := make(map[string]struct{}, len(globalSearchTypeOrder))
	for _, item := range globalSearchTypeOrder {
		allowed[item] = struct{}{}
	}

	results := make([]string, 0)
	seen := make(map[string]struct{})
	for _, item := range strings.Split(raw, ",") {
		key := strings.ToLower(strings.TrimSpace(item))
		if key == "" {
			continue
		}
		if _, ok := allowed[key]; !ok {
			continue
		}
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		results = append(results, key)
	}
	if len(results) == 0 {
		return append([]string(nil), globalSearchTypeOrder...)
	}
	return results
}

func buildAccessiblePathSet(adminID uint) map[string]struct{} {
	paths := map[string]struct{}{
		"/dashboard":             {},
		"/search/global":         {},
		"/monitor/alert-center":  {},
		"/monitor/alert-notify":  {},
		"/monitor/alert-history": {},
		"/task/config":           {},
	}
	if adminID == 0 {
		return paths
	}

	leftMenus := systemDao.QueryLeftMenuList(adminID)
	for _, item := range leftMenus {
		if item.Url != "" {
			paths[ensureRoutePath(item.Url)] = struct{}{}
		}
		children := systemDao.QueryMenuVoList(adminID, item.Id)
		for _, child := range children {
			if child.Url != "" {
				paths[ensureRoutePath(child.Url)] = struct{}{}
			}
		}
	}
	return paths
}

func ensureRoutePath(path string) string {
	path = strings.TrimSpace(path)
	if path == "" {
		return ""
	}
	if strings.HasPrefix(path, "/") {
		return path
	}
	return "/" + path
}

func isPathAccessible(path string, accessiblePaths map[string]struct{}) bool {
	path = ensureRoutePath(path)
	if path == "" {
		return true
	}
	if _, ok := accessiblePaths[path]; ok {
		return true
	}
	for candidate := range accessiblePaths {
		if candidate != "/" && strings.HasPrefix(path, candidate+"/") {
			return true
		}
	}
	return false
}

func searchTypeIndex(resourceType string) int {
	for index, item := range globalSearchTypeOrder {
		if item == resourceType {
			return index
		}
	}
	return len(globalSearchTypeOrder)
}
