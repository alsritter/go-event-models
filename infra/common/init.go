package common

import (
	commonEvent "alsritter.icu/stgo/infra/common/event"
)

var EventServer commonEvent.EventServerIface

func init() {
	var err error
	EventServer, err = commonEvent.NewMessenger(
		commonEvent.WithAddr("172.26.130.183:9876"),
		commonEvent.WithAdminAddr("172.26.130.183:10911"),
		commonEvent.WithGroupName("myGroupName"))
	if err != nil {
		panic(err)
	}
}
