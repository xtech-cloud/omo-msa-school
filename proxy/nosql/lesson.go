package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Lesson struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Weight uint32   `json:"weight" bson:"weight"`
	Name   string   `json:"name" bson:"name"`
	Remark string   `json:"remark" bson:"remark"`
	Graph  string   `json:"graph" bson:"graph"`
	Scene  string   `json:"scene" bson:"scene"`
	Cover  string   `json:"cover" bson:"cover"`
	Tags   []string `json:"tags" bson:"tags"`
	Assets []string `json:"assets" bson:"assets"`
}

func CreateLesson(info *Lesson) error {
	_, err := insertOne(TableLesson, info)
	if err != nil {
		return err
	}
	return nil
}

func GetLessonNextID() uint64 {
	num, _ := getSequenceNext(TableLesson)
	return num
}

func GetLesson(uid string) (*Lesson, error) {
	result, err := findOne(TableLesson, uid)
	if err != nil {
		return nil, err
	}
	model := new(Lesson)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetLessonsByCreator(uid string) ([]*Lesson, error) {
	var items = make([]*Lesson, 0, 100)
	msg := bson.M{"creator": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableLesson, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Lesson)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetLessonsByScene(uid string) ([]*Lesson, error) {
	var items = make([]*Lesson, 0, 100)
	msg := bson.M{"scene": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableLesson, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Lesson)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetAllLessons() ([]*Lesson, error) {
	var items = make([]*Lesson, 0, 100)
	cursor, err1 := findAll(TableLesson, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Lesson)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateLessonBase(uid, name, remark, operator string, tags []string) error {
	msg := bson.M{"name": name, "operator": operator, "remark": remark, "tags": tags, "updatedAt": time.Now()}
	_, err := updateOne(TableLesson, uid, msg)
	return err
}

func UpdateLessonAssets(uid, operator string, arr []string) error {
	msg := bson.M{"assets": arr, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableLesson, uid, msg)
	return err
}

func UpdateLessonCover(uid, operator, cover string) error {
	msg := bson.M{"cover": cover, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableLesson, uid, msg)
	return err
}

func UpdateLessonWeight(uid, operator string, weight uint32) error {
	msg := bson.M{"weight": weight, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableLesson, uid, msg)
	return err
}

func UpdateLessonGraph(uid, operator, graph string) error {
	msg := bson.M{"graph": graph, "operator": operator, "updatedAt": time.Now()}
	_, err := updateOne(TableLesson, uid, msg)
	return err
}

func RemoveLesson(uid, operator string) error {
	_, err := removeOne(TableLesson, uid, operator)
	return err
}
