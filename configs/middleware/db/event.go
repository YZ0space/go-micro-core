package db

import (
	"context"
	"fmt"
	"github.com/aka-yz/go-micro-core/configs/log"
	"os"
	"strings"
	"time"
)

const (
	UnitTestMod = "ut"
	ReleaseMod  = "release"
)

// 主要用来进行sql事件监听
type sqlEventReceiver struct {
	costThreshold int64
	logLength     int64
	mod           string
}

// TODO: 此处可以进行监控上报
func NewEventReceiver(dbname string, costThreshold int64, lenThreshold int64) *sqlEventReceiver {
	mod := ReleaseMod
	if os.Getenv("env") == "unittest" {
		mod = UnitTestMod
	}
	return &sqlEventReceiver{
		costThreshold: costThreshold,
		logLength:     lenThreshold,
		mod:           mod,
	}
}

// Event receives a simple notification when various events occur
func (s *sqlEventReceiver) Event(eventName string) {
	log.Infof(context.TODO(), "DB Event name %v", eventName)
}

// EventKv receives a notification when various events occur along with
// optional key/value data
func (s *sqlEventReceiver) EventKv(eventName string, kvs map[string]string) {
	log.Infof(context.TODO(), "DB EventKv name %v kv %v", eventName, kvs)
}

// EventErr receives a notification of an error if one occurs
func (s *sqlEventReceiver) EventErr(eventName string, err error) error {
	log.Errorf(context.TODO(), "DB EventErr name:%v err:%v", eventName, err)
	return err
}

// EventErrKv receives a notification of an error if one occurs along with
// optional key/value data
func (s *sqlEventReceiver) EventErrKv(eventName string, err error, kvs map[string]string) error {
	if err != nil && strings.Contains(err.Error(), "Duplicate entry") {
		log.Warnf(context.TODO(), "DB EventErr name:%v err:%v kvs:%v", eventName, err, kvs)
	} else {
		log.Errorf(context.TODO(), "DB EventErr name:%v err:%v kvs:%v", eventName, err, kvs)
	}
	//table, operator := table(kvs["sql"])
	return err
}

// Timing receives the time an event took to happen
func (s *sqlEventReceiver) Timing(eventName string, nanoseconds int64) {
	t := int64(time.Duration(nanoseconds) / time.Millisecond)
	if t > s.costThreshold {
		log.Infof(context.TODO(), "DB Timing name:%v cost:%v", eventName, time.Duration(nanoseconds).String())
	}
}

// TimingKv receives the time an event took to happen along with optional key/value data
func (s *sqlEventReceiver) TimingKv(eventName string, nanoseconds int64, kvs map[string]string) {
	t := int64(time.Duration(nanoseconds) / time.Millisecond)
	if t > s.costThreshold {
		for key, val := range kvs {
			if len(val) > int(s.logLength) { //控制日志长度
				kvs[key] = val[:s.logLength] + "..."
			}
		}
		log.Infof(context.TODO(), "DB TimingKv name:%v kv:%v cost:%v", eventName, kvs, time.Duration(nanoseconds).String())
	}
	if s.mod == UnitTestMod {
		log.Infof(context.TODO(), "DB TimingKv name:%v kv:%v", eventName, kvs)
	}
	//table, operator := table(kvs["sql"])

}

// SELECT * FROM {table} WHERE
// UPDATE `push_data_tab_20200401` SET `push_flag` = 1 WHERE (`user_id` IN (442547)) AND (`biz_id` = 'dc5d7e5b0efa438d97f466d66257b121')
// INSERT INTO `crm_shop_attach_tab` (`id`,`shop_id`,`calculate_buyer_time`,`buyer_num`,`extra`,`is_delete`,`ctime`,`mtime`) VALUES (0,439510,1586102400,9,'',0,1585735589,1585735589)
func table(query string) (name string, operator string) {
	qs := strings.Split(query, " ")
	for i, s := range qs {
		if (s == "FROM" || s == "from") && i < len(qs)-1 {
			if qs[0] == "DELETE" || qs[0] == "delete" {
				return strings.Replace(qs[i+1], "`", "", -1), "DELETE"
			}
			if qs[0] == "SELECT" || qs[0] == "select" {
				return qs[i+1], "SELECT"
			}
		}

		if (s == "UPDATE" || s == "update") && i < len(qs)-1 {
			return strings.Replace(qs[i+1], "`", "", -1), "UPDATE"
		}

		if (s == "INSERT" || s == "insert") && i < len(qs)-2 && (qs[i+1] == "INTO" || qs[i+1] == "into") {
			return strings.Replace(qs[i+2], "`", "", -1), "INSERT"
		}
	}

	return query, " "
}

func dbName(dataSource string) string {
	idx := strings.Index(dataSource, "/")
	if idx == -1 {
		panic(fmt.Sprintf("datasource err:%v", dataSource))
	}

	dataSource = dataSource[idx+1:]
	idx = strings.Index(dataSource, "?")
	if idx == -1 {
		panic(fmt.Sprintf("datasource err:%v", dataSource))
	}
	return dataSource[:idx]
}
