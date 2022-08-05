package regex

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_splitLabelPrefixAndName(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 string
	}{
		{"case-1", args{"aaa/bbb"}, "aaa", "bbb"},
		{"case-2", args{"aaa"}, "", "aaa"},
		{"case-3", args{"/"}, "", ""},
		{"case-4", args{"aaa/"}, "aaa", ""},
		{"case-5", args{"/aaa"}, "", "aaa"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := splitLabelPrefixAndName(tt.args.key)
			assert.Equalf(t, tt.want, got, "splitLabelPrefixAndName(%v)", tt.args.key)
			assert.Equalf(t, tt.want1, got1, "splitLabelPrefixAndName(%v)", tt.args.key)
		})
	}
}

func Test_checkLabel(t *testing.T) {
	type args struct {
		prefix string
		name   string
		value  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"invalid-prefix-1", args{"k8s.io", "test-name", "test-value"}, false},
		{"invalid-prefix-2", args{"app.kubernetes.io-1", "test_name", "test.value"}, false},
		{"invalid-prefix-3", args{"app.k8s.io..", "test", "value"}, false},
		{"invalid-name-1", args{"app.kubernetes.io", "-test", "test-value"}, false},
		{"invalid-name-2", args{"xxx", "test-", "test-value"}, false},
		{"invalid-name-3", args{"kubelet.kubernetes.io", "test@name", "test-value"}, false},
		{"invalid-name-4", args{"kubelet.kubernetes.io", "test.", "test-value"}, false},
		{"invalid-value-1", args{"kubelet.kubernetes.io", "test-name", "test#value"}, false},
		{"valid-case-1", args{"app.kubernetes.io", "test", "test-value"}, true},
		{"valid-case-2", args{"xxx", "test_name", "test.value"}, true},
		{"valid-case-3", args{"", "test", ""}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, checkLabel(tt.args.prefix, tt.args.name, tt.args.value), "checkLabel(%v, %v, %v)", tt.args.prefix, tt.args.name, tt.args.value)
		})
	}
}
