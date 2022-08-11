package targettracking

import (
	"testing"
)

func TestProcessBestContrast(t *testing.T) {
	type args struct {
		base           float64
		trackingTarget float64
		current        float64
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		{"lower than target", args{10, 100, 80}, 0},
		{"same as target", args{10, 100, 100}, 0},
		{"20% above target", args{10, 100, 120}, 2},
		{"twice as high as than targeet", args{10, 100, 300}, 20},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := processBestContrast(tt.args.base, tt.args.trackingTarget, tt.args.current); got != tt.want {
				t.Errorf("processBestContrast() = %v, want %v", got, tt.want)
			}
		})
	}
}
