package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.school/cache"
	"omo.msa.school/proxy"
)

type SchoolService struct {}

func switchSchool(info *cache.SchoolInfo) *pb.SchoolInfo {
	tmp := new(pb.SchoolInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Name = info.Name
	tmp.Scene = info.Scene
	tmp.Support = info.Support
	tmp.Grade = uint32(info.MaxGrade())
	tmp.Entity = info.Entity
	tmp.Teachers = info.Teachers()
	tmp.Honors = make([]*pb.HonorInfo, 0, len(info.Honors))
	for _, honor := range info.Honors {
		tmp.Honors = append(tmp.Honors, switchHonor(honor))
	}
	tmp.Respects = make([]*pb.HonorInfo, 0, len(info.Respects))
	for _, honor := range info.Respects {
		tmp.Respects = append(tmp.Respects, switchHonor(honor))
	}
	tmp.Subjects = make([]*pb.SubjectInfo, 0, len(info.Subjects))
	for _, item := range info.Subjects {
		tmp.Subjects = append(tmp.Subjects, &pb.SubjectInfo{Uid: item.UID, Name: item.Name, Remark: item.Remark})
	}

	return tmp
}

func switchHonor(info proxy.HonorInfo) *pb.HonorInfo {
	tmp := new(pb.HonorInfo)
	tmp.Uid = info.UID
	tmp.Name = info.Name
	tmp.Remark = info.Remark
	tmp.Parent = info.Parent
	tmp.Children = make([]*pb.HonorInfo, 0, 1)
	tmp.Bries = make([]*pb.HonorBrief, 0 ,1)
	return tmp
}

func (mine *SchoolService)AddOne(ctx context.Context, in *pb.ReqSchoolAdd, out *pb.ReplySchoolInfo) error {
	path := "school.addOne"
	inLog(path, in)

	info, err1 := cache.Context().CreateSchool(in.Name, in.Entity, in.Scene, int(in.Grade))
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchSchool(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplySchoolInfo) error {
	path := "school.getOne"
	inLog(path, in)
	var school *cache.SchoolInfo
	if in.Filter == "scene" {
		school,_ = cache.Context().GetSchoolScene(in.Value)
	}else if in.Filter == "name" {
		school,_ = cache.Context().GetSchoolByName(in.Value)
	}else if in.Filter == "class" {
		school = cache.Context().GetSchoolByClass(in.Value)
	}else if in.Filter == "user" {
		school = cache.Context().GetSchoolByUser(in.Value)
	}else{
		if len(in.Uid) > 1 {
			school,_ = cache.Context().GetSchoolByUID(in.Uid)
		}else if len(in.Parent) > 1 {
			school,_ = cache.Context().GetSchoolScene(in.Parent)
		}else {
			school,_ = cache.Context().GetSchoolByName(in.Operator)
		}
	}

	if school == nil {
		out.Status = outError(path,"not found the school by", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchSchool(school)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplySchoolList) error {
	path := "school.getByFilter"
	inLog(path, in)

	var total uint32 = 0
	var max uint32 = 0
	var list = make([]*cache.SchoolInfo, 0, 10)
	if in.Filter == "class" {
		info := cache.Context().GetSchoolByClass(in.Value)
		if info != nil {
			list = append(list, info)
		}
	}else{
		total, max, list = cache.Context().AllSchools(in.Page, in.Number)
	}

	out.List = make([]*pb.SchoolInfo, 0, len(list))
	for _, info := range list {
		tmp := switchSchool(info)
		out.List = append(out.List, tmp)
	}
	out.Total = total
	out.Pages = max
	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "school.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplySchoolList) error {
	path := "school.getList"
	inLog(path, in)

	total, max, list := cache.Context().AllSchools(in.Page, in.Number)
	out.List = make([]*pb.SchoolInfo, 0, len(list))
	for _, info := range list {
		tmp := switchSchool(info)
		out.List = append(out.List, tmp)
	}
	out.Total = total
	out.Pages = max
	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)GetArray(ctx context.Context, in *pb.RequestList, out *pb.ReplySchoolList) error {
	path := "school.getArray"
	inLog(path, in)
	out.List = make([]*pb.SchoolInfo, 0, len(in.List))
	for _, uid := range in.List {
		info,_ := cache.Context().GetSchoolByUID(uid)
		if info != nil{
			out.List = append(out.List, switchSchool(info))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *SchoolService)UpdateOne(ctx context.Context, in *pb.ReqSchoolUpdate, out *pb.ReplySchoolInfo) error {
	path := "school.updateOne"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolScene(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := school.UpdateInfo(in.Name, in.Remark, in.Operator)
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchSchool(school)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)SetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplySchoolInfo) error {
	path := "school.setByFilter"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var err error
	if in.Filter == "support" {
		err = school.UpdateSupport(in.Operator, in.Value)
	}
	if err != nil {
		out.Status = outError(path,err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchSchool(school)
	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "school.removeOne"
	inLog(path, in)


	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)UpdateSubject(ctx context.Context, in *pb.ReqSchoolSubject, out *pb.ReplySchoolSubjects) error {
	path := "school.updateSubject"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolScene(in.Scene)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err1 := school.CreateSubject(in.Name, in.Remark)
	if err1 != nil {
		out.Status = outError(path,err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Subjects = make([]*pb.SubjectInfo, 0, 20)
	for _, item := range school.Subjects {
		out.Subjects = append(out.Subjects, &pb.SubjectInfo{Uid: item.UID, Name: item.Name, Remark: item.Remark})
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *SchoolService)AppendTeacher(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "school.appendTeacher"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	teacher := cache.Context().GetTeacher(in.Uid)
	if teacher == nil {
		out.Status = outError(path,"not found the teacher", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	er := school.AppendTeacher(teacher)
	if er != nil {
		out.Status = outError(path,er.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = school.Teachers()
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *SchoolService)SubtractTeacher(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyList) error {
	path := "school.subtractTeacher"
	inLog(path, in)
	school,_ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = school.Teachers()
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}
