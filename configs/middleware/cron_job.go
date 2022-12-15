package middleware

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"strconv"
)

// for example: cron.AddFunc("* * * * * ?", func() { dosth() })

type CronJobs struct {
	c *cron.Cron
}

func (j *CronJobs) Start() {
	j.c.Start()
}

func (j *CronJobs) Stop() {
	j.c.Stop()
}

func (j *CronJobs) Remove(jobID cron.EntryID) {
	j.c.Remove(jobID)
}

func RangeField(start, end int) string {
	return fmt.Sprintf("%d-%d", start, end)
}

func GenerateCron(second, minute, hour, day, month, year int) string {
	return fmt.Sprintf("%s %s %s %s %s %s",
		strconv.Itoa(second), strconv.Itoa(minute), strconv.Itoa(hour),
		strconv.Itoa(day), strconv.Itoa(month), strconv.Itoa(year))
}

func NewCron(opts ...cron.Option) *CronJobs {
	return &CronJobs{cron.New(opts...)}
}
