package cache

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"omo.msa.school/tool"
	"sort"
	"time"
)

const (
	ClassTypeDef     ClassType = 0 // 行政班
	ClassTypeVirtual ClassType = 1 //虚拟班
)

type ClassType uint8

type ClassInfo struct {
	maxGrade uint8
	baseInfo
	School    string
	Master    string
	Assistant string
	EnrolDate proxy.DateInfo
	Number    uint16
	Type      ClassType
	Members   []proxy.ClassMember
	Teachers  []string
}

func (mine *ClassInfo) Grade() uint8 {
	return calculateGrade(mine.EnrolDate)
}

func (mine *ClassInfo) GetStatus() StudentStatus {
	if mine.Grade() > mine.maxGrade {
		return StudentFinish
	} else {
		return StudentActive
	}
}

func (mine *ClassInfo) initInfo(grade uint8, db *nosql.Class) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.maxGrade = grade
	mine.Master = db.Master
	mine.School = db.School
	mine.EnrolDate = db.EnrolDate
	mine.Number = db.Number
	mine.Assistant = db.Assistant
	mine.Members = db.Students
	mine.Type = ClassType(db.Type)
	mine.Teachers = db.Teachers
	if mine.Teachers == nil {
		mine.Teachers = make([]string, 0, 1)
		_ = nosql.UpdateClassTeachers(mine.UID, mine.Operator, mine.Teachers)
	}
	if mine.Members == nil {
		mine.Members = make([]proxy.ClassMember, 0, 1)
		_ = nosql.UpdateClassStudents(mine.UID, mine.Operator, mine.Members)
	}
}

func (mine *ClassInfo) FullName() string {
	return fmt.Sprintf("%d年级%d班", mine.Grade(), mine.Number)
}

func (mine *ClassInfo) remove(operator string) error {
	return nosql.RemoveClass(mine.UID, operator)
}

func (mine *ClassInfo) UpdateInfo(name, operator string) error {
	if name == mine.Name {
		return nil
	}
	err := nosql.UpdateClassBase(mine.UID, name, operator)
	if err == nil {
		mine.Name = name
	}
	return err
}

func (mine *ClassInfo) UpdateMaster(master, operator string) error {
	if mine.Master == master {
		return nil
	}
	err := nosql.UpdateClassMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.Operator = operator
	}
	return err
}

func (mine *ClassInfo) UpdateAssistant(master, operator string) error {
	if mine.Assistant == master {
		return nil
	}
	err := nosql.UpdateClassAssistant(mine.UID, master, operator)
	if err == nil {
		mine.Assistant = master
		mine.Operator = operator
	}
	return err
}

func (mine *ClassInfo) HadTeacher(teacher string) bool {
	for _, s := range mine.Teachers {
		if s == teacher {
			return true
		}
	}
	return false
}

func (mine *ClassInfo) AppendTeacher(teacher string) error {
	if mine.HadTeacher(teacher) {
		return nil
	}
	err := nosql.AppendClassTeacher(mine.UID, teacher)
	if err == nil {
		mine.Teachers = append(mine.Teachers, teacher)
	}
	return err
}

func (mine *ClassInfo) SubtractTeacher(teacher string) error {
	if !mine.HadTeacher(teacher) {
		return nil
	}
	err := nosql.SubtractClassTeacher(mine.UID, teacher)
	if err == nil {
		for i := 0; i < len(mine.Teachers); i += 1 {
			if mine.Teachers[i] == teacher {
				if i == len(mine.Teachers)-1 {
					mine.Teachers = append(mine.Teachers[:i])
				} else {
					mine.Teachers = append(mine.Teachers[:i], mine.Teachers[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *ClassInfo) AddStudent(info *StudentInfo) error {
	if info == nil {
		return errors.New("the student is nil")
	}
	if mine.HadStudent(info.UID) {
		return nil
	}
	uuid := fmt.Sprintf("%s-%d", mine.UID, info.ID)
	tmp := proxy.ClassMember{
		UID:     uuid,
		Student: info.UID,
		Status:  uint8(StudentActive),
		Updated: time.Now(),
	}
	err := nosql.AppendClassStudent(mine.UID, tmp)
	if err == nil {
		mine.Members = append(mine.Members, tmp)
	}
	return err
}

func (mine *ClassInfo) GetStudentsNumber() int {
	return len(mine.Members)
}

func (mine *ClassInfo) GetStudentsByStatus(st StudentStatus) []string {
	list := make([]string, 0, len(mine.Members))
	for _, item := range mine.Members {
		if item.Status == uint8(st) {
			list = append(list, item.Student)
		}
	}
	return list
}

func (mine *ClassInfo) GetStudents() []string {
	list := make([]string, 0, len(mine.Members))
	for _, item := range mine.Members {
		list = append(list, item.Student)
	}
	return list
}

func (mine *ClassInfo) HadStudent(uid string) bool {
	if uid == "" {
		return false
	}
	for _, item := range mine.Members {
		if item.Student == uid && item.Status == uint8(StudentActive) {
			return true
		}
	}
	return false
}

func (mine *ClassInfo) HadStudentByStatus(uid string, st StudentStatus) bool {
	if uid == "" {
		return false
	}
	for _, item := range mine.Members {
		if st == StudentAll {
			if item.Student == uid {
				return true
			}
		} else {
			if item.Student == uid && item.Status == uint8(st) {
				return true
			}
		}

	}
	return false
}

func (mine *ClassInfo) IsEmpty() bool {
	if mine.Members == nil || len(mine.Members) < 1 {
		return true
	} else {
		return false
	}
}

func (mine *ClassInfo) GetStudentStatus(uid string) StudentStatus {
	if mine.Grade() > mine.maxGrade {
		return StudentFinish
	}
	for _, student := range mine.Members {
		if student.Student == uid {
			return StudentStatus(student.Status)
		}
	}
	return StudentDelete
}

func (mine *ClassInfo) GetStudent(uid string) *proxy.ClassMember {
	for _, student := range mine.Members {
		if student.Student == uid {
			return &student
		}
	}
	return nil
}

func (mine *ClassInfo) RemoveStudent(uid, remark string, id uint64, st StudentStatus) error {
	if !mine.HadStudent(uid) {
		return nil
	}
	var err error
	err = nosql.SubtractClassStudent(mine.UID, uid)
	if st == StudentDelete {
		if err == nil {
			for i := 0; i < len(mine.Members); i += 1 {
				if mine.Members[i].Student == uid {
					if i == len(mine.Members)-1 {
						mine.Members = append(mine.Members[:i])
					} else {
						mine.Members = append(mine.Members[:i], mine.Members[i+1:]...)
					}
					break
				}
			}
		}
	} else if st == StudentLeave {
		tmp := proxy.ClassMember{
			UID:     fmt.Sprintf("%s-%d", mine.UID, id),
			Student: uid,
			Status:  uint8(st),
			Updated: time.Now(),
			Remark:  remark,
		}
		err = nosql.AppendClassStudent(mine.UID, tmp)
		if err == nil {
			for i := 0; i < len(mine.Members); i += 1 {
				if mine.Members[i].Student == uid {
					mine.Members[i] = tmp
					break
				}
			}
		}

	}
	return err
}

//region Class Fun
func (mine *SchoolInfo) CreateClasses(name, enrol, operator string, number, kind uint16) ([]*ClassInfo, error) {
	mine.initClasses()
	if number < 0 {
		return nil, errors.New("the number must not more than -1")
	}
	date := new(proxy.DateInfo)
	err := date.Parse(enrol)
	if err != nil {
		return nil, err
	}
	var array []*ClassInfo
	if number == 0 {
		array = make([]*ClassInfo, 0, 1)
		info, _ := mine.createClass(name, enrol, operator, 0, kind)
		if info != nil {
			array = append(array, info)
		}
	} else {
		list := mine.GetClassesByEnrol(date.Year, date.Month)
		array = make([]*ClassInfo, 0, number)
		count := len(list)
		var length int = int(number)
		if count > 0 {
			diff := int(number) - count
			if diff < 1 {
				diff = 0
				array = list
				// return nil,errors.New("the class had existed")
			}
			length = diff
		}
		for i := 0; i < length; i += 1 {
			info, _ := mine.createClass(name, enrol, operator, uint16(i+count+1), kind)
			if info != nil {
				array = append(array, info)
			}
		}
	}

	return array, nil
}

func (mine *SchoolInfo) createClass(name, enrol, operator string, number, kind uint16) (*ClassInfo, error) {
	db := new(nosql.Class)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetClassNextID()
	db.CreatedTime = time.Now()
	db.School = mine.UID
	db.Name = name
	db.Type = uint8(kind)
	db.Creator = operator
	db.EnrolDate.Parse(enrol)
	db.Number = number
	db.Teachers = make([]string, 0, 0)
	db.Students = make([]proxy.ClassMember, 0, 1)
	err1 := nosql.CreateClass(db)
	if err1 != nil {
		return nil, err1
	}
	class := new(ClassInfo)
	class.initInfo(mine.MaxGrade(), db)
	mine.classes = append(mine.classes, class)
	return class, nil
}

func (mine *SchoolInfo) hadClassByEnrol(enrol string) bool {
	for _, item := range mine.classes {
		if item.EnrolDate.String() == enrol {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) hadClass(uid string) bool {
	for _, item := range mine.classes {
		if item.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) GetClassesByEnrol(year uint16, month time.Month) []*ClassInfo {
	mine.initClasses()
	list := make([]*ClassInfo, 0, 10)
	for _, item := range mine.classes {
		if item.EnrolDate.Equal(year, month) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) GetClassesByGrade(grade uint8) []*ClassInfo {
	mine.initClasses()
	list := make([]*ClassInfo, 0, 10)
	for _, item := range mine.classes {
		if item.Grade() == grade {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) GetClasses(status StudentStatus) []*ClassInfo {
	mine.initClasses()
	list := make([]*ClassInfo, 0, 50)
	for _, item := range mine.classes {
		if status == item.GetStatus() {
			list = append(list, item)
		} else {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) GetClassesByPage(page, number uint32, st int32) (uint32, uint32, []*ClassInfo) {
	mine.initClasses()
	if number < 1 {
		number = 10
	}
	var classes []*ClassInfo
	if st > -1 {
		classes = mine.GetClasses(StudentStatus(st))
	} else {
		classes = mine.classes
	}

	total := uint32(len(classes))
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, classes
	}
	sort.Slice(classes, func(i, j int) bool {
		return classes[i].EnrolDate.Year < classes[j].EnrolDate.Year
	})
	//list := make([]*ClassInfo, 0, number)
	//for i := 0;i < len(mine.classes);i += 1{
	//	t := uint32(i) / number + 1
	//	if t == page {
	//		list = append(list, mine.classes[i])
	//	}
	//}
	total, max, list := checkPage(page, number, mine.classes)
	return total, max, list
}

func (mine *SchoolInfo) GetClass(uid string) *ClassInfo {
	if uid == "" {
		return nil
	}
	mine.initClasses()
	for _, item := range mine.classes {
		if item.UID == uid {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetClassByMaster(teacher string) *ClassInfo {
	if teacher == "" {
		return nil
	}
	mine.initClasses()
	for _, item := range mine.classes {
		if item.Master == teacher {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetClassByStudent(uid string, st StudentStatus) *ClassInfo {
	if uid == "" {
		return nil
	}
	mine.initClasses()
	for _, item := range mine.classes {
		if item.HadStudentByStatus(uid, st) {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetClassesByMaster(master string) []*ClassInfo {
	if master == "" {
		return nil
	}
	mine.initClasses()
	list := make([]*ClassInfo, 0, 2)
	for _, item := range mine.classes {
		if item.Master == master {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) GetClassesByAssistant(master string) []*ClassInfo {
	mine.initClasses()
	list := make([]*ClassInfo, 0, 2)
	for _, item := range mine.classes {
		if item.Assistant == master {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) GetClassesByTeacher(teacher string) []*ClassInfo {
	mine.initClasses()
	list := make([]*ClassInfo, 0, 2)
	for _, item := range mine.classes {
		if len(item.Teachers) > 1 && tool.HasItem(item.Teachers, teacher) {
			list = append(list, item)
		}
		if item.Master == teacher && !mine.isClassRepeated(list, item.UID) {
			list = append(list, item)
		}
		if item.Assistant == teacher && !mine.isClassRepeated(list, item.UID) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) GetClassesUIDsByTeacher(teacher string) []string {
	mine.initClasses()
	list := make([]string, 0, 2)
	for _, item := range mine.classes {
		if len(item.Teachers) > 1 && tool.HasItem(item.Teachers, teacher) {
			list = append(list, item.UID)
		}
		if item.Master == teacher && !tool.HasItem(list, item.UID) {
			list = append(list, item.UID)
		}
		if item.Assistant == teacher && !tool.HasItem(list, item.UID) {
			list = append(list, item.UID)
		}
	}
	return list
}

func (mine *SchoolInfo) isClassRepeated(list []*ClassInfo, uid string) bool {
	for _, class := range list {
		if class.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) GetClassByEntity(entity string, st StudentStatus) *ClassInfo {
	student := mine.getStudentByEntity(entity)
	if student == nil {
		return nil
	}
	return mine.GetClassByStudent(student.UID, st)
}

func (mine *SchoolInfo) checkClass(name, operator string, enrol proxy.DateInfo, class, kind uint16) *ClassInfo {
	var info *ClassInfo
	info = mine.GetClassByEnrol(enrol, class)
	if info == nil {
		_, err := mine.CreateClasses(name, enrol.String(), operator, class, kind)
		if err == nil {
			info = mine.GetClassByEnrol(enrol, class)
		}
	}
	return info
}

func (mine *SchoolInfo) GetClassByNO(grade uint8, number uint16) *ClassInfo {
	mine.initClasses()
	for _, item := range mine.classes {
		g := item.Grade()
		if g == grade && item.Number == number {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetClassByEnrol(enrol proxy.DateInfo, number uint16) *ClassInfo {
	mine.initClasses()
	for _, item := range mine.classes {
		g := item.EnrolDate.Year
		if g == enrol.Year && item.Number == number {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) RemoveClass(uid, operator string) error {
	info := mine.GetClass(uid)
	if info == nil {
		return errors.New("not found the class")
	}
	err := info.remove(operator)
	if err == nil {
		for i := 0; i < len(mine.classes); i += 1 {
			if mine.classes[i].UID == uid {
				if i == len(mine.classes)-1 {
					mine.classes = append(mine.classes[:i])
				} else {
					mine.classes = append(mine.classes[:i], mine.classes[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

//endregion
