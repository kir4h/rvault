package kv_test

import (
	"reflect"
	"testing"

	"rvault/internal/pkg/kv"

	vapi "github.com/hashicorp/vault/api"
	"github.com/spf13/viper"
)

var wantSmokeTest map[string]map[string]string

func init() {
	wantSmokeTest = map[string]map[string]string{
		"/spain/admin": {
			"admin.conf": "dsfdsflfrf43l4tlp",
		},
		"/spain/malaga/random": {
			"my.key": "d3ewf2323r21e2",
		},
		"/france/paris/key": {
			"id_rsa": "ewdfpelfr23pwlrp32l4[p23lp2k",
			"id_dsa": "fewfowefkfkwepfkewkfpweokfeowkfpk",
		},
		"/uk/london/mi5": {
			"mi5.conf": "salt, 324r23432, false",
		},
	}
}

func TestRRead(t *testing.T) {
	cluster := createTestVault(t)
	defer cluster.Cleanup()
	client := cluster.Cores[0].Client
	type args struct {
		c            *vapi.Client
		engine       string
		path         string
		concurrency  uint32
		includePaths []string
		excludePaths []string
	}
	tests := []struct {
		name       string
		args       args
		viperFlags map[string]interface{}
		want       map[string]map[string]string
		wantErr    bool
	}{
		{
			name: "Smoke Test V2",
			args: args{
				c:            client,
				engine:       engineV2,
				path:         "/",
				includePaths: []string{"*"},
			},
			viperFlags: map[string]interface{}{
				"global.kv_version": "2",
			},
			want:    wantSmokeTest,
			wantErr: false,
		},
		{
			name: "Smoke Test V1 Buffered",
			args: args{
				c:            client,
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
				concurrency:  20,
			},
			viperFlags: map[string]interface{}{
				"global.kv_version": "1",
			},
			want:    wantSmokeTest,
			wantErr: false,
		},
		{
			name: "Unknown KV Version",
			args: args{
				c:            client,
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
			},
			viperFlags: map[string]interface{}{
				"global.kv_version": "99",
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Unset KV Version",
			args: args{
				c:            client,
				engine:       engine,
				path:         "/",
				includePaths: []string{"*"},
			},
			want:    wantSmokeTest,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.viperFlags {
				viper.Set(k, v)
			}
			got, err := kv.RRead(tt.args.c, tt.args.engine, tt.args.path, tt.args.includePaths, tt.args.excludePaths,
				tt.args.concurrency)
			viper.Reset()
			if (err != nil) != tt.wantErr {
				t.Errorf("RRead() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RRead() got = %v, want %v", got, tt.want)
			}
		})
	}
}
