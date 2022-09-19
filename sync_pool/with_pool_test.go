package sync_pool

import "testing"

func TestMultiClientWithPool(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "sync pool test",
			args: args{
				uid: "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MultiClientWithPool(tt.args.uid)
		})
	}
}
