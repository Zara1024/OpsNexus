package controller

import (
	"strconv"
	"strings"

	aiModel "dodevops-api/api/ai/model"
	aiService "dodevops-api/api/ai/service"
	"dodevops-api/common/result"

	"github.com/gin-gonic/gin"
)

type AIController struct {
	service *aiService.AIService
}

func NewAIController() *AIController {
	return &AIController{
		service: aiService.NewAIService(),
	}
}

func (c *AIController) ListTemplates(ctx *gin.Context) {
	c.service.ListTemplates(ctx)
}

func (c *AIController) ListOverview(ctx *gin.Context) {
	c.service.ListOverview(ctx)
}

func (c *AIController) SuggestKnowledge(ctx *gin.Context) {
	c.service.SuggestKnowledge(ctx, strings.TrimSpace(ctx.Query("keyword")))
}

func (c *AIController) RenderPrompt(ctx *gin.Context) {
	var req aiModel.AIPromptRenderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.RenderPrompt(ctx, req)
}

func (c *AIController) ListHistory(ctx *gin.Context) {
	c.service.ListHistory(ctx)
}

func (c *AIController) GetHistoryDetail(ctx *gin.Context) {
	c.service.GetHistoryDetail(ctx, strings.TrimSpace(ctx.Param("sessionId")))
}

func (c *AIController) ChatAssistant(ctx *gin.Context) {
	var req aiModel.AIAssistantChatRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.ChatAssistant(ctx, req)
}

func (c *AIController) ListAssistantHistory(ctx *gin.Context) {
	c.service.ListAssistantHistory(ctx)
}

func (c *AIController) GetAssistantHistoryDetail(ctx *gin.Context) {
	c.service.GetAssistantHistoryDetail(ctx, strings.TrimSpace(ctx.Param("sessionId")))
}

func (c *AIController) ListInspectionTemplates(ctx *gin.Context) {
	c.service.ListInspectionTemplates(ctx)
}

func (c *AIController) ListInspectionReports(ctx *gin.Context) {
	c.service.ListInspectionReports(ctx)
}

func (c *AIController) DecideConfirmation(ctx *gin.Context) {
	var req aiModel.AIAssistantConfirmationDecisionRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	id := strings.TrimSpace(ctx.Param("id"))
	if id == "" {
		result.Failed(ctx, 400, "缺少确认ID")
		return
	}
	parsedID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		result.Failed(ctx, 400, "无效的确认ID")
		return
	}
	c.service.DecideConfirmation(ctx, uint(parsedID), req.Decision)
}

func (c *AIController) Diagnose(ctx *gin.Context) {
	var req aiModel.AIDiagnosisRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		result.Failed(ctx, 400, "参数错误: "+err.Error())
		return
	}
	c.service.Diagnose(ctx, req)
}
