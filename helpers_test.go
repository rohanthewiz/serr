package serr

import "testing"

func TestArrayToString(t *testing.T) {
	type args struct {
		strArr   []string
		delim    string
		limit    int
		listName string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{name: "Empty list",
			args: args{
				strArr:   []string{},
				delim:    ", ",
				limit:    0,
				listName: "Animals",
			},
			wantOut: "0 Animals",
		},
		{name: "Short list with no limit",
			args: args{
				strArr:   []string{"cat", "dog", "mouse"},
				delim:    ", ",
				limit:    0,
				listName: "Animals",
			},
			wantOut: "cat, dog, mouse",
		},
		{name: "List with limit",
			args: args{
				strArr: []string{"cat", "dog", "mouse", "horse", "mule", "donkey", "zebra",
					"lion", "dog", "mouse", "horse", "mule", "donkey", "zebra"},
				delim:    ", ",
				limit:    5,
				listName: "Animals",
			},
			wantOut: "14 Animals: cat, dog, mouse, horse, mule...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := ArrayToString(tt.args.strArr, tt.args.delim, tt.args.limit, tt.args.listName); gotOut != tt.wantOut {
				t.Errorf("ArrayToString() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
