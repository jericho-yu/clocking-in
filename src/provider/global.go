package provider

import (
	"clocking-in/src/provider/settingProvider"

	"github.com/jericho-yu/filesystem/filesystem"
)

var (
	Setting *settingProvider.SettingProvider
	RootDir = filesystem.NewFileSystemByRelative("../../../")
)
