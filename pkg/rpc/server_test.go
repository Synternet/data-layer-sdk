package rpc

import (
	"reflect"
	"testing"

	_ "github.com/synternet/data-layer-sdk/x/synternet/rpc"
)

func Test_extractServiceVars(t *testing.T) {
	type args struct {
		svcName string
		vars    map[string]string
	}
	tests := []struct {
		name string
		args args
		want map[string]string
	}{
		{
			name: "nil",
			args: args{svcName: "synternet.rpc.MyService", vars: nil},
			want: nil,
		},
		{
			name: "global",
			args: args{svcName: "synternet.rpc.MyService", vars: map[string]string{"var": "value"}},
			want: map[string]string{"var": "value"},
		},
		{
			name: "scoped",
			args: args{svcName: "synternet.rpc.MyService", vars: map[string]string{"var": "value", "organization.rpc.MyService/Service": "123", "synternet.rpc.MyService/variable": "654"}},
			want: map[string]string{"var": "value", "variable": "654"},
		},
		{
			name: "scoped has precedence",
			args: args{svcName: "synternet.rpc.MyService", vars: map[string]string{"synternet": "value", "organization.rpc.MyService/Service": "123", "synternet.rpc.MyService/synternet": "654"}},
			want: map[string]string{"synternet": "654"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := extractServiceVars(tt.args.svcName, tt.args.vars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("extractServiceVars() = %v, want %v", got, tt.want)
			}
		})
	}
}
