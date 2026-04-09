package service

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	systemDao "dodevops-api/api/system/dao"
	systemModel "dodevops-api/api/system/model"
	"dodevops-api/common"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
	"github.com/go-ldap/ldap/v3"
	"gorm.io/gorm"
)

const ldapConfigKey = "ldap"

type LDAPService struct{}

type LDAPAuthUser struct {
	Username string   `json:"username"`
	DN       string   `json:"dn"`
	Nickname string   `json:"nickname"`
	Email    string   `json:"email"`
	Phone    string   `json:"phone"`
	Groups   []string `json:"groups"`
}

func NewLDAPService() *LDAPService {
	return &LDAPService{}
}

func (s *LDAPService) GetConfig(c *gin.Context) {
	config, err := LoadLDAPConfig()
	if err != nil {
		result.Failed(c, 500, "获取 LDAP 配置失败: "+err.Error())
		return
	}
	result.Success(c, config)
}

func (s *LDAPService) UpdateConfig(c *gin.Context, config systemModel.LDAPConfig) {
	if err := validateLDAPConfig(config); err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	if err := SaveLDAPConfig(config); err != nil {
		result.Failed(c, 500, "保存 LDAP 配置失败: "+err.Error())
		return
	}
	result.Success(c, config)
}

func (s *LDAPService) TestConfig(c *gin.Context, config systemModel.LDAPConfig) {
	if err := validateLDAPConfig(config); err != nil {
		result.Failed(c, 400, err.Error())
		return
	}
	conn, err := dialLDAP(config)
	if err != nil {
		result.Failed(c, 500, "连接 LDAP 失败: "+err.Error())
		return
	}
	defer conn.Close()

	if strings.TrimSpace(config.BindUser) != "" {
		if err = conn.Bind(strings.TrimSpace(config.BindUser), config.BindPass); err != nil {
			result.Failed(c, 500, "LDAP 绑定失败: "+err.Error())
			return
		}
	}

	result.Success(c, gin.H{
		"connected":   true,
		"host":        config.Host,
		"port":        config.Port,
		"baseDn":      config.BaseDN,
		"groupFilter": config.GroupFilter,
	})
}

func LoadLDAPConfig() (systemModel.LDAPConfig, error) {
	defaultConfig := defaultLDAPConfig()
	item, err := systemDao.GetSysConfigByKey(ldapConfigKey)
	if err != nil {
		return defaultConfig, err
	}
	if item == nil || strings.TrimSpace(item.ConfigData) == "" {
		return defaultConfig, nil
	}

	config := defaultConfig
	if err = json.Unmarshal([]byte(item.ConfigData), &config); err != nil {
		return defaultConfig, err
	}
	if strings.TrimSpace(config.Remark) == "" {
		config.Remark = item.Remark
	}
	return config, nil
}

func SaveLDAPConfig(config systemModel.LDAPConfig) error {
	payload, err := json.Marshal(config)
	if err != nil {
		return err
	}
	return systemDao.UpsertSysConfig(&systemModel.SysConfig{
		ConfigKey:  ldapConfigKey,
		ConfigType: ldapConfigKey,
		ConfigData: string(payload),
		Status:     boolToStatus(config.Enable),
		Remark:     config.Remark,
	})
}

func AuthenticateByLDAP(username, password string) (*systemModel.SysAdmin, error) {
	config, err := LoadLDAPConfig()
	if err != nil {
		return nil, err
	}
	if !config.Enable {
		return nil, errors.New("LDAP 未启用")
	}

	authUser, err := authenticateLDAPUser(config, username, password)
	if err != nil {
		return nil, err
	}
	return syncLDAPUser(config, authUser, password)
}

func authenticateLDAPUser(config systemModel.LDAPConfig, username, password string) (*LDAPAuthUser, error) {
	if strings.TrimSpace(password) == "" {
		return nil, errors.New("LDAP 密码不能为空")
	}

	conn, err := dialLDAP(config)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	if strings.TrimSpace(config.BindUser) != "" {
		if err = conn.Bind(strings.TrimSpace(config.BindUser), config.BindPass); err != nil {
			return nil, fmt.Errorf("LDAP 管理员绑定失败: %w", err)
		}
	}

	filter := strings.TrimSpace(config.AuthFilter)
	if filter == "" {
		filter = "(uid=%s)"
	}
	filter = fmt.Sprintf(filter, ldap.EscapeFilter(username))

	req := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		1,
		10,
		false,
		filter,
		ldapSearchAttributes(config),
		nil,
	)
	resp, err := conn.Search(req)
	if err != nil {
		return nil, fmt.Errorf("LDAP 搜索失败: %w", err)
	}
	if len(resp.Entries) == 0 {
		return nil, errors.New("LDAP 用户不存在")
	}

	entry := resp.Entries[0]
	userConn, err := dialLDAP(config)
	if err != nil {
		return nil, err
	}
	defer userConn.Close()
	if err = userConn.Bind(entry.DN, password); err != nil {
		return nil, errors.New("LDAP 用户认证失败")
	}

	return &LDAPAuthUser{
		Username: username,
		DN:       entry.DN,
		Nickname: firstNonEmpty(entry.GetAttributeValue(config.Attributes.Nickname), username),
		Email:    entry.GetAttributeValue(config.Attributes.Email),
		Phone:    entry.GetAttributeValue(config.Attributes.Phone),
		Groups:   searchLDAPGroups(conn, config, username, entry.DN),
	}, nil
}

func syncLDAPUser(config systemModel.LDAPConfig, authUser *LDAPAuthUser, rawPassword string) (*systemModel.SysAdmin, error) {
	db := common.GetDB()
	admin := systemDao.GetSysAdminByUsername(authUser.Username)

	if admin.ID == 0 {
		admin = systemModel.SysAdmin{
			Username:   authUser.Username,
			Nickname:   firstNonEmpty(authUser.Nickname, authUser.Username),
			Password:   util.EncryptionMd5(rawPassword),
			Email:      authUser.Email,
			Phone:      authUser.Phone,
			Status:     1,
			CreateTime: util.HTime{Time: time.Now()},
		}
		if err := db.Create(&admin).Error; err != nil {
			return nil, err
		}
	} else {
		updates := map[string]interface{}{
			"password": util.EncryptionMd5(rawPassword),
		}
		if config.CoverAttributes || strings.TrimSpace(admin.Nickname) == "" {
			updates["nickname"] = firstNonEmpty(authUser.Nickname, admin.Nickname, authUser.Username)
		}
		if config.CoverAttributes || strings.TrimSpace(admin.Email) == "" {
			updates["email"] = firstNonEmpty(authUser.Email, admin.Email)
		}
		if config.CoverAttributes || strings.TrimSpace(admin.Phone) == "" {
			updates["phone"] = firstNonEmpty(authUser.Phone, admin.Phone)
		}
		if err := db.Model(&systemModel.SysAdmin{}).Where("id = ?", admin.ID).Updates(updates).Error; err != nil {
			return nil, err
		}
		admin = systemDao.GetSysAdminByUsername(authUser.Username)
	}

	roleIDs := collectLDAPRoleIDs(config, authUser.Groups)
	for _, roleID := range roleIDs {
		if err := assignLDAPRoleIfMissing(db, admin.ID, roleID); err != nil {
			return nil, err
		}
	}

	return &admin, nil
}

func dialLDAP(config systemModel.LDAPConfig) (*ldap.Conn, error) {
	address := net.JoinHostPort(strings.TrimSpace(config.Host), fmt.Sprintf("%d", config.Port))
	if config.TLS {
		return ldap.DialTLS("tcp", address, &tls.Config{InsecureSkipVerify: true})
	}

	conn, err := ldap.Dial("tcp", address)
	if err != nil {
		return nil, err
	}
	if config.StartTLS {
		if err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true}); err != nil {
			conn.Close()
			return nil, err
		}
	}
	return conn, nil
}

func validateLDAPConfig(config systemModel.LDAPConfig) error {
	if !config.Enable {
		return nil
	}
	if strings.TrimSpace(config.Host) == "" {
		return errors.New("LDAP host 不能为空")
	}
	if config.Port <= 0 {
		return errors.New("LDAP port 必须大于 0")
	}
	if strings.TrimSpace(config.BaseDN) == "" {
		return errors.New("LDAP baseDn 不能为空")
	}
	if strings.TrimSpace(config.AuthFilter) == "" {
		return errors.New("LDAP authFilter 不能为空")
	}
	return nil
}

func defaultLDAPConfig() systemModel.LDAPConfig {
	return systemModel.LDAPConfig{
		Enable:          false,
		Port:            389,
		AuthFilter:      "(uid=%s)",
		GroupFilter:     "(memberUid=%s)",
		GroupNameAttr:   "cn",
		CoverAttributes: true,
		TLS:             false,
		StartTLS:        false,
		Attributes: systemModel.LDAPAttributes{
			Nickname: "cn",
			Phone:    "mobile",
			Email:    "mail",
		},
	}
}

func ldapSearchAttributes(config systemModel.LDAPConfig) []string {
	attrs := []string{
		strings.TrimSpace(config.Attributes.Nickname),
		strings.TrimSpace(config.Attributes.Email),
		strings.TrimSpace(config.Attributes.Phone),
	}
	seen := make(map[string]struct{})
	result := make([]string, 0, len(attrs))
	for _, attr := range attrs {
		if attr == "" {
			continue
		}
		if _, ok := seen[attr]; ok {
			continue
		}
		seen[attr] = struct{}{}
		result = append(result, attr)
	}
	return result
}

func searchLDAPGroups(conn *ldap.Conn, config systemModel.LDAPConfig, username, userDN string) []string {
	filter := buildLDAPGroupFilter(config.GroupFilter, username, userDN)
	if filter == "" {
		return nil
	}

	groupAttr := strings.TrimSpace(config.GroupNameAttr)
	if groupAttr == "" {
		groupAttr = "cn"
	}

	req := ldap.NewSearchRequest(
		config.BaseDN,
		ldap.ScopeWholeSubtree,
		ldap.NeverDerefAliases,
		0,
		10,
		false,
		filter,
		[]string{groupAttr},
		nil,
	)
	resp, err := conn.Search(req)
	if err != nil {
		return nil
	}

	groups := make([]string, 0, len(resp.Entries))
	seen := make(map[string]struct{})
	for _, entry := range resp.Entries {
		groupName := strings.TrimSpace(entry.GetAttributeValue(groupAttr))
		if groupName == "" {
			continue
		}
		key := strings.ToLower(groupName)
		if _, ok := seen[key]; ok {
			continue
		}
		seen[key] = struct{}{}
		groups = append(groups, groupName)
	}
	return groups
}

func buildLDAPGroupFilter(template, username, userDN string) string {
	template = strings.TrimSpace(template)
	if template == "" {
		return ""
	}

	if strings.Contains(template, "{{username}}") {
		template = strings.ReplaceAll(template, "{{username}}", ldap.EscapeFilter(username))
	}
	if strings.Contains(template, "{{dn}}") {
		template = strings.ReplaceAll(template, "{{dn}}", ldap.EscapeFilter(userDN))
	}
	if strings.Contains(template, "%s") {
		template = fmt.Sprintf(template, ldap.EscapeFilter(username))
	}
	return template
}

func collectLDAPRoleIDs(config systemModel.LDAPConfig, groups []string) []uint {
	roleSet := make(map[uint]struct{})
	if config.DefaultRoleID > 0 {
		roleSet[config.DefaultRoleID] = struct{}{}
	}
	if len(groups) == 0 || len(config.RoleMappings) == 0 {
		return roleSetToSlice(roleSet)
	}

	groupSet := make(map[string]struct{}, len(groups))
	for _, group := range groups {
		groupSet[strings.ToLower(strings.TrimSpace(group))] = struct{}{}
	}

	for _, item := range config.RoleMappings {
		if item.RoleID == 0 {
			continue
		}
		if _, ok := groupSet[strings.ToLower(strings.TrimSpace(item.GroupName))]; ok {
			roleSet[item.RoleID] = struct{}{}
		}
	}
	return roleSetToSlice(roleSet)
}

func roleSetToSlice(roleSet map[uint]struct{}) []uint {
	result := make([]uint, 0, len(roleSet))
	for roleID := range roleSet {
		result = append(result, roleID)
	}
	return result
}

func assignLDAPRoleIfMissing(db *gorm.DB, adminID, roleID uint) error {
	if roleID == 0 {
		return nil
	}
	var count int64
	if err := db.Table("sys_admin_role").Where("admin_id = ? AND role_id = ?", adminID, roleID).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return db.Create(&systemModel.SysAdminRole{
		AdminId: adminID,
		RoleId:  roleID,
	}).Error
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func boolToStatus(enabled bool) int {
	if enabled {
		return 1
	}
	return 0
}
