package targettracking

import (
	"testing"
	"time"

	rebalancev1 "git.pepabo.com/akichan/rebalancer/api/v1"
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

func TestCheckScheduledValue(t *testing.T) {
	type args struct {
		scheduled []rebalancev1.Scheduled
		value     int64
		nowTime   time.Time
	}

	tests := []struct {
		name string
		args args
		want int64
	}{
		{
			"within scheduled",
			args{
				[]rebalancev1.Scheduled{
					{
						StartTime: "19:00",
						EndTime:   "23:30",
						Value:     2,
					},
				},
				0,
				time.Date(2022, 12, 20, 19, 0, 0, 0, time.Local),
			},
			2,
		},
		{
			"without scheduled",
			args{
				[]rebalancev1.Scheduled{
					{
						StartTime: "19:00",
						EndTime:   "23:30",
						Value:     2,
					},
				},
				0,
				time.Date(2022, 12, 20, 12, 0, 0, 0, time.Local),
			},
			0,
		},
		{
			"have multiple sheduled",
			args{
				[]rebalancev1.Scheduled{
					{
						StartTime: "18:00",
						EndTime:   "23:30",
						Value:     2,
					},
					{
						StartTime: "19:00",
						EndTime:   "20:00",
						Value:     3,
					},
				},
				1,
				time.Date(2022, 12, 20, 19, 0, 0, 0, time.Local),
			},
			3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := checkScheduledValue(tt.args.scheduled, tt.args.value, tt.args.nowTime); got != tt.want {
				t.Errorf("checkScheduledValue() = %v, want %v", got, tt.want)
			}
		})
	}

}
