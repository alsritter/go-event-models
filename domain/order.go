package domain

import "time"

type Order struct {
	Id                      string
	Name                    string
	Data                    string
	AutoCancelCountdownTime time.Time // 自动取消截止时间
}
