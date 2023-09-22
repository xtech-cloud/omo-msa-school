package grpc

import (
	"context"
	"fmt"
	pb "github.com/xtech-cloud/omo-msp-school/proto/school"
	pbstatus "github.com/xtech-cloud/omo-msp-status/proto/status"
	"omo.msa.school/cache"
	"omo.msa.school/proxy"
	"strconv"
)

type ClassService struct{}

func switchClass(info *cache.ClassInfo) *pb.ClassInfo {
	tmp := new(pb.ClassInfo)
	tmp.Uid = info.UID
	tmp.Id = info.ID
	tmp.Created = uint64(info.CreateTime.Unix())
	tmp.Updated = uint64(info.UpdateTime.Unix())
	tmp.Name = info.Name
	tmp.Type = uint32(info.Type)
	tmp.Operator = info.Operator
	tmp.Creator = info.Creator
	tmp.Grade = uint32(info.Grade())
	tmp.Enrol = info.EnrolDate.String()
	tmp.No = uint32(info.Number)
	tmp.Master = info.Master
	tmp.Owner = info.School
	tmp.Assistant = info.Assistant
	tmp.Teachers = info.Teachers
	tmp.Students = make([]*pb.MemberInfo, 0, len(info.Members))
	for _, member := range info.Members {
		tmp.Students = append(tmp.Students, &pb.MemberInfo{Uid: member.UID, Student: member.Student, Status: uint32(member.Status), Remark: member.Remark})
	}
	return tmp
}

func (mine *ClassService) AddOne(ctx context.Context, in *pb.ReqClassAdd, out *pb.ReplyClassList) error {
	path := "class.addOne"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Scene)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	list, err1 := school.CreateClasses(in.Name, in.Enrol, in.Operator, uint16(in.Count), uint16(in.Type))
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
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

func (mine *ClassService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyClassInfo) error {
	path := "class.getOne"
	inLog(path, in)

	info := cache.Context().GetClass(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the class", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchClass(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService) GetList(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyClassList) error {
	path := "class.getList"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var total uint32
	var max uint32
	var list []*cache.ClassInfo
	var state uint64
	if in.Filter == "status" {
		state, _ = strconv.ParseUint(in.Value, 10, 32)
		total, max, list = school.GetClassesByPage(0, 0, int32(state))
	} else {
		total, max, list = school.GetClassesByPage(0, 0, -1)
	}

	out.List = make([]*pb.ClassInfo, 0, len(list))
	for _, info := range list {
		tmp := switchClass(info)
		out.List = append(out.List, tmp)
	}
	out.Total = total
	out.Page = max
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ClassService) GetArray(ctx context.Context, in *pb.RequestList, out *pb.ReplyClassList) error {
	path := "class.getArray"
	inLog(path, in)
	out.List = make([]*pb.ClassInfo, 0, len(in.List))
	for _, uid := range in.List {
		info := cache.Context().GetClass(uid)
		if info != nil {
			out.List = append(out.List, switchClass(info))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ClassService) GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyClassList) error {
	path := "class.getByFilter"
	inLog(path, in)
	if in.Parent == "" {

	} else {
		school, _ := cache.Context().GetSchoolBy(in.Parent)
		if school == nil {
			out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		out.List = make([]*pb.ClassInfo, 0, 3)
		var classes []*cache.ClassInfo
		if in.Filter == "student" {
			class := school.GetClassByEntity(in.Value, cache.StudentActive)
			if class != nil {
				classes = make([]*cache.ClassInfo, 0, 1)
				classes = append(classes, class)
			}
		} else if in.Filter == "master" {
			classes = school.GetClassesByMaster(in.Value)
		} else if in.Filter == "teacher" {
			classes = school.GetClassesByTeacher(in.Value)
		} else if in.Filter == "assistant" {
			classes = school.GetClassesByAssistant(in.Value)
		} else if in.Filter == "enrol" {
			date := new(proxy.DateInfo)
			err := date.Parse(in.Value)
			if err != nil {
				out.Status = outError(path, err.Error(), pbstatus.ResultStatus_FormatError)
				return nil
			}
			classes = school.GetClassesByEnrol(date.Year, date.Month)
		}
		for _, class := range classes {
			out.List = append(out.List, switchClass(class))
		}
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *ClassService) GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "class.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService) UpdateOne(ctx context.Context, in *pb.ReqClassUpdate, out *pb.ReplyClassInfo) error {
	path := "class.updateOne"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info := school.GetClass(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the class", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateInfo(in.Name, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchClass(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService) SetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyClassInfo) error {
	path := "class.setByFilter"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info := school.GetClass(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the class", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchClass(info)
	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "class.removeOne"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := school.RemoveClass(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService) SetMaster(ctx context.Context, in *pb.ReqClassMaster, out *pb.ReplyInfo) error {
	path := "class.setMaster"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	oClass := school.GetClassByMaster(in.Teacher)
	if oClass != nil && oClass.UID == in.Uid {
		out.Status = outLog(path, out)
		return nil
	}
	info := school.GetClass(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the class", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateMaster(in.Teacher, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	if oClass != nil {
		_ = oClass.UpdateMaster("", in.Operator)
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService) SetAssistant(ctx context.Context, in *pb.ReqClassMaster, out *pb.ReplyInfo) error {
	path := "class.setAssistant"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info := school.GetClass(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the class", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateAssistant(in.Teacher, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *ClassService) AppendStudent(ctx context.Context, in *pb.ReqClassStudent, out *pb.ReplyClassStudents) error {
	path := "class.appendStudent"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info := school.GetClass(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the class", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	oClass, student := school.GetStudent(in.Student)
	if student == nil {
		out.Status = outError(path, "not found the student", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.AddStudent(student)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	if oClass != nil {
		_ = oClass.RemoveStudent(in.Student, "change class", student.ID, cache.StudentLeave)
	}
	out.Students = make([]*pb.MemberInfo, 0, len(info.Members))
	for _, member := range info.Members {
		out.Students = append(out.Students, &pb.MemberInfo{Uid: member.UID, Student: member.Student, Status: uint32(member.Status), Remark: member.Remark})
	}
	out.Status = outLog(path, fmt.Sprintf("the students length = %d", len(out.Students)))
	return nil
}

func (mine *ClassService) SubtractStudent(ctx context.Context, in *pb.ReqClassStudent, out *pb.ReplyClassStudents) error {
	path := "class.subtractStudent"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	class, student := school.GetStudent(in.Student)
	if student == nil {
		out.Status = outError(path, "not found the student", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	if class == nil {
		out.Status = outError(path, "not found the student class", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	err := class.RemoveStudent(in.Student, in.Remark, student.ID, cache.StudentLeave)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	_ = student.UpdateStatus(cache.StudentLeave, in.Operator)
	out.Students = make([]*pb.MemberInfo, 0, len(class.Members))
	for _, member := range class.Members {
		out.Students = append(out.Students, &pb.MemberInfo{Uid: member.UID, Student: member.Student, Status: uint32(member.Status), Remark: member.Remark})
	}
	out.Status = outLog(path, fmt.Sprintf("the students length = %d", len(out.Students)))
	return nil
}

func (mine *ClassService) AppendTeacher(ctx context.Context, in *pb.ReqClassTeacher, out *pb.ReplyList) error {
	path := "class.appendTeacher"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.School)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info := school.GetClass(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the class", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.AppendTeacher(in.Teacher)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = info.Teachers
	out.Status = outLog(path, fmt.Sprintf("the teachers length = %d", len(out.List)))
	return nil
}

func (mine *ClassService) SubtractTeacher(ctx context.Context, in *pb.ReqClassTeacher, out *pb.ReplyList) error {
	path := "class.subtractTeacher"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolBy(in.School)
	if school == nil {
		out.Status = outError(path, "not found the school by scene", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	info := school.GetClass(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the student class", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.SubtractTeacher(in.Teacher)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = info.Teachers
	out.Status = outLog(path, fmt.Sprintf("the teachers length = %d", len(out.List)))
	return nil
}
