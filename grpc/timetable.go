package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.school/cache"
	"omo.msa.school/proxy"
	"strconv"
	"time"
)

type TimetableService struct {}

func switchTimetable(info *cache.TimetableInfo) *pb.TimetableInfo {
	tmp := new(pb.TimetableInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Name = info.Name
	tmp.School = info.School
	tmp.Class = info.Class
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Year = info.Year
	tmp.Items = make([]*pb.TimetableItem, 0, len(info.Items))
	for _, item := range info.Items {
		tmp.Items = append(tmp.Items, switchTimetableItem(info.Class, &item))
	}
	return tmp
}

func switchTimetableItem(class string, info *proxy.TimetableItem) *pb.TimetableItem {
	tmp := new(pb.TimetableItem)
	tmp.Class = class
	tmp.Name = info.Name
	tmp.Weekday = uint32(info.Weekday)
	tmp.Number = uint32(info.Number)
	return tmp
}

func checkTimetable(arr []*cache.TimetableInfo, school, class string, year uint32) ([]*cache.TimetableInfo, *cache.TimetableInfo) {
	for _, info := range arr {
		if info.Class == class {
			return arr, info
		}
	}
	tmp := new(cache.TimetableInfo)
	tmp.Class = class
	tmp.School = school
	tmp.Year = year
	arr = append(arr, tmp)
	return arr, tmp
}

func (mine *TimetableService)AddOne(ctx context.Context, in *pb.ReqTimetableAdd, out *pb.ReplyTimetableInfo) error {
	path := "timetable.addOne"
	inLog(path, in)
	if in.Year < 2020 {
		out.Status = outError(path,"the year must later than 2020", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	school,_ := cache.Context().GetSchoolByUID(in.School)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	arr := make([]proxy.TimetableItem, 0, 35)
	if in.Items != nil {
		for _, item := range in.Items {
			arr = append(arr, proxy.TimetableItem{Weekday: time.Weekday(item.Weekday), Number: uint8(item.Number), Name: item.Name})
		}
	}
	info, err1 := school.CreateTimetable(in.Class, in.Operator, in.Year, arr)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchTimetable(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TimetableService)AddBatch(ctx context.Context, in *pb.ReqTimetableBatch, out *pb.ReplyTimetableList) error {
	path := "timetable.addBatch"
	inLog(path, in)
	if in.Year < 2020 {
		out.Status = outError(path,"the year must later than 2020", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	school,err := cache.Context().GetSchoolByUID(in.School)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	arr := make([]*cache.TimetableInfo, 0, 50)
	for _, item := range in.Items {
		ar, info := checkTimetable(arr, in.School, item.Class, in.Year)
		arr = ar
		info.Items = append(info.Items, proxy.TimetableItem{Weekday: time.Weekday(item.Weekday), Number: uint8(item.Number), Name: item.Name})
	}
	out.List = make([]*pb.TimetableInfo, 0, len(arr))
	for _, tmp := range arr {
		table,_ := school.GetTimetable(tmp.Class, in.Year)
		if table == nil {
			info, er := school.CreateTimetable(tmp.Class, in.Operator, in.Year, tmp.Items)
			if er == nil {
				out.List = append(out.List, switchTimetable(info))
			}
		}else{
			er := table.UpdateItems(in.Operator, tmp.Items)
			if er == nil {
				out.List = append(out.List, switchTimetable(table))
			}
		}
		school.CreateSubjects(tmp.Items)
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TimetableService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTimetableList) error {
	path := "timetable.getList"
	inLog(path, in)

	school,err := cache.Context().GetSchoolByUID(in.Parent)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	year,err := strconv.ParseUint(in.Value, 10, 32)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_FormatError)
		return nil
	}
	list, err1 := school.GetTimetablesBy(uint32(year))
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.TimetableInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchTimetable(info))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TimetableService)GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTimetableList) error {
	path := "timetable.getByFilter"
	inLog(path, in)
	if in.Number < 10 {
		in.Number = 10
	}
	school,err := cache.Context().GetSchoolByUID(in.Parent)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if in.Filter == "class" {
		year,er := strconv.ParseUint(in.Params, 10, 32)
		if er != nil {
			out.Status = outError(path,er.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
		tmp,er := school.GetTimetable(in.Value, uint32(year))
		if er != nil {
			out.Status = outError(path,er.Error(), pbstatus.ResultStatus_FormatError)
			return nil
		}
		out.List = make([]*pb.TimetableInfo, 0, 1)
		out.List = append(out.List, switchTimetable(tmp))
	}

	out.Status = outLog(path, fmt.Sprintf("the total = %d and length = %d", 1, len(out.List)))
	return nil
}

func (mine *TimetableService)GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "timetable.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *TimetableService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "timetable.removeOne"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}
