package knowledge

import (
	"dodevops-api/api/knowledge/controller"

	"github.com/gin-gonic/gin"
)

func RegisterKnowledgeRoutes(router *gin.RouterGroup) {
	ctrl := controller.NewKnowledgeController()
	router.GET("/knowledge/articles", ctrl.List)
	router.GET("/knowledge/articles/:id", ctrl.Detail)
	router.POST("/knowledge/bootstrap", ctrl.Bootstrap)
	router.POST("/knowledge/articles/bootstrap", ctrl.Bootstrap)
	router.POST("/knowledge/articles", ctrl.Create)
	router.PUT("/knowledge/articles/:id", ctrl.Update)
	router.DELETE("/knowledge/articles/:id", ctrl.Delete)
}
