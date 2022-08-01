/**
 @author: 15973
 @date: 2022/07/16
 @note:
**/
package setting

import "github.com/spf13/viper"

type Setting struct {
	vp *viper.Viper
}

func NewSetting(configs ...string) (*Setting, error) {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
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
