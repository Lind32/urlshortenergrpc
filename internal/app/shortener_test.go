package shortener

import (
	"testing"
)

func TestValidURL(t *testing.T) {
	type args struct {
		tocen string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "ValidURL_true",
			args: args{tocen: "https://www.youtube.com/"},
			want: true,
		},
		{
			name: "ValidURL_false1",
			args: args{tocen: "wwW.youtube.com/"},
			want: false,
		},
		{
			name: "ValidURL_false2",
			args: args{tocen: "http://"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ValidURL(tt.args.tocen); got != tt.want {
				t.Errorf("ValidURL() = %v, want %v", got, tt.want)
			}
		})
	}
}
