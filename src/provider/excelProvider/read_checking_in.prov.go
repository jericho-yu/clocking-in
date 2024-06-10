package excelProvider

import (
	"clocking-in/src/provider"
	"log"

	"github.com/jericho-yu/outil/excel"
)

// ReadCheckingInProvider 打卡时间Excel表
type ReadCheckingInProvider struct{ filename string }

var ReadCheckingInProv ReadCheckingInProvider

// New 实例化：打卡时间Excel表
func (ReadCheckingInProvider) New(filename string) *ReadCheckingInProvider {
	return &ReadCheckingInProvider{filename: filename}
}

// ReadClockInTime 读取数据：打卡时间
func (r *ReadCheckingInProvider) ReadClockInTime() map[uint64]map[string]string {
	startIndex := excel.ColumnTextToNumber(provider.Setting.App.CheckingIn.ClockInTime.StartColumn)
	endIndex := excel.ColumnTextToNumber(provider.Setting.App.CheckingIn.ClockInTime.EndColumn)
	title := make([]string, 0, endIndex-startIndex+1)
	for i := startIndex; i <= endIndex; i++ {
		text, err := excel.ColumnNumberToText(i)
		if err != nil {
			log.Panicf("设置表头失败：%s", err.Error())
		}
		title = append(title, text)
	}

	return r.deleteUselessData(
		excel.NewExcelReader().
			OpenFile(r.filename).
			SetOriginalRow(int(provider.Setting.App.CheckingIn.ClockInTime.StartRow)).
			SetSheetName(provider.Setting.App.CheckingIn.ClockInTime.SheetName).
			// SetFinishedRow(10).
			SetTitle(title).
			Read().
			ToMap(""),
	)
}

// ReadMonth 读取数据：月度汇总相关
func (r *ReadCheckingInProvider) ReadMonth() map[uint64]map[string]string {
	startIndex := excel.ColumnTextToNumber(provider.Setting.App.CheckingIn.Month.StartColumn)
	endIndex := excel.ColumnTextToNumber(provider.Setting.App.CheckingIn.Month.EndColumn)
	title := make([]string, 0, endIndex-startIndex+1)
	for i := startIndex; i <= endIndex; i++ {
		text, err := excel.ColumnNumberToText(i)
		if err != nil {
			log.Panicf("设置表头失败：%s", err.Error())
		}
		title = append(title, text)
	}

	return r.deleteUselessData(
		excel.NewExcelReader().
			OpenFile(r.filename).
			SetOriginalRow(int(provider.Setting.App.CheckingIn.Month.StartRow)).
			SetSheetName(provider.Setting.App.CheckingIn.Month.SheetName).
			// SetFinishedRow(10).
			SetTitle(title).
			Read().
			ToMap(""),
	)
}

// 删除不需要的数据
func (r *ReadCheckingInProvider) deleteUselessData(excelData map[uint64]map[string]string) map[uint64]map[string]string {
	for idx, excelDatum := range excelData {
		if excelDatum["B"] == "未加入考勤组" {
			delete(excelData, idx)
		}
	}
	return excelData
}
