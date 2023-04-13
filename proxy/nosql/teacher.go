package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"time"
)

type Teacher struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name      string              `json:"name" bson:"name"`
	Remark    string              `json:"remark" bson:"remark"`
	Entity    string              `json:"entity" bson:"entity"`
	User      string              `json:"user" bson:"user"`
	Classes   []string            `json:"classes" bson:"classes"`
	Subjects  []string            `json:"subjects" bson:"subjects"`
	Tags      []string            `json:"tags" bson:"tags"`
	Histories []proxy.HistoryInfo `json:"histories" bson:"histories"`
}

func CreateTeacher(info *Teacher) error {
	_, err := insertOne(TableTeacher, info)
	if err != nil {
		return err
	}
	return nil
}

func GetTeacherNextID() uint64 {
	num, _ := getSequenceNext(TableTeacher)
	return num
}

func GetTeacher(uid string) (*Teacher, error) {
	result, err := findOne(TableTeacher, uid)
	if err != nil {
		return nil, err
	}
	model := new(Teacher)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTeacherByEntity(entity string) (*Teacher, error) {
	msg := bson.M{"entity": entity}
	result, err := findOneBy(TableTeacher, msg)
	if err != nil {
		return nil, err
	}
	model := new(Teacher)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTeacherByUser(user string) (*Teacher, error) {
	msg := bson.M{"user": user}
	result, err := findOneBy(TableTeacher, msg)
	if err != nil {
		return nil, err
	}
	model := new(Teacher)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTeachersBySchool(uid string) ([]*Teacher, error) {
	var items = make([]*Teacher, 0, 100)
	msg := bson.M{"school": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTeacher, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Teacher)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetLeaveTeachers(school string) ([]*Teacher, error) {
	var items = make([]*Teacher, 0, 100)
	filter := bson.M{"histories.school": school}
	cursor, err1 := findMany(TableTeacher, filter, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Teacher)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllTeachers() ([]*Teacher, error) {
	var items = make([]*Teacher, 0, 100)
	cursor, err1 := findAll(TableTeacher, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Teacher)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateTeacherBase(uid, name, operator string, classes, subs []string) error {
	msg := bson.M{"name": name, "operator": operator, "classes": classes, "subjects": subs, "updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}

func UpdateTeacherHistories(uid, operator string, list []proxy.HistoryInfo) error {
	msg := bson.M{"operator": operator, "histories": list, "updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}

func UpdateTeacherBase2(uid, name, operator string) error {
	msg := bson.M{"name": name, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}

func UpdateTeacherSubjects(uid, operator string, array []string) error {
	msg := bson.M{"subjects": array, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableTeacher, uid, msg)
	return err
}

func RemoveTeacher(uid, operator string) error {
	_, err := removeOne(TableTeacher, uid, operator)
	return err
}

func AppendTeacherHistory(uid string, info *proxy.HistoryInfo) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"histories": info}
	_, err := appendElement(TableTeacher, uid, msg)
	return err
}

func AppendTeacherTag(uid string, tag string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"tags": tag}
	_, err := appendElement(TableTeacher, uid, msg)
	return err
}

func SubtractTeacherTag(uid string, tag string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"tags": tag}
	_, err := removeElement(TableTeacher, uid, msg)
	return err
}
