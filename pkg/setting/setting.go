package setting


//Viper 是适用于 Go 应用程序的完整配置解决方案，是目前 Go 语言中比较流行的文件配置解决方案，
//它支持处理各种不同类型的配置需求和配置格式。
import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

// NewSetting 读取配置文件，获取包含配置信息的setting
func NewSetting() (*Setting, error) {
	vp := viper.New()

	vp.SetConfigName("config")			//设定配置文件的名称为 config
	vp.AddConfigPath("configs/")			//并且设置其配置路径为相对路径 configs/
	vp.SetConfigType("yaml")				//配置类型为 yaml
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}

	return &Setting{vp}, nil
}
