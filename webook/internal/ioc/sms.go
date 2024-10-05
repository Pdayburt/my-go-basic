package ioc

import (
	"example.com/mod/webook/internal/service/sms"
	"example.com/mod/webook/internal/service/sms/memory"
)

/*func InitSMService() sms.Service {
	return memory.NewService()
}*/

func InitSMService() sms.Service {
	return memory.NewService()
}
