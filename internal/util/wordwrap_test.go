package util

import (
	"testing"
)

func TestWrapString(t *testing.T) {
	type args struct {
		s   string
		lim int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "中文测试",
			args: args{
				s:   "你好啊",
				lim: 1,
			},
			want: `你
好
啊
`,
		},
		{
			name: "中文测试",
			args: args{
				s:   "我喜欢 TiDB",
				lim: 4,
			},
			want: `我喜
欢
TiDB`,
		},
		{
			name: "english test without break line",
			args: args{
				s:   "hello",
				lim: 1,
			},
			want: "hello",
		},
		{
			name: "english test with break line",
			args: args{
				s:   "hello world",
				lim: 1,
			},
			want: "hello " +
				"world",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := String(tt.args.s, tt.args.lim); got != tt.want {
				t.Errorf("WrapString() = %v, want %v", got, tt.want)
			}
		})
	}
}
