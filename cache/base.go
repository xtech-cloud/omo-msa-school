package cache

import (
	"errors"
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/config"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"omo.msa.school/tool"
	"strings"
	"time"
)

type baseInfo struct {
	ID         uint64 `json:"-"`
	UID        string `json:"uid"`
	Name       string `json:"name"`
	Creator    string
	Operator   string
	CreateTime time.Time
	UpdateTime time.Time
}

type PairIntInfo struct {
	Value uint32
	Key   uint32
}

type cacheContext struct {
	schools  []*SchoolInfo
	teachers []*TeacherInfo
}

var cacheCtx *cacheContext

func InitData() error {
	cacheCtx = &cacheContext{}
	cacheCtx.schools = make([]*SchoolInfo, 0, 100)
	cacheCtx.teachers = make([]*TeacherInfo, 0, 100)

	err := nosql.InitDB(config.Schema.Database.IP, config.Schema.Database.Port, config.Schema.Database.Name, config.Schema.Database.Type)
	if nil != err {
		return err
	}
	//checkSequences()
	//num,_ := nosql.GetSchoolCount()
	schools, _ := nosql.GetUsableSchools()
	for _, school := range schools {
		info := new(SchoolInfo)
		info.initInfo(school)
		cacheCtx.schools = append(cacheCtx.schools, info)
	}
	logger.Infof("init schools!!! number = %d", len(schools))

	return nil
}

func Context() *cacheContext {
	return cacheCtx
}

func checkSequences() {
	arr := make([]string, 0, 6)
	arr = append(arr, "school_"+nosql.TableApply)
	arr = append(arr, "school_"+nosql.TableClass)
	arr = append(arr, "school_"+nosql.TableLesson)
	arr = append(arr, "school_"+nosql.TableStudent)
	arr = append(arr, "school_"+nosql.TableTeacher)
	all, _ := nosql.GetAllSequences()
	for _, s := range all {
		if tool.HasItem(arr, s.Name) {
			k := strings.Replace(s.Name, "school_", "", 1)
			_ = nosql.UpdateSequenceName(s.UID.Hex(), k)
		}
	}

	arr2 := make([]string, 0, 6)
	arr2 = append(arr2, nosql.TableSchool)
	arr2 = append(arr2, nosql.TableApply)
	arr2 = append(arr2, nosql.TableClass)
	arr2 = append(arr2, nosql.TableLesson)
	arr2 = append(arr2, nosql.TableStudent)
	arr2 = append(arr2, nosql.TableTeacher)
	arr2 = append(arr2, nosql.TableSequence)
	all2, _ := nosql.GetAllSequences()
	for _, s := range all2 {
		if !tool.HasItem(arr2, s.Name) {
			_ = nosql.DeleteSequence(s.UID.Hex())
		}
	}
}

func checkPage[T any](page, number uint32, all []T) (uint32, uint32, []T) {
	if len(all) < 1 {
		return 0, 0, make([]T, 0, 1)
	}
	if number < 1 {
		number = 10
	}
	total := uint32(len(all))
	if len(all) <= int(number) {
		return total, 1, all
	}
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}
	if page > maxPage {
		page = maxPage
	}

	var start = (page - 1) * number
	var end = start + number
	if end >= total {
		end = total
	}
	list := make([]T, 0, number)
	list = append(all[start:end])
	return total, maxPage, list
}

func parseDate(date string) (time.Time, error) {
	f, er := time.ParseInLocation("2006-01-02", date, time.UTC)
	if er != nil {
		return time.Now(), er
	}
	if f.IsZero() {
		return f, errors.New("the date is zero")
	}
	return f, nil
}

func calculateGrade(enrol proxy.DateInfo) uint8 {
	now := time.Now()
	diff := now.Year() - int(enrol.Year)
	if now.Month() > time.Month(7) {
		return uint8(diff + 1)
	} else {
		if diff < 1 {
			return 1
		}
		return uint8(diff)
	}
}

func (mine *cacheContext) AllSchools(page, number uint32) (uint32, uint32, []*SchoolInfo) {
	if len(mine.schools) < 1 {
		schools, _ := nosql.GetUsableSchools()
		for _, school := range schools {
			info := new(SchoolInfo)
			info.initInfo(school)
			mine.schools = append(mine.schools, info)
		}
	}
	if number < 1 {
		number = 10
	}
	if len(mine.schools) < 1 {
		return 0, 0, make([]*SchoolInfo, 0, 1)
	}
	total, max, list := checkPage(page, number, mine.schools)
	return total, max, list
}

func (mine *cacheContext) AllTeachers(page, number uint32) (uint32, uint32, []*TeacherInfo) {
	if number < 1 {
		number = 10
	}
	if len(mine.teachers) < 1 {
		return 0, 0, make([]*TeacherInfo, 0, 1)
	}
	total, max, list := checkPage(page, number, mine.teachers)
	return total, max, list
}

func parsePhones(phones string) []string {
	list := make([]string, 0, 2)
	if strings.Contains(phones, ",") {
		array := strings.Split(phones, ",")
		for _, item := range array {
			list = append(list, item)
		}
	} else {
		list = append(list, phones)
	}
	return list
}

func (mine *cacheContext) CreateSchool(name, entity, scene string, maxGrade int) (*SchoolInfo, error) {
	if scene == "" {
		return nil, errors.New("the scene uid is empty")
	}
	if entity == "" {
		return nil, errors.New("the scene entity is empty")
	}
	db := new(nosql.School)
	db.UID = primitive.NewObjectID()
	db.ID = nosql.GetSchoolNextID()
	db.CreatedTime = time.Now()
	db.Scene = scene
	db.Entity = entity
	db.Name = name
	db.Grade = uint8(maxGrade)
	if db.Grade < 6 {
		db.Grade = 6
	}
	db.Teachers = make([]string, 0, 1)
	db.Subjects = make([]proxy.SubjectInfo, 0, 1)
	db.Honors = make([]proxy.HonorInfo, 0, 1)
	db.Respects = make([]proxy.HonorInfo, 0, 1)
	err1 := nosql.CreateSchool(db)
	if err1 == nil {
		school := new(SchoolInfo)
		school.initInfo(db)

		mine.schools = append(mine.schools, school)
		return school, nil
	} else {
		return nil, err1
	}
}

func (mine *cacheContext) GetSchoolBy(uid string) (*SchoolInfo, error) {
	if len(uid) < 1 {
		return nil, errors.New("the school uid is empty")
	}

	info, er := mine.GetSchoolByScene(uid)
	if info == nil {
		return mine.getSchool(uid)
	}

	return info, er
}

func (mine *cacheContext) getSchool(uid string) (*SchoolInfo, error) {
	if len(uid) < 1 {
		return nil, errors.New("the school uid is empty")
	}
	for _, school := range mine.schools {
		if school.UID == uid {
			return school, nil
		}
	}
	db, err := nosql.GetSchool(uid)
	if err != nil {
		return nil, err
	}
	school := new(SchoolInfo)
	school.initInfo(db)
	mine.schools = append(mine.schools, school)
	return school, nil
}

func (mine *cacheContext) GetSchoolByScene(scene string) (*SchoolInfo, error) {
	if scene == "" {
		return nil, errors.New("the scene uid is empty")
	}
	for _, school := range mine.schools {
		if school.Scene == scene {
			return school, nil
		}
	}
	db, err := nosql.GetSchoolByScene(scene)
	if err != nil {
		return nil, err
	}
	school := new(SchoolInfo)
	school.initInfo(db)
	mine.schools = append(mine.schools, school)
	return school, nil
}

func (mine *cacheContext) GetSchoolByName(name string) (*SchoolInfo, error) {
	if name == "" {
		return nil, errors.New("the school uid is empty")
	}
	for i := 0; i < len(mine.schools); i += 1 {
		if mine.schools[i].Name == name {
			return mine.schools[i], nil
		}
	}
	db, _ := nosql.GetSchoolByName(name)
	if db != nil {
		school := new(SchoolInfo)
		school.initInfo(db)
		mine.schools = append(mine.schools, school)
		return school, nil
	} else {
		return nil, errors.New("not found the school by name")
	}
}

func (mine *cacheContext) GetSchoolByEntity(entity string) *SchoolInfo {
	if entity == "" {
		return nil
	}
	for i := 0; i < len(mine.schools); i += 1 {
		if mine.schools[i].Entity == entity {
			return mine.schools[i]
		}
	}
	db, _ := nosql.GetSchoolByEntity(entity)
	if db != nil {
		school := new(SchoolInfo)
		school.initInfo(db)
		mine.schools = append(mine.schools, school)
		return school
	} else {
		return nil
	}
}

func (mine *cacheContext) GetSchoolByClass(class string) *SchoolInfo {
	if class == "" {
		return nil
	}
	for i := 0; i < len(mine.schools); i += 1 {
		mine.schools[i].initClasses()
		if mine.schools[i].hadClass(class) {
			return mine.schools[i]
		}
	}
	db, _ := nosql.GetClass(class)
	if db == nil {
		return nil
	}
	tmp, _ := mine.GetSchoolBy(db.School)
	return tmp
}

func (mine *cacheContext) GetSchoolByStudent(student string, st StudentStatus) *SchoolInfo {
	if student == "" {
		return nil
	}
	for i := 0; i < len(mine.schools); i += 1 {
		if mine.schools[i].hadStudentByStatus(student, st) {
			return mine.schools[i]
		}
	}
	return nil
}

func (mine *cacheContext) GetSchoolByStudent2(student string) (*SchoolInfo, error) {
	if student == "" {
		return nil, errors.New("the student uid is empty")
	}
	db, err := nosql.GetStudent(student)
	if db != nil {
		return mine.GetSchoolBy(db.School)
	}

	return nil, err
}

func (mine *cacheContext) GetSchoolsByStudentEntity(entity string) []*SchoolInfo {
	if entity == "" {
		return nil
	}

	dbs, _ := nosql.GetStudentsByEntity(entity)
	list := make([]*SchoolInfo, 0, len(dbs))
	for _, db := range dbs {
		info, _ := mine.GetSchoolBy(db.School)
		if info != nil {
			list = append(list, info)
		}
	}

	return list
}

func (mine *cacheContext) GetSchoolByTeacher(uid string) *SchoolInfo {
	if uid == "" {
		return nil
	}
	for i := 0; i < len(mine.schools); i += 1 {
		if mine.schools[i].hadTeacher(uid) {
			return mine.schools[i]
		}
	}
	return nil
}

func (mine *cacheContext) GetSchoolByUser(uid string) *SchoolInfo {
	if uid == "" {
		return nil
	}
	for i := 0; i < len(mine.schools); i += 1 {
		if mine.schools[i].hadTeacherByUser(uid) {
			return mine.schools[i]
		}
	}
	return nil
}

func (mine *cacheContext) GetSchoolsByStudent(student string) []*SchoolInfo {
	list := make([]*SchoolInfo, 0, 5)
	if student == "" {
		return list
	}
	for i := 0; i < len(mine.schools); i += 1 {
		if mine.schools[i].hadStudentByStatus(student, StudentAll) {
			list = append(list, mine.schools[i])
		}
	}
	return list
}
