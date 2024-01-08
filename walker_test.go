package generator

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindModelSourcePath(t *testing.T) {
	type args struct {
		model string
		dir   string
	}

	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "model source file not found",
			args: args{
				model: "Block",
				dir:   "testdata",
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "model source file found",
			args: args{
				model: "Block",
				dir:   "testdata/foo",
			},
			want:    "testdata/foo/block.go",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			got, err := FindModelSourcePath(tt.args.model, tt.args.dir)
			if !tt.wantErr {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}

			require.Equal(t, tt.want, got)
		})
	}
}
