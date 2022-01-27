package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"time"
)

type School struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Grade   uint8           `json:"grade" bson:"grade"`
	Status  uint8 			`json:"status" bson:"status"`
	Name    string          `json:"name" bson:"name"`
	Cover   string          `json:"cover" bson:"cover"`
	Scene   string          `json:"scene" bson:"scene"`
	Entity  string          `json:"entity" bson:"entity"`

	Teachers []string `json:"teachers" bson:"teachers"`
	Honors []proxy.HonorInfo `json:"honors" bson:"honors"`
	Respects []proxy.HonorInfo `json:"respects" bson:"respects"`
	Subjects []proxy.SubjectInfo `json:"subjects" bson:"subjects"`
}

func CreateSchool(info *School) error {
	_, err := insertOne(TableSchool, info)
	if err != nil {
		return err
	}
	return nil
}

func GetSchoolNextID() uint64 {
	num, _ := getSequenceNext(TableSchool)
	return num
}

func GetSchoolHonorNextID() uint64 {
	num, _ := getSequenceNext("school_honor")
	return num
}

func GetSchoolSubjectNextID() uint64 {
	num, _ := getSequenceNext("school_subject")
	return num
}

func GetSchool(uid string) (*School, error) {
	result, err := findOne(TableSchool, uid)
	if err != nil {
		return nil, err
	}
	model := new(School)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSchoolCount() (int64, error) {
	return getCount(TableSchool)
}

func GetSchoolByScene(scene string) (*School, error) {
	msg := bson.M{"scene": scene}
	result, err := findOneBy(TableSchool, msg)
	if err != nil {
		return nil, err
	}
	model := new(School)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSchoolByName(name string) (*School, error) {
	msg := bson.M{"name": name}
	result, err := findOneBy(TableSchool, msg)
	if err != nil {
		return nil, err
	}
	model := new(School)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSchoolByEntity(entity string) (*School, error) {
	msg := bson.M{"entity": entity}
	result, err := findOneBy(TableSchool, msg)
	if err != nil {
		return nil, err
	}
	model := new(School)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func getOldUsableSchools() ([]*School, error) {
	var items = make([]*School, 0, 20)
	filter := bson.M{"status": bson.M{"$exists": false}}
	cursor, err1 := findMany(TableSchool, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(School)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetUsableSchools() ([]*School, error) {
	var items = make([]*School, 0, 20)
	msg := bson.M{"status": 0}
	cursor, err1 := findMany(TableSchool, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(School)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	array,er := getOldUsableSchools()
	if er == nil {
		for _, school := range array {
			items = append(items, school)
		}
	}
	return items, nil
}

func UpdateSchoolBase(uid, name, remark, operator string) error {
	msg := bson.M{"name": name, "remark": remark, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSchool, uid, msg)
	return err
}

func UpdateSchoolStatus(uid, operator string, status uint8) error {
	msg := bson.M{"status": status, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSchool, uid, msg)
	return err
}

func UpdateSchoolMaster(uid, master, operator string) error {
	msg := bson.M{"master": master, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSchool, uid, msg)
	return err
}

func UpdateSchoolTeachers(uid, operator string, list []string) error {
	msg := bson.M{"teachers": list, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSchool, uid, msg)
	return err
}

func UpdateSchoolCover(uid string, icon, operator string) error {
	msg := bson.M{"cover": icon, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSchool, uid, msg)
	return err
}

func UpdateSchoolLocal(uid string, local, operator string) error {
	msg := bson.M{"location": local, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSchool, uid, msg)
	return err
}

func UpdateSchoolGrade(uid string, grade uint8, operator string) error {
	msg := bson.M{"grade": grade, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableSchool, uid, msg)
	return err
}

func RemoveSchool(uid, operator string) error {
	_, err := removeOne(TableSchool, uid, operator)
	return err
}

func AppendSchoolTeacher(uid string, teacher string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"teachers": teacher}
	_, err := appendElement(TableSchool, uid, msg)
	return err
}

func SubtractSchoolTeacher(uid string, teacher string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"teachers": teacher}
	_, err := removeElement(TableSchool, uid, msg)
	return err
}

func AppendSchoolHonor(uid string, honor proxy.HonorInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"honors": honor}
	_, err := appendElement(TableSchool, uid, msg)
	return err
}

func SubtractSchoolHonor(uid string, honor string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"honors": bson.M{"uid": honor}}
	_, err := removeElement(TableSchool, uid, msg)
	return err
}

func AppendSchoolRespect(uid string, honor proxy.HonorInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"respects": honor}
	_, err := appendElement(TableSchool, uid, msg)
	return err
}

func SubtractSchoolRespect(uid string, honor string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"respects": bson.M{"uid": honor}}
	_, err := removeElement(TableSchool, uid, msg)
	return err
}

func AppendSchoolSubject(uid string, info proxy.SubjectInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"subjects": info}
	_, err := appendElement(TableSchool, uid, msg)
	return err
}

func SubtractSchoolSubject(uid string, subject string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"subjects": bson.M{"uid": subject}}
	_, err := removeElement(TableSchool, uid, msg)
	return err
}
