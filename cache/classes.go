package cache

import (
	"errors"
	"fmt"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"time"
)

const (
	ClassActive ClassStatus = 0  // 在读
	ClassFinish ClassStatus = 1  // 毕业
)

const (
	ClassTypeDef = 0 // 行政班
	ClassTypeVirtual = 1 //虚拟班
)

type ClassStatus uint8

type ClassInfo struct {
	maxGrade  uint8
	baseInfo
	School    string
	Master    string
	EnrolDate proxy.DateInfo
	Number    uint16
	Type      uint8
	Members   []proxy.ClassMember
	teachers  []string
}

func (mine *ClassInfo)Grade() uint8 {
	now := time.Now()
	diff := now.Year() - int(mine.EnrolDate.Year) + 1
	if mine.EnrolDate.Month > 8 {
		return uint8(diff)
	}else {
		return uint8(diff)
	}
}

func (mine *ClassInfo)GetStatus() ClassStatus {
	if mine.Grade() > mine.maxGrade {
		return ClassFinish
	}else {
		return ClassActive
	}
}

func (mine *ClassInfo)initInfo(grade uint8, db *nosql.Class)  {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.maxGrade = grade
	mine.Master = db.Master
	mine.EnrolDate = db.EnrolDate
	mine.Number = db.Number
	mine.Members = db.Students
	mine.Type = db.Type
	mine.teachers = db.Teachers
}

func (mine *ClassInfo)FullName() string {
	return fmt.Sprintf("%d年级%d班", mine.Grade(), mine.Number)
}

func (mine *ClassInfo)remove(operator string) error {
	return nosql.RemoveClass(mine.UID, operator)
}

func (mine *ClassInfo)UpdateInfo(name, operator string) error {
	if name == mine.Name {
		return nil
	}
	err := nosql.UpdateClassBase(mine.UID, name, operator)
	if err == nil {
		mine.Name = name
	}
	return err
}

func (mine *ClassInfo)UpdateMaster(master, operator string) error {
	if mine.Master == master {
		return nil
	}
	err := nosql.UpdateClassMaster(mine.UID, master, operator)
	if err == nil {
		mine.Master = master
		mine.masterInfo = Context().GetTeacherBy(master)
	}
	return err
}

func (mine *ClassInfo)HadTeacher(teacher string) bool {
	for _, s := range mine.teachers {
		if s == teacher {
			return true
		}
	}
	return false
}

func (mine *ClassInfo)AppendTeacher(teacher string) error {
	if mine.HadTeacher(teacher) {
		return nil
	}
	err := nosql.AppendClassTeacher(mine.UID, teacher)
	if err == nil {
		mine.teachers = append(mine.teachers, teacher)
	}
	return err
}

func (mine *ClassInfo)SubtractTeacher(teacher string) error {
	if !mine.HadTeacher(teacher) {
		return nil
	}
	err := nosql.SubtractClassTeacher(mine.UID, teacher)
	if err == nil {
		for i:= 0;i < len(mine.teachers);i += 1 {
			if mine.teachers[i] == teacher {
				mine.teachers = append(mine.teachers[:i], mine.teachers[i+1:]...)
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
		UID: uuid,
		Student: info.UID,
		Status: uint8(StudentActive),
	}
	err := nosql.AppendClassStudent(mine.UID, tmp)
	if err == nil {
		mine.Members = append(mine.Members, tmp)
	}
	return err
}

func (mine *ClassInfo)GetStudentsNumber() int {
	return len(mine.Members)
}

func (mine *ClassInfo)GetStudentsByStatus(st StudentStatus) []*StudentInfo {
	list := make([]*StudentInfo, 0, len(mine.Members))
	for _, student := range mine.Members {
		if student.Status == uint8(st) {
			info := mine.GetStudent(student.Student)
			if info != nil {
				list = append(list, info)
			}
		}
	}

	return list
}

func (mine *ClassInfo)HadStudent(uid string) bool {
	if uid == "" {
		return false
	}
	for _, item := range mine.Members {
		if item.Student == uid && item.Status == uint8(StudentActive){
			return true
		}
	}
	return false
}

func (mine *ClassInfo)HadStudentByStatus(uid string, st StudentStatus) bool {
	if uid == "" {
		return false
	}
	for _, item := range mine.Members {
		if st == StudentAll {
			if item.Student == uid {
				return true
			}
		}else{
			if item.Student == uid && item.Status == uint8(st) {
				return true
			}
		}

	}
	return false
}
=
func (mine *ClassInfo)IsEmpty() bool {
	if mine.Members == nil || len(mine.Members) < 1 {
		return true
	}else{
		return false
	}
}

func (mine *ClassInfo)GetStudentStatus(uid string) StudentStatus {
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

func (mine *ClassInfo)RemoveStudent(uid, remark string, st StudentStatus) error {
	info := mine.GetStudent(uid)
	if info == nil {
		return nil
	}
	uuid := fmt.Sprintf("%s-%d", mine.UID, info.ID)
	var err error
	err = nosql.SubtractClassStudent(mine.UID, uuid)
	if st == StudentDelete {
		if err == nil {
			for i:= 0;i < len(mine.members);i += 1 {
				if mine.members[i].Student == info.UID {
					mine.members = append(mine.members[:i], mine.members[i+1:]...)
					break
				}
			}
		}
	}else if st == StudentLeave {
		tmp := proxy.ClassMember{
			UID: uuid,
			Student: info.UID,
			Status: uint8(st),
			Remark: remark,
		}
		err = nosql.AppendClassStudent(mine.UID, tmp)
		if err == nil{
			for i:= 0;i < len(mine.members);i += 1 {
				if mine.members[i].Student == info.UID {
					mine.members[i] = tmp
					break
				}
			}
		}

	}
	return err
}


