package customizedtype

//gorm中重新格式化json时间数据格式返回给前端
import (
	"database/sql/driver"
	"fmt"
	"time"
)

const timeFormat = "2006-01-02 15:04:05"
const timezone = "Asia/Shanghai"

// 使用
//
//	type Category struct {
//		ID        uint   `json:"id" gorm:"primary_key"`
//		Name      string `json:"name" gorm:"type:varchar(50);not null;unique"`
//		CreatedAt CTime   `json:"created_at" gorm:"type:timestamp"`
//		UpdatedAt CTime   `json:"updated_at" gorm:"type:timestamp"`  //通过使用全局CTime格式化时间返回给前端
//	}
//
// 全局定义
type CTime time.Time

func (t CTime) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t *CTime) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = CTime(now)
	return
}

func (t CTime) IsZero() bool {
	return time.Time(t).IsZero()
}

func (t CTime) String() string {
	return time.Time(t).Format(timeFormat)
}

func (t CTime) Local() time.Time {
	loc, _ := time.LoadLocation(timezone)
	return time.Time(t).In(loc)
}

func (t CTime) Value() (driver.Value, error) {
	var zeroTime time.Time
	var ti = time.Time(t)
	if ti.UnixNano() == zeroTime.UnixNano() {
		return nil, nil
	}
	return ti, nil
}

func (t *CTime) Scan(v interface{}) error {
	value, ok := v.(time.Time)
	if ok {
		*t = CTime(value)
		return nil
	}
	return fmt.Errorf("can not convert %v to timestamp", v)
}
