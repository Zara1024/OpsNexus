// 验证码 服务层
// author xiaoRui

package service

import (
	"dodevops-api/common/util"
	"github.com/mojocn/base64Captcha"
	"image/color"
)

var store = util.RedisStore{}

// 生成验证码
func CaptMake() (id, b64s string) {
	var driver base64Captcha.Driver
	var driverString base64Captcha.DriverString
	// 配置验证码信息
	captchaConfig := base64Captcha.DriverString{
		Height:          70,
		Width:           190,
		NoiseCount:      0,
		ShowLineOptions: 0,
		Length:          4,
		Source:          "1234567890",
		BgColor: &color.RGBA{
			R: 248,
			G: 250,
			B: 252,
			A: 255,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}
	driverString = captchaConfig
	driver = driverString.ConvertFonts()
	captcha := base64Captcha.NewCaptcha(driver, store)
	lid, lb64s, _, _ := captcha.Generate()
	return lid, lb64s
}

// 验证captcha是否正确
func CaptVerify(id string, capt string) bool {
	if store.Verify(id, capt, false) {
		return true
	} else {
		return false
	}
}
