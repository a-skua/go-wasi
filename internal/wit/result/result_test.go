package result

import (
	"errors"
	"testing"

	"go.bytecodealliance.org/cm"
)

func TestHandleBool(t *testing.T) {
	tests := map[string]struct {
		result cm.BoolResult
		want   bool
	}{
		"True case": {
			result: cm.ResultOK,
			want:   true,
		},
		"False case": {
			result: cm.ResultErr,
			want:   false,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got := HandleBool(tt.result)
			if got != tt.want {
				t.Errorf("HandleBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHandle(t *testing.T) {
	tests := map[string]struct {
		result    cm.Result[string, int, string]
		wantValue int
		wantErr   string
	}{
		"Success case": {
			result:    cm.OK[cm.Result[string, int, string]](42),
			wantValue: 42,
			wantErr:   "",
		},
		"Error case": {
			result:    cm.Err[cm.Result[string, int, string]]("error occurred"),
			wantValue: 0, // zero value for int
			wantErr:   "error result: error occurred",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := Handle(tt.result)
			if value != tt.wantValue {
				t.Errorf("Handle() = %v, want %v", value, tt.wantValue)
			}
			if err != nil && err.Error() != tt.wantErr {
				t.Errorf("Handle() error = %v, want %v", err.Error(), tt.wantErr)
			} else if err == nil && tt.wantErr != "" {
				t.Errorf("Handle() expected error %v, got nil", tt.wantErr)
			}
		})
	}
}

func TestHandleErr(t *testing.T) {
	tests := map[string]struct {
		result    cm.Result[string, int, string]
		wantValue int
		wantErr   string
	}{
		"Success case": {
			result:    cm.OK[cm.Result[string, int, string]](42),
			wantValue: 42,
			wantErr:   "",
		},
		"Error case": {
			result:    cm.Err[cm.Result[string, int, string]]("error occurred"),
			wantValue: 0, // zero value for int
			wantErr:   "error occurred",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			value, err := HandleErr(tt.result, func(err string) error { return errors.New(err) })
			if value != tt.wantValue {
				t.Errorf("HandleErr() = %v, want %v", value, tt.wantValue)
			}
			if err != nil && err.Error() != tt.wantErr {
				t.Errorf("HandleErr() error = %v, want %v", err.Error(), tt.wantErr)
			} else if err == nil && tt.wantErr != "" {
				t.Errorf("HandleErr() expected error %v, got nil", tt.wantErr)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	tests := map[string]struct {
		result    cm.Result[string, int, string]
		wantValue int
		wantPanic bool
	}{
		"Success case": {
			result:    cm.OK[cm.Result[string, int, string]](42),
			wantValue: 42,
			wantPanic: false,
		},
		"Error case": {
			result:    cm.Err[cm.Result[string, int, string]]("error occurred"),
			wantValue: 0, // zero value for int
			wantPanic: true,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Errorf("Unwrap() did not panic")
					}
				}()
			}
			value := Unwrap(tt.result)
			if value != tt.wantValue {
				t.Errorf("Unwrap() = %v, want %v", value, tt.wantValue)
			}
		})
	}
}

func TestUnwrapOr(t *testing.T) {
	tests := map[string]struct {
		result       cm.Result[string, int, string]
		defaultValue int
		wantValue    int
	}{
		"Success case": {
			result:       cm.OK[cm.Result[string, int, string]](42),
			defaultValue: 0,
			wantValue:    42,
		},
		"Error case": {
			result:       cm.Err[cm.Result[string, int, string]]("error occurred"),
			defaultValue: 99,
			wantValue:    99,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			value := UnwrapOr(tt.result, tt.defaultValue)
			if value != tt.wantValue {
				t.Errorf("UnwrapOr() = %v, want %v", value, tt.wantValue)
			}
		})
	}
}
