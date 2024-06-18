package main

import (
	"log"
	"temporal-aone/backend/pkg"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// 创建 Temporal 客户端
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// 创建 Worker
	w := worker.New(c, "release-task-queue", worker.Options{})

	// 保存配置流
	w.RegisterWorkflow(pkg.ConfigWorkflow)
	w.RegisterActivity(pkg.ConfigRepoActivity)
	w.RegisterActivity(pkg.PackageActivity)

	w.RegisterActivity(pkg.CheckECSActivity)

	//机器上需要有个拉取oss配置的工具，定时任务一直拉取，

	//上传流
	w.RegisterWorkflow(pkg.BuildUploadWorkflow)
	w.RegisterActivity(pkg.UploadOSSActivity)

	//ecs处理流
	w.RegisterWorkflow(pkg.ReleaseWorkflow)
	w.RegisterActivity(pkg.UploadToECSActivity)

	w.RegisterActivity(pkg.GracefulShutdownActivity)
	w.RegisterActivity(pkg.RestartApplicationActivity)
	w.RegisterActivity(pkg.HealthCheckActivity)

	// 启动 Worker
	err = w.Start()
	if err != nil {
		log.Fatalln("Unable to start Worker", err)
	}

	log.Println("Worker started successfully")

	// 保持 Worker 运行
	select {}
}
