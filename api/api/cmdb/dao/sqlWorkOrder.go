package dao

import (
	"strings"

	"dodevops-api/api/cmdb/model"
	"gorm.io/gorm"
)

type SQLWorkOrderDao struct {
	db *gorm.DB
}

func NewSQLWorkOrderDao(db *gorm.DB) *SQLWorkOrderDao {
	return &SQLWorkOrderDao{db: db}
}

func (d *SQLWorkOrderDao) Create(order *model.CmdbSQLWorkOrder) error {
	return d.db.Create(order).Error
}

func (d *SQLWorkOrderDao) Update(order *model.CmdbSQLWorkOrder) error {
	return d.db.Save(order).Error
}

func (d *SQLWorkOrderDao) GetByID(id uint) (*model.CmdbSQLWorkOrder, error) {
	var order model.CmdbSQLWorkOrder
	err := d.db.First(&order, id).Error
	return &order, err
}

func (d *SQLWorkOrderDao) GetSummary() (model.CmdbSQLWorkOrderSummary, error) {
	summary := model.CmdbSQLWorkOrderSummary{}
	base := d.db.Model(&model.CmdbSQLWorkOrder{})

	if err := base.Count(&summary.Total).Error; err != nil {
		return summary, err
	}
	if err := d.db.Model(&model.CmdbSQLWorkOrder{}).Where("status = ?", 1).Count(&summary.Pending).Error; err != nil {
		return summary, err
	}
	if err := d.db.Model(&model.CmdbSQLWorkOrder{}).Where("status = ?", 2).Count(&summary.Approved).Error; err != nil {
		return summary, err
	}
	if err := d.db.Model(&model.CmdbSQLWorkOrder{}).Where("status = ?", 3).Count(&summary.Rejected).Error; err != nil {
		return summary, err
	}
	if err := d.db.Model(&model.CmdbSQLWorkOrder{}).Where("status = ?", 4).Count(&summary.Executing).Error; err != nil {
		return summary, err
	}
	if err := d.db.Model(&model.CmdbSQLWorkOrder{}).Where("status = ?", 5).Count(&summary.Succeeded).Error; err != nil {
		return summary, err
	}
	if err := d.db.Model(&model.CmdbSQLWorkOrder{}).Where("status = ?", 6).Count(&summary.Failed).Error; err != nil {
		return summary, err
	}
	if err := d.db.Model(&model.CmdbSQLWorkOrder{}).Where("risk_level >= ?", 2).Count(&summary.HighRisk).Error; err != nil {
		return summary, err
	}
	return summary, nil
}

func (d *SQLWorkOrderDao) List(query model.CmdbSQLWorkOrderQuery) ([]model.CmdbSQLWorkOrder, int64, error) {
	var (
		list  []model.CmdbSQLWorkOrder
		total int64
	)

	dbQuery := d.db.Model(&model.CmdbSQLWorkOrder{})
	if query.Status > 0 {
		dbQuery = dbQuery.Where("status = ?", query.Status)
	}
	if query.DatabaseID > 0 {
		dbQuery = dbQuery.Where("database_id = ?", query.DatabaseID)
	}
	if keyword := strings.TrimSpace(query.Keyword); keyword != "" {
		like := "%" + keyword + "%"
		dbQuery = dbQuery.Where(
			"order_no LIKE ? OR title LIKE ? OR database_name LIKE ? OR applicant_name LIKE ? OR approver_name LIKE ? OR sql_content LIKE ?",
			like, like, like, like, like, like,
		)
	}

	if err := dbQuery.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	err := dbQuery.Order("id DESC").
		Offset((query.Page - 1) * query.PageSize).
		Limit(query.PageSize).
		Find(&list).Error
	return list, total, err
}
