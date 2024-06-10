package main

import (
	"clocking-in/src/provider"
	"clocking-in/src/provider/excelProvider"
	"clocking-in/src/provider/settingProvider"
	"flag"
	"fmt"
	"github.com/jericho-yu/filesystem/filesystem"
	"time"
)

func main() {
	path := flag.String("path", ".", "路径")
	flag.Parse()

	provider.RootDir = filesystem.NewFileSystemByRelative(*path)

	provider.Setting = settingProvider.SingleSettingProvider(provider.RootDir.Copy().Join("settings").GetDir())

	readCheckingInProv := excelProvider.ReadCheckingInProv.New(provider.RootDir.Copy().Join(provider.Setting.App.CheckingIn.Filename).GetDir())
	clockInTimeData := readCheckingInProv.ReadClockInTime()
	monthData := readCheckingInProv.ReadMonth()

	excelProvider.AnalysisCheckingInProv.New(clockInTimeData, monthData).Analysis()

	fmt.Print("执行完毕：5秒后关闭……")

	time.Sleep(time.Second * 5)
}
