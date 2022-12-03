package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy/nosql"
	"omo.msa.school/tool"
	"time"
)

const (
	ScheduleStatusFroze = 0 //冻结
	ScheduleStatusIdle  = 1 //发布
)

type ScheduleInfo struct {
	Status    uint8
	LimitMax  uint32
	LimitMin  uint32
	StartTime uint64 //报名开始时间
	EndTime   uint64 // 报名截止时间
	baseInfo
	Name   string
	Scene  string
	Lesson string //课程
	Place  string //地址
	Date   string //日期 2022-04-04
	Times  string //期间时间 12:30-13:30

	Teachers []string
	Tags     []string
	Users    []string
}

func (mine *ScheduleInfo) initInfo(db *nosql.Schedule) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Status = db.Status
	mine.LimitMin = db.LimitMin
	mine.LimitMax = db.LimitMax

	mine.Name = db.Name
	mine.Scene = db.Scene
	mine.Lesson = db.Lesson
	mine.Place = db.Place
	mine.Date = db.Date.Format("2006-01-02")
	mine.Times = db.During
	mine.StartTime = db.StartTime
	mine.EndTime = db.EndTime
	mine.Tags = db.Tags
	mine.Teachers = db.Teachers
	mine.Users = db.Users
}

func (mine *cacheContext) CreateSchedule(scene, lesson, place, date, times, operator string, min, max uint32, teachers []string) (*ScheduleInfo, error) {
	d, er := parseDate(date)
	if er != nil {
		return nil, er
	}
	db := new(nosql.Schedule)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetScheduleNextID()
	db.CreatedTime = time.Now()
	db.Creator = operator
	db.Scene = scene
	db.Name = ""
	db.Lesson = lesson
	db.Place = place
	db.Date = d
	db.During = times
	db.LimitMin = min
	db.LimitMax = max
	db.Status = ScheduleStatusFroze

	db.Teachers = teachers
	db.Tags = make([]string, 0, 1)
	db.Users = make([]string, 0, 1)
	if db.Teachers == nil {
		db.Teachers = make([]string, 0, 1)
	}
	err := nosql.CreateSchedule(db)
	if err != nil {
		return nil, err
	}
	info := new(ScheduleInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *SchoolInfo) CreateSchedule(lesson, place, date, times, operator string, min, max uint32, teachers []string) (*ScheduleInfo, error) {
	return cacheCtx.CreateSchedule(mine.Scene, lesson, place, date, times, operator, min, max, teachers)
}

func (mine *SchoolInfo) GetSchedule(uid string) (*ScheduleInfo, error) {
	return cacheCtx.GetSchedule(uid)
}

func (mine *SchoolInfo) GetSchedules() ([]*ScheduleInfo, error) {
	return cacheCtx.GetSchedules(mine.Scene)
}

func (mine *SchoolInfo) GetSchedulesByDates(from, to string) ([]*ScheduleInfo, error) {
	return cacheCtx.GetSchedulesByDuring(mine.Scene, from, to)
}

func (mine *SchoolInfo) GetSchedulesByDate(date string) ([]*ScheduleInfo, error) {
	return cacheCtx.GetSchedulesByDate(mine.Scene, date)
}

func (mine *cacheContext) CreateSampleSchedule(scene, date, operator string) (*ScheduleInfo, error) {
	return mine.CreateSchedule(scene, "", "", date, "", operator, 0, 0, nil)
}

func (mine *cacheContext) GetSchedule(uid string) (*ScheduleInfo, error) {
	db, err := nosql.GetSchedule(uid)
	if err != nil {
		return nil, err
	}
	info := new(ScheduleInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetSchedules(scene string) ([]*ScheduleInfo, error) {
	dbs, err := nosql.GetSchedulesByScene(scene)
	if err != nil {
		return nil, err
	}
	list := make([]*ScheduleInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(ScheduleInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetSchedulesByDate(scene, date string) ([]*ScheduleInfo, error) {
	d, er := parseDate(date)
	if er != nil {
		return nil, er
	}
	dbs, err := nosql.GetSchedulesByDate(scene, d)
	if err != nil {
		return nil, err
	}
	list := make([]*ScheduleInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(ScheduleInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetSchedulesByDuring(scene, from, to string) ([]*ScheduleInfo, error) {
	f, er := parseDate(from)
	if er != nil {
		return nil, er
	}
	t, er := parseDate(to)
	if er != nil {
		return nil, er
	}
	dbs, err := nosql.GetSchedulesByDuring(scene, f, t)
	if err != nil {
		return nil, err
	}
	list := make([]*ScheduleInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(ScheduleInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list, nil
}

func (mine *ScheduleInfo) UpdateInfo(lesson, place, times, operator string, max, min uint32, teachers []string) error {
	err := nosql.UpdateScheduleBase(mine.UID, lesson, place, times, operator, max, min, teachers)
	if err == nil {
		mine.Lesson = lesson
		mine.Place = place
		mine.Times = times
		mine.Teachers = teachers
		mine.Operator = operator
		mine.LimitMax = max
		mine.LimitMin = min
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ScheduleInfo) UpdateTags(operator string, tags []string) error {
	err := nosql.UpdateScheduleTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ScheduleInfo) UpdateStatus(operator string, start, end uint64, st uint8) error {
	err := nosql.UpdateScheduleStatus(mine.UID, operator, st, start, end)
	if err == nil {
		mine.StartTime = start
		mine.EndTime = end
		mine.Operator = operator
		mine.Status = st
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ScheduleInfo) UpdateStatus2(operator string, st uint8) error {
	err := nosql.UpdateScheduleStatus(mine.UID, operator, st, 0, 0)
	if err == nil {
		mine.StartTime = 0
		mine.EndTime = 0
		mine.Operator = operator
		mine.Status = st
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ScheduleInfo) AppendUsers(operator string, users []string) error {
	arr := make([]string, 0, 10)
	arr = append(arr, mine.Users...)
	for _, user := range users {
		if !tool.HasItem(arr, user) {
			arr = append(arr, user)
		}
	}

	err := nosql.UpdateScheduleUsers(mine.UID, operator, arr)
	if err == nil {
		mine.Operator = operator
		mine.Users = arr
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ScheduleInfo) SubtractUser(operator string, users []string) error {
	arr := make([]string, 0, 10)

	for _, user := range mine.Users {
		if !tool.HasItem(users, user) {
			arr = append(arr, user)
		}
	}
	err := nosql.UpdateScheduleUsers(mine.UID, operator, arr)
	if err == nil {
		mine.Operator = operator
		mine.Users = arr
		mine.UpdateTime = time.Now()
		//for i := 0; i < len(mine.Users); i += 1 {
		//	if mine.Users[i] == user {
		//		if i == len(mine.Users)-1 {
		//			mine.Users = append(mine.Users[:i])
		//		} else {
		//			mine.Users = append(mine.Users[:i], mine.Users[i+1:]...)
		//		}
		//		break
		//	}
		//}
	}
	return err
}

func (mine *ScheduleInfo) Remove(operator string) error {
	return nosql.RemoveSchedule(mine.UID, operator)
}
