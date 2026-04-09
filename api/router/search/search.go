package search

import (
	searchController "dodevops-api/api/search/controller"

	"github.com/gin-gonic/gin"
)

func RegisterSearchRoutes(router *gin.RouterGroup) {
	ctrl := searchController.NewGlobalSearchController()
	router.GET("/search/global", ctrl.GlobalSearch)
}
