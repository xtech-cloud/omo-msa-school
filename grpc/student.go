package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.school/cache"
	"omo.msa.school/proxy"
)

type StudentService struct {}

func switchStudent(info *cache.StudentInfo, class string) *pb.StudentInfo {
	tmp := new(pb.StudentInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Name = info.Name
	tmp.Entity = info.Entity
	tmp.Enrol = info.EnrolDate.String()
	tmp.Creator = info.Creator
	tmp.Operator = info.Operator
	tmp.Card = info.IDCard
	tmp.Sn = info.SN
	tmp.Sid = info.SID
	tmp.Sex = uint32(info.Sex)
	tmp.Class = class
	tmp.School = info.School
	tmp.Custodians = make([]*pb.CustodianInfo, 0, len(info.Custodians))
	for _, custodian := range info.Custodians {
		tmp.Custodians = append(tmp.Custodians, &pb.CustodianInfo{Name: custodian.Name, Phones: custodian.Phones, Identify: custodian.Identity})
	}
	return tmp
}

func (mine *StudentService)AddOne(ctx context.Context, in *pb.ReqStudentAdd, out *pb.ReplyStudentInfo) error {
	path := "student.addOne"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info, class, err1 := school.CreateStudent(in)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchStudent(info, class)
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyStudentInfo) error {
	path := "student.getOne"
	inLog(path, in)
	var classUID = ""
	var info *cache.StudentInfo
	if len(in.Parent) > 1 {
		school,_ := cache.Context().GetSchoolByUID(in.Parent)
		if school == nil {
			out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		class, st := school.GetStudent(in.Uid)
		if st == nil {
			out.Status = outError(path,"not found the student", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		info = st
		if class != nil {
			classUID = class.UID
		}
	}else{
		st := cache.Context().GetStudent(in.Uid)
		if st == nil {
			out.Status = outError(path,"not found the student", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		info = st
	}
	out.Info = switchStudent(info, classUID)
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStudentList) error {
	path := "student.getByFilter"
	inLog(path, in)
	var list []*cache.StudentInfo
	if len(in.Parent) > 1 {
		school,_ := cache.Context().GetSchoolByUID(in.Parent)
		if school == nil {
			out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Filter == "custodian" {
			list = school.GetStudentsByCustodian(in.Value)
		}
	}else{
		if in.Filter == "entity" {
			list = cache.Context().GetStudentsByEntity(in.Value)
		}else if in.Filter == "card" {
			if in.Params == "" {
				list = cache.Context().GetStudentsByCard(in.Value)
			}else{
				list = cache.Context().GetStudentsByIDCard(in.Value, in.Params)
			}
		}
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchStudent(info,""))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *StudentService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStudentList) error {
	path := "student.getList"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	total, max, list := school.GetPageStudents(in.Page, in.Number)
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchStudent(info,""))
	}
	out.Pages = max
	out.Total = total
	out.Status = outLog(path, fmt.Sprintf("the total = %d and length = %d", total, len(out.List)))
	return nil
}

func (mine *StudentService)GetArray(ctx context.Context, in *pb.RequestList, out *pb.ReplyStudentList) error {
	path := "student.getArray"
	inLog(path, in)
	out.List = make([]*pb.StudentInfo, 0, len(in.List))
	for _, uid := range in.List {
		info := cache.Context().GetStudent(uid)
		if info != nil{
			out.List = append(out.List, switchStudent(info,""))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *StudentService)UpdateOne(ctx context.Context, in *pb.ReqStudentUpdate, out *pb.ReplyStudentInfo) error {
	path := "student.updateOne"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	_,info := school.GetStudent(in.Uid)
	if info == nil {
		out.Status = outError(path,"not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	custodians := make([]proxy.CustodianInfo, 0, len(in.Custodians))
	for _, custodian := range in.Custodians {
		custodians = append(custodians, proxy.CustodianInfo{Name: custodian.Name, Phones: custodian.Phones, Identity: custodian.Identify})
	}
	err1 := info.UpdateBase(in.Name, in.Sn, in.Card, in.Operator, uint8(in.Sex), custodians)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchStudent(info,"")
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "student.removeOne"
	inLog(path, in)
	info,_ := cache.Context().GetSchoolByUID(in.Parent)
	if info == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.RemoveStudent(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)AddBatch(ctx context.Context, in *pb.ReqStudentBatch, out *pb.ReplyStudentList) error {
	path := "student.addBatch"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(in.List))
	for _, item := range in.List {
		info,class, er := school.CreateStudent(item)
		if er == nil {
			tmp := switchStudent(info, class)
			out.List = append(out.List, tmp)
		}
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *StudentService)BindEntity(ctx context.Context, in *pb.ReqStudentBind, out *pb.ReplyStudentInfo) error {
	path := "class.bindEntity"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var info *cache.StudentInfo
	if len(in.Uid) > 1{
		_,info = school.GetStudent(in.Uid)
	}else{
		info = school.GetStudentByCard(in.Card)
	}

	if info == nil {
		out.Status = outError(path,"not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.BindEntity(in.Entity, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchStudent(info, "")
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)UpdateCustodian(ctx context.Context, in *pb.ReqStudentCustodian, out *pb.ReplyStudentInfo) error {
	path := "class.updateCustodian"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	_,info := school.GetStudent(in.Uid)

	if info == nil {
		out.Status = outError(path,"not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateCustodian(in.Name, in.Phones, in.Identify)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchStudent(info, "")
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService)UpdateTags(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "class.updateCustodian"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path,"not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	_,info := school.GetStudent(in.Uid)

	if info == nil {
		out.Status = outError(path,"not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.List, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Tags
	out.Status = outLog(path, out)
	return nil
}

