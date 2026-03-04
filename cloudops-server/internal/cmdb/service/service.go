package service

import (
	"gorm.io/gorm"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/cmdb/repository"
	"github.com/Zara1024/OpsNexus/cloudops-server/pkg/crypto"
)

// Services 聚合 cmdb 模块业务服务。
type Services struct {
	Host       *HostService
	Group      *GroupService
	Credential *CredentialService
	Terminal   *TerminalService
}

// New 创建 cmdb 模块服务集合。
func New(db *gorm.DB, aes *crypto.AESGCM) *Services {
	repos := repository.New(db)
	credentialSvc := NewCredentialService(repos.Credential, aes)
	return &Services{
		Host:       NewHostService(repos.Host, repos.Credential, credentialSvc),
		Group:      NewGroupService(repos.Group),
		Credential: credentialSvc,
		Terminal:   NewTerminalService(repos.Host, repos.Credential, repos.SSHRecord, credentialSvc),
	}
}
