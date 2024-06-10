package excelProvider

import (
	"clocking-in/src/provider"
	"encoding/json"
	"github.com/jericho-yu/filesystem/filesystem"
	"log"
	"regexp"
	"strconv"
	"strings"
)

type (
	// AnalysisCheckingInProvider 分析打卡、月度统计数据
	AnalysisCheckingInProvider struct {
		clockInTimeData map[uint64]map[string]string
		monthData       map[uint64]map[string]string
		analysisResults map[uint64]*AnalysisResult
	}

	// AnalysisResult 分析结果
	AnalysisResult struct {
		Overtimes       []string `json:"overtimes,omitempty"`         // 加班
		WrongOvertimes  []string `json:"wrong_overtimes,omitempty"`   // 错误加班
		YearLeaves      []string `json:"year_leaves,omitempty"`       // 年假
		PersonalLeaves  []string `json:"personal_leaves,omitempty"`   // 事假
		Holidays        []string `json:"holidays,omitempty"`          // 倒休
		MaternityLaves  []string `json:"maternity_laves,omitempty"`   // 产假
		MarriageLeave   []string `json:"marriage_leave,omitempty"`    // 婚假
		Miss            []string `json:"miss,omitempty"`              // 缺卡
		Absenteeism     []string `json:"absenteeism,omitempty"`       // 旷工
		LateTimes       []string `json:"late_times,omitempty"`        // 迟到
		LeaveEarlyTimes []string `json:"leave_early_times,omitempty"` // 早退
	}
)

var AnalysisCheckingInProv AnalysisCheckingInProvider

// New 实例化：分析打卡、月度统计数据
func (AnalysisCheckingInProvider) New(ClockInTimeData, MonthData map[uint64]map[string]string) *AnalysisCheckingInProvider {
	ins := &AnalysisCheckingInProvider{
		clockInTimeData: ClockInTimeData,
		monthData:       MonthData,
	}

	return ins
}

// 初始化：统计信息
func (r *AnalysisCheckingInProvider) initAnalysisResult(rowNumber uint64) *AnalysisResult {
	if val, exist := r.analysisResults[rowNumber]; !exist {
		r.analysisResults[rowNumber] = &AnalysisResult{
			Overtimes:       make([]string, 0),
			WrongOvertimes:  make([]string, 0),
			YearLeaves:      make([]string, 0),
			PersonalLeaves:  make([]string, 0),
			Holidays:        make([]string, 0),
			MaternityLaves:  make([]string, 0),
			MarriageLeave:   make([]string, 0),
			Miss:            make([]string, 0),
			Absenteeism:     make([]string, 0),
			LateTimes:       make([]string, 0),
			LeaveEarlyTimes: make([]string, 0),
		}
		return r.analysisResults[rowNumber]
	} else {
		return val
	}
}

// 获取日期名称
func (r *AnalysisCheckingInProvider) getDateName(columnText string) string {
	if dateName, exist := provider.Setting.App.CheckingIn.Dates[columnText]; exist {
		return dateName
	} else {
		log.Panicf("日期名称获取失败%s %v", columnText, provider.Setting.App.CheckingIn.ClockInTime.Overtimes)
		return ""
	}
}

// 判断时间是否是上午
func (*AnalysisCheckingInProvider) isMorning(hour int) bool {
	return hour >= 0 && hour < 12
}

// 检查加班时间是否合法
func (r *AnalysisCheckingInProvider) checkOvertime(times []string) string {
	// 获取第一位时间
	firstHour := strings.Split(times[0], ":")
	atoi, err := strconv.Atoi(firstHour[0])
	if err != nil {
		log.Panicf("获取加班时间错误")
	}
	if r.isMorning(atoi) {
		return "下班缺卡"
	} else {
		return "上班缺卡"
	}
}

// Analysis 分析
func (r *AnalysisCheckingInProvider) Analysis() {
	r.analysisOvertime()
}

// analysisOvertime 分析加班情况
func (r *AnalysisCheckingInProvider) analysisOvertime() {
	for rowNumber, _ := range r.clockInTimeData {
		// 获取加班列数据
		a := make(map[string]string)
		b := make(map[string]string)
		for _, text := range provider.Setting.App.CheckingIn.ClockInTime.Overtimes {
			a[text] = r.clockInTimeData[rowNumber][text]
			b[text] = r.monthData[rowNumber][text]
		}

		for columnText, v := range a {
			v = strings.TrimSpace(v)
			if v != "休息" {
				// 检测是否加班
				if regexp.MustCompile(`\s+`).ReplaceAllString(b[columnText], "") == "" {
					// 缺卡
					result := r.initAnalysisResult(rowNumber)
					result.WrongOvertimes = append(result.WrongOvertimes, r.getDateName(columnText)+"加班缺卡")
				} else {
					times := strings.Split(regexp.MustCompile(`\s+`).ReplaceAllString(b[columnText], " "), " ")
					if len(times) >= 2 {
						// 当日加班
						result := r.initAnalysisResult(rowNumber)
						result.Overtimes = append(result.Overtimes, r.getDateName(columnText)+"加班")
					} else {
						// 加班失败
						result := r.initAnalysisResult(rowNumber)
						result.WrongOvertimes = append(result.WrongOvertimes, r.getDateName(columnText)+r.checkOvertime(times))
					}
				}
			} else {
				// 判断是否是年假、事假、调休、婚假、产假
				result := r.initAnalysisResult(rowNumber)
				if strings.Contains(v, "年假") {
					result.YearLeaves = append(result.YearLeaves, r.getDateName(columnText)+"年假")
				}
				if strings.Contains(v, "事假") {
					result.PersonalLeaves = append(result.PersonalLeaves, r.getDateName(columnText)+v)
				}
				if strings.Contains(v, "调休") {
					result.Holidays = append(result.Holidays, r.getDateName(columnText)+v)
				}
				if strings.Contains(v, "婚假") {
					result.MarriageLeave = append(result.MarriageLeave, r.getDateName(columnText)+"婚假")
				}
				if strings.Contains(v, "产假") {
					result.MaternityLaves = append(result.MaternityLaves, r.getDateName(columnText)+"产假")
				}
				if strings.Contains(v, "缺卡") {
					result.Miss = append(result.Miss, r.getDateName(columnText)+v)
				}
				if strings.Contains(v, "旷工") {
					result.Absenteeism = append(result.Absenteeism, r.getDateName(columnText)+"旷工")
				}
				if strings.Contains(v, "迟到") {
					result.LateTimes = append(result.LateTimes, r.getDateName(columnText)+v)
				}
				if strings.Contains(v, "早退") {
					result.LeaveEarlyTimes = append(result.LeaveEarlyTimes, r.getDateName(columnText)+v)
				}
			}
		}
	}

	bytes, err := json.Marshal(r.analysisResults)
	if err != nil {
		log.Panicf("序列化json失败：%s", err.Error())
	}
	_, err = filesystem.NewFileSystemByAbs(provider.RootDir.Copy().GetDir()).Join("result.json").WriteBytes(bytes)
	if err != nil {
		log.Panicf("保存到文件失败：%s", err.Error())
	}
}
