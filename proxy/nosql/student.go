package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"omo.msa.school/tool"
	"time"
)

type Student struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name      string         `json:"name" bson:"name"`
	Entity    string         `json:"entity" bson:"entity"`
	EnrolDate proxy.DateInfo `json:"enrol" bson:"enrol"`
	Status    uint8          `json:"status" bson:"status"`
	Number    uint16         `json:"number" bson:"number"`

	Sex uint8 `json:"sex" bson:"sex"`
	//学籍号
	SID string `json:"sid" bson:"sid"`
	// 系统生成的身份序列号
	SN string `json:"sn" bson:"sn"`
	// 身份证号
	IDCard string `json:"card" bson:"card"`
	// 所属学校
	School     string                `json:"school" bson:"school"`
	Tags       []string              `json:"tags" bson:"tags"`
	Custodians []proxy.CustodianInfo `json:"custodians" bson:"custodians"`
}

func (mine *Student) HadCustodian(phone string) bool {
	for _, custodian := range mine.Custodians {
		if tool.HasItem(custodian.Phones, phone) {
			return true
		}
	}
	return false
}

func CreateStudent(info *Student) error {
	_, err := insertOne(TableStudent, info)
	if err != nil {
		return err
	}
	return nil
}

func GetStudentNextID() uint64 {
	num, _ := getSequenceNext(TableStudent)
	return num
}

func GetStudent(uid string) (*Student, error) {
	result, err := findOne(TableStudent, uid)
	if err != nil {
		return nil, err
	}
	model := new(Student)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetStudentByEntity(school, entity string) (*Student, error) {
	msg := bson.M{"school": school, "entity": entity, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableStudent, msg)
	if err != nil {
		return nil, err
	}
	model := new(Student)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetStudentsByEntity(entity string) ([]*Student, error) {
	var items = make([]*Student, 0, 10)
	msg := bson.M{"entity": entity, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentBySN(school, sn string) (*Student, error) {
	msg := bson.M{"school": school, "sn": sn, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableStudent, msg)
	if err != nil {
		return nil, err
	}
	model := new(Student)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetStudentByIDCard(school, card string) (*Student, error) {
	msg := bson.M{"school": school, "card": card, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableStudent, msg)
	if err != nil {
		return nil, err
	}
	model := new(Student)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetStudentsByCard(card string) ([]*Student, error) {
	var items = make([]*Student, 0, 5)
	msg := bson.M{"card": card, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentsBySID(sid string) ([]*Student, error) {
	var items = make([]*Student, 0, 5)
	msg := bson.M{"sid": sid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentsBySchool(uid string) ([]*Student, error) {
	var items = make([]*Student, 0, 100)
	msg := bson.M{"school": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetBindStudentsBySchool(uid string) (uint32, error) {
	msg := bson.M{"school": uid, "entity": bson.M{"$ne": ""}, "deleteAt": new(time.Time)}
	num, err1 := getCountBy(TableStudent, msg)

	return uint32(num), err1
}

func GetAllStudents() ([]*Student, error) {
	var items = make([]*Student, 0, 100)
	cursor, err1 := findAll(TableStudent, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentsByKeyword(school, key string) ([]*Student, error) {
	def := new(time.Time)
	regex := bson.M{"$regex": key}
	//filter := bson.M{"school": school, "deleteAt": def,
	//	"$or": bson.A{bson.M{"name": regex}, bson.M{"sn": regex}, bson.M{"card": regex}, bson.M{"sid": regex}}}
	filter := bson.M{"school": school, "deleteAt": def, "name": regex}
	cursor, err1 := findMany(TableStudent, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	var items = make([]*Student, 0, 20)
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(&node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentsByCustodian(school, phone string) ([]*Student, error) {
	var items = make([]*Student, 0, 10)
	//msg := bson.M{"school":school, "custodians.phone": phone}
	msg := bson.M{"school": school, "deleteAt": new(time.Time), "custodians": bson.M{"$elemMatch": bson.M{"phones": bson.M{"$elemMatch": bson.M{"$eq": phone}}}}}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentsByCustodian2(phone string) ([]*Student, error) {
	var items = make([]*Student, 0, 10)
	//msg := bson.M{"school":school, "custodians.phone": phone}
	msg := bson.M{"deleteAt": new(time.Time), "custodians": bson.M{"$elemMatch": bson.M{"phones": bson.M{"$elemMatch": bson.M{"$eq": phone}}}}}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentsByEnrol(school string, year int) ([]*Student, error) {
	var items = make([]*Student, 0, 10)
	msg := bson.M{"school": school, "deleteAt": new(time.Time), "enrol.year": year}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentsByYear(school string, year int) ([]*Student, error) {
	var items = make([]*Student, 0, 10)
	msg := bson.M{"school": school, "deleteAt": new(time.Time), "enrol.year": bson.M{"$gte": year}}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetStudentCountByStatus(school string, st uint32) uint32 {
	msg := bson.M{"school": school, "status": bson.M{"$eq": st}}
	num, _ := getCountBy(TableStudent, msg)
	return uint32(num)
}

func GetStudentsByStatus(school string, st uint32) ([]*Student, error) {
	var items = make([]*Student, 0, 10)
	msg := bson.M{"school": school, "status": bson.M{"$eq": st}}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllStudentsByStatus(st uint32) ([]*Student, error) {
	var items = make([]*Student, 0, 10)
	msg := bson.M{"status": bson.M{"$eq": st}}
	cursor, err1 := findMany(TableStudent, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Student)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateStudentBase(uid, name, sn, card, sid, operator string, sex uint8, arr []proxy.CustodianInfo) error {
	msg := bson.M{"name": name, "sn": sn, "card": card, "sid": sid, "sex": sex, "custodians": arr,
		"operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentCustodians(uid, operator string, arr []proxy.CustodianInfo) error {
	msg := bson.M{"custodians": arr, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentInfo(uid, name, sn, card, operator string, sex uint8) error {
	msg := bson.M{"name": name, "sn": sn, "card": card, "sex": sex, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentEnrol(uid, operator string, enrol proxy.DateInfo) error {
	msg := bson.M{"enrol": enrol, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentEntity(uid, entity, operator string) error {
	msg := bson.M{"entity": entity, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentState(uid, operator string, st uint8) error {
	msg := bson.M{"status": st, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentNumber(uid, operator string, num uint16) error {
	msg := bson.M{"number": num, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func UpdateStudentTags(uid, operator string, tags []string) error {
	msg := bson.M{"tags": tags, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableStudent, uid, msg)
	return err
}

func RemoveStudent(uid, operator string) error {
	_, err := deleteOne(TableStudent, uid)
	return err
}

func AppendStudentCustodian(uid string, info proxy.CustodianInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"custodians": info}
	_, err := appendElement(TableStudent, uid, msg)
	return err
}

func SubtractStudentCustodian(uid string, name string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"custodians": bson.M{"name": name}}
	_, err := removeElement(TableStudent, uid, msg)
	return err
}

func AppendStudentTag(uid string, tag string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"tags": tag}
	_, err := appendElement(TableStudent, uid, msg)
	return err
}

func SubtractStudentTag(uid string, tag string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"tags": tag}
	_, err := removeElement(TableStudent, uid, msg)
	return err
}
