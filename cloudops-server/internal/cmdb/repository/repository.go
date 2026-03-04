package repository

import "gorm.io/gorm"

// Repositories 聚合 cmdb 模块仓储。
type Repositories struct {
	Host       *HostRepository
	Group      *GroupRepository
	Credential *CredentialRepository
	SSHRecord  *SSHRecordRepository
}

// New 创建 cmdb 模块仓储集合。
func New(db *gorm.DB) *Repositories {
	return &Repositories{
		Host:       NewHostRepository(db),
		Group:      NewGroupRepository(db),
		Credential: NewCredentialRepository(db),
		SSHRecord:  NewSSHRecordRepository(db),
	}
}
