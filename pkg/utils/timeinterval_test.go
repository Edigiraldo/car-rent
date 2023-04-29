package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsValidTimeFrame(t *testing.T) {
	now := time.Now()

	type args struct {
		start time.Time
		end   time.Time
	}
	type wants struct {
		isValid bool
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "returns false when start and end times are equal",
			args: args{
				start: now,
				end:   now,
			},
			wants: wants{
				isValid: false,
			},
		},
		{
			name: "eturns false when the end time is before the start",
			args: args{
				start: now.Add(time.Second),
				end:   now,
			},
			wants: wants{
				isValid: false,
			},
		},
		{
			name: "returns true when the end time is after the start",
			args: args{
				start: now,
				end:   now.Add(time.Second),
			},
			wants: wants{
				isValid: true,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isValid := IsValidTimeFrame(test.args.start, test.args.end)

			assert.Equal(t, test.wants.isValid, isValid)
		})
	}
}
