package service

import (
	"dodevops-api/api/system/dao"
	"dodevops-api/api/system/model"
	"dodevops-api/common/result"
	"dodevops-api/common/util"
	"dodevops-api/pkg/jwt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type ISysAminService interface {
	Login(c *gin.Context, dto model.LoginDto)
	CreateSysAdmin(c *gin.Context, dto model.AddSysAdminDto)
	GetSysAdminInfo(c *gin.Context, id int)
	UpdateSysAdmin(c *gin.Context, dto model.UpdateSysAdminDto)
	DeleteSysAdminById(c *gin.Context, dto model.SysAdminIdDto)
	UpdateSysAdminStatus(c *gin.Context, dto model.UpdateSysAdminStatusDto)
	ResetSysAdminPassword(c *gin.Context, dto model.ResetSysAdminPasswordDto)
	GetSysAdminList(c *gin.Context, pageSize, pageNum int, username, status, beginTime, endTime string)
	UpdatePersonal(c *gin.Context, dto model.UpdatePersonalDto)
	UpdatePersonalPassword(c *gin.Context, dto model.UpdatePersonalPasswordDto)
}

type SysAdminServiceImpl struct{}

func SysAdminService() SysAdminServiceImpl {
	return SysAdminServiceImpl{}
}

func (s SysAdminServiceImpl) Login(c *gin.Context, dto model.LoginDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingLoginParameter), result.ApiCode.GetMessage(result.ApiCode.MissingLoginParameter))
		return
	}

	ip := c.ClientIP()
	code := util.RedisStore{}.Get(dto.IdKey, true)
	if len(code) == 0 {
		dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "captcha expired", 2)
		result.Failed(c, int(result.ApiCode.VerificationCodeHasExpired), result.ApiCode.GetMessage(result.ApiCode.VerificationCodeHasExpired))
		return
	}

	verifyRes := CaptVerify(dto.IdKey, dto.Image)
	if !verifyRes {
		dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "captcha invalid", 2)
		result.Failed(c, int(result.ApiCode.CAPTCHANOTTRUE), result.ApiCode.GetMessage(result.ApiCode.CAPTCHANOTTRUE))
		return
	}

	sysAdmin := dao.SysAdminDetail(dto)
	if sysAdmin.ID == 0 || sysAdmin.Password != util.EncryptionMd5(dto.Password) {
		ldapAdmin, ldapErr := AuthenticateByLDAP(dto.Username, dto.Password)
		if ldapErr == nil && ldapAdmin != nil {
			sysAdmin = *ldapAdmin
		} else {
			dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "password invalid", 2)
			result.Failed(c, int(result.ApiCode.PASSWORDNOTTRUE), result.ApiCode.GetMessage(result.ApiCode.PASSWORDNOTTRUE))
			return
		}
	}

	const disabledStatus = 2
	if sysAdmin.Status == disabledStatus {
		dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "account disabled", 2)
		result.Failed(c, int(result.ApiCode.STATUSISENABLE), result.ApiCode.GetMessage(result.ApiCode.STATUSISENABLE))
		return
	}

	tokenString, _ := jwt.GenerateTokenByAdmin(sysAdmin)
	dao.CreateSysLoginInfo(dto.Username, ip, util.GetRealAddressByIP(ip), util.GetBrowser(c), util.GetOs(c), "login success", 1)

	var leftMenuVo []model.LeftMenuVo
	leftMenuList := dao.QueryLeftMenuList(sysAdmin.ID)
	for _, value := range leftMenuList {
		menuSvoList := dao.QueryMenuVoList(sysAdmin.ID, value.Id)
		item := model.LeftMenuVo{
			Id:          value.Id,
			MenuName:    value.MenuName,
			Icon:        value.Icon,
			Url:         value.Url,
			MenuSvoList: menuSvoList,
		}
		leftMenuVo = append(leftMenuVo, item)
	}

	permissionList := dao.QueryPermissionList(sysAdmin.ID)
	stringList := make([]string, 0, len(permissionList))
	for _, value := range permissionList {
		stringList = append(stringList, value.Value)
	}

	result.Success(c, map[string]interface{}{
		"token":          tokenString,
		"sysAdmin":       sysAdmin,
		"leftMenuList":   leftMenuVo,
		"permissionList": stringList,
	})
}

func (s SysAdminServiceImpl) CreateSysAdmin(c *gin.Context, dto model.AddSysAdminDto) {
	err := validator.New().Struct(dto)
	if err != nil {
		result.Failed(c, int(result.ApiCode.MissingNewAdminParameter), result.ApiCode.GetMessage(result.ApiCode.MissingNewAdminParameter))
		return
	}
	ok := dao.CreateSysAdmin(dto)
	if !ok {
		result.Failed(c, int(result.ApiCode.USERNAMEALREADYEXISTS), result.ApiCode.GetMessage(result.ApiCode.USERNAMEALREADYEXISTS))
		return
	}
	result.Success(c, ok)
}

func (s SysAdminServiceImpl) GetSysAdminInfo(c *gin.Context, id int) {
	result.Success(c, dao.GetSysAdminInfo(id))
}

func (s SysAdminServiceImpl) UpdateSysAdmin(c *gin.Context, dto model.UpdateSysAdminDto) {
	result.Success(c, dao.UpdateSysAdmin(dto))
}

func (s SysAdminServiceImpl) DeleteSysAdminById(c *gin.Context, dto model.SysAdminIdDto) {
	sysAdmin := dao.GetSysAdminInfo(int(dto.Id))
	if sysAdmin.Username == "admin" || sysAdmin.Username == "root" {
		result.Failed(c, int(result.ApiCode.FAILED), "涓嶅厑璁稿垹闄dmin/root鐢ㄦ埛")
		return
	}
	dao.DeleteSysAdminById(dto)
	result.Success(c, true)
}

func (s SysAdminServiceImpl) UpdateSysAdminStatus(c *gin.Context, dto model.UpdateSysAdminStatusDto) {
	dao.UpdateSysAdminStatus(dto)
	result.Success(c, true)
}

func (s SysAdminServiceImpl) ResetSysAdminPassword(c *gin.Context, dto model.ResetSysAdminPasswordDto) {
	dao.ResetSysAdminPassword(dto)
	result.Success(c, true)
}

func (s SysAdminServiceImpl) GetSysAdminList(c *gin.Context, pageSize, pageNum int, username, status, beginTime, endTime string) {
	if pageSize < 1 {
		pageSize = 10
	}
	if pageNum < 1 {
		pageNum = 1
	}

	list, count := dao.GetSysAdminList(pageSize, pageNum, username, status, beginTime, endTime)
	pageResult := result.PageResult{
		List:     list,
		Total:    count,
		Page:     pageNum,
		PageSize: pageSize,
	}
	result.Success(c, pageResult)
}

func (s SysAdminServiceImpl) UpdatePersonal(c *gin.Context, dto model.UpdatePersonalDto) {
	result.Success(c, dao.UpdatePersonal(dto))
}

func (s SysAdminServiceImpl) UpdatePersonalPassword(c *gin.Context, dto model.UpdatePersonalPasswordDto) {
	result.Success(c, dao.UpdatePersonalPassword(dto))
}
