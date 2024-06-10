package main

import (
	"clocking-in/src/provider"
	"clocking-in/src/provider/excelProvider"
	"clocking-in/src/provider/settingProvider"
)

func main() {
	provider.Setting = settingProvider.SingleSettingProvider(provider.RootDir.Copy().Join("settings").GetDir())

	readCheckingInProv := excelProvider.ReadCheckingInProv.New(provider.RootDir.Copy().Join(provider.Setting.App.CheckingIn.Filename).GetDir())
	clockInTimeData := readCheckingInProv.ReadClockInTime()
	monthData := readCheckingInProv.ReadMonth()

	println(clockInTimeData, monthData)

	excelProvider.AnalysisCheckingInProv.New(clockInTimeData, monthData).Analysis()
}
