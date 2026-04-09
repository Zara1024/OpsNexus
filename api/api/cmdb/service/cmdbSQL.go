package service

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"errors"
	"fmt"
	"strings"
)

type CmdbSQLService struct {
	dao *dao.CmdbSQLDao
}

var errInvalidCmdbSQLParams = errors.New("invalid cmdb sql params")

const (
	cmdbSQLDefaultPage     = 1
	cmdbSQLDefaultPageSize = 10
	cmdbSQLMaxPageSize     = 100
)

func NewCmdbSQLService(dao *dao.CmdbSQLDao) *CmdbSQLService {
	return &CmdbSQLService{dao: dao}
}

func IsInvalidCmdbSQLParams(err error) bool {
	return errors.Is(err, errInvalidCmdbSQLParams)
}

// CreateDatabase 创建数据库记录
func (s *CmdbSQLService) CreateDatabase(dto model.CreateCmdbSQLDto) (*model.CmdbSQL, error) {
	db := normalizeCreateCmdbSQLRequest(dto)
	if err := validateCmdbSQL(db); err != nil {
		return nil, err
	}
	if err := s.dao.Create(&db); err != nil {
		return nil, err
	}
	decorateCmdbSQLAsset(&db)
	return &db, nil
}

// UpdateDatabase 更新数据库记录
func (s *CmdbSQLService) UpdateDatabase(dto model.UpdateCmdbSQLDto) (*model.CmdbSQL, error) {
	existing, err := s.dao.GetByID(dto.ID)
	if err != nil {
		return nil, err
	}

	db := normalizeUpdateCmdbSQLRequest(*existing, dto)
	if err = validateCmdbSQL(db); err != nil {
		return nil, err
	}

	db.CreatedAt = existing.CreatedAt
	if err = s.dao.Update(&db); err != nil {
		return nil, err
	}
	decorateCmdbSQLAsset(&db)
	return &db, nil
}

// DeleteDatabase 删除数据库记录
func (s *CmdbSQLService) DeleteDatabase(id uint) error {
	return s.dao.Delete(id)
}

// GetDatabase 获取单个数据库详情
func (s *CmdbSQLService) GetDatabase(id uint) (*model.CmdbSQL, error) {
	db, err := s.dao.GetByID(id)
	if err != nil {
		return nil, err
	}
	decorateCmdbSQLAsset(db)
	return db, nil
}

// ListDatabases 分页查询数据库列表
func (s *CmdbSQLService) ListDatabases(page, pageSize int) ([]model.CmdbSQL, int64, error) {
	page = normalizeCmdbSQLPage(page)
	pageSize = normalizeCmdbSQLPageSize(pageSize)
	list, count, err := s.dao.List(page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	decorateCmdbSQLAssets(list)
	return list, count, nil
}

// GetDatabasesByAccount 根据账号查询数据库
func (s *CmdbSQLService) GetDatabasesByAccount(accountID uint) ([]model.CmdbSQL, error) {
	list, err := s.dao.GetByAccountID(accountID)
	if err != nil {
		return nil, err
	}
	decorateCmdbSQLAssets(list)
	return list, nil
}

// GetDatabasesByGroup 根据业务组查询数据库
func (s *CmdbSQLService) GetDatabasesByGroup(groupID uint) ([]model.CmdbSQL, error) {
	list, err := s.dao.GetByGroupID(groupID)
	if err != nil {
		return nil, err
	}
	decorateCmdbSQLAssets(list)
	return list, nil
}

// GetDatabasesByName 根据名称查询数据库
func (s *CmdbSQLService) GetDatabasesByName(name string) ([]model.CmdbSQL, error) {
	list, err := s.dao.GetByName(name)
	if err != nil {
		return nil, err
	}
	decorateCmdbSQLAssets(list)
	return list, nil
}

func (s *CmdbSQLService) ResolveDatabase(databaseID uint, databaseName string) (*model.CmdbSQL, error) {
	db, _, err := resolveCmdbSQLTarget(databaseID, databaseName, s.dao.GetByID, s.dao.GetByName)
	return db, err
}

func (s *CmdbSQLService) ResolveDatabaseTarget(databaseID uint, databaseName string) (*model.CmdbSQL, string, error) {
	return resolveCmdbSQLTarget(databaseID, databaseName, s.dao.GetByID, s.dao.GetByName)
}

// GetDatabasesByType 根据类型查询数据库
func (s *CmdbSQLService) GetDatabasesByType(dbType int) ([]model.CmdbSQL, error) {
	if dbType < 1 || dbType > 5 {
		return nil, errors.New("无效的数据库类型")
	}
	list, err := s.dao.GetByType(dbType)
	if err != nil {
		return nil, err
	}
	decorateCmdbSQLAssets(list)
	return list, nil
}

func ResolveCmdbSQLSchemaName(db model.CmdbSQL, requested string) string {
	return resolveCmdbSQLSchemaName(db, requested)
}

func resolveCmdbSQLSchemaName(db model.CmdbSQL, requested string) string {
	return firstNonEmptyCmdbSQL(
		strings.TrimSpace(requested),
		strings.TrimSpace(db.DefaultDatabase),
		strings.TrimSpace(db.Name),
	)
}

func resolveCmdbSQLAssetLookup(
	databaseID uint,
	databaseName string,
	getByID func(uint) (*model.CmdbSQL, error),
	getByName func(string) ([]model.CmdbSQL, error),
) (*model.CmdbSQL, error) {
	switch {
	case databaseID > 0:
		db, err := getByID(databaseID)
		if err != nil {
			return nil, err
		}
		decorateCmdbSQLAsset(db)
		if !db.IsActive {
			return nil, fmt.Errorf("database asset %q is inactive", strings.TrimSpace(db.Name))
		}
		return db, nil
	case strings.TrimSpace(databaseName) != "":
		name := strings.TrimSpace(databaseName)
		list, err := getByName(name)
		if err != nil {
			return nil, err
		}
		switch len(list) {
		case 0:
			return nil, fmt.Errorf("database asset %q not found", name)
		case 1:
			db := list[0]
			decorateCmdbSQLAsset(&db)
			if !db.IsActive {
				return nil, fmt.Errorf("database asset %q is inactive", strings.TrimSpace(db.Name))
			}
			return &db, nil
		default:
			return nil, fmt.Errorf("database asset %q is ambiguous", name)
		}
	default:
		return nil, errors.New("database id or name is required")
	}
}

func resolveCmdbSQLTarget(
	databaseID uint,
	databaseName string,
	getByID func(uint) (*model.CmdbSQL, error),
	getByName func(string) ([]model.CmdbSQL, error),
) (*model.CmdbSQL, string, error) {
	requested := strings.TrimSpace(databaseName)
	// Request contract:
	// - when databaseId is present, databaseName is treated as a schema override
	// - when databaseId is absent, databaseName is treated as the asset name and schema comes from asset defaults
	switch {
	case databaseID > 0:
		db, err := resolveCmdbSQLAssetLookup(databaseID, "", getByID, getByName)
		if err != nil {
			return nil, "", err
		}
		schemaName := resolveCmdbSQLSchemaName(*db, requested)
		if schemaName == "" {
			return nil, "", errors.New("database schema name is required")
		}
		return db, schemaName, nil
	case requested != "":
		db, err := resolveCmdbSQLAssetLookup(0, requested, getByID, getByName)
		if err != nil {
			return nil, "", err
		}
		schemaName := resolveCmdbSQLSchemaName(*db, "")
		if schemaName == "" {
			return nil, "", errors.New("database schema name is required")
		}
		return db, schemaName, nil
	default:
		return nil, "", errors.New("database id or name is required")
	}
}

func normalizeCreateCmdbSQLRequest(dto model.CreateCmdbSQLDto) model.CmdbSQL {
	remark := firstNonEmptyCmdbSQL(strings.TrimSpace(dto.Remark), strings.TrimSpace(dto.Description))
	name := strings.TrimSpace(dto.Name)

	db := model.CmdbSQL{
		Name:            name,
		Address:         strings.TrimSpace(dto.Address),
		Platform:        strings.TrimSpace(dto.Platform),
		DefaultDatabase: firstNonEmptyCmdbSQL(strings.TrimSpace(dto.DefaultDatabase), name),
		Type:            dto.Type,
		AccountID:       dto.AccountID,
		GroupID:         dto.GroupID,
		ProtocolGroup:   normalizeCmdbSQLProtocolGroup(dto.ProtocolGroup, dto.Type),
		Tags:            strings.TrimSpace(dto.Tags),
		IsActive:        normalizeCmdbSQLActive(dto.IsActive),
		Remark:          remark,
		Description:     remark,
	}
	return db
}

func normalizeUpdateCmdbSQLRequest(existing model.CmdbSQL, dto model.UpdateCmdbSQLDto) model.CmdbSQL {
	db := normalizeCreateCmdbSQLRequest(model.CreateCmdbSQLDto{
		Name:            dto.Name,
		Address:         dto.Address,
		Platform:        dto.Platform,
		DefaultDatabase: dto.DefaultDatabase,
		Type:            dto.Type,
		AccountID:       dto.AccountID,
		GroupID:         dto.GroupID,
		ProtocolGroup:   dto.ProtocolGroup,
		Tags:            dto.Tags,
		IsActive:        dto.IsActive,
		Remark:          dto.Remark,
		Description:     dto.Description,
	})

	db.ID = existing.ID
	db.CreatedAt = existing.CreatedAt
	if strings.TrimSpace(dto.DefaultDatabase) == "" {
		db.DefaultDatabase = firstNonEmptyCmdbSQL(strings.TrimSpace(existing.DefaultDatabase), db.Name)
	}
	if dto.IsActive == nil {
		db.IsActive = existing.IsActive
	}
	if strings.TrimSpace(dto.ProtocolGroup) == "" {
		db.ProtocolGroup = normalizeCmdbSQLProtocolGroup(existing.ProtocolGroup, db.Type)
	}
	if strings.TrimSpace(dto.Remark) == "" && strings.TrimSpace(dto.Description) == "" {
		db.Remark = firstNonEmptyCmdbSQL(strings.TrimSpace(existing.Remark), strings.TrimSpace(existing.Description))
	}
	db.Description = db.Remark
	return db
}

func validateCmdbSQL(db model.CmdbSQL) error {
	if strings.TrimSpace(db.Name) == "" {
		return invalidCmdbSQLParams("name cannot be empty")
	}
	if strings.TrimSpace(db.Address) == "" {
		return invalidCmdbSQLParams("address cannot be empty")
	}
	if db.Type < 1 || db.Type > 5 {
		return invalidCmdbSQLParams("invalid database type")
	}
	if db.AccountID == 0 {
		return invalidCmdbSQLParams("accountId cannot be empty")
	}
	if db.GroupID == 0 {
		return invalidCmdbSQLParams("groupId cannot be empty")
	}
	if resolveCmdbSQLSchemaName(db, "") == "" {
		return invalidCmdbSQLParams("defaultDatabase cannot be empty")
	}
	return nil
}

func normalizeCmdbSQLProtocolGroup(value string, dbType int) string {
	parts := strings.Split(value, ",")
	normalized := make([]string, 0, len(parts))
	seen := make(map[string]struct{}, len(parts))
	for _, part := range parts {
		item := strings.ToLower(strings.TrimSpace(part))
		if item == "" {
			continue
		}
		if _, ok := seen[item]; ok {
			continue
		}
		seen[item] = struct{}{}
		normalized = append(normalized, item)
	}
	if len(normalized) > 0 {
		return strings.Join(normalized, ",")
	}
	return defaultCmdbSQLProtocolGroup(dbType)
}

func defaultCmdbSQLProtocolGroup(dbType int) string {
	switch dbType {
	case 1:
		return "mysql"
	case 2:
		return "postgresql"
	case 3:
		return "redis"
	case 4:
		return "mongodb"
	case 5:
		return "elasticsearch"
	default:
		return ""
	}
}

func normalizeCmdbSQLActive(value *bool) bool {
	if value == nil {
		return true
	}
	return *value
}

func decorateCmdbSQLAsset(db *model.CmdbSQL) {
	if db == nil {
		return
	}
	db.Description = strings.TrimSpace(db.Remark)
}

func decorateCmdbSQLAssets(list []model.CmdbSQL) {
	for i := range list {
		decorateCmdbSQLAsset(&list[i])
	}
}

func invalidCmdbSQLParams(message string) error {
	return fmt.Errorf("%w: %s", errInvalidCmdbSQLParams, message)
}

func firstNonEmptyCmdbSQL(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}

func normalizeCmdbSQLPage(page int) int {
	if page < 1 {
		return cmdbSQLDefaultPage
	}
	return page
}

func normalizeCmdbSQLPageSize(pageSize int) int {
	if pageSize < 1 {
		return cmdbSQLDefaultPageSize
	}
	if pageSize > cmdbSQLMaxPageSize {
		return cmdbSQLMaxPageSize
	}
	return pageSize
}
