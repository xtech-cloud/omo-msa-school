package proxy

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type DateInfo struct {
	Name string `json:"name" bson:"name"`
	Year uint16 `json:"year" bson:"year"`
	Month time.Month `json:"month" bson:"month"`
	Day uint8 `json:"day" bson:"day"`
}

func (mine *DateInfo)String() string {
	return fmt.Sprintf("%d/%d/%d", mine.Year, mine.Month, mine.Day)
}

func (mine *DateInfo)Parse(msg string) error {
	if len(msg) < 1 {
		return errors.New("the date is empty")
	}
	mine.Name = msg
	array := strings.Split(msg, "/")
	if array != nil && len(array) > 2 {
		year,_ := strconv.ParseUint(array[0], 10, 32)
		mine.Year = uint16(year)
		month,_ := strconv.ParseUint(array[1], 10, 32)
		mine.Month = time.Month(month)
		day,_ := strconv.ParseUint(array[2], 10, 32)
		mine.Day = uint8(day)
		return nil
	}else{
		return errors.New("the split date format is error")
	}
}