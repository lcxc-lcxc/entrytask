package setting

import "github.com/spf13/viper"

/**
对viper对封装使用
*/

type Setting struct {
	vp *viper.Viper
}

// NewSetting 参数对configs为 从命令行flag输进来的configs路径。
func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config") //设置默认配置文件名为config
	vp.SetConfigType("yaml")   // 设置默认配置文件后缀为.yaml
	for _, config := range configs {
		if config != "" {
			vp.AddConfigPath(config)
		}
	}
	err := vp.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return &Setting{vp}, nil
}
