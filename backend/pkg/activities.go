package pkg

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	gosundheit "github.com/AppsFlyer/go-sundheit"
	"github.com/AppsFlyer/go-sundheit/checks"
	"github.com/cloudflare/tableflip"
	"github.com/go-git/go-git/v5"
	gitHttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	RepoURL  string
	UserName string
	Token    string
	Tag      string

	BinaryPath     string
	ConfigFilePath string
	Version        string
	LocalPath      string

	ECSUploadPath string
	ECSUser       string
	ECSServer     string

	AccessKeyID     string
	AccessKeySecret string
	RegionID        string
	ECSIPAddress    string

	HealthCheckURL string
}

func generateFolderName(repoURL, tag string) string {
	repoName := fmt.Sprintf("%s_%s_%s", getRepoName(repoURL), time.Now().Format("20060102150405"), tag)
	return fmt.Sprintf("reposity/%s", repoName)
}

func getRepoName(repoURL string) string {
	parts := strings.Split(repoURL, "/")
	return parts[len(parts)-1]
}

// validateGitRepo 验证指定的Git仓库地址是否可以连接
func validateGitRepo(repoURL, username, token, cloneDir string) error {
	_, err := git.PlainClone(cloneDir, false, &git.CloneOptions{
		URL: repoURL,
		Auth: &gitHttp.BasicAuth{
			Username: username, // 通常是你的用户名或令牌
			Password: token,    // 密码或是个人访问令牌
		},
		Progress: os.Stdout,
	})
	return err
}

func ConfigRepoActivity(ctx context.Context, config Config) error {
	localPath := generateFolderName(config.RepoURL, config.Tag)

	if err := validateGitRepo(config.RepoURL, config.Tag, config.Token, localPath); err != nil {
		return err
	}

	// 将配置信息持久化到数据库，

	//db, err := gorm.Open(sqlite.Open("config_info.db"), &gorm.Config{})
	//if err != nil {
	//	return config, err
	//}

	//err = db.AutoMigrate(&Config{})
	//if err != nil {
	//	return config, fmt.Errorf("database migration error: %v", err)
	//}
	//err = shared.GetDB().Create(&config).Error
	//if err != nil {
	//	return config, fmt.Errorf("database insert error: %v", err)
	//}

	return nil
}

func UploadOSSActivity(ctx context.Context, config Config) error {
	//将指定目录可执行文件以及配置文件文件上传到oss上

	return nil
}

func CheckECSActivity(ctx context.Context, config Config) error {

	return nil
}

func BuildActivity(ctx context.Context, config Config) error {
	fmt.Println("Building the project...")

	var stdOut, stdErr bytes.Buffer
	cmd := exec.Command("go", "build", "-o", config.BinaryPath, config.LocalPath)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("build error: %v - stderr: %s", err, stdErr.String())
	}

	fmt.Println("Build output:", stdOut.String())
	return nil
}

func TestActivity(ctx context.Context, config Config) error {
	fmt.Println("Running tests...")

	var stdOut, stdErr bytes.Buffer
	cmd := exec.Command("go", "test", "./...", config.LocalPath)
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("test error: %v - stderr: %s", err, stdErr.String())
	}

	fmt.Println("Test output:", stdOut.String())
	return nil
}

func PackageActivity(ctx context.Context, config Config) error {
	fmt.Println("Packaging the project...")

	var stdOut, stdErr bytes.Buffer
	tarBinaryPath := fmt.Sprintf("%s_binary.tar.gz", config.LocalPath)
	tarConfigPath := fmt.Sprintf("%s_config.tar.gz", config.LocalPath)
	cmdBinary := exec.Command("tar", "-czvf", tarBinaryPath, config.BinaryPath)
	cmdConfig := exec.Command("tar", "-czvf", tarConfigPath, config.ConfigFilePath)
	cmdBinary.Stdout = &stdOut
	cmdBinary.Stderr = &stdErr
	cmdConfig.Stdout = &stdOut
	cmdConfig.Stderr = &stdErr

	err := cmdBinary.Run()
	if err != nil {
		return fmt.Errorf("package error: %v - stderr: %s", err, stdErr.String())
	}

	err = cmdConfig.Run()
	if err != nil {
		return fmt.Errorf("package error: %v - stderr: %s", err, stdErr.String())
	}

	fmt.Println("Package output:", stdOut.String())
	return nil
}

func UploadToECSActivity(ctx context.Context, config Config) error {
	fmt.Println("Uploading packages to ECS...")

	conn, err := net.Dial("tcp", config.ECSServer)
	if err != nil {
		return fmt.Errorf("failed to connect to ECS server: %v", err)
	}
	defer conn.Close()

	uploadFile := func(filePath string) error {
		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %v", err)
		}
		defer file.Close()

		_, err = io.Copy(conn, file)
		if err != nil {
			return fmt.Errorf("failed to upload file: %v", err)
		}

		return nil
	}

	tarBinaryPath := fmt.Sprintf("%s_binary.tar.gz", config.LocalPath)
	tarConfigPath := fmt.Sprintf("%s_config.tar.gz", config.LocalPath)

	if err := uploadFile(tarBinaryPath); err != nil {
		return err
	}
	if err := uploadFile(tarConfigPath); err != nil {
		return err
	}

	fmt.Println("Upload completed")
	return nil
}

func HealthCheckActivity(ctx context.Context, config Config) error {
	fmt.Println("Running health checks...")

	// 创建健康检查
	health := gosundheit.New()

	// 定义 HTTP 检查配置
	httpCheckConf := checks.HTTPCheckConfig{
		CheckName: config.ECSServer + ".health.check",
		Timeout:   1 * time.Second,
		URL:       config.HealthCheckURL,
	}

	// 创建 HTTP 健康检查
	httpCheck, err := checks.NewHTTPCheck(httpCheckConf)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// 注册健康检查
	err = health.RegisterCheck(
		httpCheck,
		gosundheit.InitialDelay(time.Second),       // 初始延迟
		gosundheit.ExecutionPeriod(10*time.Second), // 健康检查执行间隔
	)
	if err != nil {
		return fmt.Errorf("registering health check error: %v", err)
	}

	// 模拟等待一段时间以执行健康检查
	time.Sleep(30 * time.Second)

	// 获取并检查健康检查结果
	results, _ := health.Results()
	if result, ok := results[config.ECSServer+".health.check"]; ok {
		if result.Error != nil {
			return fmt.Errorf("health check failed: %v", result.Error)
		}
	} else {
		return fmt.Errorf("health check result not found for %s", config.ECSServer)
	}

	fmt.Println("Health check passed")
	return nil
}

func GracefulShutdownActivity(ctx context.Context, config Config) error {
	fmt.Println("Shutting down the application gracefully...")

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", config.ECSUser, config.ECSServer), "pkill -SIGTERM myapp")
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("graceful shutdown error: %v - stderr: %s", err, stdErr.String())
	}

	fmt.Println("Graceful shutdown successful:", stdOut.String())
	return nil
}

func RestartApplicationActivity(ctx context.Context, config Config) error {
	fmt.Println("Restarting the application...")

	upgrader, err := tableflip.New(tableflip.Options{})
	if err != nil {
		return fmt.Errorf("restart application error: %v", err)
	}
	defer upgrader.Stop()

	cmd := exec.Command("ssh", fmt.Sprintf("%s@%s", config.ECSUser, config.ECSServer), "nohup /path/to/deployed/binary &")
	cmd.Env = append(os.Environ(), fmt.Sprintf("GIT_ASKPASS=echo '%s'", config.Token))
	var stdOut, stdErr bytes.Buffer
	cmd.Stdout = &stdOut
	cmd.Stderr = &stdErr

	err = cmd.Run()
	if err != nil {
		return fmt.Errorf("restart application error: %v - stderr: %s", err, stdErr.String())
	}

	fmt.Println("Application restart successful:", stdOut.String())
	return nil
}
