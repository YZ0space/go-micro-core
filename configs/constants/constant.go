package constants

// refer: https://help.aliyun.com/document_detail/64769.html
const (
	// CronJobAllSign 所有可能的值
	CronJobAllSign = "*"
	// CronJobListEnumSign 列出枚举值
	CronJobListEnumSign = ","
	// CronJobRangeSign 范围
	CronJobRangeSign = "-"
	// CronJobNoSign 不指定值，仅日期和星期域支持
	CronJobNoSign = "?"
	// CronJobLastSign 仅日期和星期域支持
	CronJobLastSign = "L"
	// CronJobWeekSign 除周末以外的有效工作日，在离指定日期的最近的有效工作日触发事件。
	CronJobWeekSign = "W"
)
