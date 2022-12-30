package command

import "testing"

func Test_joinPath(t *testing.T) {
	type args struct {
		cur  string
		path string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"1", args{"/a", "/b"}, "/b"},
		{"2", args{"/a", "/b/c"}, "/b/c"},
		{"3", args{"/a", "b/c"}, "/a/b/c"},
		{"4", args{"a", "b/c"}, "/a/b/c"},
		{"5", args{"a", "b/c/"}, "/a/b/c"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := joinPath(tt.args.cur, tt.args.path); got != tt.want {
				t.Errorf("joinPath() = %v, want %v", got, tt.want)
			}
		})
	}
}
