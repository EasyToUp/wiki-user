package response

import "wiki-user/server/config"

type SysConfigResponse struct {
	Config config.Server `json:"config"`
}
