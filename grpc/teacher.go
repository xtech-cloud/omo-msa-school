package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.school/cache"
	"omo.msa.school/proxy"
)

type TeacherService struct{}

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
	school := cache.Context().GetSchoolByTeacher(info.UID)
	if school != nil {
		tmp.Owner = school.UID
		tmp.Classes = school.GetClassesUIDsByTeacher(info.UID)
	}
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

func (mine *TeacherService) AddOne(ctx context.Context, in *pb.ReqTeacherAdd, out *pb.ReplyTeacherInfo) error {
	path := "teacher.addOne"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info, err1 := school.CreateTeacher(in.Name, in.Entity, in.User, in.Operator, in.Classes, in.Subjects)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	for _, uid := range in.Classes {
		class := school.GetClass(uid)
		if class != nil {
			_ = class.AppendTeacher(info.UID)
		}
	}
	out.Info = switchTeacher(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyTeacherInfo) error {
	path := "teacher.getOne"
	inLog(path, in)
	var info *cache.TeacherInfo
	if len(in.Parent) > 1 {
		school, err := cache.Context().GetSchoolByUID(in.Parent)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Filter == "" {
			info = school.GetTeacher(in.Value)
		} else if in.Filter == "entity" {
			info = school.GetTeacherByEntity(in.Value)
		} else if in.Filter == "user" {
			info = school.GetTeacherByUser(in.Value)
		}
	} else {
		if in.Filter == "" {
			info = cache.Context().GetTeacher(in.Value)
		} else if in.Filter == "entity" {
			info = cache.Context().GetTeacherByEntity(in.Value)
		} else if in.Filter == "user" {
			info = cache.Context().GetTeacherByUser(in.Value)
		}

	}

	if info == nil {
		out.Status = outError(path, "not found the teacher", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchTeacher(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService) GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTeacherList) error {
	path := "teacher.getList"
	inLog(path, in)

	var list []*cache.TeacherInfo
	var total uint32 = 0
	var max uint32 = 0
	if len(in.Parent) > 1 {
		school, err := cache.Context().GetSchoolByUID(in.Parent)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
			return nil
		}
		total, max, list = school.GetTeachersByPage(in.Page, in.Number)
	} else {
		total, max, list = cache.Context().AllTeachers(in.Page, in.Number)
	}

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

func (mine *TeacherService) GetArray(ctx context.Context, in *pb.RequestList, out *pb.ReplyTeacherList) error {
	path := "teacher.getArray"
	inLog(path, in)
	out.List = make([]*pb.TeacherInfo, 0, len(in.List))
	for _, uid := range in.List {
		info := cache.Context().GetTeacher(uid)
		if info != nil {
			out.List = append(out.List, switchTeacher(info))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeacherService) GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTeacherList) error {
	path := "teacher.getByFilter"
	inLog(path, in)
	if in.Number < 10 {
		in.Number = 10
	}
	out.List = make([]*pb.TeacherInfo, 0, 5)
	if len(in.Parent) > 1 {
		school, err := cache.Context().GetSchoolByUID(in.Parent)
		if err != nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Filter == "entity" {
			info := school.GetTeacherByEntity(in.Value)
			if info != nil {
				out.List = append(out.List, switchTeacher(info))
			}
		} else if in.Filter == "name" {
			list := school.GetTeachersByName(in.Value)
			for _, info := range list {
				out.List = append(out.List, switchTeacher(info))
			}
		} else if in.Filter == "leave" {
			list := cache.Context().GetLeaveTeachers(school.UID)
			for _, info := range list {
				out.List = append(out.List, switchTeacher(info))
			}
		}
	} else {
		if in.Filter == "name" {
			info := cache.Context().GetTeacherByName(in.Value)
			if info != nil {
				out.List = append(out.List, switchTeacher(info))
			}
		} else if in.Filter == "user" {
			info := cache.Context().GetTeacherByUser(in.Value)
			if info != nil {
				out.List = append(out.List, switchTeacher(info))
			}
		}
	}

	out.Pages = 0
	out.Page = in.Page
	out.Total = 1
	out.Status = outLog(path, fmt.Sprintf("the total = %d and length = %d", 1, len(out.List)))
	return nil
}

func (mine *TeacherService) GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "teacher.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService) UpdateOne(ctx context.Context, in *pb.ReqTeacherUpdate, out *pb.ReplyTeacherInfo) error {
	path := "teacher.updateOne"
	inLog(path, in)
	info := cache.Context().GetTeacher(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the teacher by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err1 := info.UpdateBase(in.Name, in.Operator, in.Classes, in.Subjects)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchTeacher(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService) SetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyTeacherInfo) error {
	path := "teacher.setByFilter"
	inLog(path, in)
	info := cache.Context().GetTeacher(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the teacher by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchTeacher(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "teacher.removeOne"
	inLog(path, in)
	info, _ := cache.Context().GetSchoolByUID(in.Parent)
	if info == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.RemoveTeacherByUID(in.Uid, in.Value)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *TeacherService) AddBatch(ctx context.Context, in *pb.ReqTeacherBatch, out *pb.ReplyTeacherList) error {
	path := "teacher.addBatch"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = make([]*pb.TeacherInfo, 0, len(in.List))
	for _, item := range in.List {
		info, er := school.CreateTeacher(item.Name, item.Entity, item.User, in.Operator, item.Classes, item.Subjects)
		if er == nil {
			tmp := switchTeacher(info)
			out.List = append(out.List, tmp)
			for _, uid := range item.Classes {
				class := school.GetClass(uid)
				if class != nil {
					_ = class.AppendTeacher(info.UID)
				}
			}
		}
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *TeacherService) UpdateTags(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "teacher.updateTags"
	inLog(path, in)
	//school := cache.Context().GetSchool(in.Parent)
	//if school == nil {
	//	out.Status = outError(path,"not found the school by scene", pbstatus.ResultStatus_NotExisted)
	//	return nil
	//}
	info := cache.Context().GetTeacher(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the teacher", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.Operator, in.List)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.List = in.List
	out.Status = outLog(path, out)
	return nil
}
