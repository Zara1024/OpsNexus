package model

import "dodevops-api/common/util"

type KnowledgeBase struct {
	ID         uint       `gorm:"column:id;primaryKey" json:"id"`
	Type       string     `gorm:"column:type" json:"type"`
	Category   string     `gorm:"column:category" json:"category"`
	Title      string     `gorm:"column:title" json:"title"`
	Content    string     `gorm:"column:content" json:"content"`
	Keywords   string     `gorm:"column:keywords" json:"keywords"`
	Tags       string     `gorm:"column:tags" json:"tags"`
	Score      float64    `gorm:"column:score" json:"score"`
	UseCount   int64      `gorm:"column:use_count" json:"useCount"`
	Enabled    int        `gorm:"column:enabled" json:"enabled"`
	CreateTime util.HTime `gorm:"column:create_time" json:"createTime"`
	UpdateTime util.HTime `gorm:"column:update_time" json:"updateTime"`
}

func (KnowledgeBase) TableName() string {
	return "knowledge_base"
}

type KnowledgeQuery struct {
	Page     int
	PageSize int
	Keyword  string
	Type     string
	Category string
	Enabled  int
}
