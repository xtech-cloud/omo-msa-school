package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.school/cache"
	"omo.msa.school/proxy"
)

type TeacherService struct {}

func switchTeacher(info *cache.TeacherInfo) *pb.TeacherInfo {
	tmp := new(pb.TeacherInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Name = info.Name
	tmp.Entity = info.Entity
	tmp.User = info.User
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Classes = info.Classes
	tmp.Subjects = info.Subjects
	tmp.Owner = info.Owner
	tmp.Histories = make([]*pb.HistoryInfo, 0, len(info.Histories))
	for _, history := range info.Histories {
		tmp.Histories = append(tmp.Histories, switchHistory(&history))
	}
	return tmp
}

func switchHistory(info *proxy.HistoryInfo) *pb.HistoryInfo {
	tmp := new(pb.HistoryInfo)
	tmp.School = info.School
	tmp.Uid = info.UID
	tmp.Remark = info.Remark
	tmp.Enrol = info.Enrol
	tmp.Grade = uint32(info.Grade)
	tmp.Class = uint32(info.Class)
	return tmp
}

func (mine *TeacherService)AddOne(ctx context.Context, in *pb.ReqTeacherAdd, out *pb.ReplyTeacherInfo) error {
	path := "teacher.addOne"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info, err1 := school.CreateTeacher(in.Name, in.Entity, in.User, in.Operator)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchTeacher(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeacherInfo) error {
	path := "teacher.getOne"
	inLog(path, in)
	var info *cache.TeacherInfo
	if len(in.Parent) > 1 {
		info = cache.Context().GetTeacherByEntity(in.Parent)
	}else if len(in.Operator) > 1 {
		info = cache.Context().GetTeacherByUser(in.Operator)
	}else{
		info = cache.Context().GetTeacher(in.Uid)
	}

	if info == nil {
		out.Status = outError(path,"not found the teacher", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchTeacher(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTeacherList) error {
	path := "teacher.getList"
	inLog(path, in)

	total, max, list := cache.Context().AllTeachers(in.Page, in.Number)
	out.List = make([]*pb.TeacherInfo, 0, len(list))
	out.Pages = max
	for _, info := range list {
		out.List = append(out.List, switchTeacher(info))
	}
	out.Page = in.Page
	out.Total = total
	out.Status = outLog(path, fmt.Sprintf("the total = %d and length = %d", total, len(out.List)))
	return nil
}

func (mine *TeacherService)GetArray(ctx context.Context, in *pb.RequestList, out *pb.ReplyTeacherList) error {
	path := "teacher.getArray"
	inLog(path, in)
	out.List = make([]*pb.TeacherInfo, 0, len(in.List))
	for _, uid := range in.List {
		info := cache.Context().GetTeacher(uid)
		if info != nil{
			out.List = append(out.List, switchTeacher(info))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeacherService)GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTeacherList) error {
	path := "teacher.getByFilter"
	inLog(path, in)

	total, max, list := cache.Context().AllTeachers(in.Page, in.Number)
	out.List = make([]*pb.TeacherInfo, 0, len(list))
	out.Pages = max
	for _, info := range list {
		out.List = append(out.List, switchTeacher(info))
	}
	out.Page = in.Page
	out.Total = total
	out.Status = outLog(path, fmt.Sprintf("the total = %d and length = %d", total, len(out.List)))
	return nil
}

func (mine *TeacherService)UpdateOne(ctx context.Context, in *pb.ReqTeacherUpdate, out *pb.ReplyTeacherInfo) error {
	path := "teacher.updateOne"
	inLog(path, in)
	info := cache.Context().GetTeacher(in.Uid)
	if info == nil {
		out.Status = outError(path,"not found the teacher by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err1 := info.UpdateBase(in.Name, in.Operator, in.Classes, in.Subjects)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchTeacher(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "teacher.removeOne"
	inLog(path, in)
	info,_ := cache.Context().GetSchoolByUID(in.Parent)
	if info == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.RemoveTeacherByUID(in.Uid, "")
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)AddBatch(ctx context.Context, in *pb.ReqTeacherBatch, out *pb.ReplyTeacherList) error {
	path := "teacher.addBatch"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = make([]*pb.TeacherInfo, 0, len(in.List))
	for _, item := range in.List {
		info, er := school.CreateTeacher(item.Name, item.Entity, item.User, in.Operator)
		if er == nil {
			tmp := switchTeacher(info)
			out.List = append(out.List, tmp)
		}
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeacherService)UpdateTags(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "teacher.updateTags"
	inLog(path, in)
	//school := cache.Context().GetSchool(in.Parent)
	//if school == nil {
	//	out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
	//	return nil
	//}
	info := cache.Context().GetTeacher(in.Uid)
	if info == nil {
		out.Status = outError(path,"not found the teacher", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.List = in.List
	out.Status = outLog(path, out)
	return nil
}