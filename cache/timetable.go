package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"time"
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

func (mine *TimetableInfo) Delete() error {
	return nosql.RemoveTimetable(mine.UID, "")
}

func (mine *TimetableInfo) UpdateItems(operator string, list []proxy.TimetableItem) error {
	err := nosql.UpdateTimetableItems(mine.UID, operator, list)
	if err == nil {
		mine.Items = list
	}
	return err
}

func (mine *SchoolInfo) CreateTimetable(class, operator string, year uint32, items []proxy.TimetableItem) (*TimetableInfo, error) {
	db := new(nosql.Timetable)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetTimetableNextID()
	db.CreatedTime = time.Now()
	db.Name = ""
	db.Creator = operator
	db.School = mine.UID
	db.Class = class
	db.Year = year
	db.Items = items
	if db.Items == nil {
		db.Items = make([]proxy.TimetableItem, 0, 1)
	}
	err := nosql.CreateTimetable(db)
	if err != nil {
		return nil, err
	}

	info := new(TimetableInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *SchoolInfo) GetTimetablesBy(year uint32) ([]*TimetableInfo, error) {
	list, err := nosql.GetTimetablesBy(mine.UID, year)
	if err != nil {
		return nil, err
	}
	arr := make([]*TimetableInfo, 0, len(list))
	for _, item := range list {
		info := new(TimetableInfo)
		info.initInfo(item)
		arr = append(arr, info)
	}
	return arr, nil
}

func (mine *SchoolInfo) GetTimetable(class string, year uint32) (*TimetableInfo, error) {
	db, err := nosql.GetTimetable(mine.UID, class, year)
	if err != nil {
		return nil, err
	}
	info := new(TimetableInfo)
	info.initInfo(db)
	return info, nil
}
