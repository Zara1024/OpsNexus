package controller

import (
	"dodevops-api/api/cmdb/dao"
	"dodevops-api/api/cmdb/model"
	"dodevops-api/api/cmdb/service"
	"dodevops-api/common"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"errors"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CmdbDeviceController struct {
	service *service.CmdbDeviceService
}

func NewCmdbDeviceController() *CmdbDeviceController {
	return &CmdbDeviceController{
		service: service.NewCmdbDeviceService(dao.NewCmdbDeviceDao(common.GetDB())),
	}
}

func (c *CmdbDeviceController) CreateDevice(ctx *gin.Context) {
	var dto model.CreateCmdbDeviceDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}

	device, err := c.service.CreateDevice(dto)
	if err != nil {
		result.Failed(ctx, resolveCmdbDeviceErrorCode(err, constant.CMDB_DEVICE_CREATE_FAILED), err.Error())
		return
	}

	result.Success(ctx, device)
}

func (c *CmdbDeviceController) UpdateDevice(ctx *gin.Context) {
	var dto model.UpdateCmdbDeviceDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}

	device, err := c.service.UpdateDevice(dto)
	if err != nil {
		result.Failed(ctx, resolveCmdbDeviceErrorCode(err, constant.CMDB_DEVICE_UPDATE_FAILED), err.Error())
		return
	}

	result.Success(ctx, device)
}

func (c *CmdbDeviceController) DeleteDevice(ctx *gin.Context) {
	id, err := resolveCmdbDeviceID(ctx)
	if err != nil || id == 0 {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}

	if err := c.service.DeleteDevice(id); err != nil {
		result.Failed(ctx, resolveCmdbDeviceErrorCode(err, constant.CMDB_DEVICE_DELETE_FAILED), err.Error())
		return
	}

	result.Success(ctx, nil)
}

func (c *CmdbDeviceController) ListDevices(ctx *gin.Context) {
	page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(ctx.DefaultQuery("pageSize", "10"))

	devices, total, err := c.service.ListDevices(page, pageSize)
	if err != nil {
		result.Failed(ctx, constant.CMDB_DEVICE_QUERY_FAILED, err.Error())
		return
	}

	result.Success(ctx, gin.H{
		"list":  devices,
		"total": total,
	})
}

func (c *CmdbDeviceController) GetDevice(ctx *gin.Context) {
	id, err := strconv.Atoi(ctx.Query("id"))
	if err != nil || id == 0 {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}

	device, err := c.service.GetDevice(uint(id))
	if err != nil {
		result.Failed(ctx, resolveCmdbDeviceErrorCode(err, constant.CMDB_DEVICE_QUERY_FAILED), err.Error())
		return
	}

	result.Success(ctx, device)
}

func (c *CmdbDeviceController) BatchTestDeviceConnectivity(ctx *gin.Context) {
	var dto model.BatchCmdbDeviceConnectivityDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		result.Failed(ctx, constant.INVALID_PARAMS, "参数错误")
		return
	}

	summary, err := c.service.BatchTestDeviceConnectivity(dto)
	if err != nil {
		result.Failed(ctx, resolveCmdbDeviceErrorCode(err, constant.CMDB_DEVICE_QUERY_FAILED), err.Error())
		return
	}

	result.Success(ctx, summary)
}

func resolveCmdbDeviceID(ctx *gin.Context) (uint, error) {
	if queryID := ctx.Query("id"); queryID != "" {
		id, err := strconv.Atoi(queryID)
		if err != nil {
			return 0, err
		}
		return uint(id), nil
	}

	var dto model.CmdbDeviceIDDto
	if err := ctx.ShouldBindJSON(&dto); err != nil {
		return 0, err
	}
	return dto.ID, nil
}

func resolveCmdbDeviceErrorCode(err error, defaultCode int) int {
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		return constant.CMDB_DEVICE_NOT_FOUND
	case service.IsInvalidCmdbDeviceParams(err):
		return constant.INVALID_PARAMS
	default:
		return defaultCode
	}
}
