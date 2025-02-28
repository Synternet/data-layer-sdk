package service

import (
	"reflect"
	"testing"
)

func TestNewSubject(t *testing.T) {
	type args struct {
		prefix   string
		suffixes []string
	}
	tests := []struct {
		name string
		args args
		want Subject
	}{
		{
			name: "simple",
			args: args{prefix: "prefix", suffixes: []string{"a", "b"}},
			want: "prefix.a.b",
		},
		{
			name: "empty prefix",
			args: args{prefix: "", suffixes: []string{"a", "b"}},
			want: "a.b",
		},
		{
			name: "empty suffix",
			args: args{prefix: "prefix", suffixes: []string{"", "b"}},
			want: "prefix.b",
		},
		{
			name: "spaces",
			args: args{prefix: " prefix ", suffixes: []string{"", "b"}},
			want: "prefix.b",
		},
		{
			name: "only prefix",
			args: args{prefix: " prefix ", suffixes: []string{}},
			want: "prefix",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewSubject(tt.args.prefix, tt.args.suffixes...); got != tt.want {
				t.Errorf("NewSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubject_Tokens(t *testing.T) {
	tokens := Subject("a.b.c.>").Tokens()
	want := []string{"a", "b", "c", ">"}

	if !reflect.DeepEqual(tokens, want) {
		t.Errorf("Subject.Tokens() = %v, want %v", tokens, want)
	}
}

func TestSubject_Validate(t *testing.T) {
	tests := []struct {
		name    string
		s       Subject
		wantErr bool
	}{
		{"good", Subject("a.*.c.>"), false},
		{"bad?", Subject("a.?"), true},
		{"empty token", Subject("a."), true},
		{"bad,", Subject("a.,s"), true},
		{"bad whitespace", Subject("hellow world"), true},
		{"bad whitespace1", Subject("hellow. world"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.s.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Subject.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSubject_Match(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		pattern Subject
		subject Subject
		want    bool
	}{
		{"longer", Subject("a.b"), Subject("a.b.c"), false},
		{"shorter", Subject("a.b"), Subject("a"), false},
		{"shorter *", Subject("a.b.*"), Subject("a.b"), false},
		{"even longer", Subject("a.b.*"), Subject("a.b.c.d"), false},
		{"exact", Subject("a.b"), Subject("a.b"), true},
		{"false", Subject("a.b"), Subject("a.c"), false},
		{"all exact", Subject("a.b.>"), Subject("a.b.>"), true},
		{"all longer", Subject("a.b.>"), Subject("a.b.c.d.e"), true},
		{"all subj reverse", Subject("a.b.c"), Subject("a.b.>"), false},
		{"all almost reverse", Subject("a.b.*"), Subject("a.b.>"), true}, // NOTE: This is a weird case
		{"all almost", Subject("a.b.>"), Subject("a.b.*"), true},
		{"all specific", Subject("a.b.>"), Subject("a.b.c"), true},
		{"middle reverse", Subject("a.b.c"), Subject("a.*.c"), false},
		{"middle1", Subject("a.*.c"), Subject("a.b.c"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pattern.Match(tt.subject); got != tt.want {
				t.Errorf("Subject.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubject_SymmetricMatch(t *testing.T) {
	type args struct {
	}
	tests := []struct {
		name    string
		pattern Subject
		subject Subject
		want    bool
	}{
		{"exact", Subject("a.b"), Subject("a.b"), true},
		{"false", Subject("a.b"), Subject("a.c"), false},
		{"all exact", Subject("a.b.>"), Subject("a.b.>"), true},
		{"all subj reverse", Subject("a.b.c"), Subject("a.b.>"), true},
		{"all almost", Subject("a.b.>"), Subject("a.b.*"), true},
		{"all almost reverse", Subject("a.b.*"), Subject("a.b.>"), true}, // NOTE: This is a weird case
		{"all specific", Subject("a.b.>"), Subject("a.b.c"), true},
		{"middle reverse", Subject("a.b.c"), Subject("a.*.c"), true},
		{"middle1", Subject("a.*.c"), Subject("a.b.c"), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.pattern.SymmetricMatch(tt.subject); got != tt.want {
				t.Errorf("SymmetricSubject.Match() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSubject_RemainingTokens(t *testing.T) {
	tests := []struct {
		name    string
		s       Subject
		pattern []string
		want    []string
	}{
		{
			"simple",
			Subject("a.b.c.d"),
			[]string{"a", "b"},
			[]string{"c", "d"},
		},
		{
			"empty subject",
			Subject(""),
			[]string{"a", "b"},
			[]string{},
		},
		{
			"empty pattern",
			Subject("a.b.c"),
			[]string{},
			[]string{"a", "b", "c"},
		},
		{
			"pattern.>",
			Subject("a.b.c.d.e"),
			[]string{"a", "b", ">"},
			[]string{"c", "d", "e"},
		},
		{
			"pattern.*",
			Subject("a.b.c.d.e"),
			[]string{"a", "b", "*"},
			[]string{"c"},
		},
		{
			"pattern.*.-",
			Subject("a.b.c.d.e"),
			[]string{"a", "b", "*", "-"},
			[]string{"c", "d", "e"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.s.RemainingTokens(tt.pattern); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Subject.RemainingTokens() = %v, want %v", got, tt.want)
			}
		})
	}
}
