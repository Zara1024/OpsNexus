package controller

import (
	"strconv"

	searchModel "dodevops-api/api/search/model"
	searchService "dodevops-api/api/search/service"

	"github.com/gin-gonic/gin"
)

type GlobalSearchController struct {
	service searchService.GlobalSearchService
}

func NewGlobalSearchController() *GlobalSearchController {
	return &GlobalSearchController{
		service: searchService.NewGlobalSearchService(),
	}
}

func (c *GlobalSearchController) GlobalSearch(ctx *gin.Context) {
	query := searchModel.GlobalSearchQuery{
		Keyword: ctx.Query("keyword"),
		Types:   ctx.Query("types"),
		Limit:   parseInt(ctx.Query("limit"), 8),
	}
	c.service.GlobalSearch(ctx, query)
}

func parseInt(value string, fallback int) int {
	if value == "" {
		return fallback
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return fallback
	}
	return parsed
}
