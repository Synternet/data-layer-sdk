package rpc

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

func Test_pascalToDotCase(t *testing.T) {
	tests := []struct {
		name string
		args string
		want []string
	}{
		{
			name: "regular",
			args: "simplename",
			want: []string{"simplename"},
		},
		{
			name: "Regular",
			args: "Simplename",
			want: []string{"simplename"},
		},
		{
			name: "XRegular",
			args: "XSimplename",
			want: []string{"x", "simplename"},
		},
		{
			name: "camelCase",
			args: "simpleName",
			want: []string{"simple", "name"},
		},
		{
			name: "PascalCase",
			args: "SimpleName",
			want: []string{"simple", "name"},
		},
		{
			name: "PKGPascalCase",
			args: "PREFIXSimpleName",
			want: []string{"prefix", "simple", "name"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := splitPascalCase(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_deriveServiceTokens(t *testing.T) {
	tests := []struct {
		name        string
		serviceName protoreflect.FullName
		want        []string
	}{
		{
			name:        "regular",
			serviceName: protoreflect.FullName("types.rpc.MySpecialService"),
			want:        []string{"service", "rpc", "my", "special"},
		},
		{
			name:        "local",
			serviceName: protoreflect.FullName("MyService"),
			want:        []string{"service", "my"},
		},
		{
			name:        "semi local",
			serviceName: protoreflect.FullName("pkg.MyService"),
			want:        []string{"service", "pkg", "my"},
		},
		{
			name:        "empty",
			serviceName: protoreflect.FullName(""),
			want:        []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := deriveServiceTokens(tt.serviceName)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_getCustomServicePrefix(t *testing.T) {
	serviceDesc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName("synternet.rpc.TestService"))
	require.NoError(t, err)
	require.NotNil(t, serviceDesc)

	got := getCustomServicePrefix(serviceDesc)
	assert.Equal(t, got, "override.test")
}

func Test_getCustomMethodSuffix(t *testing.T) {
	serviceDesc, err := protoregistry.GlobalFiles.FindDescriptorByName(protoreflect.FullName("synternet.rpc.TestService"))
	require.NoError(t, err)
	require.NotNil(t, serviceDesc)

	methodDesc := serviceDesc.(protoreflect.ServiceDescriptor).Methods().ByName(protoreflect.Name("TestStreamOnly"))
	require.NotNil(t, methodDesc)

	got := getCustomMethodSuffix(methodDesc)
	assert.Equal(t, got, "override.test.stream.data")
}

func Test_splitAndParametrizeTokens(t *testing.T) {
	type args struct {
		tokens []string
		vars   map[string]string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple",
			args: args{tokens: []string{"a", "b.c"}, vars: nil},
			want: []string{"a", "b", "c"},
		},
		{
			name: "no vars",
			args: args{tokens: []string{"a", "b.{c}"}, vars: nil},
			want: []string{"a", "b", "*"},
		},
		{
			name: "vars",
			args: args{tokens: []string{"a", "b.{c}"}, vars: map[string]string{"c": "123"}},
			want: []string{"a", "b", "123"},
		},
		{
			name: "split vars",
			args: args{tokens: []string{"a", "b.{c.d}"}, vars: map[string]string{"c": "123"}},
			want: []string{"a", "b", "123", "*"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitAndParametrizeTokens(tt.args.tokens, tt.args.vars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("splitAndParametrizeTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
