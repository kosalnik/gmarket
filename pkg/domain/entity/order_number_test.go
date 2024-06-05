package entity

import "testing"

func TestOrderNumber_Valid(t *testing.T) {
	tests := map[string]struct {
		b    OrderNumber
		want bool
	}{
		"positive": {
			b:    OrderNumber("5062821234567892"),
			want: true,
		},
		"positive2": {
			b:    OrderNumber("18"),
			want: true,
		},
		"negative empty": {
			b:    OrderNumber(""),
			want: false,
		},
		"negative one digit": {
			b:    OrderNumber("1"),
			want: false,
		},
		"negative no number": {
			b:    OrderNumber("asdf123"),
			want: false,
		},
		"negative": {
			b:    OrderNumber("5062821234567893"),
			want: false,
		},
		"negative2": {
			b:    OrderNumber("5062821734567892"),
			want: false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tt.b.Valid(); got != tt.want {
				t.Errorf("Valid() = %v, want %v", got, tt.want)
			}
		})
	}
}
