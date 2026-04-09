package dao

import (
	"strings"

	"dodevops-api/api/k8s/model"
	"gorm.io/gorm"
)

type KubeClusterDao struct {
	DB *gorm.DB
}

func NewKubeClusterDao(db *gorm.DB) *KubeClusterDao {
	return &KubeClusterDao{DB: db}
}

func (d *KubeClusterDao) Create(cluster *model.KubeCluster) error {
	return d.DB.Create(cluster).Error
}

func (d *KubeClusterDao) GetByID(id uint) (*model.KubeCluster, error) {
	var cluster model.KubeCluster
	err := d.DB.Where("id = ?", id).First(&cluster).Error
	return &cluster, err
}

func (d *KubeClusterDao) GetByName(name string) (*model.KubeCluster, error) {
	var cluster model.KubeCluster
	err := d.DB.Where("name = ?", name).First(&cluster).Error
	return &cluster, err
}

func (d *KubeClusterDao) List(page, size int, name, version string, status *int) ([]model.KubeCluster, int64, error) {
	var clusters []model.KubeCluster
	var total int64

	query := d.DB.Model(&model.KubeCluster{})
	if strings.TrimSpace(name) != "" {
		query = query.Where("name LIKE ?", "%"+strings.TrimSpace(name)+"%")
	}
	if strings.TrimSpace(version) != "" {
		query = query.Where("version LIKE ?", "%"+strings.TrimSpace(version)+"%")
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * size
	err := query.Offset(offset).
		Limit(size).
		Order("created_at DESC").
		Find(&clusters).Error

	return clusters, total, err
}

func (d *KubeClusterDao) Update(id uint, updates map[string]interface{}) error {
	return d.DB.Model(&model.KubeCluster{}).Where("id = ?", id).Updates(updates).Error
}

func (d *KubeClusterDao) Delete(id uint) error {
	return d.DB.Where("id = ?", id).Delete(&model.KubeCluster{}).Error
}

func (d *KubeClusterDao) UpdateStatus(id uint, status int) error {
	return d.DB.Model(&model.KubeCluster{}).Where("id = ?", id).Update("status", status).Error
}

func (d *KubeClusterDao) UpdateCredential(id uint, credential string) error {
	return d.DB.Model(&model.KubeCluster{}).Where("id = ?", id).Update("credential", credential).Error
}

func (d *KubeClusterDao) IsClusterNameExists(name string) (bool, error) {
	var count int64
	err := d.DB.Model(&model.KubeCluster{}).Where("name = ?", name).Count(&count).Error
	return count > 0, err
}

func (d *KubeClusterDao) GetClusterCountByStatus(status int) (int64, error) {
	var count int64
	err := d.DB.Model(&model.KubeCluster{}).Where("status = ?", status).Count(&count).Error
	return count, err
}
