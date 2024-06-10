package excelProvider

import (
	"clocking-in/src/provider"
	"encoding/json"
	"github.com/jericho-yu/filesystem/filesystem"
	"github.com/jericho-yu/outil/array"
	"github.com/jericho-yu/outil/dict"
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
		Name            string   `json:"姓名"`             // 姓名
		NeedWrite       bool     `json:"-"`              // 是否需要记载
		Overtimes       []string `json:"加班,omitempty"`   // 加班
		OvertimesDays   uint64   `json:"加班次数,omitempty"` // 加班天数
		WrongOvertimes  []string `json:"加班错误,omitempty"` // 加班错误
		YearLeaves      []string `json:"年假,omitempty"`   // 年假
		PersonalLeaves  []string `json:"事假,omitempty"`   // 事假
		Holidays        []string `json:"调休,omitempty"`   // 调休
		MaternityLaves  []string `json:"产假,omitempty"`   // 产假
		MarriageLeave   []string `json:"婚假,omitempty"`   // 婚假
		Miss            []string `json:"缺卡,omitempty"`   // 缺卡
		Absenteeism     []string `json:"旷工,omitempty"`   // 旷工
		LateTimes       []string `json:"迟到,omitempty"`   // 迟到
		LeaveEarlyTimes []string `json:"早退,omitempty"`   // 早退
	}
)

var AnalysisCheckingInProv AnalysisCheckingInProvider

// New 实例化：分析打卡、月度统计数据
func (AnalysisCheckingInProvider) New(ClockInTimeData, MonthData map[uint64]map[string]string) *AnalysisCheckingInProvider {
	ins := &AnalysisCheckingInProvider{
		clockInTimeData: ClockInTimeData,
		monthData:       MonthData,
		analysisResults: make(map[uint64]*AnalysisResult),
	}

	return ins
}

// 初始化：统计信息
func (r *AnalysisCheckingInProvider) initAnalysisResult(rowNumber uint64) *AnalysisResult {
	if val, exist := r.analysisResults[rowNumber]; !exist {
		r.analysisResults[rowNumber] = &AnalysisResult{
			Overtimes:       make([]string, 0),
			OvertimesDays:   0,
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
func (r *AnalysisCheckingInProvider) checkOvertime(time string) string {
	if time == "缺卡" {
		return "加班缺卡"
	}
	firstHour := strings.Split(time, ":")
	atoi, err := strconv.Atoi(firstHour[0])
	if err != nil {
		log.Panicf("获取加班时间错误：%v", time)
	}
	if r.isMorning(atoi) {
		return "加班下班缺卡"
	} else {
		return "加班上班缺卡"
	}
}

// Analysis 分析
func (r *AnalysisCheckingInProvider) Analysis() {
	var (
		err   error
		bytes []byte
		fs    *filesystem.FileSystem
	)
	r.analysisOvertime()
	r.analysisLeave()

	if len(r.analysisResults) > 0 {
		bytes, err = json.MarshalIndent(r.analysisResults, "", "	")
		if err != nil {
			log.Panicf("序列化json失败：%s", err.Error())
		}
		fs = filesystem.NewFileSystemByAbs(provider.RootDir.Copy().GetDir()).Join("result.json")
		if fs.IsExist {
			err = fs.Delete()
			if err != nil {
				log.Panicf("删除文件错误：%s", err.Error())
			}
		}
		_, err = fs.WriteBytes(bytes)
		if err != nil {
			log.Panicf("保存到文件失败：%s", err.Error())
		}
	}
}

// 分析加班情况
func (r *AnalysisCheckingInProvider) analysisOvertime() {
	for rowNumber, _ := range r.clockInTimeData {
		// 获取加班列数据
		a := make(map[string]string)
		b := make(map[string]string)
		for _, text := range provider.Setting.App.CheckingIn.ClockInTime.Overtimes {
			a[text] = r.monthData[rowNumber][text]
			b[text] = r.clockInTimeData[rowNumber][text]
		}

		result := r.initAnalysisResult(rowNumber)
		result.Name = r.clockInTimeData[rowNumber]["A"]

		// 分析加班情况
		for columnText, v := range a {
			v = strings.TrimSpace(v)
			log.Println(result.Name, v, columnText)
			if v != "休息" {
				// 检测是否加班
				if regexp.MustCompile(`\s+`).ReplaceAllString(b[columnText], "") == "" {
					// 缺卡
					result.NeedWrite = true
					result.WrongOvertimes = append(result.WrongOvertimes, r.getDateName(columnText)+"加班缺卡")
				} else {
					times := strings.Split(regexp.MustCompile(`\s+`).ReplaceAllString(b[columnText], " "), " ")
					if len(times) >= 2 {
						// 当日加班
						result.NeedWrite = true
						result.Overtimes = append(result.Overtimes, r.getDateName(columnText)+"加班")
						if array.In[string](columnText, provider.Setting.App.CheckingIn.ClockInTime.Day3overtimes) {
							// 三薪加班
							result.OvertimesDays += 3
						} else {
							// 普通加班
							result.OvertimesDays++
						}
					} else if len(times) == 1 {
						// 加班失败
						result.NeedWrite = true
						result.WrongOvertimes = append(result.WrongOvertimes, r.getDateName(columnText)+r.checkOvertime(times[0]))
					} else {
						// 加班失败
						result.NeedWrite = true
						result.WrongOvertimes = append(result.WrongOvertimes, r.getDateName(columnText)+"加班缺卡")
					}
				}
			}
		}
	}
}

// 分析请假
func (r *AnalysisCheckingInProvider) analysisLeave() {
	for rowNumber, _ := range r.monthData {
		result := r.initAnalysisResult(rowNumber)
		result.Name = r.monthData[rowNumber]["A"]

		for _, columnText := range dict.GetKeys[string](provider.Setting.App.CheckingIn.Dates) {
			v := strings.TrimSpace(r.monthData[rowNumber][columnText])

			// 判断是否是年假、事假、调休、婚假、产假
			log.Println(result.Name, v, columnText)
			if strings.Contains(v, "年假") {
				result.NeedWrite = true
				result.YearLeaves = append(result.YearLeaves, r.getDateName(columnText)+"年假")
			}
			if strings.Contains(v, "事假") {
				result.NeedWrite = true
				result.PersonalLeaves = append(result.PersonalLeaves, r.getDateName(columnText)+v)
			}
			if strings.Contains(v, "调休") {
				result.NeedWrite = true
				result.Holidays = append(result.Holidays, r.getDateName(columnText)+v)
			}
			if strings.Contains(v, "婚假") {
				result.NeedWrite = true
				result.MarriageLeave = append(result.MarriageLeave, r.getDateName(columnText)+"婚假")
			}
			if strings.Contains(v, "产假") {
				result.NeedWrite = true
				result.MaternityLaves = append(result.MaternityLaves, r.getDateName(columnText)+"产假")
			}
			if strings.Contains(v, "缺卡") {
				result.NeedWrite = true
				result.Miss = append(result.Miss, r.getDateName(columnText)+v)
			}
			if strings.Contains(v, "旷工") {
				result.NeedWrite = true
				result.Absenteeism = append(result.Absenteeism, r.getDateName(columnText)+"旷工")
			}
			if strings.Contains(v, "迟到") {
				result.NeedWrite = true
				result.LateTimes = append(result.LateTimes, r.getDateName(columnText)+v)
			}
			if strings.Contains(v, "早退") {
				result.NeedWrite = true
				result.LeaveEarlyTimes = append(result.LeaveEarlyTimes, r.getDateName(columnText)+v)
			}
		}

	}
}
