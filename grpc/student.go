package grpc

import (
	"context"
	"omo.msa.school/cache"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
)

type StudentService struct {}

func switchStudent(info *cache.StudentInfo) *pb.StudentInfo {
	tmp := new(pb.StudentInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Type = uint32(info.Type)

	return tmp
}

func (mine *StudentService)AddOne(ctx context.Context, in *pb.ReqStudentAdd, out *pb.ReplyStudentInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateStudentes(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		tmp := switchStudent(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyStudentList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateStudentes(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		tmp := switchStudent(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStudentList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateStudentes(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		tmp := switchStudent(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)UpdateOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyStudentInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateStudentes(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		tmp := switchStudent(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateStudentes(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		tmp := switchStudent(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)AddBach(ctx context.Context, in *pb.ReqStudentBatch, out *pb.ReplyStudentList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateStudentes(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		tmp := switchStudent(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)BindEntity(ctx context.Context, in *pb.ReqStudentBind, out *pb.ReplyStudentInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateStudentes(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		tmp := switchStudent(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

