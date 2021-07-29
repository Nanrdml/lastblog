package global

import (
	"lastblog/pkg/logger"
	"lastblog/pkg/setting"
)

//在读取了文件的配置信息后，还是不够的，因为我们需要将配置信息和应用程序关联起来，我们才能够去使用它，

//在main函数中绑定的实例，这将作为全局实例
var (
	ServerSetting   *setting.ServerSettingS
	AppSetting      *setting.AppSettingS
	DatabaseSetting *setting.DatabaseSettingS
	Logger 			*logger.Logger
)
