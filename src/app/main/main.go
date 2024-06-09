package main

import (
	"clocking-in/src/provider/excel"
	"clocking-in/src/provider/setting"

	"github.com/jericho-yu/filesystem/filesystem"
)

func main() {
	rootDir := filesystem.NewFileSystemByAbs(".").Join("../../..")
	setting := setting.SingleSettingProvider(rootDir.Copy().Join("settings").GetDir())

	excel.ReadCheckingIn(rootDir.Copy().Join(setting.App.ClockIn.Filename).GetDir())

	// println(setting.App.ClockIn.Filename)
}
