package cache

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy/nosql"
	"time"
)

type LessonInfo struct {
	Weight uint32
	baseInfo
	Remark string
	Scene  string
	Cover  string
	Tags   []string
	Assets []string
}

func (mine *LessonInfo) initInfo(db *nosql.Lesson) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Weight = db.Weight
	mine.Remark = db.Remark
	mine.Cover = db.Cover
	mine.Scene = db.Scene
	mine.Assets = db.Assets
	mine.Tags = db.Tags
}

func (mine *SchoolInfo) CreateLesson(name, remark, cover, operator string, tags []string) (*LessonInfo, error) {
	db := new(nosql.Lesson)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetLessonNextID()
	db.CreatedTime = time.Now()
	db.Scene = mine.Scene
	db.Name = name
	db.Remark = remark
	db.Cover = cover
	db.Creator = operator
	db.Tags = tags
	err := nosql.CreateLesson(db)
	if err != nil {
		return nil, err
	}
	info := new(LessonInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *SchoolInfo) GetLesson(uid string) (*LessonInfo, error) {
	return cacheCtx.GetLesson(uid)
}

func (mine *cacheContext) GetLesson(uid string) (*LessonInfo, error) {
	db, err := nosql.GetLesson(uid)
	if err != nil {
		return nil, err
	}
	info := new(LessonInfo)
	info.initInfo(db)
	return info, nil
}

func (mine *cacheContext) GetLessons(scene string) ([]*LessonInfo, error) {
	dbs, err := nosql.GetLessonsByScene(scene)
	if err != nil {
		return nil, err
	}
	list := make([]*LessonInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(LessonInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list, nil
}

func (mine *cacheContext) GetLessonsByCreator(operator string) ([]*LessonInfo, error) {
	dbs, err := nosql.GetLessonsByCreator(operator)
	if err != nil {
		return nil, err
	}
	list := make([]*LessonInfo, 0, len(dbs))
	for _, db := range dbs {
		info := new(LessonInfo)
		info.initInfo(db)
		list = append(list, info)
	}
	return list, nil
}

func (mine *LessonInfo) UpdateInfo(name, remark, operator string, tags []string) error {
	err := nosql.UpdateLessonBase(mine.UID, name, remark, operator, tags)
	if err == nil {
		mine.Name = name
		mine.Remark = remark
		mine.Operator = operator
		mine.Tags = tags
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *LessonInfo) UpdateCover(operator, cover string) error {
	err := nosql.UpdateLessonCover(mine.UID, operator, cover)
	if err == nil {
		mine.Cover = cover
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *LessonInfo) UpdateWeight(operator string, weight uint32) error {
	err := nosql.UpdateLessonWeight(mine.UID, operator, weight)
	if err == nil {
		mine.Weight = weight
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *LessonInfo) UpdateAssets(operator string, arr []string) error {
	err := nosql.UpdateLessonAssets(mine.UID, operator, arr)
	if err == nil {
		mine.Assets = arr
		mine.Operator = operator
		mine.UpdateTime = time.Now()
	}
	return err
}

func (mine *LessonInfo) Remove(operator string) error {
	return nosql.RemoveLesson(mine.UID, operator)
}
