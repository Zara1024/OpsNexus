package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/monitor/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/monitor/repository"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/crypto"
	apperrors "github.com/Zara1024/OpsNexus/cloudops-server/pkg/errors"
	"gorm.io/gorm"
)

// Services 聚合 monitor 模块服务。
type Services struct {
	AlertRule *AlertRuleService
	Alert     *AlertService
	Channel   *ChannelService
}

// New 创建服务集合。
func New(db *gorm.DB, aes *crypto.AESGCM) *Services {
	repos := repository.New(db)
	notifySvc := NewNotifyService()
	return &Services{
		AlertRule: NewAlertRuleService(repos.AlertRule),
		Alert:     NewAlertService(repos.Alert),
		Channel:   NewChannelService(repos.Channel, aes, notifySvc),
	}
}

type AlertRuleService struct {
	repo *repository.AlertRuleRepository
}

func NewAlertRuleService(repo *repository.AlertRuleRepository) *AlertRuleService {
	return &AlertRuleService{repo: repo}
}

type AlertRuleCreateRequest struct {
	Name           string `json:"name" binding:"required,max=128"`
	Description    string `json:"description"`
	Severity       string `json:"severity" binding:"required,oneof=P1 P2 P3 P4"`
	Expression     string `json:"expression" binding:"required"`
	Duration       string `json:"duration"`
	Labels         string `json:"labels"`
	Annotations    string `json:"annotations"`
	Enabled        *bool  `json:"enabled"`
	NotifyChannels string `json:"notify_channels"`
}

type AlertRuleUpdateRequest = AlertRuleCreateRequest

func (s *AlertRuleService) List() ([]model.AlertRule, error) { return s.repo.List() }

func (s *AlertRuleService) Create(req *AlertRuleCreateRequest) error {
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	return s.repo.Create(&model.AlertRule{
		Name:           strings.TrimSpace(req.Name),
		Description:    strings.TrimSpace(req.Description),
		Severity:       strings.TrimSpace(req.Severity),
		Expression:     strings.TrimSpace(req.Expression),
		Duration:       defaultString(req.Duration, "5m"),
		Labels:         defaultJSON(req.Labels, "{}"),
		Annotations:    defaultJSON(req.Annotations, "{}"),
		Enabled:        enabled,
		NotifyChannels: defaultJSON(req.NotifyChannels, "[]"),
	})
}

func (s *AlertRuleService) Update(id int64, req *AlertRuleUpdateRequest) error {
	if _, err := s.repo.GetByID(id); err != nil {
		return apperrors.ErrNotFound
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	return s.repo.Update(id, map[string]interface{}{
		"name":            strings.TrimSpace(req.Name),
		"description":     strings.TrimSpace(req.Description),
		"severity":        strings.TrimSpace(req.Severity),
		"expression":      strings.TrimSpace(req.Expression),
		"duration":        defaultString(req.Duration, "5m"),
		"labels":          defaultJSON(req.Labels, "{}"),
		"annotations":     defaultJSON(req.Annotations, "{}"),
		"enabled":         enabled,
		"notify_channels": defaultJSON(req.NotifyChannels, "[]"),
	})
}

func (s *AlertRuleService) Delete(id int64) error { return s.repo.Delete(id) }

func (s *AlertRuleService) Toggle(id int64, enabled bool) error {
	return s.repo.Update(id, map[string]interface{}{"enabled": enabled})
}

type AlertService struct {
	repo *repository.AlertEventRepository
}

func NewAlertService(repo *repository.AlertEventRepository) *AlertService {
	return &AlertService{repo: repo}
}

func (s *AlertService) ListActive() ([]model.AlertEvent, error)  { return s.repo.ListActive() }
func (s *AlertService) ListHistory() ([]model.AlertEvent, error) { return s.repo.ListHistory() }
func (s *AlertService) Ack(id, userID int64) error               { return s.repo.Ack(id, userID) }
func (s *AlertService) Silence(id int64, minute int) error {
	if minute <= 0 {
		minute = 30
	}
	return s.repo.Silence(id, time.Now().Add(time.Duration(minute)*time.Minute))
}

type ChannelService struct {
	repo      *repository.ChannelRepository
	aes       *crypto.AESGCM
	notifySvc *NotifyService
}

func NewChannelService(repo *repository.ChannelRepository, aes *crypto.AESGCM, notifySvc *NotifyService) *ChannelService {
	return &ChannelService{repo: repo, aes: aes, notifySvc: notifySvc}
}

type ChannelCreateRequest struct {
	Name    string            `json:"name" binding:"required,max=128"`
	Type    string            `json:"type" binding:"required,max=32"`
	Config  map[string]string `json:"config" binding:"required"`
	Enabled *bool             `json:"enabled"`
}

type ChannelItem struct {
	ID        int64             `json:"id"`
	Name      string            `json:"name"`
	Type      string            `json:"type"`
	Enabled   bool              `json:"enabled"`
	Config    map[string]string `json:"config"`
	CreatedAt time.Time         `json:"created_at"`
}

func (s *ChannelService) List() ([]ChannelItem, error) {
	list, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	res := make([]ChannelItem, 0, len(list))
	for _, item := range list {
		plain, _ := s.aes.Decrypt(item.ConfigEncrypted)
		cfg := map[string]string{}
		_ = json.Unmarshal([]byte(plain), &cfg)
		maskSecrets(cfg)
		res = append(res, ChannelItem{ID: item.ID, Name: item.Name, Type: item.Type, Enabled: item.Enabled, Config: cfg, CreatedAt: item.CreatedAt})
	}
	return res, nil
}

func (s *ChannelService) Create(req *ChannelCreateRequest) error {
	raw, _ := json.Marshal(req.Config)
	cipher, err := s.aes.Encrypt(string(raw))
	if err != nil {
		return err
	}
	enabled := true
	if req.Enabled != nil {
		enabled = *req.Enabled
	}
	return s.repo.Create(&model.NotificationChannel{Name: strings.TrimSpace(req.Name), Type: strings.TrimSpace(req.Type), ConfigEncrypted: cipher, Enabled: enabled})
}

func (s *ChannelService) Test(id int64) error {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return apperrors.ErrNotFound
	}
	plain, err := s.aes.Decrypt(item.ConfigEncrypted)
	if err != nil {
		return err
	}
	cfg := map[string]string{}
	if err = json.Unmarshal([]byte(plain), &cfg); err != nil {
		return err
	}
	return s.notifySvc.SendTest(item.Type, cfg)
}

func defaultString(v, d string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return d
	}
	return v
}

func defaultJSON(v, d string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return d
	}
	return v
}

func maskSecrets(cfg map[string]string) {
	for k, v := range cfg {
		lk := strings.ToLower(k)
		if strings.Contains(lk, "password") || strings.Contains(lk, "secret") || strings.Contains(lk, "token") {
			if v != "" {
				cfg[k] = "******"
			}
		}
	}
}
