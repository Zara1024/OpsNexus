package service

import (
	"strings"
	"time"

	"dodevops-api/api/knowledge/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	. "dodevops-api/pkg/db"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type KnowledgeService struct{}

func NewKnowledgeService() *KnowledgeService {
	return &KnowledgeService{}
}

func (s *KnowledgeService) List(c *gin.Context, query model.KnowledgeQuery) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}

	db := Db.Model(&model.KnowledgeBase{})
	if strings.TrimSpace(query.Keyword) != "" {
		like := "%" + strings.TrimSpace(query.Keyword) + "%"
		db = db.Where("title LIKE ? OR content LIKE ? OR keywords LIKE ? OR tags LIKE ?", like, like, like, like)
	}
	if strings.TrimSpace(query.Type) != "" {
		db = db.Where("type = ?", strings.TrimSpace(query.Type))
	}
	if strings.TrimSpace(query.Category) != "" {
		db = db.Where("category = ?", strings.TrimSpace(query.Category))
	}
	if query.Enabled >= 0 {
		db = db.Where("enabled = ?", query.Enabled)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		result.Failed(c, 500, "获取知识库列表失败: "+err.Error())
		return
	}

	var list []model.KnowledgeBase
	if err := db.Order("update_time DESC, id DESC").
		Offset((query.Page - 1) * query.PageSize).
		Limit(query.PageSize).
		Find(&list).Error; err != nil {
		result.Failed(c, 500, "获取知识库列表失败: "+err.Error())
		return
	}

	result.Success(c, gin.H{
		"list":     list,
		"total":    total,
		"page":     query.Page,
		"pageSize": query.PageSize,
	})
}

func (s *KnowledgeService) Detail(c *gin.Context, id uint) {
	var item model.KnowledgeBase
	if err := Db.First(&item, id).Error; err != nil {
		result.Failed(c, 404, "知识库文章不存在")
		return
	}
	Db.Model(&model.KnowledgeBase{}).Where("id = ?", id).UpdateColumn("use_count", gorm.Expr("use_count + ?", 1))
	item.UseCount++
	result.Success(c, item)
}

func (s *KnowledgeService) Create(c *gin.Context, item *model.KnowledgeBase) {
	now := util.HTime{Time: time.Now()}
	item.Type = strings.TrimSpace(item.Type)
	item.Category = strings.TrimSpace(item.Category)
	item.Title = strings.TrimSpace(item.Title)
	item.Content = strings.TrimSpace(item.Content)
	item.Keywords = strings.TrimSpace(item.Keywords)
	item.Tags = strings.TrimSpace(item.Tags)
	if item.Score <= 0 {
		item.Score = 0.5
	}
	if item.Enabled != 0 && item.Enabled != 1 {
		item.Enabled = 1
	}
	item.CreateTime = now
	item.UpdateTime = now

	if item.Title == "" || item.Content == "" || item.Type == "" {
		result.Failed(c, 400, "标题、类型和内容不能为空")
		return
	}
	if err := Db.Create(item).Error; err != nil {
		result.Failed(c, 500, "创建知识库文章失败: "+err.Error())
		return
	}
	result.Success(c, item)
}

func (s *KnowledgeService) Update(c *gin.Context, id uint, item *model.KnowledgeBase) {
	updates := map[string]interface{}{
		"type":        strings.TrimSpace(item.Type),
		"category":    strings.TrimSpace(item.Category),
		"title":       strings.TrimSpace(item.Title),
		"content":     strings.TrimSpace(item.Content),
		"keywords":    strings.TrimSpace(item.Keywords),
		"tags":        strings.TrimSpace(item.Tags),
		"score":       item.Score,
		"enabled":     item.Enabled,
		"update_time": util.HTime{Time: time.Now()},
	}
	if updates["type"] == "" || updates["title"] == "" || updates["content"] == "" {
		result.Failed(c, 400, "标题、类型和内容不能为空")
		return
	}
	if item.Score <= 0 {
		updates["score"] = 0.5
	}
	if item.Enabled != 0 && item.Enabled != 1 {
		updates["enabled"] = 1
	}
	if err := Db.Model(&model.KnowledgeBase{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		result.Failed(c, 500, "更新知识库文章失败: "+err.Error())
		return
	}
	s.Detail(c, id)
}

func (s *KnowledgeService) Delete(c *gin.Context, id uint) {
	if err := Db.Delete(&model.KnowledgeBase{}, id).Error; err != nil {
		result.Failed(c, 500, "删除知识库文章失败: "+err.Error())
		return
	}
	result.Success(c, gin.H{"id": id})
}

func (s *KnowledgeService) Bootstrap(c *gin.Context) {
	seeds := []model.KnowledgeBase{
		buildKnowledgeSeed(
			"runbook",
			"K8s",
			"K8s 终端审计回放 SOP",
			"terminal audit,pod,kubectl,runbook",
			"k8s,terminal,audit",
			0.9,
			"# K8s 终端审计回放 SOP\n\n## 适用场景\n- Pod Web terminal 异常回放\n- kubectl terminal 误操作复盘\n- 需要从 `sys_session_recording` 与 `sys_command_audit` 交叉核对时\n\n## 操作步骤\n1. 在终端录像审计页按会话 ID、命令关键字或风险等级筛选。\n2. 优先关注 `录像状态=可回放` 的会话。\n3. 先看命令审计，再切到录像回放核对输入/输出顺序。\n4. 对 `command_only` 历史会话，只能走命令聚合复盘。\n\n## 排障提示\n- `文件缺失`：通常是录制文件未落盘或已被清理。\n- `空文件`：通常是会话异常中断或终端刚连接即退出。\n- `窗口变化`：常用于判断用户是否在复制粘贴大段内容。\n",
		),
		buildKnowledgeSeed(
			"change",
			"数据库",
			"SQL 变更工单回滚检查清单",
			"sql,rollback,change,workorder",
			"sql,rollback,change",
			0.95,
			"# SQL 变更工单回滚检查清单\n\n## 提交前\n- 明确变更目标库、schema、影响表\n- 确认是否带 `WHERE` / `LIMIT`\n- 评估是否需要先备份旧数据快照\n\n## 审批时重点\n- `DELETE` / `TRUNCATE` / `DROP` 默认按高风险处理\n- `UPDATE` 必须确认影响范围与回滚方案\n- 对 DDL 变更，优先要求结构备份或变更脚本回退方案\n\n## 执行前\n1. 先检查工单里的 `回滚建议`\n2. 对 UPDATE / DELETE 尽量保留执行前快照\n3. 执行后核对影响行数和结果日志\n",
		),
		buildKnowledgeSeed(
			"incident",
			"告警",
			"AlertManager 告警分诊清单",
			"alertmanager,alert,incident,triage",
			"alertmanager,alert,triage",
			0.85,
			"# AlertManager 告警分诊清单\n\n## 一级判断\n- 是否为重复告警\n- 是否存在静默策略\n- 是否影响生产链路\n\n## 二级排查\n- 先看告警中心事件详情\n- 再查关联主机 / K8s 节点 / 应用发布记录\n- 若涉及终端操作，联动终端审计回放\n\n## 处置建议\n- 需要人工操作的故障，先沉淀操作过程到知识库或工单\n- 恢复后补充根因、处置时间线和后续优化项\n",
		),
		buildKnowledgeSeed(
			"faq",
			"AI",
			"AI 诊断上下文整理模板",
			"ai,diagnosis,context,knowledge",
			"ai,diagnosis,prompt",
			0.8,
			"# AI 诊断上下文整理模板\n\n## 最小上下文\n- 发生时间\n- 影响对象\n- 当前现象\n- 最近变更\n- 已执行排查动作\n\n## 推荐补充\n- 告警详情\n- 终端审计命令轨迹\n- SQL 工单内容与执行结果\n- K8s workload / pod / event 信息\n\n## 目标\n让 AI 输出更聚焦：\n1. 故障摘要\n2. 可能根因\n3. 优先级排序后的排查建议\n4. 可沉淀成知识条目的结论\n",
		),
		buildKnowledgeSeed(
			"inspection",
			"巡检",
			"周度巡检报告复盘模板",
			"inspection,report,巡检,复盘",
			"inspection,runbook,ai",
			0.88,
			"# 周度巡检报告复盘模板\n\n## 巡检范围\n- 主机资源与可用性\n- 数据库健康与连接情况\n- SSL 证书剩余有效期\n- 近期告警与发布变更\n\n## 输出要求\n1. 总结本周新增风险项\n2. 标记需要立刻跟进的阻断问题\n3. 识别可沉淀为 SOP / 知识的重复问题\n4. 给出下周治理建议与责任人\n",
		),
	}

	inserted := make([]model.KnowledgeBase, 0)
	for _, seed := range seeds {
		var count int64
		if err := Db.Model(&model.KnowledgeBase{}).Where("title = ?", seed.Title).Count(&count).Error; err != nil {
			result.Failed(c, 500, "初始化知识库失败: "+err.Error())
			return
		}
		if count > 0 {
			continue
		}
		item := seed
		if err := Db.Create(&item).Error; err != nil {
			result.Failed(c, 500, "初始化知识库失败: "+err.Error())
			return
		}
		inserted = append(inserted, item)
	}

	result.Success(c, gin.H{
		"inserted": len(inserted),
		"list":     inserted,
	})
}

func buildKnowledgeSeed(articleType, category, title, keywords, tags string, score float64, content string) model.KnowledgeBase {
	now := util.HTime{Time: time.Now()}
	return model.KnowledgeBase{
		Type:       articleType,
		Category:   category,
		Title:      title,
		Content:    content,
		Keywords:   keywords,
		Tags:       tags,
		Score:      score,
		Enabled:    1,
		CreateTime: now,
		UpdateTime: now,
	}
}
