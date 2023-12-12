package cache

import (
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"omo.msa.school/tool"
	"strconv"
	"strings"
	"time"
)

type SchoolInfo struct {
	maxGrade uint8
	Status   uint8
	baseInfo
	Scene   string
	Cover   string
	Support string

	Entity        string
	Honors        []proxy.HonorInfo // 学生荣誉
	Respects      []proxy.HonorInfo // 教师荣誉
	Subjects      []proxy.SubjectInfo
	teacherList   []string
	classes       []*ClassInfo
	isInitClasses bool
}

func (mine *SchoolInfo) initInfo(db *nosql.School) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Cover = db.Cover
	mine.Scene = db.Scene
	mine.Support = db.Support
	mine.Honors = db.Honors
	mine.Respects = db.Respects
	mine.Subjects = db.Subjects
	mine.Entity = db.Entity
	mine.Status = db.Status
	mine.maxGrade = db.Grade
	if mine.maxGrade < 1 {
		mine.maxGrade = 6
	}
	mine.teacherList = db.Teachers
	mine.isInitClasses = false
	if mine.teacherList == nil {
		mine.teacherList = make([]string, 0, 1)
		_ = nosql.UpdateSchoolTeachers(mine.UID, mine.Operator, mine.teacherList)
	}
}

func (mine *SchoolInfo) initClasses() {
	if mine.isInitClasses {
		return
	}
	classes, err := nosql.GetClassesBySchool(mine.UID)
	if err == nil {
		mine.classes = make([]*ClassInfo, 0, len(classes))
		for _, item := range classes {
			tmp := new(ClassInfo)
			tmp.initInfo(mine.MaxGrade(), item)
			if tmp.Grade() <= mine.MaxGrade() {
				mine.classes = append(mine.classes, tmp)
			}
		}
	} else {
		mine.classes = make([]*ClassInfo, 0, 1)
	}
	mine.isInitClasses = true
}

func (mine *SchoolInfo) MaxGrade() uint8 {
	return mine.maxGrade
}

func (mine *SchoolInfo) UpdateInfo(name, remark, operator string) error {
	err1 := nosql.UpdateSchoolBase(mine.UID, name, remark, operator)
	if err1 != nil {
		return err1
	}
	mine.Name = name
	mine.Operator = operator
	return nil
}

func (mine *SchoolInfo) UpdateGrade(grade uint8, operator string) error {
	if grade < 6 {
		grade = 6
	}
	if mine.maxGrade == grade {
		return nil
	}
	err := nosql.UpdateSchoolGrade(mine.UID, grade, operator)
	if err != nil {
		return err
	}
	mine.maxGrade = grade
	return nil
}

func (mine *SchoolInfo) UpdateSupport(operator, support string) error {
	if mine.Support == support {
		return nil
	}
	err := nosql.UpdateSchoolSupport(mine.UID, operator, support)
	if err != nil {
		return err
	}
	mine.Support = support
	return nil
}

func (mine *SchoolInfo) UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateSchoolStatus(mine.UID, operator, st)
	if err != nil {
		return err
	}
	mine.Status = st
	return nil
}

func (mine *SchoolInfo) IsCustodian(phone string) bool {
	for _, info := range mine.AllStudents() {
		if info.HadCustodian(phone) {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) CreateStudentHonor(name, remark, parent string) error {
	for _, item := range mine.Honors {
		if item.Name == name {
			return errors.New("the name had exist")
		}
	}
	uuid := fmt.Sprintf("%s-%s%d", mine.UID, "s", nosql.GetSchoolHonorNextID())
	honor := proxy.HonorInfo{
		UID:    uuid,
		Name:   name,
		Remark: remark,
		Parent: parent,
	}
	err := nosql.AppendSchoolHonor(mine.UID, honor)
	if err == nil {
		mine.Honors = append(mine.Honors, honor)
	}
	return err
}

func (mine *SchoolInfo) GetHonor(student bool, uid string) *proxy.HonorInfo {
	if student {
		for _, honor := range mine.Honors {
			if honor.UID == uid {
				return &honor
			}
		}
	} else {
		for _, honor := range mine.Respects {
			if honor.UID == uid {
				return &honor
			}
		}
	}
	return nil
}

func (mine *SchoolInfo) CreateTeacherHonor(name, remark, parent string) error {
	for _, item := range mine.Respects {
		if item.Name == name {
			return errors.New("the name had exist")
		}
	}
	uuid := fmt.Sprintf("%s-%s%d", mine.UID, "t", nosql.GetSchoolHonorNextID())
	honor := proxy.HonorInfo{
		UID:    uuid,
		Name:   name,
		Remark: remark,
		Parent: parent,
	}
	err := nosql.AppendSchoolRespect(mine.UID, honor)
	if err == nil {
		mine.Respects = append(mine.Respects, honor)
	}
	return err
}

func (mine *SchoolInfo) RemoveHonor(uid string, kind pb.TargetType) error {
	var err error
	if kind == pb.TargetType_TStudent {
		err = nosql.SubtractSchoolHonor(mine.UID, uid)
		if err == nil {
			for i := 0; i < len(mine.Honors); i += 1 {
				if mine.Honors[i].UID == uid {
					if i == len(mine.Honors)-1 {
						mine.Honors = append(mine.Honors[:i])
					} else {
						mine.Honors = append(mine.Honors[:i], mine.Honors[i+1:]...)
					}
					break
				}
			}
		}
	} else {
		err = nosql.SubtractSchoolRespect(mine.UID, uid)
		if err == nil {
			for i := 0; i < len(mine.Respects); i += 1 {
				if mine.Respects[i].UID == uid {
					mine.Respects = append(mine.Respects[:i], mine.Respects[i+1:]...)
					break
				}
			}
		}
	}
	return err
}

func (mine *SchoolInfo) GetSubject(uid string) *proxy.SubjectInfo {
	for _, item := range mine.Subjects {
		if item.UID == uid {
			return &item
		}
	}
	return nil
}

func (mine *SchoolInfo) CreateSubject(name, remark string) error {
	for _, item := range mine.Subjects {
		if item.Name == name {
			return nil
		}
	}
	uuid := fmt.Sprintf("%s-%d", mine.UID, nosql.GetSchoolSubjectNextID())
	info := proxy.SubjectInfo{
		UID:    uuid,
		Name:   name,
		Remark: remark,
	}
	err := nosql.AppendSchoolSubject(mine.UID, info)
	if err == nil {
		mine.Subjects = append(mine.Subjects, info)
	}
	return err
}

func (mine *SchoolInfo) CreateSubjects(items []proxy.TimetableItem) {
	for _, item := range items {
		_ = mine.CreateSubject(item.Name, item.Name)
	}
}

func (mine *SchoolInfo) RemoveSubject(uid string) error {
	var err error
	err = nosql.SubtractSchoolSubject(mine.UID, uid)
	if err == nil {
		for i := 0; i < len(mine.Subjects); i += 1 {
			if mine.Subjects[i].UID == uid {
				mine.Subjects = append(mine.Subjects[:i], mine.Subjects[i+1:]...)
				break
			}
		}
	}
	return err
}

//region Statistic
func (mine *SchoolInfo) GetGradeStudents() []*PairIntInfo {
	list := make([]*PairIntInfo, 0, 6)
	for i := 0; i < 6; i += 1 {
		pair := new(PairIntInfo)
		pair.Key = uint32(i + 1)
		num := 0
		classes := mine.GetClassesByGrade(uint8(i + 1))
		for _, class := range classes {
			num = num + class.GetStudentsNumber()
		}
		pair.Value = uint32(num)
		list = append(list, pair)
	}
	return list
}

//endregion

//region Student Fun
func (mine *SchoolInfo) CreateStudent(data *pb.ReqStudentAdd) (*StudentInfo, *ClassInfo, error) {
	list := make([]proxy.CustodianInfo, 0, 2)
	if data.Custodians != nil {
		for _, custodian := range data.Custodians {
			list = append(list, proxy.CustodianInfo{
				Name:     custodian.Name,
				Phones:   custodian.Phones,
				Identity: custodian.Identify,
			})
		}
	}

	enrol := new(proxy.DateInfo)
	er := enrol.Parse(data.Enrol)
	if er != nil {
		enrol.Year = uint16(time.Now().Year())
		enrol.Month = time.September
		enrol.Day = 1
	}
	student, err := cacheCtx.createStudent(mine.UID, data.Name, data.Sn, data.Card, data.Operator, enrol, uint8(data.Sex), StudentStatus(data.Status), list)
	if err != nil {
		return nil, nil, err
	}
	_ = student.UpdateClassNumber(uint16(data.Number), data.Operator)
	class := mine.GetClass(data.Class)
	if class != nil {
		_ = class.AddStudent(student)
		_ = student.UpdateEnrol(class.EnrolDate, data.Operator)
		_ = student.UpdateClassNumber(class.Number, data.Operator)
	} else {
		class, _ = cacheCtx.GetClassByEnrol(mine.UID, enrol, uint16(data.Number))
		if class != nil {
			_ = student.UpdateEnrol(class.EnrolDate, data.Operator)
		}
	}
	if len(data.Entity) > 0 {
		_ = student.BindEntity(data.Entity, data.Operator)
	}
	return student, class, nil
}

func (mine *SchoolInfo) GetStudentByEntity(entity string) *StudentInfo {
	if entity == "" {
		return nil
	}
	db, err := nosql.GetStudentByEntity(mine.UID, entity)
	if err == nil {
		info := new(StudentInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *SchoolInfo) GetStudentClassByEntity(entity string) (*ClassInfo, *StudentInfo) {
	if entity == "" {
		return nil, nil
	}
	db, err := nosql.GetStudentByEntity(mine.UID, entity)
	if err == nil {
		info := new(StudentInfo)
		info.initInfo(db)
		cla := mine.GetClassByStudent(info.UID, StudentActive)
		return cla, info
	}
	return nil, nil
}

func (mine *SchoolInfo) GetStudentBySN(sn string) *StudentInfo {
	if sn == "" {
		return nil
	}
	db, err := nosql.GetStudentBySN(mine.UID, sn)
	if err == nil {
		info := new(StudentInfo)
		info.initInfo(db)
		return info
	}
	return nil
}

func (mine *SchoolInfo) GetClassAndStudent(uid string) (*ClassInfo, *StudentInfo) {
	if uid == "" {
		return nil, nil
	}
	student := cacheCtx.GetStudent(uid)
	if student == nil {
		student = mine.GetStudentByEntity(uid)
		if student == nil {
			return nil, nil
		}
	}
	mine.initClasses()
	for _, class := range mine.classes {
		students := class.GetStudentsByStatus(StudentActive)
		for _, studentUid := range students {
			if studentUid == uid {
				return class, student
			}
		}
	}
	return nil, student
}

func (mine *SchoolInfo) GetStudentsByCustodian(phone, name string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 2)
	if phone == "" {
		return list
	}
	array, err := nosql.GetStudentsByCustodian(mine.UID, phone)
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

func (mine *SchoolInfo) GetStudentsByEnrol(enrol string, num uint16) []*StudentInfo {
	list := make([]*StudentInfo, 0, 2)
	if enrol == "" {
		return list
	}
	year, _ := strconv.Atoi(enrol)
	array, err := nosql.GetStudentsByEnrol(mine.UID, year)
	if err != nil {
		return list
	}
	for _, student := range array {
		if num < 1 {
			info := new(StudentInfo)
			info.initInfo(student)
			list = append(list, info)
		} else {
			if student.Number == num {
				info := new(StudentInfo)
				info.initInfo(student)
				list = append(list, info)
			}
		}
	}
	return list
}

func (mine *SchoolInfo) GetBindStudents(grades []string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 2)
	max := int(mine.maxGrade)
	if max == 0 {
		max = 6
	}
	now := time.Now()
	year := now.Year() - max
	//month := 8
	//if now.Month() >= 8 {
	//	year += 1
	//}
	array, err := nosql.GetStudentsByYear(mine.UID, year)
	if err != nil {
		return list
	}
	length := len(grades)
	had := false
	if length == 0 || tool.HasItem(grades, mine.Scene) {
		had = true
	}
	for _, student := range array {
		grade := calculateGrade(student.EnrolDate)
		add := had
		if !add && tool.HasItem(grades, strconv.Itoa(int(grade))) {
			add = true
		}
		if studentAlive(student) && add {
			info := new(StudentInfo)
			info.initInfo(student)
			list = append(list, info)
		}
	}
	return list
}

func studentAlive(db *nosql.Student) bool {
	if db.Status == uint8(StudentActive) || db.Status == uint8(StudentUnknown) {
		if len(db.Entity) > 1 {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) GetStudentsByName(name string) []*StudentInfo {
	if name == "" {
		return nil
	}
	list := make([]*StudentInfo, 0, 10)
	if name == "" {
		return list
	}
	all := mine.AllStudents()
	for _, info := range all {
		if strings.Contains(info.Name, name) {
			list = append(list, info)
		}
	}
	return list
}

func (mine *SchoolInfo) CreateSimpleStudent(name, entity, sn, card, operator string, enrol proxy.DateInfo, sex uint8) (*StudentInfo, error) {
	db := new(nosql.Student)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetStudentNextID()
	db.CreatedTime = time.Now()
	db.Name = name
	db.Creator = operator
	db.EnrolDate = enrol
	db.Tags = make([]string, 0, 1)
	db.Entity = entity
	db.Sex = sex
	db.SN = sn
	db.Status = uint8(StudentActive)
	if len(card) == 19 {
		db.IDCard = card[1:]
		db.SID = card
	} else if len(card) == 18 {
		db.SID = "G" + card
		db.IDCard = card
	}

	db.School = mine.UID
	db.Custodians = make([]proxy.CustodianInfo, 0, 1)
	err := nosql.CreateStudent(db)
	if err != nil {
		return nil, err
	}

	student := new(StudentInfo)
	student.initInfo(db)
	return student, nil
}

//
//func (mine *SchoolInfo) createStudent(operator string, data *StudentTemp) (*StudentInfo, error) {
//	list := make([]proxy.CustodianInfo, 0, 1)
//	if data.Custodian.Name != "" {
//		list = append(list, data.Custodian)
//	}
//	date := proxy.DateInfo{
//		Year:  uint16(time.Now().Year() - int(data.Grade) + 1),
//		Month: time.January,
//		Day:   1,
//	}
//	student,err := cacheCtx.createStudent(mine.UID, data.Name, data.SN, data.Card, operator, date, data.Sex, list)
//	if err != nil {
//		return nil, err
//	}
//	mine.appendStudent(student, operator, data.Grade, data.Class, 0)
//	return student, nil
//}

func (mine *SchoolInfo) appendStudent(student *StudentInfo, operator string, grade uint8, num uint16, kind ClassType) {
	if student == nil {
		return
	}
	date := &proxy.DateInfo{
		Year:  uint16(time.Now().Year() - int(grade) + 1),
		Month: time.January,
		Day:   1,
	}
	class := mine.checkClass("", operator, date, num, kind)
	if class != nil {
		_ = class.AddStudent(student)
	}
}

func (mine *SchoolInfo) HadStudentBySN(sn string) bool {
	if sn == "" {
		return false
	}
	for _, info := range mine.AllStudents() {
		if info.SN == sn {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) RemoveStudent(uid, operator string) error {
	if uid == "" {
		return errors.New("the student uid is empty")
	}
	class, info := mine.GetClassAndStudent(uid)
	if info == nil {
		return errors.New("not found the student")
	}
	if info.Remove(operator) {
		if class != nil {
			_ = class.RemoveStudent(uid, "the admin delete student", info.ID, StudentDelete)
		}
	}
	return nil
}

func (mine *SchoolInfo) hadStudentByStatus(uid string, st StudentStatus) bool {
	if uid == "" {
		return false
	}
	for _, class := range mine.classes {
		if class.HadStudentByStatus(uid, st) {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) hadStudent(uid string) bool {
	if uid == "" {
		return false
	}
	all := mine.AllStudents()
	for _, info := range all {
		if info.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) AllStudents() []*StudentInfo {
	students, err := nosql.GetStudentsBySchool(mine.UID)
	if err == nil {
		list := make([]*StudentInfo, 0, len(students))
		for _, db := range students {
			info := new(StudentInfo)
			info.initInfo(db)
			list = append(list, info)
		}
		return list
	} else {
		return make([]*StudentInfo, 0, 1)
	}
}

func (mine *SchoolInfo) AllActEntities() []*StudentInfo {
	mine.initClasses()
	list := make([]*StudentInfo, 0, 200)
	for _, class := range mine.classes {
		for _, item := range class.Members {
			if item.Status == uint8(StudentActive) {
				student := cacheCtx.GetStudent(item.Student)
				if student != nil && len(student.Entity) > 2 {
					list = append(list, student)
				}
			}
		}
	}
	return list
}

func (mine *SchoolInfo) GetAllStudentsByStatus(st StudentStatus, bind bool) []*StudentInfo {
	arr, err := nosql.GetStudentsByStatus(mine.UID, uint32(st))
	list := make([]*StudentInfo, 0, len(arr))
	if err != nil {
		return list
	}
	for _, db := range arr {
		if bind {
			if len(db.Entity) > 2 {
				info := new(StudentInfo)
				info.initInfo(db)
				list = append(list, info)
			}
		} else {
			info := new(StudentInfo)
			info.initInfo(db)
			list = append(list, info)
		}
	}
	return list
}

func (mine *SchoolInfo) SearchStudents(flag string, act bool) []*StudentInfo {
	list := make([]*StudentInfo, 0, 100)
	dbs, err := nosql.GetStudentsByKeyword(mine.UID, flag)
	var sts []uint
	if act {
		sts = []uint{uint(StudentActive), uint(StudentUnknown)}
	} else {
		sts = []uint{uint(StudentLeave), uint(StudentDelete), uint(StudentDelete)}
	}
	if err == nil {
		for _, item := range dbs {
			if tool.HasItemByUint(sts, uint(item.Status)) {
				tmp := new(StudentInfo)
				tmp.initInfo(item)
				list = append(list, tmp)
			}
		}
	}

	return list
}

func (mine *SchoolInfo) GetStudentsByClass(uid string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 100)
	mine.initClasses()
	for _, class := range mine.classes {
		if class.UID == uid {
			array := class.GetStudentsByStatus(StudentActive)
			if len(array) > 0 {
				for _, item := range array {
					if !mine.hadStudentIn(list, item) {
						info := cacheCtx.GetStudent(item)
						if info != nil {
							list = append(list, info)
						}
					}
				}
			}
			break
		}
	}
	return list
}

func (mine *SchoolInfo) hadStudentIn(list []*StudentInfo, uid string) bool {
	for _, info := range list {
		if info.UID == uid {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) GetStudentByCard(sn string) *StudentInfo {
	if sn == "" {
		return nil
	}
	all := mine.AllStudents()
	for _, info := range all {
		if info.IDCard == sn {
			return info
		}
	}
	return nil
}

func (mine *SchoolInfo) GetStudentByCustodian(phone, name string) *StudentInfo {
	if phone == "" {
		return nil
	}
	all := mine.AllStudents()
	for _, info := range all {
		if info.HadCustodian(phone) {
			if name == "" {
				return info
			} else {
				if name == info.Name {
					return info
				}
			}
		}
	}
	return nil
}

func (mine *SchoolInfo) GetAllStudentsByPage(page, number uint32) (uint32, uint32, []*StudentInfo) {
	if number < 1 {
		number = 10
	}
	all := mine.AllStudents()
	if len(all) < 1 {
		return 0, 0, make([]*StudentInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, all)

	return total, maxPage, list
}

func (mine *SchoolInfo) GetStudents(page, number uint32, st StudentStatus) (uint32, uint32, []*StudentInfo) {
	if number < 1 {
		number = 10
	}
	all := mine.GetAllStudentsByStatus(st, false)
	if len(all) < 1 {
		return 0, 0, make([]*StudentInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, all)

	return total, maxPage, list
}

func (mine *SchoolInfo) GetLeaveStudents(page, number uint32) (uint32, uint32, []*StudentInfo) {
	if number < 1 {
		number = 10
	}
	list1 := mine.GetAllStudentsByStatus(StudentDelete, false)
	list2 := mine.GetAllStudentsByStatus(StudentLeave, false)
	list3 := mine.GetAllStudentsByStatus(StudentFinish, false)
	all := make([]*StudentInfo, 0, len(list1)+len(list2)+len(list3))
	all = append(all, list1...)
	all = append(all, list2...)
	all = append(all, list3...)
	if len(all) < 1 {
		return 0, 0, make([]*StudentInfo, 0, 1)
	}
	total, maxPage, arr := checkPage(page, number, all)
	return total, maxPage, arr
}

func (mine *SchoolInfo) GetActiveStudents(page, number uint32) (uint32, uint32, []*StudentInfo) {
	if number < 1 {
		number = 10
	}
	list1 := mine.GetAllStudentsByStatus(StudentActive, false)
	list2 := mine.GetAllStudentsByStatus(StudentUnknown, false)
	all := make([]*StudentInfo, 0, len(list1)+len(list2))
	all = append(all, list1...)
	all = append(all, list2...)
	if len(all) < 1 {
		return 0, 0, make([]*StudentInfo, 0, 1)
	}
	total, maxPage, arr := checkPage(page, number, all)
	return total, maxPage, arr
}

func (mine *SchoolInfo) GetActiveBindStudents(page, number uint32) (uint32, uint32, []*StudentInfo) {
	if number < 1 {
		number = 10
	}
	list1 := mine.GetAllStudentsByStatus(StudentActive, true)
	list2 := mine.GetAllStudentsByStatus(StudentUnknown, true)
	all := make([]*StudentInfo, 0, len(list1)+len(list2))
	all = append(all, list1...)
	all = append(all, list2...)
	if len(all) < 1 {
		return 0, 0, make([]*StudentInfo, 0, 1)
	}
	total, maxPage, arr := checkPage(page, number, all)
	return total, maxPage, arr
}

//endregion
