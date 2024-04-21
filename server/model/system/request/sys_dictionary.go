package request

import (
	"wiki-user/server/model/common/request"
	"wiki-user/server/model/system"
)

type SysDictionarySearch struct {
	system.SysDictionary
	request.PageInfo
}
