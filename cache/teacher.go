package cache

import (
	"errors"
	"fmt"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"time"
)

type TeacherInfo struct {
	baseInfo
	entity     string
	user string
	Classes    []string
	Subjects   []string
	Tags       []string
	Histories  []proxy.HistoryInfo
}

func (mine *TeacherInfo)initInfo(db *nosql.Teacher) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.user = db.User
	mine.entity = db.Entity
	mine.Subjects = db.Subjects
	mine.Classes = db.Classes
	mine.Tags = db.Tags
	mine.Histories = db.Histories
	if mine.Histories == nil {
		mine.Histories = make([]proxy.HistoryInfo, 0, 5)
	}
}

func (mine *TeacherInfo)createHistory(school, remark string) *proxy.HistoryInfo {
	info := new(proxy.HistoryInfo)
	uuid := fmt.Sprintf("%s-%d", mine.UID, len(mine.Histories)+1)
	info.UID = uuid
	info.School = school
	info.Remark = remark
	info.Created = time.Now()
	return info
}

func (mine *TeacherInfo)Remove(school, remark string) error {
	info := mine.createHistory(school, remark)
	err := nosql.AppendTeacherHistory(mine.UID, info)
	if err == nil {
		mine.Histories = append(mine.Histories, *info)
	}
	return err
}

func (mine *TeacherInfo)IsActive(school string) bool {
	for _, history := range mine.Histories {
		if history.School == school {
			return false
		}
	}
	return true
}

func (mine *TeacherInfo) UpdateBase(name, cover, operator string, classes, subs []string, props []*pb.PropertyInfo) error {
	var err error
	err = nosql.UpdateTeacherBase(mine.UID,name, operator, classes, subs)
	if err == nil {
		mine.Name = name
		mine.Classes = classes
		mine.Subjects = subs
	}
	return err
}

func (mine *TeacherInfo)HadHonor(honor string) bool {
	return mine.hadTag(honor)
}

func (mine *TeacherInfo)hadSubject(sub string) bool {
	if sub == "" {
		return false
	}
	for _, s := range mine.Subjects {
		if s == sub {
			return true
		}
	}
	return false
}

func (mine *TeacherInfo)hadClass(class string) bool {
	if class == "" {
		return false
	}
	for _, s := range mine.Classes {
		if s == class {
			return true
		}
	}
	return false
}

func (mine *TeacherInfo)hadTag(tag string) bool {
	if tag == "" {
		return false
	}
	for _, s := range mine.Tags {
		if s == tag {
			return true
		}
	}
	return false
}

func (mine *TeacherInfo) appendTag(tag string) error {
	if tag == "" {
		return errors.New("the tag is empty")
	}
	if mine.hadTag(tag) {
		return errors.New("the tag had existed")
	}
	err := nosql.AppendTeacherTag(mine.UID, tag)
	if err == nil {
		mine.Tags = append(mine.Tags, tag)
	}
	return err
}

func (mine *TeacherInfo) subtractTag(tag string) error {
	if tag == "" {
		return errors.New("the tag is empty")
	}
	if !mine.hadTag(tag) {
		return errors.New("the tag not existed")
	}
	err := nosql.SubtractTeacherTag(mine.UID, tag)
	if err == nil {
		for i := 0;i < len(mine.Tags);i += 1 {
			if mine.Tags[i] == tag {
				mine.Tags = append(mine.Tags[:i], mine.Tags[i+1:]...)
				break
			}
		}
	}
	return err
}