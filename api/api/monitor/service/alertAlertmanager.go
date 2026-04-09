package service

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"dodevops-api/api/monitor/dao"
	"dodevops-api/api/monitor/model"
	"dodevops-api/common/result"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const alertManagerSourceType = 4

var alertManagerHTTPClient = &http.Client{
	Timeout: 15 * time.Second,
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

type alertManagerUpstreamError struct {
	StatusCode int
	Message    string
}

func (e *alertManagerUpstreamError) Error() string {
	if strings.TrimSpace(e.Message) == "" {
		return fmt.Sprintf("AlertManager upstream status=%d", e.StatusCode)
	}
	return e.Message
}

type alertManagerStatusResponse struct {
	Cluster struct {
		Name   string `json:"name"`
		Status string `json:"status"`
	} `json:"cluster"`
	VersionInfo struct {
		Version   string `json:"version"`
		Revision  string `json:"revision"`
		BuildDate string `json:"buildDate"`
	} `json:"versionInfo"`
	Uptime string `json:"uptime"`
}

type alertManagerSilenceResponse struct {
	ID       string `json:"id"`
	Matchers []struct {
		Name    string `json:"name"`
		Value   string `json:"value"`
		IsRegex bool   `json:"isRegex"`
		IsEqual bool   `json:"isEqual"`
	} `json:"matchers"`
	StartsAt  time.Time `json:"startsAt"`
	EndsAt    time.Time `json:"endsAt"`
	CreatedBy string    `json:"createdBy"`
	Comment   string    `json:"comment"`
	Status    struct {
		State string `json:"state"`
	} `json:"status"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type alertManagerReceiverResponse struct {
	Name         string                                         `json:"name"`
	Active       bool                                           `json:"active"`
	Integrations []model.MonitorAlertManagerReceiverIntegration `json:"integrations"`
}

type alertManagerCreateSilencePayload struct {
	Matchers  []model.MonitorAlertManagerMatcher `json:"matchers"`
	StartsAt  time.Time                          `json:"startsAt"`
	EndsAt    time.Time                          `json:"endsAt"`
	CreatedBy string                             `json:"createdBy"`
	Comment   string                             `json:"comment"`
}

func (s *MonitorAlertServiceImpl) GetAlertManagerStatus(c *gin.Context, query model.MonitorAlertManagerQuery) {
	source, err := resolveAlertManagerSource(query)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	var statusResp alertManagerStatusResponse
	if err = doAlertManagerRequest(source, http.MethodGet, "/api/v2/status", nil, &statusResp); err != nil {
		writeAlertManagerError(c, err, "get alertmanager status failed")
		return
	}

	result.Success(c, model.MonitorAlertManagerStatus{
		SourceID:      source.ID,
		SourceName:    source.Name,
		Endpoint:      source.APIBaseURL,
		Available:     true,
		ClusterName:   statusResp.Cluster.Name,
		ClusterStatus: statusResp.Cluster.Status,
		Version:       statusResp.VersionInfo.Version,
		Revision:      statusResp.VersionInfo.Revision,
		BuildDate:     statusResp.VersionInfo.BuildDate,
		Uptime:        normalizeAlertManagerTimeText(statusResp.Uptime),
	})
}

func (s *MonitorAlertServiceImpl) GetAlertManagerSilenceList(c *gin.Context, query model.MonitorAlertManagerQuery) {
	source, err := resolveAlertManagerSource(query)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	var items []alertManagerSilenceResponse
	if err = doAlertManagerRequest(source, http.MethodGet, "/api/v2/silences", nil, &items); err != nil {
		writeAlertManagerError(c, err, "get alertmanager silences failed")
		return
	}

	list := make([]model.MonitorAlertManagerSilence, 0, len(items))
	for _, item := range items {
		list = append(list, convertAlertManagerSilence(source, item))
	}
	sort.SliceStable(list, func(i, j int) bool {
		leftPriority := alertManagerSilencePriority(list[i].Status)
		rightPriority := alertManagerSilencePriority(list[j].Status)
		if leftPriority != rightPriority {
			return leftPriority < rightPriority
		}
		return list[i].UpdatedAt > list[j].UpdatedAt
	})

	result.Success(c, list)
}

func (s *MonitorAlertServiceImpl) CreateAlertManagerSilence(c *gin.Context, req model.MonitorAlertManagerSilenceCreateRequest) {
	source, err := resolveAlertManagerSource(model.MonitorAlertManagerQuery{SourceID: req.SourceID})
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	payload, err := buildAlertManagerSilencePayload(c, req)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	var createResp struct {
		SilenceID string `json:"silenceID"`
	}
	if err = doAlertManagerRequest(source, http.MethodPost, "/api/v2/silences", payload, &createResp); err != nil {
		writeAlertManagerError(c, err, "create alertmanager silence failed")
		return
	}

	result.Success(c, gin.H{
		"sourceId":   source.ID,
		"sourceName": source.Name,
		"silenceId":  createResp.SilenceID,
	})
}

func (s *MonitorAlertServiceImpl) DeleteAlertManagerSilence(c *gin.Context, query model.MonitorAlertManagerQuery, silenceID string) {
	source, err := resolveAlertManagerSource(query)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	if err = doAlertManagerRequest(source, http.MethodDelete, "/api/v2/silence/"+url.PathEscape(silenceID), nil, nil); err != nil {
		writeAlertManagerError(c, err, "delete alertmanager silence failed")
		return
	}

	result.Success(c, gin.H{
		"sourceId":   source.ID,
		"sourceName": source.Name,
		"silenceId":  silenceID,
	})
}

func (s *MonitorAlertServiceImpl) GetAlertManagerReceiverList(c *gin.Context, query model.MonitorAlertManagerQuery) {
	source, err := resolveAlertManagerSource(query)
	if err != nil {
		result.Failed(c, 400, err.Error())
		return
	}

	var items []alertManagerReceiverResponse
	if err = doAlertManagerRequest(source, http.MethodGet, "/api/v2/receivers", nil, &items); err != nil {
		writeAlertManagerError(c, err, "get alertmanager receivers failed")
		return
	}

	list := make([]model.MonitorAlertManagerReceiver, 0, len(items))
	for _, item := range items {
		list = append(list, model.MonitorAlertManagerReceiver{
			SourceID:     source.ID,
			SourceName:   source.Name,
			Name:         item.Name,
			Active:       item.Active,
			Integrations: item.Integrations,
		})
	}
	sort.SliceStable(list, func(i, j int) bool {
		return strings.ToLower(list[i].Name) < strings.ToLower(list[j].Name)
	})

	result.Success(c, list)
}

func resolveAlertManagerSource(query model.MonitorAlertManagerQuery) (*model.MonitorAlertSource, error) {
	if query.SourceID > 0 {
		source, err := dao.GetMonitorAlertSourceByID(query.SourceID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("alertmanager source not found")
			}
			return nil, err
		}
		if source.Type != alertManagerSourceType {
			return nil, errors.New("selected source is not alertmanager")
		}
		if source.Status != 1 {
			return nil, errors.New("selected alertmanager source is disabled")
		}
		if strings.TrimSpace(source.APIBaseURL) == "" {
			return nil, errors.New("alertmanager source apiBaseUrl is empty")
		}
		return source, nil
	}

	sources, err := dao.GetMonitorAlertSourcesByType(alertManagerSourceType, true)
	if err != nil {
		return nil, err
	}
	if len(sources) == 0 {
		return nil, errors.New("no enabled alertmanager source found")
	}
	if strings.TrimSpace(sources[0].APIBaseURL) == "" {
		return nil, errors.New("alertmanager source apiBaseUrl is empty")
	}
	return &sources[0], nil
}

func buildAlertManagerSilencePayload(c *gin.Context, req model.MonitorAlertManagerSilenceCreateRequest) (*alertManagerCreateSilencePayload, error) {
	matchers := make([]model.MonitorAlertManagerMatcher, 0, len(req.Matchers))
	for _, matcher := range req.Matchers {
		name := strings.TrimSpace(matcher.Name)
		value := strings.TrimSpace(matcher.Value)
		if name == "" || value == "" {
			return nil, errors.New("silence matcher name and value are required")
		}
		matchers = append(matchers, model.MonitorAlertManagerMatcher{
			Name:    name,
			Value:   value,
			IsRegex: matcher.IsRegex,
			IsEqual: !matcher.IsRegex,
		})
	}
	if len(matchers) == 0 {
		return nil, errors.New("at least one silence matcher is required")
	}

	now := time.Now()
	startAt, err := parseAlertManagerTime(req.StartsAt)
	if err != nil {
		return nil, err
	}
	if startAt.IsZero() {
		startAt = now
	}

	endAt, err := parseAlertManagerTime(req.EndsAt)
	if err != nil {
		return nil, err
	}
	if endAt.IsZero() {
		endAt = startAt.Add(2 * time.Hour)
	}
	if !endAt.After(startAt) {
		return nil, errors.New("silence end time must be later than start time")
	}

	createdBy := strings.TrimSpace(req.CreatedBy)
	if createdBy == "" {
		if username, err := jwt.GetAdminName(c); err == nil && strings.TrimSpace(username) != "" {
			createdBy = strings.TrimSpace(username)
		}
	}
	if createdBy == "" {
		createdBy = "OpsNexus"
	}

	comment := strings.TrimSpace(req.Comment)
	if comment == "" {
		return nil, errors.New("silence comment is required")
	}

	return &alertManagerCreateSilencePayload{
		Matchers:  matchers,
		StartsAt:  startAt,
		EndsAt:    endAt,
		CreatedBy: createdBy,
		Comment:   comment,
	}, nil
}

func parseAlertManagerTime(value string) (time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return time.Time{}, nil
	}

	layouts := []string{
		time.RFC3339,
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05",
		"2006-01-02 15:04",
		"2006-01-02T15:04",
	}
	for _, layout := range layouts {
		if parsed, err := time.Parse(layout, value); err == nil {
			return parsed, nil
		}
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed, nil
		}
	}
	return time.Time{}, fmt.Errorf("invalid time format: %s", value)
}

func doAlertManagerRequest(source *model.MonitorAlertSource, method, apiPath string, requestBody interface{}, responseBody interface{}) error {
	endpoint, err := buildAlertManagerURL(source.APIBaseURL, apiPath)
	if err != nil {
		return err
	}

	var bodyReader io.Reader
	if requestBody != nil {
		payload, marshalErr := json.Marshal(requestBody)
		if marshalErr != nil {
			return marshalErr
		}
		bodyReader = bytes.NewReader(payload)
	}

	req, err := http.NewRequest(method, endpoint, bodyReader)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	if requestBody != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	applyAlertManagerAuth(req, source.AppKey)

	resp, err := alertManagerHTTPClient.Do(req)
	if err != nil {
		return &alertManagerUpstreamError{
			StatusCode: http.StatusBadGateway,
			Message:    "alertmanager request failed: " + err.Error(),
		}
	}
	defer resp.Body.Close()

	body, readErr := io.ReadAll(io.LimitReader(resp.Body, 1<<20))
	if readErr != nil {
		return readErr
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		message := strings.TrimSpace(string(body))
		if message == "" {
			message = resp.Status
		}
		return &alertManagerUpstreamError{
			StatusCode: resp.StatusCode,
			Message:    fmt.Sprintf("alertmanager upstream error: %s", message),
		}
	}

	if responseBody == nil || len(bytes.TrimSpace(body)) == 0 {
		return nil
	}
	if err = json.Unmarshal(body, responseBody); err != nil {
		return fmt.Errorf("decode alertmanager response failed: %w", err)
	}
	return nil
}

func buildAlertManagerURL(baseURL, apiPath string) (string, error) {
	parsed, err := url.Parse(strings.TrimSpace(baseURL))
	if err != nil {
		return "", fmt.Errorf("invalid alertmanager url: %w", err)
	}
	if parsed.Scheme == "" || parsed.Host == "" {
		return "", errors.New("invalid alertmanager url")
	}

	parsed.Path = strings.TrimRight(parsed.Path, "/") + "/" + strings.TrimLeft(apiPath, "/")
	parsed.RawQuery = ""
	parsed.Fragment = ""
	return parsed.String(), nil
}

func applyAlertManagerAuth(req *http.Request, appKey string) {
	appKey = strings.TrimSpace(appKey)
	if appKey == "" {
		return
	}
	if strings.HasPrefix(appKey, "Bearer ") || strings.HasPrefix(appKey, "Basic ") {
		req.Header.Set("Authorization", appKey)
		return
	}
	req.Header.Set("Authorization", "Bearer "+appKey)
}

func convertAlertManagerSilence(source *model.MonitorAlertSource, item alertManagerSilenceResponse) model.MonitorAlertManagerSilence {
	matchers := make([]model.MonitorAlertManagerMatcher, 0, len(item.Matchers))
	for _, matcher := range item.Matchers {
		matchers = append(matchers, model.MonitorAlertManagerMatcher{
			Name:    matcher.Name,
			Value:   matcher.Value,
			IsRegex: matcher.IsRegex,
			IsEqual: matcher.IsEqual,
		})
	}

	return model.MonitorAlertManagerSilence{
		SourceID:   source.ID,
		SourceName: source.Name,
		ID:         item.ID,
		Matchers:   matchers,
		StartsAt:   formatAlertManagerTime(item.StartsAt),
		EndsAt:     formatAlertManagerTime(item.EndsAt),
		CreatedBy:  item.CreatedBy,
		Comment:    item.Comment,
		Status:     normalizeAlertManagerSilenceStatus(item.Status.State, item.EndsAt),
		UpdatedAt:  formatAlertManagerTime(item.UpdatedAt),
	}
}

func normalizeAlertManagerSilenceStatus(status string, endsAt time.Time) string {
	status = strings.TrimSpace(strings.ToLower(status))
	if status != "" {
		return status
	}
	if !endsAt.IsZero() && endsAt.Before(time.Now()) {
		return "expired"
	}
	return "active"
}

func formatAlertManagerTime(value time.Time) string {
	if value.IsZero() {
		return ""
	}
	return value.In(time.Local).Format("2006-01-02 15:04:05")
}

func normalizeAlertManagerTimeText(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return ""
	}
	if parsed, err := time.Parse(time.RFC3339, value); err == nil {
		return formatAlertManagerTime(parsed)
	}
	return value
}

func alertManagerSilencePriority(status string) int {
	switch strings.ToLower(strings.TrimSpace(status)) {
	case "active":
		return 0
	case "pending":
		return 1
	case "expired":
		return 2
	default:
		return 3
	}
}

func writeAlertManagerError(c *gin.Context, err error, fallback string) {
	var upstreamErr *alertManagerUpstreamError
	if errors.As(err, &upstreamErr) {
		result.Failed(c, http.StatusBadGateway, upstreamErr.Error())
		return
	}

	message := strings.TrimSpace(err.Error())
	if message == "" {
		message = fallback
	}
	result.Failed(c, 500, message)
}
