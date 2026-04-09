package dao

import (
	"strings"

	"dodevops-api/api/k8s/model"
	. "dodevops-api/pkg/db"
)

// CreateWorkloadCapacitySuggestionHistory stores one workload capacity suggestion snapshot.
func CreateWorkloadCapacitySuggestionHistory(item *model.WorkloadCapacitySuggestionHistoryEntity) error {
	return Db.Create(item).Error
}

// ListWorkloadCapacitySuggestionHistory returns workload capacity suggestion snapshots by target.
func ListWorkloadCapacitySuggestionHistory(
	clusterID uint,
	namespaceName string,
	workloadType string,
	workloadName string,
	pageNum int,
	pageSize int,
) (list []model.WorkloadCapacitySuggestionHistoryEntity, count int64, err error) {
	db := Db.Model(&model.WorkloadCapacitySuggestionHistoryEntity{}).
		Where(
			"cluster_id = ? AND namespace_name = ? AND workload_type = ? AND workload_name = ?",
			clusterID,
			strings.TrimSpace(namespaceName),
			strings.ToLower(strings.TrimSpace(workloadType)),
			strings.TrimSpace(workloadName),
		)

	if err = db.Count(&count).Error; err != nil {
		return list, count, err
	}

	err = db.Order("created_at DESC, id DESC").
		Limit(pageSize).
		Offset((pageNum - 1) * pageSize).
		Find(&list).Error
	return list, count, err
}
