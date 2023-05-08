package cmd

import (
	"testing"
	"time"
)

func TestTotalWorkTime(t *testing.T) {
	type args struct {
		workTimes map[string]time.Duration
	}
	tests := []struct {
		name string
		args args
		want time.Duration
	}{
		{
			name: "When workTimes does not have work times",
			args: args{
				workTimes: map[string]time.Duration{},
			},
			want: 0 * time.Minute,
		},
		{
			name: "When workTimes have three work times",
			args: args{
				workTimes: map[string]time.Duration{
					"/path/to/directory1": 3 * time.Hour,
					"/path/to/directory2": 20 * time.Minute,
					"/path/to/directory3": 45 * time.Minute,
				},
			},
			want: 4*time.Hour + 5*time.Minute,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := TotalWorkTime(tt.args.workTimes); got != tt.want {
				t.Errorf("TotalWorkTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
