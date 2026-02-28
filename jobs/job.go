package jobs

import (
	"context"
	"dbgod_workOrder_job/config"
	"dbgod_workOrder_job/loggerzap"
	"dbgod_workOrder_job/service"
	"fmt"
	"github.com/darkit/cron"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

// DbJob 实现Job接口的简单任务
type DbJob struct {
	JobName   string //每个job的名称
	Counter   *int64
	DsnSource string //源数据库连接字符串
	SqlStr    string //查询源库的sql
	Token     string
	ApiGet    string
	ApiPost   string
	ApiPut    string
}

func (j *DbJob) Name() string {
	return j.JobName
}

func (j *DbJob) Run(ctx context.Context) error {
	// 创建客户端
	client := service.NewDBClient(
		j.ApiGet,
		service.WithTimeout(15*time.Second),
		service.WithHeader("token", j.Token),
	)

	// 获取工单数据
	ctx2 := context.Background()
	workOrders, err := client.GetWorkOrdersList(ctx2)
	// 输出get请求工单列表结果
	if err != nil {
		loggerzap.L().Error(err.Error())
		return err
	} else {
		loggerzap.L().Info(fmt.Sprintf("未审核工单记录数:%d Msg:%s", workOrders.Data.Count, workOrders.Msg))
	}

	// 遍历工单列表响应,每一行就是一个工单
	for _, order := range workOrders.Data.List {
		// 创建客户端（设置自定义超时和请求头）
		client = service.NewDBClient(
			j.ApiPost,
			service.WithTimeout(15*time.Second),
			service.WithHeader("token", j.Token),
		)

		// 构造请求参数-创建数据库
		req := &service.CreateDBRequest{
			WorkOrderId:   order.Id,
			DbType:        order.DbTypeInt,
			DbName:        order.DbName,
			Department:    order.DepartmentInt,
			ApplicantUser: order.ApplicantUser,
		}

		// 发送请求-创建数据库
		resp, err := client.CreateDatabase(context.Background(), req)
		if err != nil {
			loggerzap.L().Error(fmt.Sprintf("Error: %v\n", err))
			return err
		}

		// 处理响应-创建数据库
		loggerzap.L().Info(fmt.Sprintf("Response: %+v\n", resp))

		// 创建数据库之后等待几秒，有些数据库建库比较慢
		//time.Sleep(5 * time.Second)

		// 更新下工单状态，1改2
		client = service.NewDBClient(
			j.ApiPut,
			service.WithTimeout(15*time.Second),
			service.WithHeader("token", j.Token),
		)
		// 构造请求参数-更新工单状态
		req2 := &service.UpdateWorkorderRequest{
			Status:        2,
			DbName:        order.DbName,
			ID:            order.Id,
			UsageInt:      order.UsageInt,
			TimeLimitInt:  order.TimeLimitInt,
			DepartmentInt: order.DepartmentInt,
			DbTypeInt:     order.DbTypeInt,
		}
		// 发送请求-更新工单状态
		resp2, err2 := client.UpdateWorkorder(context.Background(), req2)
		if err2 != nil {
			loggerzap.L().Info(fmt.Sprintf("Error: %v\n", err2))
		}

		// 处理响应-更新工单状态
		loggerzap.L().Info(fmt.Sprintf("Response: %+v\n", resp2))

	}

	count := atomic.AddInt64(j.Counter, 1)

	loggerzap.L().Info(fmt.Sprintf("%s executed (count: %d)\n", j.JobName, count))
	return nil
}

func StartJob() {
	// 创建调度器
	scheduler := cron.New(cron.WithLogger(&config.CustomLogger{}))
	// 计数器
	counter1 := int64(0)
	// 创建调度器的job
	// 遍历到单个instance
	for _, job := range config.GetConfig().Jobs {
		loggerzap.L().Info("current run instance:" + job.Instance)

		var dsnSource string
		//源库的连接
		//for _, connection := range job.DsnSource {
		//	dsnSource = connection
		//}

		// 遍历配置文件queries一个或多个name与query属性
		for _, query := range job.Queries {
			loggerzap.L().Info("job name:" + query.Name)

			err := scheduler.ScheduleJobByName(job.Interval, &DbJob{
				JobName:   query.Name,
				DsnSource: dsnSource,
				SqlStr:    query.Query,
				Token:     job.Token,
				ApiGet:    job.ApiGet,
				ApiPost:   job.ApiPost,
				ApiPut:    job.ApiPut,
				Counter:   &counter1,
			}, cron.JobOptions{
				Async:         true,
				MaxConcurrent: 1,
				Timeout:       30 * time.Second,
			})
			if err != nil {
				loggerzap.L().Fatal("添加优雅任务失败:" + err.Error())
			}
		}

	}

	// 启动调度器
	err := scheduler.Start()
	if err != nil {
		loggerzap.L().Fatal("启动调度器失败:" + err.Error())
	}

	// 等待信号
	loggerzap.L().Info(fmt.Sprintf("▶️  调度器已启动，按 Ctrl+C 停止...\n\n"))

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	loggerzap.L().Info(fmt.Sprintf("🛑 正在停止调度器..."))

	scheduler.Stop()
	loggerzap.L().Info(fmt.Sprintf("✅ 调度器已停止"))

	// 显示最终统计
	loggerzap.L().Info(fmt.Sprintf("📊 最终统计:"))
	loggerzap.L().Info(fmt.Sprintf("- 接口任务执行: %d次\n", atomic.LoadInt64(&counter1)))
}
