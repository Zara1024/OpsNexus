package dao

import (
	"dodevops-api/api/cmdb/model"

	"gorm.io/gorm"
)

type CmdbDeviceDao struct {
	db *gorm.DB
}

func NewCmdbDeviceDao(db *gorm.DB) *CmdbDeviceDao {
	return &CmdbDeviceDao{db: db}
}

func (d *CmdbDeviceDao) Create(device *model.CmdbDevice) error {
	return d.db.Create(device).Error
}

func (d *CmdbDeviceDao) Update(device *model.CmdbDevice) error {
	return d.db.Save(device).Error
}

func (d *CmdbDeviceDao) Delete(id uint) error {
	return d.db.Delete(&model.CmdbDevice{}, id).Error
}

func (d *CmdbDeviceDao) GetByID(id uint) (*model.CmdbDevice, error) {
	var device model.CmdbDevice
	err := d.db.First(&device, id).Error
	return &device, err
}

func (d *CmdbDeviceDao) List(page, pageSize int) ([]model.CmdbDevice, int64, error) {
	var devices []model.CmdbDevice
	var count int64

	if err := d.db.Model(&model.CmdbDevice{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	err := d.db.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&devices).Error
	return devices, count, err
}

func (d *CmdbDeviceDao) GetByIDs(ids []uint) ([]model.CmdbDevice, error) {
	var devices []model.CmdbDevice
	err := d.db.Where("id IN ?", ids).Find(&devices).Error
	return devices, err
}
