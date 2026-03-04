package service

import (
	"strings"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/model"
	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/repository"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/crypto"
)

// CredentialService 提供凭据管理能力。
type CredentialService struct {
	repo *repository.CredentialRepository
	aes  *crypto.AESGCM
}

// CredentialCreateRequest 定义创建凭据请求。
type CredentialCreateRequest struct {
	Name       string `json:"name" binding:"required,max=128"`
	Type       string `json:"type" binding:"required,oneof=password key"`
	Username   string `json:"username" binding:"required,max=128"`
	Password   string `json:"password"`
	PrivateKey string `json:"private_key"`
	Passphrase string `json:"passphrase"`
}

// CredentialUpdateRequest 定义更新凭据请求。
type CredentialUpdateRequest = CredentialCreateRequest

// CredentialView 定义凭据展示结构。
type CredentialView struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
}

// NewCredentialService 创建凭据服务。
func NewCredentialService(repo *repository.CredentialRepository, aes *crypto.AESGCM) *CredentialService {
	return &CredentialService{repo: repo, aes: aes}
}

// CreateCredential 创建凭据。
func (s *CredentialService) CreateCredential(req *CredentialCreateRequest) error {
	pwd, err := s.aes.Encrypt(strings.TrimSpace(req.Password))
	if err != nil {
		return err
	}
	key, err := s.aes.Encrypt(strings.TrimSpace(req.PrivateKey))
	if err != nil {
		return err
	}
	passphrase, err := s.aes.Encrypt(strings.TrimSpace(req.Passphrase))
	if err != nil {
		return err
	}
	return s.repo.Create(&model.Credential{
		Name:                strings.TrimSpace(req.Name),
		Type:                strings.TrimSpace(req.Type),
		Username:            strings.TrimSpace(req.Username),
		PasswordEncrypted:   pwd,
		PrivateKeyEncrypted: key,
		PassphraseEncrypted: passphrase,
	})
}

// UpdateCredential 更新凭据。
func (s *CredentialService) UpdateCredential(id int64, req *CredentialUpdateRequest) error {
	pwd, err := s.aes.Encrypt(strings.TrimSpace(req.Password))
	if err != nil {
		return err
	}
	key, err := s.aes.Encrypt(strings.TrimSpace(req.PrivateKey))
	if err != nil {
		return err
	}
	passphrase, err := s.aes.Encrypt(strings.TrimSpace(req.Passphrase))
	if err != nil {
		return err
	}
	return s.repo.Update(id, map[string]interface{}{
		"name":                  strings.TrimSpace(req.Name),
		"type":                  strings.TrimSpace(req.Type),
		"username":              strings.TrimSpace(req.Username),
		"password_encrypted":    pwd,
		"private_key_encrypted": key,
		"passphrase_encrypted":  passphrase,
	})
}

// DeleteCredential 删除凭据。
func (s *CredentialService) DeleteCredential(id int64) error {
	return s.repo.SoftDelete(id)
}

// ListCredentials 查询凭据列表。
func (s *CredentialService) ListCredentials() ([]CredentialView, error) {
	list, err := s.repo.List()
	if err != nil {
		return nil, err
	}
	result := make([]CredentialView, 0, len(list))
	for _, item := range list {
		result = append(result, CredentialView{
			ID:        item.ID,
			Name:      item.Name,
			Type:      item.Type,
			Username:  item.Username,
			CreatedAt: item.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return result, nil
}

// ResolvePlainCredential 解密并返回可用于连接的凭据。
func (s *CredentialService) ResolvePlainCredential(id int64) (username, password, privateKey, passphrase string, err error) {
	item, err := s.repo.GetByID(id)
	if err != nil {
		return "", "", "", "", err
	}
	password, err = s.aes.Decrypt(item.PasswordEncrypted)
	if err != nil {
		return "", "", "", "", err
	}
	privateKey, err = s.aes.Decrypt(item.PrivateKeyEncrypted)
	if err != nil {
		return "", "", "", "", err
	}
	passphrase, err = s.aes.Decrypt(item.PassphraseEncrypted)
	if err != nil {
		return "", "", "", "", err
	}
	return item.Username, password, privateKey, passphrase, nil
}
