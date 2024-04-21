package initialize

import (
	"wiki-user/server/global"
	"wiki-user/server/model/example"
)

// CheckSetting 检查设置表是否有初始化
func CheckSetting() {
	var s example.ExaFileUploadAndDownload
	global.WK_DB.Where("id = 1").First(&s)
	if s.ID == 0 {
		s.ID = 1
		global.WK_DB.Create(&s)
	}
}
