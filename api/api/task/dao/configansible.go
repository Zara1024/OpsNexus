package dao

import (
	"dodevops-api/api/task/model"

	"gorm.io/gorm"
)

type ConfigAnsibleDao struct {
	DB *gorm.DB
}

func NewConfigAnsibleDao(db *gorm.DB) *ConfigAnsibleDao {
	return &ConfigAnsibleDao{DB: db}
}

func (d *ConfigAnsibleDao) Create(config *model.ConfigAnsible) error {
	return d.DB.Create(config).Error
}

func (d *ConfigAnsibleDao) Update(config *model.ConfigAnsible) error {
	return d.DB.Save(config).Error
}

func (d *ConfigAnsibleDao) Delete(id uint) error {
	return d.DB.Delete(&model.ConfigAnsible{}, id).Error
}

func (d *ConfigAnsibleDao) GetByID(id uint) (*model.ConfigAnsible, error) {
	var config model.ConfigAnsible
	err := d.DB.First(&config, id).Error
	return &config, err
}

type ConfigAnsibleListResponse struct {
	List  []model.ConfigAnsible `json:"list"`
	Total int64                 `json:"total"`
}

func (d *ConfigAnsibleDao) List(page, size int, name string, configType int) (*ConfigAnsibleListResponse, error) {
	var configs []model.ConfigAnsible
	var total int64

	query := d.DB.Model(&model.ConfigAnsible{})
	if name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if configType > 0 {
		query = query.Where("type = ?", configType)
	}

	if err := query.Count(&total).
		Offset((page - 1) * size).
		Limit(size).
		Order("id DESC").
		Find(&configs).Error; err != nil {
		return nil, err
	}

	return &ConfigAnsibleListResponse{
		List:  configs,
		Total: total,
	}, nil
}
