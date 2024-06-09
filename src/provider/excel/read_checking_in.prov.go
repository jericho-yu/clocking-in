package excel

import (
	"fmt"

	outilExcel "github.com/jericho-yu/outil/excel"
)

func ReadCheckingIn(filename string) {
	excel := outilExcel.NewExcelReader().AutoReadBySheetName(filename, "月度汇总 (2)").ToList()
	for _, v := range excel {
		fmt.Printf("%+v\n", v)
	}
}
