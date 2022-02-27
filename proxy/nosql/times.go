package nosql

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"omo.msa.school/proxy"
	"time"
)

type Timetable struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Name   string                `json:"name" bson:"name"`
	Year   uint32                `json:"year" bson:"year"`
	School string                `json:"school" bson:"school"`
	Class  string                `json:"class" bson:"class"`
	Items  []proxy.TimetableItem `json:"items" bson:"items"`
}

func CreateTimetable(info *Timetable) error {
	_, err := insertOne(TableTimes, info)
	if err != nil {
		return err
	}
	return nil
}

func GetTimetableNextID() uint64 {
	num, _ := getSequenceNext(TableTimes)
	return num
}

func GetTimetableByUID(uid string) (*Timetable, error) {
	result, err := findOne(TableTimes, uid)
	if err != nil {
		return nil, err
	}
	model := new(Timetable)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTimetable(school, class string, year uint32) (*Timetable, error) {
	msg := bson.M{"school": school, "class": class, "year": year, "deleteAt": new(time.Time)}
	result, err := findOneBy(TableTimes, msg)
	if err != nil {
		return nil, err
	}
	model := new(Timetable)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetTimetablesBy(school string, year uint32) ([]*Timetable, error) {
	var items = make([]*Timetable, 0, 10)
	msg := bson.M{"school": school, "year": year, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableTimes, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Timetable)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateTimetableItems(uid, operator string, list []proxy.TimetableItem) error {
	msg := bson.M{"operator": operator, "items": list, "updatedAt": time.Now()}
	_, err := updateOne(TableTimes, uid, msg)
	return err
}

func RemoveTimetable(uid, operator string) error {
	_, err := deleteOne(TableTimes, uid)
	return err
}
