package grpc

import (
	"context"
	"errors"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.school/cache"
	"strconv"
)

type ScheduleService struct{}

func switchSchedule(info *cache.ScheduleInfo) *pb.ScheduleInfo {
	tmp := new(pb.ScheduleInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Name = info.Name
	tmp.Lesson = info.Lesson
	tmp.Place = info.Place
	tmp.Date = info.Date
	tmp.During = info.Times
	tmp.Max = info.LimitMax
	tmp.Min = info.LimitMin
	tmp.Status = uint32(info.Status)
	tmp.Scene = info.Scene
	tmp.Start = uint64(info.StartTime)
	tmp.End = uint64(info.EndTime)

	tmp.Tags = info.Tags
	tmp.Teachers = info.Teachers
	tmp.Users = info.Users
	return tmp
}

func (mine *ScheduleService) AddOne(ctx context.Context, in *pb.ReqScheduleAdd, out *pb.ReplyScheduleInfo) error {
	path := "schedule.addOne"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolScene(in.Scene)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	info, err1 := school.CreateSchedule(in.Lesson, in.Place, in.Date, in.During, in.Operator, in.Min, in.Max, in.Teachers)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchSchedule(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyScheduleInfo) error {
	path := "schedule.getOne"
	inLog(path, in)
	var info *cache.ScheduleInfo
	var err error
	if len(in.Parent) > 1 {
		scene, er := cache.Context().GetSchoolScene(in.Parent)
		if er != nil {
			out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
			return nil
		}
		info, err = scene.GetSchedule(in.Uid)

	} else {

	}

	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchSchedule(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyScheduleList) error {
	path := "schedule.getByFilter"
	inLog(path, in)
	if in.Number < 10 {
		in.Number = 10
	}
	var err error
	var list []*cache.ScheduleInfo
	scene, er := cache.Context().GetSchoolScene(in.Parent)
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if in.Filter == "" {
		list, err = scene.GetSchedules()
	} else if in.Filter == "dates" {
		if len(in.List) < 2 {
			err = errors.New("the filter of dates length < 2")
		}
		from := in.List[0]
		to := in.List[1]
		list, err = scene.GetSchedulesByDates(from, to)
	} else if in.Filter == "date" {
		list, err = scene.GetSchedulesByDate(in.Value)
	} else {
		err = errors.New("the filter not defined")
	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = make([]*pb.ScheduleInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchSchedule(info))
	}
	out.Pages = 0
	out.Page = in.Page
	out.Total = 1
	out.Status = outLog(path, fmt.Sprintf("the total = %d and length = %d", 1, len(out.List)))
	return nil
}

func (mine *ScheduleService) GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "schedule.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) UpdateOne(ctx context.Context, in *pb.ReqScheduleUpdate, out *pb.ReplyScheduleInfo) error {
	path := "schedule.updateOne"
	inLog(path, in)
	info, err := cache.Context().GetSchedule(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err1 := info.UpdateInfo(in.Lesson, in.Place, in.During, in.Operator, in.Max, in.Min, in.Teachers)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchSchedule(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) SetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyScheduleInfo) error {
	path := "schedule.setByFilter"
	inLog(path, in)
	info, err := cache.Context().GetSchedule(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var er error
	if in.Filter == "status" {
		st, er1 := strconv.ParseInt(in.Value, 10, 32)
		if er1 == nil {
			if len(in.List) < 2 {
				er = info.UpdateStatus2(in.Operator, uint8(st))
			} else {
				start, _ := strconv.ParseInt(in.List[0], 10, 64)
				end, _ := strconv.ParseInt(in.List[1], 10, 64)
				er = info.UpdateStatus(in.Operator, start, end, uint8(st))
			}
		}
	} else if in.Filter == "tags" {
		er = info.UpdateTags(in.Operator, in.List)
	} else {
		er = errors.New("the filter not defined")
	}
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
		return nil
	}
	out.Info = switchSchedule(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "schedule.removeOne"
	inLog(path, in)
	info, err := cache.Context().GetSchedule(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.Remove(in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) AppendUser(ctx context.Context, in *pb.ReqScheduleUser, out *pb.ReplyScheduleUsers) error {
	path := "schedule.appendUsers"
	inLog(path, in)
	info, err := cache.Context().GetSchedule(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.AppendUsers(in.Operator, in.Users)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ScheduleService) SubtractUser(ctx context.Context, in *pb.ReqScheduleUser, out *pb.ReplyScheduleUsers) error {
	path := "schedule.subtractUsers"
	inLog(path, in)
	info, err := cache.Context().GetSchedule(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err = info.SubtractUser(in.Operator, in.Users)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}
