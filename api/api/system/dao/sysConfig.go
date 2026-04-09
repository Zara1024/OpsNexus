package dao

import (
	"errors"
	"time"

	"dodevops-api/api/system/model"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"

	"gorm.io/gorm"
)

func GetSysConfigByKey(configKey string) (*model.SysConfig, error) {
	var item model.SysConfig
	err := Db.Where("config_key = ?", configKey).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &item, nil
}

func UpsertSysConfig(item *model.SysConfig) error {
	existing, err := GetSysConfigByKey(item.ConfigKey)
	if err != nil {
		return err
	}

	now := util.HTime{Time: time.Now()}
	if existing == nil {
		item.CreateTime = now
		item.UpdateTime = now
		return Db.Create(item).Error
	}

	item.ID = existing.ID
	item.CreateTime = existing.CreateTime
	item.UpdateTime = now
	return Db.Model(&model.SysConfig{}).
		Where("id = ?", existing.ID).
		Updates(map[string]interface{}{
			"config_type": item.ConfigType,
			"config_data": item.ConfigData,
			"status":      item.Status,
			"remark":      item.Remark,
			"update_time": item.UpdateTime,
		}).Error
}
