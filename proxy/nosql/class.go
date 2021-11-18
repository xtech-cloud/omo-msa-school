package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"time"
)

type Class struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name      string              `json:"name" bson:"name"`
	School    string              `json:"school" bson:"school"`
	Master    string 			  `json:"master" bson:"master"`
	EnrolDate proxy.DateInfo      `json:"enrol" bson:"enrol"`
	Type      uint8 			  `json:"type" bson:"type"`
	Number    uint16              `json:"number" bson:"number"`
	Teachers  []string 			  `json:"teachers" bson:"teachers"`
	Students  []proxy.ClassMember `json:"students" bson:"students"`
}

func CreateClass(info *Class) error {
	_, err := insertOne(TableClass, info)
	if err != nil {
		return err
	}
	return nil
}

func GetClassNextID() uint64 {
	num, _ := getSequenceNext(TableClass)
	return num
}

func GetClass(uid string) (*Class, error) {
	result, err := findOne(TableClass, uid)
	if err != nil {
		return nil, err
	}
	model := new(Class)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetClassesBySchool(uid string) ([]*Class, error) {
	var items = make([]*Class, 0, 100)
	msg := bson.M{"school": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableClass, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Class)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllClasses() ([]*Class, error) {
	var items = make([]*Class, 0, 100)
	cursor, err1 := findAll(TableClass, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Class)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateClassBase(uid, name, operator string) error {
	msg := bson.M{"name": name, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableClass, uid, msg)
	return err
}

func UpdateClassMaster(uid, master, operator string) error {
	msg := bson.M{"master": master, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableClass, uid, msg)
	return err
}

func RemoveClass(uid, operator string) error {
	_, err := removeOne(TableClass, uid, operator)
	return err
}

func AppendClassStudent(uid string, info proxy.ClassMember) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"students": info}
	_, err := appendElement(TableClass, uid, msg)
	return err
}

func SubtractClassStudent(uid string, uuid string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"students": bson.M{"uid": uuid}}
	_, err := removeElement(TableClass, uid, msg)
	return err
}

func AppendClassTeacher(uid, teacher string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"students": teacher}
	_, err := appendElement(TableClass, uid, msg)
	return err
}

func SubtractClassTeacher(uid, teacher string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"students": teacher}
	_, err := removeElement(TableClass, uid, msg)
	return err
}
