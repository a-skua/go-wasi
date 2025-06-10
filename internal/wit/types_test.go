package wit

import (
	"go.bytecodealliance.org/cm"
	"testing"
)

func TestHandleResult(t *testing.T) {
	t.Run("Success case", func(t *testing.T) {
		result := cm.OK[cm.Result[string, string, struct{}]]("success value")
		ok, err := HandleResult(result)

		if ok == nil {
			t.Error("Expected OK value, got nil")
		}
		if *ok != "success value" {
			t.Errorf("Expected 'success value', got %s", *ok)
		}
		if err != nil {
			t.Error("Expected nil error, got non-nil")
		}
	})

	t.Run("Error case", func(t *testing.T) {
		result := cm.Err[cm.Result[string, string, string]]("error message")
		ok, err := HandleResult(result)

		if ok != nil {
			t.Error("Expected nil OK value, got non-nil")
		}
		if err == nil {
			t.Error("Expected error value, got nil")
		}
		if *err != "error message" {
			t.Errorf("Expected 'error message', got %s", *err)
		}
	})
}
