package cache

import (
	"errors"
	"omo.msa.school/proxy"
	"omo.msa.school/proxy/nosql"
	"omo.msa.school/tool"
)

const (
	StudentDelete StudentStatus = 0 //管理员删除
	StudentActive StudentStatus = 1 // 在读
	StudentFinish StudentStatus = 2 // 毕业
	StudentLeave  StudentStatus = 3 // 中途离开，转校
	StudentAll StudentStatus = 99 //全部记录
)

type StudentStatus uint8

type StudentInfo struct {
	Sex uint8
	Status uint8
	baseInfo
	Entity     string
	SN         string //学号
	IDCard     string //身份证
	SID        string //学籍号
	School     string
	EnrolDate  proxy.DateInfo
	Tags       []string
	Custodians []proxy.CustodianInfo
}

func (mine *StudentInfo) initInfo(db *nosql.Student) {
	mine.UID = db.UID.Hex()
	mine.ID = db.ID
	mine.UpdateTime = db.UpdatedTime
	mine.CreateTime = db.CreatedTime
	mine.Creator = db.Creator
	mine.Operator = db.Operator
	mine.Name = db.Name
	mine.Entity = db.Entity
	mine.Tags = db.Tags
	mine.SN = db.SN
	mine.Sex = db.Sex
	mine.SID = db.SID
	mine.IDCard = db.IDCard
	mine.EnrolDate = db.EnrolDate
	mine.School = db.School
	mine.Status = db.Status
	mine.Custodians = db.Custodians
	if mine.Custodians == nil {
		mine.Custodians = make([]proxy.CustodianInfo, 0, 1)
	}
}

func (mine *StudentInfo) Birthday() string {
	if len(mine.SID) == 19 {
		return mine.SID[7:15]
	}
	if len(mine.IDCard) == 18 {
		return mine.SID[6:14]
	}
	return ""
}

func (mine *StudentInfo) UpdateCustodian(name, phones, identify string) error {
	if len(phones) < 1 {
		return errors.New("the custodian phone is empty")
	}
	if len(name) < 2 {
		name = "default"
	}
	list := parsePhones(phones)
	if mine.hadCustodian(name) {
		_ = nosql.SubtractStudentCustodian(mine.UID, name)
	}
	info := proxy.CustodianInfo{Name: name, Phones: list, Identity: identify}
	err := nosql.AppendStudentCustodian(mine.UID, info)
	if err == nil {
		mine.Custodians = append(mine.Custodians, info)
	}
	return err
}

func (mine *StudentInfo) HadCustodian(phone string) bool {
	if phone == ""{
		return false
	}
	for _, custodian := range mine.Custodians {
		if tool.HasItem(custodian.Phones, phone) {
			return true
		}
	}
	return false
}

func (mine *StudentInfo) hadCustodian(name string) bool {
	if name == ""{
		return false
	}
	for _, custodian := range mine.Custodians {
		if custodian.Name == name {
			return true
		}
	}
	return false
}

func (mine *StudentInfo) UpdateBase(name, sn, card, operator string, sex uint8, arr []proxy.CustodianInfo) error {
	var err error
	var sid = mine.SID
	if card == "" {
		card = mine.IDCard
	}
	if card != mine.IDCard {
		sid = "G"+card
	}
	if name == "" {
		name = mine.Name
	}
	err = nosql.UpdateStudentBase(mine.UID, name, sn, card, sid, operator, sex, arr)
	if err == nil {
		mine.Name = name
		mine.Custodians = arr
		mine.SN = sn
		mine.IDCard = card
		mine.SID = sid
		mine.Sex = sex
		mine.Operator = operator
	}
	return err
}

func (mine *StudentInfo) UpdateSelf(name, sn,card, operator string, sex uint8) error {
	var err error
	err = nosql.UpdateStudentInfo(mine.UID, name, sn, card, operator, sex)
	if err == nil {
		mine.Name = name
		mine.IDCard = card
		mine.Sex = sex
		mine.SN = sn
		mine.Operator = operator
	}
	return err
}

func (mine *StudentInfo) UpdateEnrol(enrol proxy.DateInfo, operator string) error {
	err := nosql.UpdateStudentEnrol(mine.UID, operator, enrol)
	if err == nil {
		mine.EnrolDate = enrol
		mine.Operator = operator
	}
	return err
}

func (mine *StudentInfo) UpdateTags(tags []string, operator string) error {
	err := nosql.UpdateStudentTags(mine.UID, operator, tags)
	if err == nil {
		mine.Tags = tags
		mine.Operator = operator
	}
	return err
}

func (mine *StudentInfo) UpdateStatus(st uint8, operator string) error {
	err := nosql.UpdateStudentState(mine.UID, operator, st)
	if err == nil {
		mine.Status = st
		mine.Operator = operator
	}
	return err
}

func (mine *StudentInfo) BindEntity(entity, operator string) error {
	err := nosql.UpdateStudentEntity(mine.UID, entity, operator)
	if err == nil {
		mine.Entity = entity
		mine.Operator = operator
	}
	return err
}

func (mine *StudentInfo) Remove(operator string) bool {
	if mine.Entity != "" {
		return false
	}
	er := nosql.RemoveStudent(mine.UID, operator)
	if er == nil {
		return true
	}
	return false
}

func (mine *StudentInfo)hadTag(tag string) bool {
	if tag == ""{
		return false
	}
	for _, s := range mine.Tags {
		if s == tag {
			return true
		}
	}
	return false
}

func (mine *StudentInfo) appendTag(tag string) error {
	if tag == ""{
		return errors.New("the tag is empty")
	}
	if mine.hadTag(tag) {
		return errors.New("the tag had existed")
	}
	err := nosql.AppendStudentTag(mine.UID, tag)
	if err == nil {
		mine.Tags = append(mine.Tags, tag)
	}
	return err
}

func (mine *StudentInfo) subtractTag(tag string) error {
	if tag == ""{
		return errors.New("the tag is empty")
	}
	if !mine.hadTag(tag) {
		return errors.New("the tag not existed")
	}
	err := nosql.SubtractStudentTag(mine.UID, tag)
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