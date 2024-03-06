package main

import (
	"crypto/md5"
	"encoding/hex"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/logger"
	_ "github.com/micro/go-plugins/registry/consul/v2"
	_ "github.com/micro/go-plugins/registry/etcdv3/v2"
	"github.com/robfig/cron/v3"
	proto "github.com/xtech-cloud/omo-msp-school/proto/school"
	"io"
	"omo.msa.school/cache"
	"omo.msa.school/config"
	"omo.msa.school/grpc"
	"os"
	"path/filepath"
	"time"
)

var (
	BuildVersion string
	BuildTime    string
	CommitID     string
)

func main() {
	config.Setup()
	err := cache.InitData()
	if err != nil {
		panic(err)
	}
	// New Service
	service := micro.NewService(
		micro.Name("omo.msa.school"),
		micro.Version(BuildVersion),
		micro.RegisterTTL(time.Second*time.Duration(config.Schema.Service.TTL)),
		micro.RegisterInterval(time.Second*time.Duration(config.Schema.Service.Interval)),
		micro.Address(config.Schema.Service.Address),
	)
	// Initialise service
	service.Init()
	// Register Handler
	_ = proto.RegisterClassesServiceHandler(service.Server(), new(grpc.ClassService))
	_ = proto.RegisterSchoolServiceHandler(service.Server(), new(grpc.SchoolService))
	_ = proto.RegisterStudentServiceHandler(service.Server(), new(grpc.StudentService))
	_ = proto.RegisterTeacherServiceHandler(service.Server(), new(grpc.TeacherService))
	_ = proto.RegisterTimetableServiceHandler(service.Server(), new(grpc.TimetableService))
	_ = proto.RegisterLessonServiceHandler(service.Server(), new(grpc.LessonService))
	_ = proto.RegisterScheduleServiceHandler(service.Server(), new(grpc.ScheduleService))

	go checkTimer()

	app, _ := filepath.Abs(os.Args[0])

	logger.Info("-------------------------------------------------------------")
	logger.Info("- Micro Service Agent -> Run")
	logger.Info("-------------------------------------------------------------")
	logger.Infof("- version      : %s", BuildVersion)
	logger.Infof("- application  : %s", app)
	logger.Infof("- md5          : %s", md5hex(app))
	logger.Infof("- build        : %s", BuildTime)
	logger.Infof("- commit       : %s", CommitID)
	logger.Info("-------------------------------------------------------------")
	// Run service
	if err := service.Run(); err != nil {
		logger.Fatal(err)
	}
}

func checkTimer() {
	time.Sleep(time.Second * 5)
	cache.Context().CheckStudentFinish()
	cache.Context().CheckStudentError()
	cli := cron.New()
	_, er := cli.AddFunc("1 22 1 * *", func() {
		cache.Context().CheckStudentFinish()
	})
	if er != nil {
		logger.Warn("start cron failed that err = " + er.Error())
		return
	}
	cli.Start()

	//std := new(grpc.StudentService)
	//in := &proto.RequestPage{
	//	Parent: "642b90575fd5b1224731c00a",
	//	Filter: "bind",
	//	List:   []string{"1", "2"},
	//}
	//out := &proto.ReplyStudentList{}
	//std.GetByFilter(context.Background(), in, out)
}

func md5hex(_file string) string {
	h := md5.New()

	f, err := os.Open(_file)
	if err != nil {
		return ""
	}
	defer f.Close()

	io.Copy(h, f)

	return hex.EncodeToString(h.Sum(nil))
}
