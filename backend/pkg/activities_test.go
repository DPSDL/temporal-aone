package pkg

import (
	"context"
	"testing"
)

func TestConfigRepoActivity(t *testing.T) {
	type args struct {
		ctx    context.Context
		config Config
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				ctx: nil,
				config: Config{
					RepoURL:        "https://github.com/DPSDL/HaoNing.git",
					UserName:       "DPSDL",
					Token:          "",
					Tag:            "v1.0",
					BinaryPath:     "",
					ConfigFilePath: "",
					Version:        "",
					LocalPath:      "",
					ECSUploadPath:  "",
					ECSServer:      "",
					ECSUser:        "",
					HealthCheckURL: "",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ConfigRepoActivity(tt.args.ctx, tt.args.config); (err != nil) != tt.wantErr {
				t.Errorf("ConfigRepoActivity() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
