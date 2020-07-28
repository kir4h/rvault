package api_test

import (
	"testing"

	"rvault/internal/pkg/api"
)

func TestGetBasePath(t *testing.T) {
	type args struct {
		engine    string
		kvVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Smoke V1",
			args: args{
				engine:    "secret",
				kvVersion: "1",
			},
			want:    "secret",
			wantErr: false,
		},
		{
			name: "Smoke V2",
			args: args{
				engine:    "secret",
				kvVersion: "2",
			},
			want:    "secret/metadata",
			wantErr: false,
		},
		{
			name: "Fail on Unsupported Version",
			args: args{
				engine:    "secret",
				kvVersion: "99",
			},
			want:    "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := api.GetListBasePath(tt.args.engine, tt.args.kvVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetListBasePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetListBasePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetReadBasePath(t *testing.T) {
	type args struct {
		engine    string
		kvVersion string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Smoke Test V1",
			args: args{
				engine:    "secret",
				kvVersion: "1",
			},
			want:    "secret",
			wantErr: false,
		},
		{
			name: "Smoke Test V2",
			args: args{
				engine:    "secret",
				kvVersion: "2",
			},
			want:    "secret/data",
			wantErr: false,
		},
		{
			name: "Fail on Unsupported Version",
			args: args{
				engine:    "secret",
				kvVersion: "99",
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := api.GetReadBasePath(tt.args.engine, tt.args.kvVersion)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetReadBasePath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetReadBasePath() got = %v, want %v", got, tt.want)
			}
		})
	}
}
