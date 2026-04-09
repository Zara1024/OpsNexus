package service

import (
	"time"

	"dodevops-api/api/configcenter/dao"
	"dodevops-api/api/configcenter/model"
	"dodevops-api/common/constant"
	"dodevops-api/common/result"
	"dodevops-api/common/util"

	"github.com/gin-gonic/gin"
)

type EcsAuthServiceInterface interface {
	GetEcsAuthList(c *gin.Context)
	GetEcsAuthListWithPage(c *gin.Context, page, pageSize int, name string)
	GetEcsAuthByName(c *gin.Context, name string)
	GetEcsAuthById(c *gin.Context, id uint)
	CreateEcsAuth(c *gin.Context, dto *model.CreateEcsPasswordAuthDto)
	UpdateEcsAuth(c *gin.Context, id uint, dto *model.CreateEcsPasswordAuthDto)
	DeleteEcsAuth(c *gin.Context, id uint)
}

type EcsAuthServiceImpl struct {
	dao dao.EcsAuthDao
}

func (s *EcsAuthServiceImpl) GetEcsAuthList(c *gin.Context) {
	list := s.dao.GetEcsAuthList()
	var vos []model.EcsAuthVo
	for _, auth := range list {
		vos = append(vos, model.EcsAuthVo{
			ID:         auth.ID,
			Name:       auth.Name,
			Type:       auth.Type,
			Username:   auth.Username,
			Password:   auth.Password,
			PublicKey:  auth.PublicKey,
			Port:       auth.Port,
			CreateTime: auth.CreateTime,
			Remark:     auth.Remark,
		})
	}
	result.Success(c, vos)
}

func (s *EcsAuthServiceImpl) GetEcsAuthListWithPage(c *gin.Context, page, pageSize int, name string) {
	list, total := s.dao.GetEcsAuthListWithPage(page, pageSize, name)
	var vos []model.EcsAuthVo
	for _, auth := range list {
		vos = append(vos, model.EcsAuthVo{
			ID:         auth.ID,
			Name:       auth.Name,
			Type:       auth.Type,
			Username:   auth.Username,
			Password:   auth.Password,
			PublicKey:  auth.PublicKey,
			Port:       auth.Port,
			CreateTime: auth.CreateTime,
			Remark:     auth.Remark,
		})
	}

	pageResult := result.PageResult{
		List:     vos,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}
	result.Success(c, pageResult)
}

func (s *EcsAuthServiceImpl) CreateEcsAuth(c *gin.Context, dto *model.CreateEcsPasswordAuthDto) {
	if s.dao.CheckNameExists(dto.Name) {
		result.FailedWithCode(c, constant.ECS_AUTH_NAME_EXISTS, "凭据名称已存在")
		return
	}

	auth := model.EcsAuth{
		Name:       dto.Name,
		Username:   dto.Username,
		Password:   dto.Password,
		Port:       dto.Port,
		CreateTime: util.HTime{Time: time.Now()},
		Remark:     dto.Remark,
		Type:       dto.Type,
		PublicKey:  dto.PublicKey,
	}
	if err := s.dao.CreateEcsAuth(&auth); err != nil {
		result.FailedWithCode(c, constant.ECS_AUTH_CREATE_FAILED, err.Error())
		return
	}
	result.Success(c, true)
}

func (s *EcsAuthServiceImpl) UpdateEcsAuth(c *gin.Context, id uint, dto *model.CreateEcsPasswordAuthDto) {
	auth := model.EcsAuth{
		Name:      dto.Name,
		Username:  dto.Username,
		Password:  dto.Password,
		Port:      dto.Port,
		Remark:    dto.Remark,
		Type:      dto.Type,
		PublicKey: dto.PublicKey,
	}
	if err := s.dao.UpdateEcsAuth(id, &auth); err != nil {
		result.FailedWithCode(c, constant.ECS_AUTH_UPDATE_FAILED, err.Error())
		return
	}
	result.Success(c, true)
}

func (s *EcsAuthServiceImpl) DeleteEcsAuth(c *gin.Context, id uint) {
	auth, err := s.dao.GetEcsAuthById(id)
	if err == nil && auth.Name == "免密认证" {
		result.FailedWithCode(c, constant.ECS_AUTH_DELETE_FAILED, "不允许删除免密认证凭据")
		return
	}
	if err := s.dao.DeleteEcsAuth(id); err != nil {
		result.FailedWithCode(c, constant.ECS_AUTH_DELETE_FAILED, err.Error())
		return
	}
	result.Success(c, true)
}

func (s *EcsAuthServiceImpl) GetEcsAuthByName(c *gin.Context, name string) {
	auth, err := s.dao.GetEcsAuthByName(name)
	if err != nil {
		result.FailedWithCode(c, constant.ECS_AUTH_NOT_FOUND, "凭据不存在")
		return
	}

	vo := model.EcsAuthVo{
		ID:         auth.ID,
		Name:       auth.Name,
		Type:       auth.Type,
		Username:   auth.Username,
		Password:   auth.Password,
		PublicKey:  auth.PublicKey,
		Port:       auth.Port,
		CreateTime: auth.CreateTime,
		Remark:     auth.Remark,
	}
	result.Success(c, vo)
}

func (s *EcsAuthServiceImpl) GetEcsAuthById(c *gin.Context, id uint) {
	auth, err := s.dao.GetEcsAuthById(id)
	if err != nil {
		result.FailedWithCode(c, constant.ECS_AUTH_NOT_FOUND, "凭据不存在")
		return
	}

	vo := model.EcsAuthVo{
		ID:         auth.ID,
		Name:       auth.Name,
		Type:       auth.Type,
		Username:   auth.Username,
		Password:   auth.Password,
		PublicKey:  auth.PublicKey,
		Port:       auth.Port,
		CreateTime: auth.CreateTime,
		Remark:     auth.Remark,
	}
	result.Success(c, vo)
}

func GetEcsAuthService() EcsAuthServiceInterface {
	return &EcsAuthServiceImpl{
		dao: dao.NewEcsAuthDao(),
	}
}
