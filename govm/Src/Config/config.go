package Config

import (
	"SparkEven/govm/Src/Common/Def"
	"github.com/spf13/viper"
)

var (
	configMap      = make(map[string]interface{}, 8)
	isLoadFinished = false
)

func UpdateConfig(k string, v interface{}) {
	if isLoadFinished {
		panic("")
	}

	configMap[k] = v
}

func GetValue[T int32 | int64 | string | bool](k string) (v T, exists bool) {
	temp, exists := configMap[k]
	if exists {
		v = temp.(T)
	}

	return
}

func Load(path string) error {
	viper.SetConfigName(Def.Config_FileName)
	viper.SetConfigType("toml")
	if len(path) == 0 {
		path = "."
	}

	viper.AddConfigPath(path)

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	for k, v := range viper.AllSettings() {
		configMap[k] = v
	}

	return err
}
