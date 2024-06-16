package main

import (
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
	"log"
	"temporal-aone/backend/pkg"
)

func main() {
	// 创建 Temporal 客户端
	c, err := client.NewClient(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create Temporal client", err)
	}
	defer c.Close()

	// 创建 Worker
	w := worker.New(c, "release-task-queue", worker.Options{})

	// 注册工作流和活动
	w.RegisterWorkflow(pkg.ConfigWorkflow)
	w.RegisterWorkflow(pkg.BuildUploadWorkflow)
	w.RegisterWorkflow(pkg.ReleaseWorkflow)
	w.RegisterActivity(pkg.ConfigRepoActivity)
	w.RegisterActivity(pkg.BuildActivity)
	w.RegisterActivity(pkg.TestActivity)
	w.RegisterActivity(pkg.PackageActivity)
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
