// options package is a package that stores publisher's configuration and has convenience functions to configure it.
package options

import (
	"testing"
)

func TestOptions_ServiceSubject(t *testing.T) {
	tests := []struct {
		name          string
		SubjectPrefix string
		Name          string
		prefix        string
		parts         []string
		streamId      string
		want          string
	}{
		{"original subject", "test", "case", ".", []string{"foo", "bar", ">"}, "", "test.case.foo.bar.>"},
		{"original subject w/o prefix", "test", "case", "", []string{"foo", "bar", ">"}, "", "foo.bar.>"},
		{"original subject w/ prefix", "test", "case", "prefix.t", []string{"foo", "bar", ">"}, "", "prefix.t.foo.bar.>"},
		{"subject with stream id", "test", "case", ".", []string{"foo", "bar", ">"}, "stream.id", "test.case.foo.bar.stream.id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				Prefix: tt.SubjectPrefix,
				Name:   tt.Name,
			}
			if got := o.subject(tt.prefix, tt.parts, tt.streamId); got != tt.want {
				t.Errorf("Options.ServiceSubject() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOptions_Subject(t *testing.T) {
	tests := []struct {
		name          string
		SubjectPrefix string
		Name          string
		suffixes      []string
		want          string
	}{
		{"simple case", "test", "case", []string{"foo", "*", "bar", ">"}, "test.case.foo.*.bar.>"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			o := Options{
				Prefix: tt.SubjectPrefix,
				Name:   tt.Name,
			}
			if got := o.Subject(tt.suffixes...); got != tt.want {
				t.Errorf("Options.Subject() = %v, want %v", got, tt.want)
			}
		})
	}
}
