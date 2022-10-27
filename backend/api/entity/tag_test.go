package entity

import "testing"

func TestTagsUnique(t *testing.T) {
	tests := []struct {
		name   string
		target Tags
		want   Tags
	}{
		{
			name:   "unique Tags",
			target: Tags{{Name: "a"}, {Name: "b"}, {Name: "c"}},
			want:   Tags{{Name: "a"}, {Name: "b"}, {Name: "c"}},
		},
		{
			name:   "not unique Tags",
			target: Tags{{Name: "a"}, {Name: "b"}, {Name: "b"}},
			want:   Tags{{Name: "a"}, {Name: "b"}},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.target.Unique()
			if len(got) != len(tt.want) {
				t.Errorf("want %v, but %v", tt.want, got)
				return
			}
			for i := range got {
				if got[i].Name != tt.want[i].Name {
					t.Errorf("want %v, but %v", tt.want, got)
					return
				}
			}
		})
	}
}
