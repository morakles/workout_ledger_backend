package postgres

import "testing"

func TestIsUniqueViolation(t *testing.T) {
	cases := []struct {
		name    string
		message string
		want    bool
	}{
		{name: "duplicate key", message: "ERROR: duplicate key value violates unique constraint", want: true},
		{name: "unique constraint", message: "unique constraint violated", want: true},
		{name: "violates unique constraint", message: "violates UNIQUE constraint", want: true},
		{name: "other error", message: "some other database error", want: false},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := &fakeError{msg: tc.message}
			if got := isUniqueViolation(err); got != tc.want {
				t.Fatalf("expected %v for %q, got %v", tc.want, tc.message, got)
			}
		})
	}
}

type fakeError struct {
	msg string
}

func (f *fakeError) Error() string {
	return f.msg
}
