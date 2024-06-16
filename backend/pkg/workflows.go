package pkg

import (
	"go.temporal.io/sdk/workflow"
	"time"
)

func generateVersion() string {
	return time.Now().Format("20060102150405")
}

func ConfigWorkflow(ctx workflow.Context, config Config) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)
	config.Version = generateVersion()

	// 执行ConfigRepoActivity
	err := workflow.ExecuteActivity(ctx, ConfigRepoActivity, config).Get(ctx, &config)
	if err != nil {
		logger.Error("ConfigRepoActivity failed.", "Error", err)
		return err
	}

	logger.Info("Config workflow completed successfully", "Version", config.Version)
	return nil
}

func BuildUploadWorkflow(ctx workflow.Context, config Config) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	// 执行BuildActivity
	err := workflow.ExecuteActivity(ctx, BuildActivity, config).Get(ctx, nil)
	if err != nil {
		logger.Error("BuildActivity failed.", "Error", err)
		return err
	}

	// 执行TestActivity
	err = workflow.ExecuteActivity(ctx, TestActivity, config).Get(ctx, nil)
	if err != nil {
		logger.Error("TestActivity failed.", "Error", err)
		return err
	}

	// 执行PackageActivity
	err = workflow.ExecuteActivity(ctx, PackageActivity, config).Get(ctx, nil)
	if err != nil {
		logger.Error("PackageActivity failed.", "Error", err)
		return err
	}

	// 执行UploadToECSActivity
	err = workflow.ExecuteActivity(ctx, UploadToECSActivity, config).Get(ctx, nil)
	if err != nil {
		logger.Error("UploadToECSActivity failed.", "Error", err)
		return err
	}

	logger.Info("Build and upload workflow completed successfully")
	return nil
}

func ReleaseWorkflow(ctx workflow.Context, config Config) error {
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	logger := workflow.GetLogger(ctx)

	// 优雅关停
	err := workflow.ExecuteActivity(ctx, GracefulShutdownActivity, config).Get(ctx, nil)
	if err != nil {
		logger.Error("GracefulShutdownActivity failed.", "Error", err)
		return err
	}

	// 重启应用
	err = workflow.ExecuteActivity(ctx, RestartApplicationActivity, config).Get(ctx, nil)
	if err != nil {
		logger.Error("RestartApplicationActivity failed.", "Error", err)
		return err
	}

	// 健康检查
	err = workflow.ExecuteActivity(ctx, HealthCheckActivity, config).Get(ctx, nil)
	if err != nil {
		logger.Error("HealthCheckActivity failed.", "Error", err)
		return err
	}

	logger.Info("Release workflow completed successfully", "Version", config.Version)
	return nil
}
