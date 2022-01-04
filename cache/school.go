package cache

import (
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"sort"
	"time"
)

type SchoolInfo struct {
	maxGrade uint8
	Status   uint8
	baseInfo
	Scene    string
	Cover    string

	Entity     string
	Honors     []proxy.HonorInfo // 学生荣誉
	Respects   []proxy.HonorInfo // 教师荣誉
	Subjects   []proxy.SubjectInfo
	classes    []*ClassInfo
	teachers   []*TeacherInfo
}

func (mine *SchoolInfo)initInfo(db *nosql.School) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Cover = db.Cover
	mine.Scene = db.Scene
	mine.Honors = db.Honors
	mine.Respects = db.Respects
	mine.Subjects = db.Subjects
	mine.Entity = db.Entity
	mine.Status = db.Status
	classes,err := nosql.GetClassesBySchool(mine.UID)
	if err == nil {
		mine.classes = make([]*ClassInfo, 0, len(classes))
		for _, item := range classes {
			tmp := new(ClassInfo)
			tmp.initInfo(mine.MaxGrade(), item)
			mine.classes = append(mine.classes, tmp)
		}
	}else{
		mine.classes = make([]*ClassInfo, 0, 1)
	}

	mine.teachers = make([]*TeacherInfo, 0, len(db.Teachers))
	for _, item := range db.Teachers {
		tmp := Context().GetTeacher(item)
		if tmp != nil && !mine.HadTeacher(tmp.UID){
			mine.teachers = append(mine.teachers, tmp)
		}
	}
}

func (mine *SchoolInfo)MaxGrade() uint8 {
	if mine.maxGrade == 0 {
		return 6
	}
	return mine.maxGrade
}

func (mine *SchoolInfo)UpdateInfo(name, remark, operator string) error {
	err1 := nosql.UpdateSchoolBase(mine.UID, name, remark, operator)
	if err1 != nil {
		return err1
	}
	mine.Name = name
	mine.Operator = operator
	return nil
}

func (mine *SchoolInfo)UpdateGrade(grade uint8, operator string) error {
	if grade < 6 {
		grade = 6
	}
	err := nosql.UpdateSchoolGrade(mine.UID, grade, operator)
	if err != nil {
		return err
	}
	mine.maxGrade = grade
	return nil
}

func (mine *SchoolInfo)UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateSchoolStatus(mine.UID, operator, st)
	if err != nil {
		return err
	}
	mine.Status = st
	return nil
}

func (mine *SchoolInfo)IsCustodian(phone string) bool {
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
		UID: uuid,
		Name: name,
		Remark: remark,
		Parent: parent,
	}
	err := nosql.AppendSchoolHonor(mine.UID, honor)
	if err == nil {
		mine.Honors = append(mine.Honors, honor)
	}
	return err
}

func (mine *SchoolInfo)GetHonor(student bool, uid string) *proxy.HonorInfo {
	if student {
		for _, honor := range mine.Honors {
			if honor.UID == uid {
				return &honor
			}
		}
	}else{
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
		UID: uuid,
		Name: name,
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
			for i := 0;i < len(mine.Honors);i += 1 {
				if mine.Honors[i].UID == uid {
					mine.Honors = append(mine.Honors[:i], mine.Honors[i+1:]...)
					break
				}
			}
		}
	}else{
		err = nosql.SubtractSchoolRespect(mine.UID, uid)
		if err == nil {
			for i := 0;i < len(mine.Respects);i += 1 {
				if mine.Respects[i].UID == uid {
					mine.Respects = append(mine.Respects[:i], mine.Respects[i+1:]...)
					break
				}
			}
		}
	}
	return err
}

func (mine *SchoolInfo)GetSubject(uid string) *proxy.SubjectInfo {
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
			return errors.New("the name had exist")
		}
	}
	uuid := fmt.Sprintf("%s-%d", mine.UID, nosql.GetSchoolSubjectNextID())
	info := proxy.SubjectInfo{
		UID: uuid,
		Name: name,
		Remark: remark,
	}
	err := nosql.AppendSchoolSubject(mine.UID, info)
	if err == nil {
		mine.Subjects = append(mine.Subjects, info)
	}
	return err
}

func (mine *SchoolInfo) RemoveSubject(uid string) error {
	var err error
	err = nosql.SubtractSchoolSubject(mine.UID, uid)
	if err == nil {
		for i := 0;i < len(mine.Subjects);i += 1 {
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
		pair.Key = uint32(i+1)
		num := 0
		classes := mine.GetClassesByGrade(uint8(i+1))
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
func (mine *SchoolInfo) CreateStudent(data *pb.ReqStudentAdd) (*StudentInfo, string, error) {
	list := make([]proxy.CustodianInfo, 0, 2)
	if data.Custodians != nil {
		for _, custodian := range data.Custodians {
			if custodian.Name != "" {
				list = append(list, proxy.CustodianInfo{
					Name:     custodian.Name,
					Phones:    custodian.Phones,
					Identity: custodian.Identify,
				})
			}
		}
	}
	date := proxy.DateInfo{
		Year:  0,
		Month: time.January,
		Day:   1,
	}
	student,err := cacheCtx.createStudent(mine.UID, data.Name, data.Sn, data.Card, data.Operator, date, uint8(data.Sex), list)
	if err != nil {
		return nil, "", err
	}
	class := mine.GetClass(data.Class)
	var classUID = ""
	if class != nil {
		classUID = class.UID
		_ = class.AddStudent(student)
		_ = student.UpdateEnrol(class.EnrolDate, data.Operator)
	}
	return student, classUID, nil
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

func (mine *SchoolInfo) GetStudent(uid string) (*ClassInfo, *StudentInfo) {
	if uid == "" {
		return nil, nil
	}
	for _, class := range mine.classes {
		students := class.GetStudentsByStatus(StudentActive)
		for _, studentUid := range students {
			if studentUid == uid {
				return class, mine.getStudent(uid)
			}
		}
	}
	return nil, cacheCtx.GetStudent(uid)
}

func (mine *SchoolInfo) GetStudentsByCustodian(phone string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 2)
	if phone == "" {
		return list
	}
	array,err := nosql.GetStudentsByCustodian(mine.UID, phone)
	if err != nil {
		return list
	}
	for _, student := range array {
		info := new(StudentInfo)
		info.initInfo(student)
		list = append(list, info)
	}
	return list
}

func (mine *SchoolInfo) GetStudentsByName(name string) []*StudentInfo {
	list := make([]*StudentInfo, 0, 10)
	if name == "" {
		return list
	}
	all := mine.AllStudents()
	for _, info := range all {
		if info.Name == name {
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
	if len(card) == 19 {
		db.IDCard = card[1:]
		db.SID = card
	}else if len(card) == 18 {
		db.SID = "G"+card
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

func (mine *SchoolInfo) appendStudent(student *StudentInfo, operator string, grade uint8, num,kind uint16) {
	if student == nil {
		return
	}
	date := proxy.DateInfo{
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
	class,info := mine.GetStudent(uid)
	if info == nil {
		return errors.New("not found the student")
	}
	if info.Entity != "" {
		return errors.New("the student had bind")
	}
	info.Remove(operator)
	if class != nil {
		_ = class.RemoveStudent(uid, "the admin delete student", info.ID, StudentDelete)
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

func (mine *SchoolInfo) getStudentByEntity(uid string) *StudentInfo {
	if uid == "" {
		return nil
	}
	all := mine.AllStudents()
	for _, info := range all {
		if info.Entity == uid {
			return info
		}
	}
	return nil
}

func (mine *SchoolInfo) getStudent(uid string) *StudentInfo {
	if uid == "" {
		return nil
	}
	all := mine.AllStudents()
	for _, info := range all {
		if info.UID == uid {
			return info
		}
	}
	return nil
}

func (mine *SchoolInfo) AllStudents() []*StudentInfo {
	students, err := nosql.GetStudentsBySchool(mine.UID)
	if err == nil {
		list := make([]*StudentInfo, 0, len(students))
		for _, db := range students {
			info:=new(StudentInfo)
			info.initInfo(db)
			list = append(list, info)
		}
		return list
	}else{
		return make([]*StudentInfo, 0, 1)
	}
}

func (mine *SchoolInfo) GetStudentsByStatus(st StudentStatus) []*StudentInfo {
	list := make([]*StudentInfo, 0, 100)
	for _, class := range mine.classes {
		array := class.GetStudentsByStatus(st)
		if len(array) > 0 {
			for _, uid := range array {
				if !mine.hadStudentIn(list, uid) {
					info := mine.getStudent(uid)
					if info != nil {
						list = append(list, info)
					}
				}
			}
		}
	}
	return list
}

func (mine *SchoolInfo)hadStudentIn(list []*StudentInfo, uid string) bool {
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

func (mine *SchoolInfo) GetPageStudents(page, number uint32) (uint32, uint32, []*StudentInfo) {
	if number < 1 {
		number = 10
	}
	all := mine.AllStudents()
	if len(all) < 1 {
		return 0, 0, make([]*StudentInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, all)

	return total, maxPage, list.([]*StudentInfo)
}

func (mine *SchoolInfo)GetActiveStudents(page, number uint32) (uint32, uint32, []*StudentInfo) {
	if number < 1 {
		number = 10
	}
	all := mine.GetStudentsByStatus(StudentActive)
	if len(all) < 1 {
		return 0, 0, make([]*StudentInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, all)

	return total, maxPage, list.([]*StudentInfo)
}
//endregion

//region Teacher Fun
func (mine *SchoolInfo)AllTeachers() []*TeacherInfo {
	return mine.teachers
}

func (mine *SchoolInfo)Teachers() []string {
	list := make([]string, 0, len(mine.teachers))
	for _, teacher := range mine.teachers {
		list = append(list, teacher.UID)
	}
	return list
}

func (mine *SchoolInfo) GetTeacherByEntity(entity string) *TeacherInfo {
	if entity == "" {
		return nil
	}
	for _, item := range mine.teachers {
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
	for _, item := range mine.teachers {
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
	for _, item := range mine.teachers {
		if item.UID == uid {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetTeacherByName(name string) *TeacherInfo {
	if name == "" {
		return nil
	}
	for _, item := range mine.teachers {
		if item.Name == name {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetTeachersBySub(subject string) []*TeacherInfo {
	list := make([]*TeacherInfo, 0, 10)
	if subject == "" {
		return list
	}
	for _, item := range mine.teachers {
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
	for _, item := range mine.teachers {
		if item.hadClass(class) {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo) HadTeacher(uid string) bool {
	if uid == "" {
		return false
	}
	for _, item := range mine.teachers {
		if item.UID == uid{
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) HadTeacherByUser(uid string) bool {
	if uid == "" {
		return false
	}
	for _, item := range mine.teachers {
		if item.User == uid && item.IsActive(mine.UID){
			return true
		}
	}
	return false
}

func (mine *SchoolInfo) CreateTeacher(name, entity, user, operator string) (*TeacherInfo, error) {
	if mine.HadTeacherByUser(user) {
		return mine.GetTeacherByUser(user),nil
	}
	teacher, err := Context().createTeacher(name, operator, mine.Scene, entity, user,nil, nil)
	if err != nil {
		return nil, err
	}
	mine.AppendTeacher(teacher)
	return teacher,nil
}

func (mine *SchoolInfo)AppendTeacher(info *TeacherInfo) {
	if mine.HadTeacherByUser(info.UID) {
		return
	}
	err := nosql.AppendSchoolTeacher(mine.UID, info.UID)
	if err == nil {
		mine.teachers = append(mine.teachers, info)
	}
}

func (mine *SchoolInfo)HadTeacherByEntity(entity string) bool {
	for _, teacher := range mine.teachers {
		if teacher.Entity == entity {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo)GetTeachersByPage(page, number uint32) (uint32, uint32, []*TeacherInfo) {
	if number < 1 {
		number = 10
	}
	if len(mine.teachers) < 1 {
		return 0, 0, make([]*TeacherInfo, 0, 1)
	}
	total, maxPage, list := checkPage(page, number, mine.teachers)
	return total, maxPage, list.([]*TeacherInfo)
}

func (mine *SchoolInfo)RemoveTeacher(entity, remark string) error {
	info := mine.GetTeacherByEntity(entity)
	if info == nil {
		return errors.New("not found the teacher")
	}
	_ = info.Remove(mine.UID, remark)
	err :=  nosql.SubtractSchoolTeacher(mine.UID, info.UID)
	if err == nil{
		for i:= 0;i < len(mine.teachers);i += 1 {
			if mine.teachers[i].UID == entity {
				mine.teachers = append(mine.teachers[:i], mine.teachers[i+1:]...)
				break
			}
		}
	}
	return err
}

func (mine *SchoolInfo)RemoveTeacherByUID(uid, remark string) error {
	info := mine.GetTeacher(uid)
	if info == nil {
		return errors.New("not found the teacher")
	}
	_ = info.Remove(mine.UID, remark)
	err :=  nosql.SubtractSchoolTeacher(mine.UID, info.UID)
	if err == nil{
		for i:= 0;i < len(mine.teachers);i += 1 {
			if mine.teachers[i].UID == uid {
				mine.teachers = append(mine.teachers[:i], mine.teachers[i+1:]...)
				break
			}
		}
	}
	return err
}
//endregion

//region Class Fun
func (mine *SchoolInfo)CreateClasses(name, enrol, operator string, number, kind uint16) ([]*ClassInfo, error) {
	if number < 1 {
		return nil, errors.New("the number must not more than 0")
	}
	date := new(proxy.DateInfo)
	err := date.Parse(enrol)
	if err != nil {
		return nil, err
	}
	list := mine.GetClassesByEnrol(date.Year, date.Month)
	array := make([]*ClassInfo, 0, number)
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
	for i := 0;i < length;i += 1 {
		info,_ := mine.createClass(name, enrol, operator, uint16(i + count + 1), kind)
		if info != nil {
			array = append(array, info)
		}
	}
	return array,nil
}

func (mine *SchoolInfo) createClass(name, enrol, operator string, number,kind uint16) (*ClassInfo, error) {
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

func (mine *SchoolInfo)HadClassByEnrol(enrol string) bool {
	for _, item := range mine.classes {
		if item.EnrolDate.String() == enrol {
			return true
		}
	}
	return false
}

func (mine *SchoolInfo)GetClassesByEnrol(year uint16, month time.Month) []*ClassInfo {
	list := make([]*ClassInfo, 0, 10)
	for _, item := range mine.classes {
		if  item.EnrolDate.Year == year && item.EnrolDate.Month == month {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo)GetClassesByGrade(grade uint8) []*ClassInfo {
	list := make([]*ClassInfo, 0, 10)
	for _, item := range mine.classes {
		if  item.Grade() == grade {
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo)GetClasses(status ClassStatus) []*ClassInfo {
	list := make([]*ClassInfo, 0, 50)
	for _, item := range mine.classes {
		if  status == ClassFinish && item.Grade() > mine.MaxGrade() {
			list = append(list, item)
		}else{
			list = append(list, item)
		}
	}
	return list
}

func (mine *SchoolInfo)GetClassesByPage(page, number uint32, st int32) (uint32, uint32, []*ClassInfo) {
	if number < 1 {
		number = 10
	}
	var classes []*ClassInfo
	if st > -1 {
		classes = mine.GetClasses(ClassStatus(st))
	}else{
		classes = mine.classes
	}

	total := uint32(len(classes))
	maxPage := total/ number + 1
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
	return total, max ,list.([]*ClassInfo)
}

func (mine *SchoolInfo)GetClass(uid string) *ClassInfo {
	for _, item := range mine.classes {
		if item.UID == uid {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo)GetClassByStudent(uid string, st StudentStatus) *ClassInfo {
	for _, item := range mine.classes {
		if item.HadStudentByStatus(uid, st) {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo)GetClassByEntity(entity string, st StudentStatus) *ClassInfo {
	student := mine.getStudentByEntity(entity)
	if student == nil {
		return nil
	}
	return mine.GetClassByStudent(student.UID, st)
}

func (mine *SchoolInfo)checkClass(name, operator string, enrol proxy.DateInfo, class, kind uint16) *ClassInfo {
	var info *ClassInfo
	info = mine.GetClassByEnrol(enrol, class)
	if info == nil {
		_,err := mine.CreateClasses(name, enrol.String(), operator, class, kind)
		if err == nil {
			info = mine.GetClassByEnrol(enrol, class)
		}
	}
	return info
}

func (mine *SchoolInfo) GetClassByNO(grade uint8, number uint16) *ClassInfo {
	for _, item := range mine.classes {
		g := item.Grade()
		if g == grade && item.Number == number {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo) GetClassByEnrol(enrol proxy.DateInfo, number uint16) *ClassInfo {
	for _, item := range mine.classes {
		g := item.EnrolDate.Year
		if g == enrol.Year && item.Number == number {
			return item
		}
	}
	return nil
}

func (mine *SchoolInfo)RemoveClass(uid, operator string) error {
	info := mine.GetClass(uid)
	if info == nil {
		return errors.New("not found the class")
	}
	err := info.remove(operator)
	if err == nil{
		for i:= 0;i < len(mine.classes);i += 1 {
			if mine.classes[i].UID == uid {
				mine.classes = append(mine.classes[:i], mine.classes[i+1:]...)
				break
			}
		}
	}
	return err
}
//endregion
