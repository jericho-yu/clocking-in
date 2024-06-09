package setting

import (
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v2"
)

// AppSetting App配置
type AppSetting struct {
	ClockIn struct { // 打卡文件
		Filename string `yaml:"filename"`  // Excel文件路径
		StartRow int    `yaml:"start-row"` // 数据起始行号
	} `yaml:"clock-in"`
	Collect struct { // 汇总文件
		Filename string `yaml:"filename"`  // Excel文件路径
		StartRow int    `yaml:"start-row"` // 数据起始行号
	} `yaml:"collect"`
}

// NewAppSetting 实例化：App配置
func NewAppSetting(filename string) *AppSetting {
	var (
		file       []byte
		err        error
		appSetting *AppSetting
	)
	file, err = os.ReadFile(path.Join(filename, "app.yaml"))
	if err != nil {
		println(fmt.Sprintf("读取配置文件（app.yaml）失败：%s", err.Error()))
		panic(fmt.Sprintf("读取配置文件（app.yaml）失败：%s", err.Error()))
	}

	err = yaml.Unmarshal(file, &appSetting)
	if err != nil {
		panic(fmt.Sprintf("解析配置文件（app.yaml）失败：%s", err.Error()))
	}

	return appSetting
}
