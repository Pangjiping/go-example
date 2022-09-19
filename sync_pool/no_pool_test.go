package sync_pool

import "testing"

func TestMultiClientWithoutPool(t *testing.T) {
	type args struct {
		uid string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "no pool test",
			args: args{
				uid: "123456",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			MultiClientWithoutPool(tt.args.uid)
		})
	}
}
