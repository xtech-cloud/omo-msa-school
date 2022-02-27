package cache

import (
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
)

type TimetableInfo struct {
	Year uint32
	baseInfo
	School string
	Class  string
	Items  []proxy.TimetableItem
}

func (mine *TimetableInfo) initInfo(db *nosql.Timetable) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Year = db.Year
	mine.School = db.School
	mine.Class = db.Class
	mine.Items = db.Items
	if mine.Items == nil {
		mine.Items = make([]proxy.TimetableItem, 0, 1)
	}
}

func (mine *TimetableInfo)Delete() error {
	return nosql.RemoveTimetable(mine.UID, "")
}

func (mine *TimetableInfo)UpdateItems(operator string, list []proxy.TimetableItem) error {
	err := nosql.UpdateTimetableItems(mine.UID, operator, list)
	if err == nil {
		mine.Items = list
	}
	return err
}

