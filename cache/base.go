package cache

import (
	"errors"
	"github.com/micro/go-micro/v2/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/config"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"reflect"
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
	Key uint32
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

	//num,_ := nosql.GetSchoolCount()
	schools,_ := nosql.GetUsableSchools()
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

func checkPage( page, number uint32, all interface{}) (uint32, uint32, interface{}) {
	if number < 1 {
		number = 10
	}
	array := reflect.ValueOf(all)
	total := uint32(array.Len())
	maxPage := total/number + 1
	if page < 1 {
		return total, maxPage, all
	}

	var start = (page - 1) * number
	var end = start + number
	if end > total {
		end = total
	}

	list := array.Slice(int(start), int(end))
	return total, maxPage, list.Interface()
}

func (mine *cacheContext)AllSchools(page, number uint32) []*SchoolInfo {
	if len(mine.schools) < 1 {
		schools,_ := nosql.GetUsableSchools()
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
		return make([]*SchoolInfo, 0, 1)
	}
	_, _, list := checkPage(page, number, mine.schools)
	return list.([]*SchoolInfo)
}

func (mine *cacheContext)createSchoolInfo(name, entity, scene string, maxGrade int) (*SchoolInfo,error) {
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
		return school,nil
	}else{
		return nil,err1
	}
}

func (mine *cacheContext) GetSchoolByUID(uid string) (*SchoolInfo,error) {
	if len(uid) < 1 {
		return nil,errors.New("the school uid is empty")
	}
	for i := 0;i < len(mine.schools);i += 1 {
		if mine.schools[i].UID == uid {
			return mine.schools[i], nil
		}
	}
	db,err := nosql.GetSchool(uid)
	if err != nil {
		return nil,err
	}
	school := new(SchoolInfo)
	school.initInfo(db)
	mine.schools = append(mine.schools, school)
	return school,nil
}

func (mine *cacheContext) GetSchool(scene string) *SchoolInfo {
	if scene == "" {
		return nil
	}
	for i := 0;i < len(mine.schools);i += 1 {
		if mine.schools[i].Scene == scene {
			return mine.schools[i]
		}
	}
	return nil
}

func (mine *cacheContext) GetSchoolByName(name string) *SchoolInfo {
	if name == "" {
		return nil
	}
	for i := 0;i < len(mine.schools);i += 1 {
		if mine.schools[i].Name == name {
			return mine.schools[i]
		}
	}
	db,_ := nosql.GetSchoolByName(name)
	if db != nil {
		school := new(SchoolInfo)
		school.initInfo(db)
		mine.schools = append(mine.schools, school)
		return school
	}else{
		return nil
	}
}

func (mine *cacheContext) GetSchoolByEntity(entity string) *SchoolInfo {
	if entity == "" {
		return nil
	}
	for i := 0;i < len(mine.schools);i += 1 {
		if mine.schools[i].entity == entity {
			return mine.schools[i]
		}
	}
	db,_ := nosql.GetSchoolByEntity(entity)
	if db != nil {
		school := new(SchoolInfo)
		school.initInfo(db)
		mine.schools = append(mine.schools, school)
		return school
	}else{
		return nil
	}
}

func (mine *cacheContext) GetSchoolByStudent(student string, st StudentStatus) *SchoolInfo {
	if student == "" {
		return nil
	}
	for i := 0;i < len(mine.schools);i += 1 {
		if mine.schools[i].hadStudentByStatus(student, st) {
			return mine.schools[i]
		}
	}
	return nil
}

func (mine *cacheContext) GetSchoolByStudent2(student string) *SchoolInfo {
	if student == "" {
		return nil
	}
	for i := 0;i < len(mine.schools);i += 1 {
		if mine.schools[i].hadStudent(student) {
			return mine.schools[i]
		}
	}
	return nil
}

func (mine *cacheContext) GetSchoolsByStudentEntity(entity string) ([]*SchoolInfo,[]*ClassInfo) {
	list := make([]*SchoolInfo, 0, 5)
	classes := make([]*ClassInfo, 0, 10)
	if entity == "" {
		return list, classes
	}
	for i := 0;i < len(mine.schools);i += 1 {
		stu := mine.schools[i].getStudentByEntity(entity)
		if stu != nil {
			list = append(list, mine.schools[i])
			class := mine.schools[i].GetClassByStudent(stu.UID, StudentActive)
			if class != nil {
				classes = append(classes, class)
			}
		}
	}
	return list, classes
}

func (mine *cacheContext) GetSchoolByTeacher(user string) *SchoolInfo {
	if user == "" {
		return nil
	}
	for i := 0;i < len(mine.schools);i += 1 {
		if mine.schools[i].HadTeacherByUser(user) {
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
	for i := 0;i < len(mine.schools);i += 1 {
		if mine.schools[i].hadStudentByStatus(student, StudentAll) {
			list = append(list, mine.schools[i])
		}
	}
	return list
}
