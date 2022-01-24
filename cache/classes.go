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
	ClassTypeDef ClassType = 0 // 行政班
	ClassTypeVirtual ClassType = 1 //虚拟班
)

const (
	DeviceTypeOther DeviceType = 0
	DeviceTypeWall DeviceType = 1 //榜样墙（6屏）
	DeviceTypeWide DeviceType = 2 //榜样栏（双屏）
	DeviceTypeBoard DeviceType = 3 //班级驿站
	DeviceTypeCard DeviceType = 4 // 班牌
)

type ClassStatus uint8

type DeviceType uint8

type ClassType uint8

type ClassInfo struct {
	maxGrade  uint8
	baseInfo
	School    string
	Master    string
	Assistant string
	EnrolDate proxy.DateInfo
	Number    uint16
	Type      ClassType
	Members   []proxy.ClassMember
	Teachers  []string
	Devices   []proxy.DeviceInfo
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
	mine.School = db.School
	mine.EnrolDate = db.EnrolDate
	mine.Number = db.Number
	mine.Assistant = db.Assistant
	mine.Members = db.Students
	mine.Type = ClassType(db.Type)
	mine.Teachers = db.Teachers
	mine.Devices = db.Devices
	if mine.Devices == nil {
		mine.Devices = make([]proxy.DeviceInfo, 0, 2)
	}
	if mine.Teachers == nil {
		mine.Teachers = make([]string, 0, 1)
	}
	if mine.Members == nil {
		mine.Members = make([]proxy.ClassMember, 0, 1)
	}
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
		mine.Operator = operator
	}
	return err
}

func (mine *ClassInfo)UpdateAssistant(master, operator string) error {
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


func (mine *ClassInfo)HadTeacher(teacher string) bool {
	for _, s := range mine.Teachers {
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
		mine.Teachers = append(mine.Teachers, teacher)
	}
	return err
}

func (mine *ClassInfo)SubtractTeacher(teacher string) error {
	if !mine.HadTeacher(teacher) {
		return nil
	}
	err := nosql.SubtractClassTeacher(mine.UID, teacher)
	if err == nil {
		for i:= 0;i < len(mine.Teachers);i += 1 {
			if mine.Teachers[i] == teacher {
				if i == len(mine.Teachers) - 1 {
					mine.Teachers = append(mine.Teachers[:i])
				}else{
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

func (mine *ClassInfo)GetStudentsByStatus(st StudentStatus) []string {
	list := make([]string, 0, len(mine.Members))
	for _, student := range mine.Members {
		if student.Status == uint8(st) {
			info := mine.GetStudent(student.Student)
			if info != nil {
				list = append(list, info.Student)
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

func (mine *ClassInfo)GetStudent(uid string) *proxy.ClassMember {
	for _, student := range mine.Members {
		if student.Student == uid {
			return &student
		}
	}
	return nil
}

func (mine *ClassInfo)HadDevice(uid string) bool {
	for _, device := range mine.Devices {
		if device.UID == uid {
			return true
		}
	}
	return false
}

func (mine *ClassInfo)HadProduct(tp uint8) bool {
	for _, device := range mine.Devices {
		if device.Type == tp {
			return true
		}
	}
	return false
}

func (mine *ClassInfo)BindDevice(uid, remark string, tp DeviceType) error {
	if mine.HadDevice(uid) {
		return nil
	}
	info := proxy.DeviceInfo{
		UID: uid,
		Remark: remark,
		Type: uint8(tp),
	}
	err := nosql.AppendClassDevice(mine.UID, info)
	if err == nil {
		mine.Devices = append(mine.Devices, info)
	}
	return err
}

func (mine *ClassInfo)UnbindDevice(uid string) error {
	if !mine.HadDevice(uid) {
		return nil
	}

	err := nosql.SubtractClassDevice(mine.UID, uid)
	if err == nil {
		for i:= 0;i < len(mine.Devices);i += 1 {
			if mine.Devices[i].UID == uid {
				if i == len(mine.Devices) - 1 {
					mine.Devices = append(mine.Devices[:i])
				}else{
					mine.Devices = append(mine.Devices[:i], mine.Devices[i+1:]...)
				}
				break
			}
		}
	}
	return err
}

func (mine *ClassInfo)RemoveStudent(uid, remark string, id uint64, st StudentStatus) error {
	if mine.HadStudent(uid) {
		return nil
	}
	uuid := fmt.Sprintf("%s-%d", mine.UID, id)
	var err error
	err = nosql.SubtractClassStudent(mine.UID, uuid)
	if st == StudentDelete {
		if err == nil {
			for i:= 0;i < len(mine.Members);i += 1 {
				if mine.Members[i].Student == uid {
					if i == len(mine.Members) - 1 {
						mine.Members = append(mine.Members[:i])
					}else{
						mine.Members = append(mine.Members[:i], mine.Members[i+1:]...)
					}
					break
				}
			}
		}
	}else if st == StudentLeave {
		tmp := proxy.ClassMember{
			UID: uuid,
			Student: uid,
			Status: uint8(st),
			Remark: remark,
		}
		err = nosql.AppendClassStudent(mine.UID, tmp)
		if err == nil{
			for i:= 0;i < len(mine.Members);i += 1 {
				if mine.Members[i].Student == uid {
					mine.Members[i] = tmp
					break
				}
			}
		}

	}
	return err
}


