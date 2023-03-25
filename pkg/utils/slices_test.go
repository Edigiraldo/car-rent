package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsInSlice(t *testing.T) {
	type args struct {
		s     []interface{}
		value interface{}
	}
	type wants struct {
		isInSlice bool
	}
	tests := []struct {
		name  string
		args  args
		wants wants
	}{
		{
			name: "returns true when character 'g' is in slice",
			args: args{
				s:     []interface{}{"a", "s", "d", "f", "g", "h", "j", "k", "l", ";"},
				value: "g",
			},
			wants: wants{
				isInSlice: true,
			},
		},
		{
			name: "returns false when number 3 is not in slice",
			args: args{
				s:     []interface{}{21, 6, 60, 8, 2, 0, 8, 4, 98, 4},
				value: 3,
			},
			wants: wants{
				isInSlice: false,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			isInSlice := IsInSlice(test.args.s, test.args.value)

			assert.Equal(t, test.wants.isInSlice, isInSlice)
		})
	}
}
