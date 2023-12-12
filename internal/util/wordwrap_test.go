// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
