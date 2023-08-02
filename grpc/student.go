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

type StudentService struct{}

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
	tmp.Status = uint32(info.Status)
	tmp.Sex = uint32(info.Sex)
	if class == "" {
		tmp.Class = info.Class
	} else {
		tmp.Class = class
	}
	tmp.School = info.School
	tmp.Custodians = make([]*pb.CustodianInfo, 0, len(info.Custodians))
	for _, custodian := range info.Custodians {
		tmp.Custodians = append(tmp.Custodians, &pb.CustodianInfo{Name: custodian.Name, Phones: custodian.Phones, Identify: custodian.Identity})
	}
	return tmp
}

func (mine *StudentService) AddOne(ctx context.Context, in *pb.ReqStudentAdd, out *pb.ReplyStudentInfo) error {
	path := "student.addOne"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	tmp := school.GetStudentByCard(in.Card)
	if tmp == nil {
		info, class, err1 := school.CreateStudent(in)
		if err1 != nil {
			out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
			return nil
		}
		out.Info = switchStudent(info, class)
	} else {
		class := school.GetClassByStudent(tmp.UID, cache.StudentAll)
		if class != nil {
			out.Info = switchStudent(tmp, class.UID)
		} else {
			cla := school.GetClass(in.Class)
			if cla != nil {
				_ = cla.AddStudent(tmp)
				out.Info = switchStudent(tmp, cla.UID)
			} else {
				out.Info = switchStudent(tmp, in.Class)
			}
		}
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) GetOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyStudentInfo) error {
	path := "student.getOne"
	inLog(path, in)
	var classUID = ""
	var info *cache.StudentInfo
	var class *cache.ClassInfo
	if len(in.Parent) > 1 {
		school, _ := cache.Context().GetSchoolByUID(in.Parent)
		if school == nil {
			out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if len(in.Filter) > 0 {
			if in.Filter == "card" {
				info = school.GetStudentByCard(in.Value)
			} else if in.Filter == "sn" {
				info = school.GetStudentBySN(in.Value)
			} else if in.Filter == "entity" {
				info = school.GetStudentByEntity(in.Value)
			}
		} else {
			class, info = school.GetStudent(in.Uid)
		}
		if info == nil {
			out.Status = outError(path, "not found the student", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if class != nil {
			classUID = class.UID
		}
	} else {
		st := cache.Context().GetStudent(in.Uid)
		if st == nil {
			out.Status = outError(path, "not found the student", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		info = st
	}
	out.Info = switchStudent(info, classUID)
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) GetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStudentList) error {
	path := "student.getByFilter"
	inLog(path, in)
	var list = make([]*cache.StudentInfo, 0, 10)
	if len(in.Parent) > 1 {
		school, _ := cache.Context().GetSchoolByUID(in.Parent)
		if school == nil {
			out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
			return nil
		}
		if in.Filter == "custodian" {
			list = school.GetStudentsByCustodian(in.Value)
		} else if in.Filter == "card" {
			student := school.GetStudentByCard(in.Value)
			if student != nil {
				list = append(list, student)
			}
		} else if in.Filter == "entity" {
			student := school.GetStudentByEntity(in.Value)
			if student != nil {
				list = append(list, student)
			}
		} else if in.Filter == "name" {
			list = school.GetStudentsByName(in.Value)
		} else if in.Filter == "class" {
			list = school.GetStudentsByClass(in.Value)
		} else if in.Filter == "sn" {
			tmp := school.GetStudentBySN(in.Value)
			if tmp != nil {
				list = append(list, tmp)
			}
		} else if in.Filter == "search" {
			list = school.SearchStudents(in.Value)
		} else if in.Filter == "enrol" {
			list = school.GetStudentsByEnrol(in.Value)
		}
	} else {
		if in.Filter == "entity" {
			list = cache.Context().GetStudentsByEntity(in.Value)
		} else if in.Filter == "entities" {
			for _, uid := range in.List {
				arr := cache.Context().GetStudentsByEntity(uid)
				if len(arr) > 0 {
					list = append(list, arr...)
				}
			}
		} else if in.Filter == "card" {
			if in.Params == "" {
				list = cache.Context().GetStudentsByCard(in.Value)
			} else {
				list = cache.Context().GetStudentsByIDCard(in.Value, in.Params)
			}
		} else if in.Filter == "custodian" {
			list = cache.Context().GetStudentsByCustodian(in.Value)
		}
	}
	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		class := cache.Context().GetClassByStudent(info.UID)
		cid := ""
		if class != nil {
			cid = class.UID
		}
		out.List = append(out.List, switchStudent(info, cid))
	}
	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *StudentService) GetList(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStudentList) error {
	path := "student.getList"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var total uint32 = 0
	var max uint32 = 0
	var list []*cache.StudentInfo
	if in.Filter == "entities" {
		total, max, list = school.GetPageStudentEntities(in.Page, in.Number)
	} else if in.Filter == "type" {
		tp, er := strconv.ParseInt(in.Value, 10, 32)
		if er == nil {
			total, max, list = school.GetStudentsByType(in.Page, in.Number, cache.StudentStatus(tp))
		}
	} else if in.Filter == "status" {
		st, er := strconv.ParseInt(in.Value, 10, 32)
		if er == nil {
			total, max, list = school.GetStudents(in.Page, in.Number, cache.StudentStatus(st))
		}
	} else {
		total, max, list = school.GetStudents(in.Page, in.Number, cache.StudentActive)
	}

	out.List = make([]*pb.StudentInfo, 0, len(list))
	for _, info := range list {
		out.List = append(out.List, switchStudent(info, ""))
	}
	out.Pages = max
	out.Total = total
	out.Status = outLog(path, fmt.Sprintf("the total = %d and length = %d", total, len(out.List)))
	return nil
}

func (mine *StudentService) GetArray(ctx context.Context, in *pb.RequestList, out *pb.ReplyStudentList) error {
	path := "student.getArray"
	inLog(path, in)
	out.List = make([]*pb.StudentInfo, 0, len(in.List))
	if len(in.Parent) < 1 {
		for _, uid := range in.List {
			info := cache.Context().GetStudent(uid)
			if info != nil {
				out.List = append(out.List, switchStudent(info, ""))
			}
		}
	} else {
		school, err := cache.Context().GetSchoolByUID(in.Parent)
		if school == nil {
			out.Status = outError(path, err.Error(), pbstatus.ResultStatus_NotExisted)
			return nil
		}
		for _, uid := range in.List {
			class, info := school.GetStudent(uid)
			if info != nil {
				if class == nil {
					out.List = append(out.List, switchStudent(info, ""))
				} else {
					out.List = append(out.List, switchStudent(info, class.UID))
				}
			}
		}
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *StudentService) GetStatistic(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStatistic) error {
	path := "student.getStatistic"
	inLog(path, in)

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) UpdateOne(ctx context.Context, in *pb.ReqStudentUpdate, out *pb.ReplyStudentInfo) error {
	path := "student.updateOne"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	_, info := school.GetStudent(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	custodians := make([]proxy.CustodianInfo, 0, len(in.Custodians))
	for _, custodian := range in.Custodians {
		custodians = append(custodians, proxy.CustodianInfo{Name: custodian.Name, Phones: custodian.Phones, Identity: custodian.Identify})
	}
	err1 := info.UpdateBase(in.Name, in.Sn, in.Card, in.Operator, uint8(in.Sex), custodians)
	if err1 != nil {
		out.Status = outError(path, err1.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Info = switchStudent(info, "")
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) SetByFilter(ctx context.Context, in *pb.RequestPage, out *pb.ReplyStudentInfo) error {
	path := "student.setByFilter"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	_, info := school.GetStudent(in.Uid)
	if info == nil {
		out.Status = outError(path, "not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}

	out.Info = switchStudent(info, "")
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) RemoveOne(ctx context.Context, in *pb.RequestInfo, out *pb.ReplyInfo) error {
	path := "student.removeOne"
	inLog(path, in)
	info, _ := cache.Context().GetSchoolByUID(in.Parent)
	if info == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.RemoveStudent(in.Uid, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}

	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) AddBatch(ctx context.Context, in *pb.ReqStudentBatch, out *pb.ReplyStudentList) error {
	path := "student.addBatch"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	out.List = make([]*pb.StudentInfo, 0, len(in.List))
	for _, item := range in.List {
		info, class, er := school.CreateStudent(item)
		if er == nil {
			tmp := switchStudent(info, class)
			out.List = append(out.List, tmp)
		}
	}

	out.Status = outLog(path, fmt.Sprintf("the length = %d", len(out.List)))
	return nil
}

func (mine *StudentService) BindEntity(ctx context.Context, in *pb.ReqStudentBind, out *pb.ReplyStudentInfo) error {
	path := "class.bindEntity"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	var info *cache.StudentInfo
	if len(in.Uid) > 1 {
		_, info = school.GetStudent(in.Uid)
	} else {
		info = school.GetStudentByCard(in.Card)
	}

	if info == nil {
		out.Status = outError(path, "not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.BindEntity(in.Entity, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchStudent(info, "")
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) UpdateCustodian(ctx context.Context, in *pb.ReqStudentCustodian, out *pb.ReplyStudentInfo) error {
	path := "class.updateCustodian"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Owner)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	_, info := school.GetStudent(in.Uid)

	if info == nil {
		out.Status = outError(path, "not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateCustodian(in.Name, in.Phones, in.Identify)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Info = switchStudent(info, "")
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) UpdateTags(ctx context.Context, in *pb.RequestList, out *pb.ReplyList) error {
	path := "class.tags"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	_, info := school.GetStudent(in.Uid)

	if info == nil {
		out.Status = outError(path, "not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateTags(in.List, in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.List = info.Tags
	out.Status = outLog(path, out)
	return nil
}

func (mine *StudentService) UpdateStatus(ctx context.Context, in *pb.RequestState, out *pb.ReplyInfo) error {
	path := "class.updateCustodian"
	inLog(path, in)
	school, _ := cache.Context().GetSchoolByUID(in.Parent)
	if school == nil {
		out.Status = outError(path, "not found the school by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	_, info := school.GetStudent(in.Flag)
	if info == nil {
		out.Status = outError(path, "not found the student by uid", pbstatus.ResultStatus_NotExisted)
		return nil
	}
	err := info.UpdateStatus(uint8(in.State), in.Operator)
	if err != nil {
		out.Status = outError(path, err.Error(), pbstatus.ResultStatus_DBException)
		return nil
	}
	out.Status = outLog(path, out)
	return nil
}
