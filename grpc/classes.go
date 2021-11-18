package grpc

import (
	"context"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.school/cache"
)

type ClassService struct {}

func switchClass(info *cache.ClassInfo) *pb.ClassInfo {
	tmp := new(pb.ClassInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = info.CreateTime.Unix()
	tmp.Updated = info.UpdateTime.Unix()
	tmp.Name = info.Name
	tmp.Type = uint32(info.Type)

	return tmp
}

func (mine *ClassService)AddOne(ctx context.Context, in *pb.ReqClassAdd, out *pb.ReplyClassInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyClassList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyClassList) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService)UpdateOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyClassInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService)SetMaster(ctx context.Context, in *pb.ReqClassMaster, out *pb.ReplyInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService)StudentJoin(ctx context.Context, in *pb.ReqClassJoin, out *pb.ReplyInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService)StudentKick(ctx context.Context, in *pb.ReqClassKick, out *pb.ReplyInfo) error {
	path := "class.addOne"
	inLog(path, in)
	school := cache.Context().GetSchool(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, operator, uint16(in.Count), 0)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}

	out.Status = outLog(path, out)
	return nil
}