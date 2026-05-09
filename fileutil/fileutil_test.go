package fileutil

import "testing"

func TestReplaceExtension(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		newExt   string
		want     string
	}{
		{
			name:     "json to xlsx",
			filename: "spec.json",
			newExt:   ".xlsx",
			want:     "spec.xlsx",
		},
		{
			name:     "yaml to xlsx",
			filename: "maturity.yaml",
			newExt:   ".xlsx",
			want:     "maturity.xlsx",
		},
		{
			name:     "no extension",
			filename: "myfile",
			newExt:   ".xlsx",
			want:     "myfile.xlsx",
		},
		{
			name:     "path with directories",
			filename: "/path/to/spec.json",
			newExt:   ".xlsx",
			want:     "/path/to/spec.xlsx",
		},
		{
			name:     "relative path",
			filename: "../data/spec.json",
			newExt:   ".xlsx",
			want:     "../data/spec.xlsx",
		},
		{
			name:     "double extension",
			filename: "spec.backup.json",
			newExt:   ".xlsx",
			want:     "spec.backup.xlsx",
		},
		{
			name:     "hidden file with extension",
			filename: ".config.json",
			newExt:   ".xlsx",
			want:     ".config.xlsx",
		},
		{
			name:     "hidden file no extension",
			filename: ".gitignore",
			newExt:   ".bak",
			want:     ".bak", // filepath.Ext(".gitignore") = ".gitignore", so it gets replaced
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ReplaceExtension(tt.filename, tt.newExt)
			if got != tt.want {
				t.Errorf("ReplaceExtension(%q, %q) = %q, want %q",
					tt.filename, tt.newExt, got, tt.want)
			}
		})
	}
}
