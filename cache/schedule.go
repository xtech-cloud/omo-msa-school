package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy/nosql"
	"omo.msa.school/tool"
	"time"
)

const (
	ScheduleStatusFroze = 0 //冻结编辑中
	ScheduleStatusOrder = 1 //发布
	ScheduleStatusIdle  = 2 //开课
)

type ScheduleInfo struct {
	Status    uint8
	LimitMax  uint32
	LimitMin  uint32
	StartTime int64  //报名开始时间
	EndTime   int64  // 报名截止时间
	Date      string //日期 2022-04-04
	baseInfo
	Name   string
	Remark string
	Scene  string
	Lesson string //课程
	Place  string //地址
	Times  string //期间时间 12:30-13:30
	Reason string //取消原因

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
	mine.Remark = db.Remark
	mine.Scene = db.Scene
	mine.Lesson = db.Lesson
	mine.Place = db.Place
	mine.Date = time.Unix(db.Date, 0).Format("2006-01-02")
	mine.Times = db.During
	mine.StartTime = db.StartTime
	mine.EndTime = db.EndTime
	mine.Reason = db.Reason
	mine.Tags = db.Tags
	mine.Teachers = db.Teachers
	mine.Users = db.Users
}

func (mine *cacheContext) CreateSchedule(scene, remark, lesson, place, date, times, operator string, min, max uint32, teachers []string) (*ScheduleInfo, error) {
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
	db.Date = d.Unix()
	db.During = times
	db.LimitMin = min
	db.LimitMax = max
	db.Remark = remark
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

func (mine *SchoolInfo) CreateSchedule(remark, lesson, place, date, times, operator string, min, max uint32, teachers []string) (*ScheduleInfo, error) {
	return cacheCtx.CreateSchedule(mine.Scene, remark, lesson, place, date, times, operator, min, max, teachers)
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

func (mine *cacheContext) CreateSampleSchedule(scene, remark, date, operator string) (*ScheduleInfo, error) {
	return mine.CreateSchedule(scene, remark, "", "", date, "", operator, 0, 0, nil)
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
	dbs, err := nosql.GetSchedulesByDate(scene, d.Unix())
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
	dbs, err := nosql.GetSchedulesByDuring(scene, f.Unix(), t.Unix())
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

func (mine *ScheduleInfo) UpdateInfo(remark, lesson, place, times, operator string, max, min uint32, teachers []string) error {
	err := nosql.UpdateScheduleBase(mine.UID, remark, lesson, place, times, operator, max, min, teachers)
	if err == nil {
		mine.Remark = remark
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

func (mine *ScheduleInfo) UpdateRemark(operator, remark string) error {
	err := nosql.UpdateScheduleRemark(mine.UID, operator, remark)
	if err == nil {
		mine.Remark = remark
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ScheduleInfo) UpdateStatus(operator, reason string, start, end int64, st uint8) error {
	err := nosql.UpdateScheduleStatus(mine.UID, operator, reason, st, start, end)
	if err == nil {
		mine.StartTime = start
		mine.EndTime = end
		mine.Operator = operator
		mine.Status = st
		mine.Reason = reason
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *ScheduleInfo) UpdateStatus2(operator, reason string, st uint8) error {
	err := nosql.UpdateScheduleStatus(mine.UID, operator, reason, st, 0, 0)
	if err == nil {
		mine.StartTime = 0
		mine.EndTime = 0
		mine.Operator = operator
		mine.Status = st
		mine.Reason = reason
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
