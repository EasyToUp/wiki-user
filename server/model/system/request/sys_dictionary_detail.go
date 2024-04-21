package request

import (
	"wiki-user/server/model/common/request"
	"wiki-user/server/model/system"
)

type SysDictionaryDetailSearch struct {
	system.SysDictionaryDetail
	request.PageInfo
}
