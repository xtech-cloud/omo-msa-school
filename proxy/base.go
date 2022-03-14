package proxy

import "time"

type HonorInfo struct {
	UID    string `json:"uid" bson:"uid"`
	Name   string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
	Parent string `json:"parent" bson:"parent"`
}

type SubjectInfo struct {
	UID    string `json:"uid" bson:"uid"`
	Name   string `json:"name" bson:"name"`
	Remark string `json:"remark" bson:"remark"`
}

// 监护人信息
type CustodianInfo struct {
	Name     string   `json:"name" bson:"name"`
	Phones   []string `json:"phones" bson:"phones"`
	Identity string   `json:"identity" bson:"identity"`
}

type HistoryInfo struct {
	UID     string `json:"uid" bson:"uid"`
	School  string `json:"school" bson:"school"`
	Grade   uint8  `json:"grade" bson:"grade"`
	Class   uint16 `json:"class" bson:"class"`
	Remark  string `json:"remark" bson:"remark"`
	Enrol   string `json:"enrol" bson:"enrol"`
	Created uint64 `json:"created" bson:"created"`
}

type ClassMember struct {
	UID     string `bson:"uid"`
	Student string `bson:"student"`
	Status  uint8  `bson:"status"`
	Remark  string `bson:"remark"`
	Updated time.Time `bson:"updated"`
}

type DeviceInfo struct {
	Type   uint8  `json:"type" bson:"type"`
	UID    string `json:"uid" bson:"uid"`
	Remark string `json:"remark" bson:"remark"`
}

type TimetableItem struct {
	Weekday   time.Weekday `json:"weekday" bson:"weekday"`
	Number uint8        `json:"number" bson:"number"`
	Name   string       `json:"name" bson:"name"`
}
