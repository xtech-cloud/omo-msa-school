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
	Entity    string
	User      string
	Classes   []string
	Subjects  []string
	Tags      []string
	Histories []proxy.HistoryInfo
}

func (mine *TeacherInfo) initInfo(db *nosql.Teacher) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.User = db.User
	mine.Entity = db.Entity
	mine.Subjects = db.Subjects
	mine.Classes = db.Classes
	mine.Tags = db.Tags
	mine.Histories = db.Histories
	if mine.Histories == nil {
		mine.Histories = make([]proxy.HistoryInfo, 0, 5)
		_ = nosql.UpdateTeacherHistories(mine.UID, mine.Operator, mine.Histories)
	}
}

func (mine *TeacherInfo) createHistory(school, remark string) *proxy.HistoryInfo {
	info := new(proxy.HistoryInfo)
	uuid := fmt.Sprintf("%s-%d", mine.UID, len(mine.Histories)+1)
	info.UID = uuid
	info.School = school
	info.Remark = remark
	info.Created = uint64(time.Now().Unix())
	return info
}

func (mine *TeacherInfo) Remove(school, remark string) error {
	info := mine.createHistory(school, remark)
	err := nosql.AppendTeacherHistory(mine.UID, info)
	if err == nil {
		mine.Histories = append(mine.Histories, *info)
	}
	return err
}

func (mine *TeacherInfo) IsActive(school string) bool {
	//for _, history := range mine.Histories {
	//	if history.School == school {
	//		return false
	//	}
	//}
	//return true
	return true
}

func (mine *TeacherInfo) UpdateTags(operator string, tags []string) error {
	var err error
	return err
}

func (mine *TeacherInfo) UpdateBase(name, operator string, classes, subs []string) error {
	var err error
	err = nosql.UpdateTeacherBase(mine.UID, name, operator, classes, subs)
	if err == nil {
		mine.Name = name
		mine.Classes = classes
		mine.Subjects = subs
		mine.Operator = operator
	}
	return err
}

func (mine *TeacherInfo) HadHonor(honor string) bool {
	return mine.hadTag(honor)
}

func (mine *TeacherInfo) hadSubject(sub string) bool {
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

func (mine *TeacherInfo) hadClass(class string) bool {
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

func (mine *TeacherInfo) hadTag(tag string) bool {
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
		for i := 0; i < len(mine.Tags); i += 1 {
			if mine.Tags[i] == tag {
				mine.Tags = append(mine.Tags[:i], mine.Tags[i+1:]...)
				break
			}
		}
	}
	return err
}

//region Teacher Fun
func (mine *SchoolInfo) AllTeachers() []*TeacherInfo {
	teachers := make([]*TeacherInfo, 0, len(mine.teacherList))
	for _, item := range mine.teacherList {
		tmp := Context().GetTeacher(item)
		if tmp != nil {
			teachers = append(teachers, tmp)
		}
	}
	return teachers
}

func (mine *SchoolInfo) Teachers() []string {
	return mine.teacherList
}

func (mine *SchoolInfo) GetTeacherByEntity(entity string) *TeacherInfo {
	if entity == "" {
		return nil
	}
	teachers := mine.AllTeachers()
	for _, item := range teachers {
		if item.Entity == entity {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetTeacherByUser(user string) *TeacherInfo {
	if user == "" {
		return nil
	}
	teachers := mine.AllTeachers()
	for _, item := range teachers {
		if item.User == user {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetTeacher(uid string) *TeacherInfo {
	if uid == "" {
		return nil
	}
	return Context().GetTeacher(uid)
}

func (mine *SchoolInfo) GetTeachersByName(name string) []*TeacherInfo {
	list := make([]*TeacherInfo, 0, 10)
	if name == "" {
		return list
	}
	teachers := mine.AllTeachers()
	for _, item := range teachers {
		if item.Name == name {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) GetTeachersBySub(subject string) []*TeacherInfo {
	list := make([]*TeacherInfo, 0, 10)
	if subject == "" {
		return list
	}
	teachers := mine.AllTeachers()
	for _, item := range teachers {
		if item.hadSubject(subject) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) GetTeachersByClass(class string) []*TeacherInfo {
	list := make([]*TeacherInfo, 0, 10)
	if class == "" {
		return list
	}
	info := mine.GetClass(class)
	if info == nil {
		return list
	}
	for _, item := range info.Teachers {
		tmp := cacheCtx.GetTeacher(item)
		if tmp != nil {
			list = append(list, tmp)
		}
	}
	return list
}

func (mine *SchoolInfo) hadTeacher(uid string) bool {
	if uid == "" {
		return false
	}
	for _, item := range mine.teacherList {
		if item == uid {
			return true
		}
	}
	return false
}

//func (mine *SchoolInfo) hadActTeacher(uid string) bool {
//	if uid == "" {
//		return false
//	}
//	mine.AllTeachers()
//	for _, item := range mine.teachers {
//		if item.UID == uid && item.IsActive(mine.UID){
//			return true
//		}
//	}
//	return false
//}

func (mine *SchoolInfo) hadTeacherByUser(uid string) bool {
	if uid == "" {
		return false
	}
	teachers := mine.AllTeachers()
	for _, item := range teachers {
		if item.User == uid {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) CreateTeacher(name, entity, user, operator string, classes, subs []string) (*TeacherInfo, error) {
	if mine.hadTeacherByUser(user) {
		return mine.GetTeacherByUser(user), nil
	}
	teacher, err := Context().createTeacher(name, operator, mine.Scene, entity, user, classes, subs)
	if err != nil {
		return nil, err
	}
	mine.AppendTeacher(teacher)
	return teacher, nil
}

func (mine *SchoolInfo) AppendTeacher(info *TeacherInfo) error {
	if mine.hadTeacher(info.UID) {
		return nil
	}
	err := nosql.AppendSchoolTeacher(mine.UID, info.UID)
	if err == nil {
		mine.teacherList = append(mine.teacherList, info.UID)
	}
	return err
}

func (mine *SchoolInfo) HadTeacherByEntity(entity string) bool {
	teachers := mine.AllTeachers()
	for _, teacher := range teachers {
		if teacher.Entity == entity {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) GetTeachersByPage(page, number uint32) (uint32, uint32, []*TeacherInfo) {
	if number < 1 {
		number = 10
	}
	teachers := mine.AllTeachers()
	if len(teachers) < 1 {
		return 0, 0, make([]*TeacherInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, teachers)
	return total, maxPage, list
}

func (mine *SchoolInfo) RemoveTeacher(entity, remark string) error {
	mine.AllTeachers()
	info := mine.GetTeacherByEntity(entity)
	if info == nil {
		return errors.New("not found the teacher")
	}
	_ = info.Remove(mine.UID, remark)
	err := nosql.SubtractSchoolTeacher(mine.UID, info.UID)
	if err == nil {
		mine.removeTeacherUID(info.UID)
	}
	return err
}

func (mine *SchoolInfo) removeTeacherUID(uid string) {
	for i := 0; i < len(mine.teacherList); i += 1 {
		if mine.teacherList[i] == uid {
			if i == len(mine.teacherList)-1 {
				mine.teacherList = append(mine.teacherList[:i])
			} else {
				mine.teacherList = append(mine.teacherList[:i], mine.teacherList[i+1:]...)
			}
			break
		}
	}
}

func (mine *SchoolInfo) RemoveTeacherByUID(uid, remark string) error {
	mine.AllTeachers()
	info := mine.GetTeacher(uid)
	if info == nil {
		return errors.New("not found the teacher")
	}
	_ = info.Remove(mine.UID, remark)
	err := nosql.SubtractSchoolTeacher(mine.UID, info.UID)
	if err == nil {
		mine.removeTeacherUID(uid)
	}
	return err
}

//endregion
