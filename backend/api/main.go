package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"temporal-aone/backend/pkg"
	"temporal-aone/backend/shared"
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	healthhttp "github.com/AppsFlyer/go-sundheit/http"
	"github.com/cloudflare/tableflip"
	"github.com/gin-gonic/gin"
	"github.com/google/gops/agent"
	"go.temporal.io/sdk/client"
)

// 定义请求结构体
type WorkflowRequest struct {
	RepoURL        string `json:"repo_url"`
	Token          string `json:"token"`
	BinaryPath     string `json:"binary_path"`
	ConfigFilePath string `json:"config_file_path"`
	Version        string `json:"version"`
	ECSUploadPath  string `json:"ecs_upload_path"`
	ECSServer      string `json:"ecs_server"`
	ECSUser        string `json:"ecs_user"`
	ECSPassword    string `json:"ecs_password"`
	HealthCheckURL string `json:"health_check_url"`
}

func startConfigWorkflow(c *gin.Context) {
	var req WorkflowRequest
	// Parse the request body
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create Temporal client
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer temporalClient.Close()
	options := client.StartWorkflowOptions{
		ID:        "config-workflow",
		TaskQueue: "release-task-queue",
	}
	config := pkg.Config{
		RepoURL:        req.RepoURL,
		Token:          req.Token,
		BinaryPath:     req.BinaryPath,
		ConfigFilePath: req.ConfigFilePath,
		Version:        req.Version,
		ECSUploadPath:  req.ECSUploadPath,
		HealthCheckURL: req.HealthCheckURL,
	}
	we, err := temporalClient.ExecuteWorkflow(context.Background(), options, pkg.ConfigWorkflow, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
	})
}

// / Handler for starting build and upload workflow
func startBuildUploadWorkflow(c *gin.Context) {
	var req WorkflowRequest
	// Parse the request body
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create Temporal client
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer temporalClient.Close()
	options := client.StartWorkflowOptions{
		ID:        "build-upload-workflow",
		TaskQueue: "release-task-queue",
	}
	config := pkg.Config{
		RepoURL:        req.RepoURL,
		Token:          req.Token,
		BinaryPath:     req.BinaryPath,
		ConfigFilePath: req.ConfigFilePath,
		Version:        req.Version,
		ECSUploadPath:  req.ECSUploadPath,
		HealthCheckURL: req.HealthCheckURL,
	}
	we, err := temporalClient.ExecuteWorkflow(context.Background(), options, pkg.BuildUploadWorkflow, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
	})
}

// Handler for starting release workflow
func startReleaseWorkflow(c *gin.Context) {
	var req WorkflowRequest
	// Parse the request body
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Create Temporal client
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer temporalClient.Close()
	options := client.StartWorkflowOptions{
		ID:        "release-workflow",
		TaskQueue: "release-task-queue",
	}
	config := pkg.Config{
		RepoURL:        req.RepoURL,
		Token:          req.Token,
		BinaryPath:     req.BinaryPath,
		ConfigFilePath: req.ConfigFilePath,
		Version:        req.Version,
		ECSUploadPath:  req.ECSUploadPath,
		HealthCheckURL: req.HealthCheckURL,
	}
	we, err := temporalClient.ExecuteWorkflow(context.Background(), options, pkg.ReleaseWorkflow, config)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"workflow_id": we.GetID(),
		"run_id":      we.GetRunID(),
	})
}

func main() {
	// 初始化配置、日志和数据库
	shared.InitConfig()
	shared.InitLogger()
	shared.InitDatabase()

	// 设置 Gin 路由
	r := gin.Default()

	r.POST("/api/start-config", startConfigWorkflow)
	r.POST("/api/start-build-upload", startBuildUploadWorkflow)
	r.POST("/api/start-release", startReleaseWorkflow)

	// 创建健康检查实例
	health := gosundheit.New()

	// 定义一个 HTTP 依赖检查
	httpCheckConf := checks.HTTPCheckConfig{
		CheckName: "httpbin.url.check",
		Timeout:   1 * time.Second,
		URL:       "http://httpbin.org/status/200,300",
	}
	// 创建 HTTP 检查
	httpCheck, err := checks.NewHTTPCheck(httpCheckConf)
	if err != nil {
		shared.Logger.Infof("Failed to create HTTP check: %v", err)
	}
	// 注册 HTTP 检查
	err = health.RegisterCheck(
		httpCheck,
		gosundheit.InitialDelay(time.Second),
		gosundheit.ExecutionPeriod(10*time.Second),
	)
	if err != nil {
		shared.Logger.Infof("Failed to register check: %v", err)
	}

	// 注册健康检查端点
	r.GET("/admin/health.json", gin.WrapH(healthhttp.HandleHealthJSON(health)))

	// 设置 Tableflip
	upg, err := tableflip.New(tableflip.Options{})
	if err != nil {
		shared.Logger.Infof("tableflip setup failed: %v", err)
	}

	// 启动 HTTP 服务
	ln, err := upg.Fds.Listen("tcp", ":"+strconv.Itoa(shared.Config.Server.Port))
	if err != nil {
		shared.Logger.Fatalf("failed to listen on port: %v", err)
	}

	srv := &http.Server{
		Handler: r,
	}

	go func() {
		if err := srv.Serve(ln); err != nil && err != http.ErrServerClosed {
			shared.Logger.Fatalf("serve failed: %s", err)
		}
	}()

	// Unix 信号处理
	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
		for {
			switch <-sig {
			case syscall.SIGHUP:
				if err := upg.Upgrade(); err != nil {
					shared.Logger.Printf("upgrade failed: %v", err)
				}
			default:
				if err := srv.Shutdown(context.Background()); err != nil {
					shared.Logger.Printf("shutdown error: %v", err)
				}
				return
			}
		}
	}()

	if err := agent.Listen(agent.Options{}); err != nil {
		shared.Logger.Fatal(err)
	}

	// 升级准备就绪
	if err := upg.Ready(); err != nil {
		shared.Logger.Fatalf("upgrade ready error: %v", err)
	}

	// 等待退出
	<-upg.Exit()
	shared.Logger.Println("reloaded")
}
