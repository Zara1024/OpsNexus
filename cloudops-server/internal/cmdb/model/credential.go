package model

import "time"

// Credential 定义主机登录凭据。
type Credential struct {
	ID                  int64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	Name                string     `gorm:"column:name;size:128;not null" json:"name"`
	Type                string     `gorm:"column:type;size:16;not null" json:"type"`
	Username            string     `gorm:"column:username;size:128;not null;default:''" json:"username"`
	PasswordEncrypted   string     `gorm:"column:password_encrypted;type:text;not null;default:''" json:"password_encrypted"`
	PrivateKeyEncrypted string     `gorm:"column:private_key_encrypted;type:text;not null;default:''" json:"private_key_encrypted"`
	PassphraseEncrypted string     `gorm:"column:passphrase_encrypted;type:text;not null;default:''" json:"passphrase_encrypted"`
	CreatedAt           time.Time  `gorm:"column:created_at" json:"created_at"`
	UpdatedAt           time.Time  `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt           *time.Time `gorm:"column:deleted_at" json:"-"`
}

// TableName 指定表名。
func (Credential) TableName() string {
	return "cmdb_credentials"
}
