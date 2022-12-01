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

type LessonService struct{}

func switchLesson(info *cache.LessonInfo) *pb.LessonInfo {
	tmp := new(pb.LessonInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Scene = info.Scene
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Tags = info.Tags
	tmp.Assets = info.Assets
	tmp.Cover = info.Cover
	return tmp
}

func (mine *LessonService) AddOne(ctx context.Context, in *pb.ReqLessonAdd, out *pb.ReplyLessonInfo) error {
	path := "lesson.addOne"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolScene(in.Scene)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info, err1 := school.CreateLesson(in.Name, in.Remark, in.Cover, in.Operator, in.Tags)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchLesson(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *LessonService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyLessonInfo) error {
	path := "lesson.getOne"
	inLog(path, in)
	var info *cache.LessonInfo
	var er error
	if len(in.Parent) > 1 {
		scene, err := cache.Context().GetSchoolByUID(in.Parent)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Filter == "" {
			info, er = scene.GetLesson(in.Uid)
		}
	} else {
		if in.Filter == "" {
			info, er = cache.Context().GetLesson(in.Uid)
		}

	}

	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchLesson(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *LessonService) GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyLessonList) error {
	path := "lesson.getByFilter"
	inLog(path, in)
	if in.Number < 10 {
		in.Number = 10
	}

	var list []*cache.LessonInfo
	var err error
	if len(in.Parent) > 1 {
		list, err = cache.Context().GetLessons(in.Parent)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
			return nil
		}

	} else {

	}
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = make([]*pb.LessonInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchLesson(info))
	}
	out.Pages = 0
	out.Page = in.Page
	out.Total = 1
	out.Status = outLog(path, fmt.Sprintf("the total = %d and length = %d", 1, len(out.List)))
	return nil
}

func (mine *LessonService) GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "lesson.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *LessonService) UpdateOne(ctx context.Context, in *pb.ReqLessonUpdate, out *pb.ReplyLessonInfo) error {
	path := "lesson.updateOne"
	inLog(path, in)
	info, err := cache.Context().GetLesson(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err1 := info.UpdateInfo(in.Name, in.Remark, in.Operator, in.Tags)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchLesson(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *LessonService) SetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyLessonInfo) error {
	path := "lesson.setByFilter"
	inLog(path, in)
	info, err := cache.Context().GetLesson(in.Uid)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var er error
	if in.Filter == "weight" {
		w, er1 := strconv.ParseInt(in.Value, 10, 32)
		if er1 == nil {
			er = info.UpdateWeight(in.Operator, uint32(w))
		}
	} else if in.Filter == "cover" {
		er = info.UpdateCover(in.Operator, in.Value)
	} else if in.Filter == "assets" {
		er = info.UpdateAssets(in.Operator, in.List)
	} else {
		er = errors.New("the filter not defined")
	}
	if er != nil {
		out.Status = outError(path, er.Error(), pbstatus.ResultStatus_FormatError)
		return nil
	}

	out.Info = switchLesson(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *LessonService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "lesson.removeOne"
	inLog(path, in)
	info, err := cache.Context().GetLesson(in.Uid)
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
