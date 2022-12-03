package nosql

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Schedule struct {
	UID         primitive.ObjectID `bson:"_id"`
	ID          uint64             `json:"id" bson:"id"`
	CreatedTime time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedTime time.Time          `json:"updatedAt" bson:"updatedAt"`
	DeleteTime  time.Time          `json:"deleteAt" bson:"deleteAt"`
	Creator     string             `json:"creator" bson:"creator"`
	Operator    string             `json:"operator" bson:"operator"`

	Status    uint8     `json:"status" bson:"status"`
	LimitMax  uint32    `json:"max" bson:"max"`
	LimitMin  uint32    `json:"min" bson:"min"`
	StartTime uint64    `json:"startTime" bson:"startTime"` //报名开始时间
	EndTime   uint64    `json:"endTime" bson:"endTime"`     //报名截止时间
	Date      time.Time `json:"date" bson:"date"`           //日期
	Name      string    `json:"name" bson:"name"`
	Scene     string    `json:"scene" bson:"scene"`
	Lesson    string    `json:"lesson" bson:"lesson"` //课程
	Place     string    `json:"place" bson:"place"`   //地址
	During    string    `json:"during" bson:"during"`

	Teachers []string `json:"teachers" bson:"teachers"`
	Tags     []string `json:"tags" bson:"tags"`
	Users    []string `json:"users" bson:"users"`
}

func CreateSchedule(info *Schedule) error {
	_, err := insertOne(TableSchedules, info)
	if err != nil {
		return err
	}
	return nil
}

func GetScheduleNextID() uint64 {
	num, _ := getSequenceNext(TableSchedules)
	return num
}

func GetSchedule(uid string) (*Schedule, error) {
	result, err := findOne(TableSchedules, uid)
	if err != nil {
		return nil, err
	}
	model := new(Schedule)
	err1 := result.Decode(model)
	if err1 != nil {
		return nil, err1
	}
	return model, nil
}

func GetSchedulesByCreator(uid string) ([]*Schedule, error) {
	var items = make([]*Schedule, 0, 100)
	msg := bson.M{"creator": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableSchedules, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Schedule)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSchedulesByScene(uid string) ([]*Schedule, error) {
	var items = make([]*Schedule, 0, 100)
	msg := bson.M{"scene": uid, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableSchedules, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Schedule)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSchedulesByDate(scene string, date time.Time) ([]*Schedule, error) {
	var items = make([]*Schedule, 0, 100)
	msg := bson.M{"scene": scene, "date": date, "deleteAt": new(time.Time)}
	cursor, err1 := findMany(TableSchedules, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Schedule)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func GetSchedulesByDuring(scene string, from, to time.Time) ([]*Schedule, error) {
	var items = make([]*Schedule, 0, 100)
	msg := bson.M{"scene": scene, "deleteAt": new(time.Time), "$and": bson.A{bson.M{"date": bson.M{"$gte": from}}, bson.M{"date": bson.M{"$lte": to}}}}
	cursor, err1 := findMany(TableSchedules, msg, 0)
	if err1 != nil {
		return nil, err1
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var node = new(Schedule)
		if err := cursor.Decode(node); err != nil {
			return nil, err
		} else {
			items = append(items, node)
		}
	}
	return items, nil
}

func UpdateScheduleBase(uid, lesson, place, times, operator string, max, min uint32, teachers []string) error {
	msg := bson.M{"lesson": lesson, "place": place, "during": times, "operator": operator,
		"teachers": teachers, "max": max, "min": min, "updatedAt": time.Now()}
	_, err := updateOne(TableSchedules, uid, msg)
	return err
}

func UpdateScheduleStatus(uid, operator string, st uint8, start, end uint64) error {
	msg := bson.M{"operator": operator, "status": st, "startTime": start, "endTime": end, "updatedAt": time.Now()}
	_, err := updateOne(TableSchedules, uid, msg)
	return err
}

func UpdateScheduleTags(uid, operator string, tags []string) error {
	msg := bson.M{"operator": operator, "tags": tags, "updatedAt": time.Now()}
	_, err := updateOne(TableSchedules, uid, msg)
	return err
}

func UpdateScheduleUsers(uid, operator string, users []string) error {
	msg := bson.M{"operator": operator, "users": users, "updatedAt": time.Now()}
	_, err := updateOne(TableSchedules, uid, msg)
	return err
}

func AppendScheduleUser(uid, user string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"users": user}
	_, err := appendElement(TableSchedules, uid, msg)
	return err
}

func SubtractScheduleUser(uid, user string) error {
	if len(uid) < 1 {
		return errors.New("the uid is empty")
	}
	msg := bson.M{"users": user}
	_, err := removeElement(TableSchedules, uid, msg)
	return err
}

func RemoveSchedule(uid, operator string) error {
	_, err := removeOne(TableSchedules, uid, operator)
	return err
}
