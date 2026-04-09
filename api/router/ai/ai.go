package ai

import (
	"dodevops-api/api/ai/controller"

	"github.com/gin-gonic/gin"
)

func RegisterAIRoutes(router *gin.RouterGroup) {
	ctrl := controller.NewAIController()
	router.GET("/ai/overview", ctrl.ListOverview)
	router.GET("/ai/templates", ctrl.ListTemplates)
	router.GET("/ai/knowledge/suggest", ctrl.SuggestKnowledge)
	router.POST("/ai/render", ctrl.RenderPrompt)
	router.POST("/ai/diagnose", ctrl.Diagnose)
	router.GET("/ai/history", ctrl.ListHistory)
	router.GET("/ai/history/:sessionId", ctrl.GetHistoryDetail)
	router.POST("/ai/assistant/chat", ctrl.ChatAssistant)
	router.GET("/ai/assistant/history", ctrl.ListAssistantHistory)
	router.GET("/ai/assistant/history/:sessionId", ctrl.GetAssistantHistoryDetail)
	router.GET("/ai/assistant/templates", ctrl.ListInspectionTemplates)
	router.GET("/ai/assistant/reports", ctrl.ListInspectionReports)
	router.POST("/ai/assistant/confirm/:id", ctrl.DecideConfirmation)
}
