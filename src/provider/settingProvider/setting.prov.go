package settingProvider

import "sync"

// SettingProvider 配置提供者
type SettingProvider struct {
	App *AppSetting
}

var (
	settingProviderOnce sync.Once        // 单例模式
	settingProvider     *SettingProvider // 对象：配置提供者
)

// SingleSettingProvider 单利化：配置提供者
func SingleSettingProvider(filename string) *SettingProvider {
	settingProviderOnce.Do(func() {
		settingProvider = &SettingProvider{App: NewAppSetting(filename)}
	})
	return settingProvider
}
