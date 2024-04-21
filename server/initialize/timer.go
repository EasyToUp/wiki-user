package initialize

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"wiki-user/server/config"
	"wiki-user/server/global"
	utils "wiki-user/server/util"
)

func Timer() {
	if global.WK_CONFIG.Timer.Start {
		for i := range global.WK_CONFIG.Timer.Detail {
			go func(detail config.Detail) {
				var option []cron.Option
				if global.WK_CONFIG.Timer.WithSeconds {
					option = append(option, cron.WithSeconds())
				}
				_, err := global.WK_Timer.AddTaskByFunc("ClearDB", global.WK_CONFIG.Timer.Spec, func() {
					err := utils.ClearTable(global.WK_DB, detail.TableName, detail.CompareField, detail.Interval)
					if err != nil {
						fmt.Println("timer error:", err)
					}
				}, option...)
				if err != nil {
					fmt.Println("add timer error:", err)
				}
			}(global.WK_CONFIG.Timer.Detail[i])
		}
	}
}
