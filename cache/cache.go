package cache

import (
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"omo.msa.school/tool"
	"strings"
	"time"
)

func (mine *cacheContext) GetClass(uid string) *ClassInfo {
	if len(uid) < 1 {
		return nil
	}
	for _, item := range mine.schools {
		t := item.GetClass(uid)
		if t != nil {
			return t
		}
	}
	return nil
}

func (mine *cacheContext) GetClassByEnrol(school string, enrol *proxy.DateInfo, num uint16) (*ClassInfo, error) {
	if num < 1 {
		return nil, errors.New("the class number is 0")
	}
	info, err := mine.GetSchoolBy(school)
	if err != nil {
		return nil, err
	}
	tmp := info.GetClassByEnrol(enrol, uint16(num))
	if tmp == nil {
		msg := enrol.String()
		name := fmt.Sprintf("%s-%d", msg, num)
		tmp, err = info.createClass(name, msg, info.Operator, uint16(num), ClassTypeDef)
	}
	return tmp, err
}

func (mine *cacheContext) GetClassByStudent(uid string) *ClassInfo {
	if uid == "" {
		return nil
	}
	for _, item := range mine.schools {
		t := item.GetClassByStudent(uid, StudentActive)
		if t != nil {
			return t
		}
	}
	return nil
}

func (mine *cacheContext) GetTeachersByPage(page, number uint32) (uint32, uint32, []*TeacherInfo) {
	if number < 1 {
		number = 10
	}
	if len(mine.teachers) < 1 {
		return 0, 0, make([]*TeacherInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, mine.teachers)
	return total, maxPage, list
}

func (mine *cacheContext) createTeacher(name, operator, scene, entity, user string, classes, subs []string) (*TeacherInfo, error) {
	db := new(nosql.Teacher)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetTeacherNextID()
	db.CreatedTime = time.Now()
	db.Entity = entity
	db.Name = strings.TrimRight(name, " ")
	db.Classes = classes
	db.Subjects = subs
	db.User = user
	db.Tags = make([]string, 0, 1)
	if db.Subjects == nil {
		db.Subjects = make([]string, 0, 1)
	}
	if db.Classes == nil {
		db.Classes = make([]string, 0, 1)
	}
	db.Histories = make([]proxy.HistoryInfo, 0, 1)
	err1 := nosql.CreateTeacher(db)
	if err1 != nil {
		return nil, err1
	}

	teacher := new(TeacherInfo)
	teacher.initInfo(db)
	mine.teachers = append(mine.teachers, teacher)
	return teacher, nil
}

func (mine *cacheContext) GetTeacher(uid string) *TeacherInfo {
	if uid == "" {
		return nil
	}
	for _, item := range mine.teachers {
		if item.UID == uid {
			return item
		}
	}
	db, err := nosql.GetTeacher(uid)
	if err == nil {
		info := new(TeacherInfo)
		info.initInfo(db)
		mine.teachers = append(mine.teachers, info)
		return info
	}
	return nil
}

func (mine *cacheContext) GetTeacherByEntity(entity string) *TeacherInfo {
	if entity == "" {
		return nil
	}
	for _, item := range mine.teachers {
		if item.Entity == entity {
			return item
		}
	}
	db, err := nosql.GetTeacherByEntity(entity)
	if err == nil {
		info := new(TeacherInfo)
		info.initInfo(db)
		mine.teachers = append(mine.teachers, info)
		return info
	}
	return nil
}

func (mine *cacheContext) GetTeachersByEntity(entity string) *TeacherInfo {
	if entity == "" {
		return nil
	}
	for _, item := range mine.teachers {
		if item.Entity == entity {
			return item
		}
	}
	db, err := nosql.GetTeacherByEntity(entity)
	if err == nil {
		info := new(TeacherInfo)
		info.initInfo(db)
		mine.teachers = append(mine.teachers, info)
		return info
	}
	return nil
}

func (mine *cacheContext) GetTeacherByUser(user string) *TeacherInfo {
	if user == "" {
		return nil
	}
	for _, item := range mine.teachers {
		if item.User == user {
			return item
		}
	}
	db, err := nosql.GetTeacherByUser(user)
	if err == nil {
		info := new(TeacherInfo)
		info.initInfo(db)
		mine.teachers = append(mine.teachers, info)
		return info
	}
	return nil
}

func (mine *cacheContext) GetLeaveTeachers(school string) []*TeacherInfo {
	list := make([]*TeacherInfo, 0, 100)
	if school == "" {
		return list
	}

	dbs, err := nosql.GetLeaveTeachers(school)
	if err == nil {
		for _, db := range dbs {
			info := new(TeacherInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *cacheContext) GetTeacherByName(name string) *TeacherInfo {
	if name == "" {
		return nil
	}
	for _, item := range mine.teachers {
		if item.Name == name {
			return item
		}
	}
	db, err := nosql.GetTeacherByUser(name)
	if err == nil {
		info := new(TeacherInfo)
		info.initInfo(db)
		mine.teachers = append(mine.teachers, info)
		return info
	}
	return nil
}

func (mine *cacheContext) CheckTeacher(entity string) *SchoolInfo {
	for _, school := range mine.schools {
		if school.HadTeacherByEntity(entity) {
			return school
		}
	}
	return nil
}

func (mine *cacheContext) createStudent(owner, name, sn, sid, operator string, enrol *proxy.DateInfo, sex uint8, st StudentStatus, custodians []proxy.CustodianInfo) (*StudentInfo, error) {
	db := new(nosql.Student)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetStudentNextID()
	db.CreatedTime = time.Now()
	db.Entity = ""
	db.Name = strings.TrimRight(name, " ")
	db.Creator = operator
	db.EnrolDate = enrol.Clone()
	db.School = owner
	db.Tags = make([]string, 0, 1)
	if st == StudentDelete {
		st = StudentActive
	}
	db.Status = uint8(st)
	db.Sex = sex
	db.SN = sn
	if len(sid) == 19 {
		db.IDCard = sid[1:]
		db.SID = sid
	} else if len(sid) == 18 {
		db.IDCard = sid
		db.SID = "G" + sid
	}

	if custodians != nil {
		db.Custodians = make([]proxy.CustodianInfo, 0, len(custodians))
		for _, custodian := range custodians {
			if len(custodian.Phones) > 0 {
				db.Custodians = append(db.Custodians, proxy.CustodianInfo{
					Name:     custodian.Name,
					Phones:   custodian.Phones,
					Identity: custodian.Identity,
				})
			}
		}
	} else {
		db.Custodians = make([]proxy.CustodianInfo, 0, 1)
	}

	err := nosql.CreateStudent(db)
	if err != nil {
		return nil, err
	}

	student := new(StudentInfo)
	student.initInfo(db)
	return student, nil
}

func (mine *cacheContext) GetStudent(uid string) *StudentInfo {
	if uid == "" {
		return nil
	}
	db, err := nosql.GetStudent(uid)
	if err == nil {
		info := new(StudentInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *cacheContext) GetStudentsByIDCard(card, phone string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 4)
	if len(card) < 1 {
		return list
	}
	array, err := nosql.GetStudentsByCard(card)
	if err == nil {
		for _, item := range array {
			if item.HadCustodian(phone) {
				info := new(StudentInfo)
				info.initInfo(item)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetStudentsByCustodian(phone, name string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 2)
	if phone == "" {
		return list
	}
	array, err := nosql.GetStudentsByCustodian2(phone)
	if err != nil {
		return list
	}
	for _, db := range array {
		if name == "" {
			info := new(StudentInfo)
			info.initInfo(db)
			list = append(list, info)
		} else {
			if db.Name == name {
				info := new(StudentInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetStudentsByCard(card string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 4)
	if len(card) < 1 {
		return list
	}
	//array, err := nosql.GetStudentsByCard(card)
	//if err == nil {
	//	for _, item := range array {
	//		info := new(StudentInfo)
	//		info.initInfo(item)
	//		list = append(list, info)
	//	}
	//}
	for _, school := range mine.schools {
		for _, student := range school.AllStudents() {
			if student.IDCard == card {
				list = append(list, student)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetStudentsByEntity(uid string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 4)
	if len(uid) < 1 {
		return list
	}
	array, err := nosql.GetStudentsByEntity(uid)
	tArray := make([]string, 0, len(array))
	if err == nil {
		for _, item := range array {
			if !tool.HasItem(tArray, item.UID.Hex()) {
				info := new(StudentInfo)
				info.initInfo(item)
				list = append(list, info)
				tArray = append(tArray, info.UID)
			}
		}
	}
	return list
}

func (mine *cacheContext) GetActiveStudentByEntity(uid string) (*ClassInfo, *StudentInfo) {
	for _, school := range mine.schools {
		student := school.GetStudentByEntity(uid)
		if student != nil {
			class, _ := school.GetStudent(student.UID)
			if class != nil {
				return class, student
			}
		}
	}
	return nil, nil
}

func (mine *cacheContext) GetStudents(array []string) []*StudentInfo {
	if array == nil {
		return nil
	}
	list := make([]*StudentInfo, 0, len(array))
	for _, s := range array {
		item := mine.GetStudent(s)
		if item != nil {
			list = append(list, item)
		}
	}
	return list
}

func (mine *cacheContext) CheckStudentFinish() {
	dbs, _ := nosql.GetAllStudentsByStatus(uint32(StudentActive))
	for _, db := range dbs {
		student := new(StudentInfo)
		student.initInfo(db)
		school, _ := mine.GetSchoolBy(student.School)
		if school != nil {
			if student.Grade() > school.maxGrade {
				_ = student.UpdateStatus(StudentFinish, student.Operator)
			}
		}
	}
}

func (mine *cacheContext) CheckStudentError() {
	db2s, _ := nosql.GetAllStudentsByStatus(uint32(StudentFinish))
	for _, db := range db2s {
		student := new(StudentInfo)
		student.initInfo(db)
		school, _ := mine.GetSchoolBy(student.School)
		if school != nil {
			if student.Grade() < school.maxGrade {
				_ = student.UpdateStatus(StudentActive, student.Operator)
			}
		}
	}
}
