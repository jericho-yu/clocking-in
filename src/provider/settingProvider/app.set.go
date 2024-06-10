package settingProvider

import (
	"log"
	"os"

	"github.com/jericho-yu/filesystem/filesystem"
	"gopkg.in/yaml.v2"
)

// AppSetting App配置
type AppSetting struct {
	// CheckingIn 打卡文件相关配置
	CheckingIn struct {
		Filename    string            `yaml:"filename"` // Excel文件路径
		Dates       map[string]string // 日期相关
		ClockInTime struct {
			StartRow    uint     `yaml:"start-row"`    // 数据起始行
			StartColumn string   `yaml:"start-column"` // 起始列
			EndColumn   string   `yaml:"end-column"`   // 终止列
			SheetName   string   `yaml:"sheet-name"`   // 打卡时间sheet名称
			Overtimes   []string `yaml:"overtimes"`    // 假期所在列
		} `yaml:"clock-in-time"`
		Month struct {
			StartRow    uint64   `yaml:"start-row"`    // 数据起始行号
			StartColumn string   `yaml:"start-column"` // 起始列
			EndColumn   string   `yaml:"end-column"`   // 终止列
			SheetName   string   `yaml:"sheet-name"`   // 月度汇总sheet名称
			Overtimes   []string `yaml:"overtimes"`    // 假期所在列
		} `yaml:"month"` // 月度汇总相关配置
	} `yaml:"checking-in"`

	// Collect 汇总表相关配置
	Collect struct {
		Filename string `yaml:"filename"`  // Excel文件路径
		StartRow uint64 `yaml:"start-row"` // 数据起始行号
	} `yaml:"collect"`
}

// NewAppSetting 实例化：App配置
func NewAppSetting(filename string) *AppSetting {
	var (
		file       []byte
		err        error
		appSetting *AppSetting
	)
	file, err = os.ReadFile(filesystem.NewFileSystemByAbs(filename).Join("app.yaml").GetDir())
	if err != nil {
		log.Panicf("读取配置文件（app.yaml）失败：%s", err.Error())
	}

	err = yaml.Unmarshal(file, &appSetting)
	if err != nil {
		log.Panicf("解析配置文件（app.yaml）失败：%s", err.Error())
	}

	return appSetting
}
