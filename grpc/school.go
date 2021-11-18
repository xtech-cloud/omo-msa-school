package grpc

import (
	"context"
	"omo.msa.school/cache"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
)

type SchoolService struct {}

func switchSchool(info *cache.SchoolInfo) *pb.SchoolInfo {
	tmp := new(pb.SchoolInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Type = uint32(info.Type)

	return tmp
}

func (mine *SchoolService)AddOne(ctx context.Context, in *pb.ReqSchoolAdd, out *pb.ReplySchoolInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateSchooles(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.SchoolInfo, 0, len(list))
	for _, info := range list {
		tmp := switchSchool(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySchoolInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateSchooles(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.SchoolInfo, 0, len(list))
	for _, info := range list {
		tmp := switchSchool(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplySchoolList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateSchooles(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.SchoolInfo, 0, len(list))
	for _, info := range list {
		tmp := switchSchool(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)UpdateOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySchoolInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateSchooles(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.SchoolInfo, 0, len(list))
	for _, info := range list {
		tmp := switchSchool(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateSchooles(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.SchoolInfo, 0, len(list))
	for _, info := range list {
		tmp := switchSchool(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)UpdateSubject(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateSchooles(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.SchoolInfo, 0, len(list))
	for _, info := range list {
		tmp := switchSchool(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}
