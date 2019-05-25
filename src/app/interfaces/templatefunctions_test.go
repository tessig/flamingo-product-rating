package interfaces_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/tessig/flamingo-product-rating/src/app/interfaces"
)

func TestRandomIntFunc_Func(t *testing.T) {
	type fields struct {
		DebugMode bool
	}
	type args struct {
		min int
		max int
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "debug mode",
			fields: fields{
				DebugMode: true,
			},
			args: args{
				min: 7,
				max: 15,
			},
		},
		{
			name: "debug mode off",
			fields: fields{
				DebugMode: false,
			},
			args: args{
				min: 7,
				max: 15,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &interfaces.RandomIntFunc{}
			f.Inject(
				&struct {
					DebugMode bool `inject:"config:debug.mode"`
				}{
					DebugMode: tt.fields.DebugMode,
				},
			)
			tf := f.Func(context.Background()).(func(int, int) int)
			for i := 0; i < 100; i++ {
				if got := tf(tt.args.min, tt.args.max); got < tt.args.min || got >= tt.args.max {
					t.Errorf("%d: ForFunc.Func() = %v, not between %d and %d", i, got, tt.args.min, tt.args.max)
				}
			}
		})
	}
}

func TestForFunc_Func(t *testing.T) {
	type args struct {
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "simple case",
			args: args{
				start: 1,
				end:   7,
			},
			want: []string{"1", "2", "3", "4", "5", "6", "7"},
		},
		{
			name: "negative case",
			args: args{
				start: -2,
				end:   3,
			},
			want: []string{"-2", "-1", "0", "1", "2", "3"},
		},
		{
			name: "empty case",
			args: args{
				start: 2,
				end:   2,
			},
			want: []string{"2"},
		},
		{
			name: "impossible case",
			args: args{
				start: 5,
				end:   2,
			},
			want: []string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &interfaces.ForFunc{}
			tf := f.Func(context.Background()).(func(int, int) []string)
			if got := tf(tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBarTypeFunc_Func(t *testing.T) {
	type args struct {
		stars int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "0",
			args: args{
				stars: 0,
			},
			want: "success",
		},
		{
			name: "1",
			args: args{
				stars: 1,
			},
			want: "danger",
		},
		{
			name: "2",
			args: args{
				stars: 2,
			},
			want: "warning",
		},
		{
			name: "3",
			args: args{
				stars: 3,
			},
			want: "info",
		},
		{
			name: "4",
			args: args{
				stars: 4,
			},
			want: "primary",
		},
		{
			name: "5",
			args: args{
				stars: 5,
			},
			want: "success",
		},
		{
			name: "711",
			args: args{
				stars: 711,
			},
			want: "success",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &interfaces.BarTypeFunc{}
			tf := f.Func(context.Background()).(func(int) string)
			if got := tf(tt.args.stars); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ForFunc.Func() = %v, want %v", got, tt.want)
			}
		})
	}
}
