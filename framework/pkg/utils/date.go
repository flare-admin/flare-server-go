package utils

import (
	"sync"
	"time"
)

type TimeZoneManager struct {
	mu  sync.RWMutex
	loc *time.Location
}

// 单例
var manager = &TimeZoneManager{}

func init() {
	// 初始化默认时区
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err)
	}
	manager.mu.Lock()
	manager.loc = loc
	manager.mu.Unlock()
}

// SetTimeZone 设置时区（动态切换）
func SetTimeZone(zone string) error {
	if zone == "" {
		return nil
	}
	loc, err := time.LoadLocation(zone)
	if err != nil {
		return err
	}
	manager.mu.Lock()
	manager.loc = loc
	manager.mu.Unlock()
	return nil
}

// GetTimeNow 获取当前时间
// 返回: 当前时区下的 time.Time 对象
func GetTimeNow() time.Time {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	if manager.loc == nil {
		return time.Now() // 防护
	}
	return time.Now().In(manager.loc)
}

// In 获取指定时间在当前系统时区的时间
func In(t time.Time) time.Time {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	if manager.loc == nil {
		return t
	}
	return t.In(manager.loc)
}

// GetAfterYearData 获取指定年数后的时间戳
// year: 年数
// 返回: 指定年数后的 Unix 时间戳
func GetAfterYearData(year int) int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.AddDate(year, 0, 0).Unix()
}

// GetAfterDayData 获取指定天数后的时间戳
// day: 天数
// 返回: 指定天数后的 Unix 时间戳
func GetAfterDayData(day int64) int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.AddDate(0, 0, int(day)).Unix()
}

// GetNowTimeStr 获取当前时间的格式化字符串
// 返回: 格式为 "2006-01-02 15:04:05" 的时间字符串
func GetNowTimeStr() string {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.Format("2006-01-02 15:04:05")
}

// GetNowDataStr 获取当前时间的格式化字符串
// 返回: 格式为 "2006-01-02" 的时间字符串
func GetNowDataStr() string {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.Format("2006-01-02")
}

// GetStartAndEndOfDay 获取当天的开始和结束时间戳
// 返回: (当天开始时间戳, 当天结束时间戳)
func GetStartAndEndOfDay() (int64, int64) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location())
	return startOfDay.Unix(), endOfDay.Unix()
}

// GetStartAndEndOfBeforeDay 获取指定天数前的开始和结束时间戳
// day: 天数
// 返回: (指定天数前的开始时间戳, 指定天数前的结束时间戳)
func GetStartAndEndOfBeforeDay(day int) (int64, int64) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	startOfYesterday := time.Date(now.Year(), now.Month(), now.Day()-day, 0, 0, 0, 0, now.Location())
	endOfYesterday := time.Date(now.Year(), now.Month(), now.Day()-day, 23, 59, 59, 999999999, now.Location())
	return startOfYesterday.Unix(), endOfYesterday.Unix()
}

// GetTimeBeforeDays 获取指定天数前的时间戳
// days: 天数
// 返回: 指定天数前的 Unix 时间戳
func GetTimeBeforeDays(days int) int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.AddDate(0, 0, -days).Unix()
}

// GeTimestampOfSpecifiedTimeDays 获取指定天数前特定时间的时间戳
// days: 天数
// hour: 小时
// minutes: 分钟
// seconds: 秒
// 返回: 指定时间点的 Unix 时间戳
func GeTimestampOfSpecifiedTimeDays(days int, hour, minutes, seconds int) int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	tim := now.AddDate(0, 0, -days)
	return GeTimestampOfSpecifiedTime(tim, hour, minutes, seconds)
}

// GeTimestampOfSpecifiedTime 获取指定时间点的 Unix 时间戳
// tim: 基准时间
// hour: 小时
// minutes: 分钟
// seconds: 秒
// 返回: 指定时间点的 Unix 时间戳
func GeTimestampOfSpecifiedTime(tim time.Time, hour, minutes, seconds int) int64 {
	today := time.Date(tim.Year(), tim.Month(), tim.Day(), hour, minutes, seconds, 0, tim.Location())
	return today.Unix()
}

// GeTimestampOfSpecifiedTimeDay 获取当天指定时间的时间戳
// hour: 小时
// minutes: 分钟
// seconds: 秒
// 返回: 指定时间点的 Unix 时间戳
func GeTimestampOfSpecifiedTimeDay(hour, minutes, seconds int) int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), hour, minutes, seconds, 0, now.Location())
	return today.Unix()
}

// GetDateStr 获取当前时间的格式化字符串（不含分隔符）
// 返回: 格式为 "20060102150405" 的时间字符串
func GetDateStr() string {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.Format("20060102150405")
}

// GetDateUnix 获取当前时间的 Unix 时间戳（秒级）
// 返回: 当前时间的 Unix 时间戳
func GetDateUnix() int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.Unix()
}

// GetDateUnixMilli 获取当前时间的毫秒级时间戳
// 返回: 当前时间的毫秒级时间戳
func GetDateUnixMilli() int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.UnixMilli()
}

// GetDateUnixMilliStr 获取当前时间的毫秒级时间戳字符串
// 返回: 格式为 "20060102150405000" 的时间字符串
func GetDateUnixMilliStr() string {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	return now.Format("20060102150405000")
}

// IsYesterday 判断时间戳是否为昨天
// timestamp: Unix 时间戳
// 返回: 如果是昨天返回 true，否则返回 false
func IsYesterday(timestamp int64) bool {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, manager.loc)
	yesterdayStart := today.AddDate(0, 0, -1)
	yesterdayEnd := today
	t := time.Unix(timestamp, 0).In(manager.loc)
	return t.After(yesterdayStart) && t.Before(yesterdayEnd)
}

// GetAreaStartAndEndOfDay 获取当天的开始和结束时间戳
// 返回: (当天开始时间戳, 当天结束时间戳)
func GetAreaStartAndEndOfDay() (int64, int64) {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, manager.loc)
	endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, manager.loc)
	return startOfDay.Unix(), endOfDay.Unix()
}

// GetCurrentQuarter 获取当前季度
// 返回: 当前季度（1-4）
func GetCurrentQuarter() int {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	month := now.Month()
	switch {
	case month >= 1 && month <= 3:
		return 1
	case month >= 4 && month <= 6:
		return 2
	case month >= 7 && month <= 9:
		return 3
	case month >= 10 && month <= 12:
		return 4
	default:
		return 0
	}
}

// IsToday 判断时间戳是否为今天
// timestamp: Unix 时间戳
// 返回: 如果是今天返回 true，否则返回 false
func IsToday(timestamp int64) bool {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	t := time.Unix(timestamp, 0).In(manager.loc)
	now := time.Now().In(manager.loc)
	return now.YearDay() == t.YearDay() && now.Year() == t.Year()
}

// GetDateUnixByTime 将指定时间戳转换为当天0点的时间戳
// timestamp: 指定的时间戳
// 返回: 当天0点的时间戳
func GetDateUnixByTime(timestamp int64) int64 {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	t := time.Unix(timestamp, 0).In(manager.loc)
	startOfDay := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, manager.loc)
	return startOfDay.Unix()
}

// IsWorkday 判断指定时间戳是否为工作日
// 工作日为周一至周五，不包括法定节假日
func IsWorkday(timestamp int64) bool {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	t := time.Unix(timestamp, 0).In(manager.loc)

	// 获取星期几（0-6，0表示周日）
	weekday := t.Weekday()

	// 如果是周六或周日，则不是工作日
	if weekday == time.Saturday || weekday == time.Sunday {
		return false
	}

	// TODO: 这里可以添加法定节假日的判断逻辑
	// 目前仅判断周末，后续可以根据需求添加法定节假日的判断

	return true
}

// IsInTimeRange 判断当前时间是否在允许提现时间段
func IsInTimeRange(start, end string) bool {
	if start == "" || end == "" {
		return true
	}
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	now := time.Now().In(manager.loc)
	startTime, _ := time.Parse("15:04:05", start)
	endTime, _ := time.Parse("15:04:05", end)
	cur := now.Format("15:04:05")
	curTime, _ := time.Parse("15:04:05", cur)
	return (curTime.After(startTime) || curTime.Equal(startTime)) && (curTime.Before(endTime) || curTime.Equal(endTime))
}

func getLocation() *time.Location {
	manager.mu.RLock()
	defer manager.mu.RUnlock()
	return manager.loc
}

// LuaDataFunc 提供给Lua脚本使用的函数
var LuaDataFunc = map[string]interface{}{
	// 获取当前时间戳（秒）
	"now": func() int64 {
		return GetDateUnix()
	},
	// 获取当前时间戳（毫秒）
	"now_ms": func() int64 {
		return GetDateUnixMilli()
	},
	// 获取今天的开始时间戳
	"today_start": func() int64 {
		start, _ := GetStartAndEndOfDay()
		return start
	},
	// 获取今天的结束时间戳
	"today_end": func() int64 {
		_, end := GetStartAndEndOfDay()
		return end
	},
	// 获取本周的开始时间戳
	"week_start": func() int64 {
		now := GetTimeNow()
		offset := int(now.Weekday())
		if offset == 0 {
			offset = 7
		}
		return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, -offset+1).Unix()
	},
	// 获取本周的结束时间戳
	"week_end": func() int64 {
		now := GetTimeNow()
		offset := int(now.Weekday())
		if offset == 0 {
			offset = 7
		}
		return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 999999999, now.Location()).AddDate(0, 0, 7-offset).Unix()
	},
	// 获取本月的开始时间戳
	"month_start": func() int64 {
		now := GetTimeNow()
		return time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location()).Unix()
	},
	// 获取本月的结束时间戳
	"month_end": func() int64 {
		now := GetTimeNow()
		return time.Date(now.Year(), now.Month()+1, 0, 23, 59, 59, 999999999, now.Location()).Unix()
	},
	// 判断是否是工作日
	"is_workday": func() bool {
		return IsWorkday(GetDateUnix())
	},
	// 判断是否在时间范围内
	"is_in_time_range": func(start, end string) bool {
		return IsInTimeRange(start, end)
	},
	// 获取指定时间戳的年月日
	"get_date": func(timestamp int64) map[string]int {
		t := time.Unix(timestamp, 0).In(getLocation())
		return map[string]int{
			"year":  t.Year(),
			"month": int(t.Month()),
			"day":   t.Day(),
			"hour":  t.Hour(),
			"min":   t.Minute(),
			"sec":   t.Second(),
		}
	},
	// 计算两个时间戳之间的天数差
	"days_between": func(t1, t2 int64) int {
		d1 := time.Unix(t1, 0).In(getLocation())
		d2 := time.Unix(t2, 0).In(getLocation())
		return int(d2.Sub(d1).Hours() / 24)
	},
	// 判断是否是今天
	"is_today": func(timestamp int64) bool {
		return IsToday(timestamp)
	},
	// 判断是否是昨天
	"is_yesterday": func(timestamp int64) bool {
		return IsYesterday(timestamp)
	},
	// 获取当前季度
	"get_quarter": func() int {
		return GetCurrentQuarter()
	},
	// 获取指定天数前的时间戳
	"get_before_days": func(days int) int64 {
		return GetTimeBeforeDays(days)
	},
	// 获取指定天数前的开始和结束时间戳
	"get_before_day_range": func(days int) (int64, int64) {
		return GetStartAndEndOfBeforeDay(days)
	},
	// 获取指定时间点的Unix时间戳
	"get_specified_time": func(hour, minutes, seconds int) int64 {
		return GeTimestampOfSpecifiedTimeDay(hour, minutes, seconds)
	},
	// 获取指定天数前特定时间的时间戳
	"get_before_specified_time": func(days, hour, minutes, seconds int) int64 {
		return GeTimestampOfSpecifiedTimeDays(days, hour, minutes, seconds)
	},
	// 获取当前时间的格式化字符串
	"get_time_str": func() string {
		return GetNowTimeStr()
	},
	// 获取当前时间的格式化字符串（不含分隔符）
	"get_date_str": func() string {
		return GetDateStr()
	},
	// 获取当前时间的毫秒级时间戳字符串
	"get_date_ms_str": func() string {
		return GetDateUnixMilliStr()
	},
	// 获取指定年数后的时间戳
	"get_after_year": func(year int) int64 {
		return GetAfterYearData(year)
	},
	// 获取指定天数后的时间戳
	"get_after_day": func(day int64) int64 {
		return GetAfterDayData(day)
	},
	// 将时间戳转换为当天0点的时间戳
	"get_day_start": func(timestamp int64) int64 {
		return GetDateUnixByTime(timestamp)
	},
}
