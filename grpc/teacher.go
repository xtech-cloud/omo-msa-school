package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	"omo.msa.school/cache"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
)

type TeacherService struct {}

func switchTeacher(info *cache.TeacherInfo) *pb.TeacherInfo {
	tmp := new(pb.TeacherInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Type = uint32(info.Type)

	return tmp
}

func (mine *TeacherService)AddOne(ctx context.Context, in *pb.ReqTeacherAdd, out *pb.ReplyTeacherInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateTeacheres(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.TeacherInfo, 0, len(list))
	for _, info := range list {
		tmp := switchTeacher(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeacherList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateTeacheres(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.TeacherInfo, 0, len(list))
	for _, info := range list {
		tmp := switchTeacher(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTeacherList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateTeacheres(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.TeacherInfo, 0, len(list))
	for _, info := range list {
		tmp := switchTeacher(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)UpdateOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeacherInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateTeacheres(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.TeacherInfo, 0, len(list))
	for _, info := range list {
		tmp := switchTeacher(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateTeacheres(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.TeacherInfo, 0, len(list))
	for _, info := range list {
		tmp := switchTeacher(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService)AddBath(ctx context.Context, in *pb.ReqTeacherBatch, out *pb.ReplyTeacherList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateTeacheres(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.TeacherInfo, 0, len(list))
	for _, info := range list {
		tmp := switchTeacher(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}