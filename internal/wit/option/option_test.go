package option

import (
	"testing"

	"go.bytecodealliance.org/cm"
)

func TestHandle(t *testing.T) {
	tests := map[string]struct {
		option    cm.Option[int]
		wantValue int
		wantOk    bool
	}{
		"Some value": {
			option:    cm.Some(42),
			wantValue: 42,
			wantOk:    true,
		},
		"None value": {
			option:    cm.None[int](),
			wantValue: 0, // zero value for int
			wantOk:    false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			value, ok := Handle(tt.option)
			if value != tt.wantValue || ok != tt.wantOk {
				t.Errorf("Handle() = (%v, %v), want (%v, %v)", value, ok, tt.wantValue, tt.wantOk)
			}
		})
	}
}

func TestUnwrap(t *testing.T) {
	tests := map[string]struct {
		option    cm.Option[int]
		wantValue int
		wantPanic bool
	}{
		"Some value": {
			option:    cm.Some(42),
			wantValue: 42,
			wantPanic: false,
		},
		"None value": {
			option:    cm.None[int](),
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
			value := Unwrap(tt.option)
			if value != tt.wantValue {
				t.Errorf("Unwrap() = %v, want %v", value, tt.wantValue)
			}
		})
	}
}

func TestUnwrapOr(t *testing.T) {
	tests := map[string]struct {
		option       cm.Option[int]
		defaultValue int
		wantValue    int
	}{
		"Some value": {
			option:       cm.Some(42),
			defaultValue: 0,
			wantValue:    42,
		},
		"None value with default": {
			option:       cm.None[int](),
			defaultValue: 100,
			wantValue:    100,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			value := UnwrapOr(tt.option, tt.defaultValue)
			if value != tt.wantValue {
				t.Errorf("UnwrapOr() = %v, want %v", value, tt.wantValue)
			}
		})
	}
}
