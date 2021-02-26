package assert

import (
	"context"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"
)

func mustBeFunction(v interface{}) {
	if reflect.TypeOf(v).Kind() != reflect.Func {
		panic("Value is not a function")
	}
}

func describe(v interface{}) string {
	value := reflect.ValueOf(v)

	if v == nil {
		return "[nil] nil"
	}

	return fmt.Sprintf("[%s] %v", value.Type(), v)
}

//Fail the current test case with given message
func Fail(msg string, args ...interface{}) {
	if currentT == nil {
		panic("Did you forget to call RegisterT(t)?")
	}
	currentT.Errorf(msg, args...)
}

var currentT *testing.T
var envVariables map[string]string

//RegisterT saves current testing.T for further usage by Expect
func RegisterT(t *testing.T) {
	if currentT == nil {
		copyEnv()
	}

	currentT = t
	restartEnv()
}

func copyEnv() {
	envVariables = make(map[string]string)
	for _, e := range os.Environ() {
		key := strings.Split(e, "=")[0]
		envVariables[key] = os.Getenv(key)
	}
}

func restartEnv() {
	for k, v := range envVariables {
		os.Setenv(k, v)
	}
}

// AnyAssertions is used to assert any kind of value
type AnyAssertions struct {
	actual interface{}
}

// Expect starts new assertions on given value
func Expect(actual interface{}) *AnyAssertions {
	if currentT == nil {
		panic("Did you forget to call RegisterT(t)?")
	}
	return &AnyAssertions{
		actual: actual,
	}
}

// Equals asserts that actual value equals expected value
func (a *AnyAssertions) Equals(expected interface{}) bool {
	if reflect.DeepEqual(expected, a.actual) {
		return true
	}
	err := fmt.Errorf("Equals assertion failed. \n Expected: \n\t\t %s\n Actual: \n\t\t %s", describe(expected), describe(a.actual))
	currentT.Error(err)
	return false
}

// ContainsString asserts that actual value contains given string
func (a *AnyAssertions) ContainsString(substr string) bool {
	if strings.Contains(a.actual.(string), substr) {
		return true
	}
	err := fmt.Errorf("ContainsString assertion failed. \n String: \n\t\t %s\n Actual: \n\t\t %s", substr, describe(a.actual))
	currentT.Error(err)
	return false
}

// NotEquals asserts that actual value is different than given value
func (a *AnyAssertions) NotEquals(other interface{}) bool {
	if !reflect.DeepEqual(other, a.actual) {
		return true
	}
	err := fmt.Errorf("NotEquals assertion failed. \n Other: \n\t\t %s\n Actual: \n\t\t %s", describe(other), describe(a.actual))
	currentT.Error(err)
	return false
}

// IsTrue asserts that actual value is true
func (a *AnyAssertions) IsTrue() bool {
	return a.Equals(true)
}

// IsFalse asserts that actual value is false
func (a *AnyAssertions) IsFalse() bool {
	return a.Equals(false)
}

// IsEmpty asserts that actual value is an empty string
func (a *AnyAssertions) IsEmpty() bool {
	return a.Equals("")
}

// IsNotEmpty asserts that actual value is not an empty string
func (a *AnyAssertions) IsNotEmpty() bool {
	return a.NotEquals("")
}

// IsNotNil asserts that actual value is not nil
func (a *AnyAssertions) IsNotNil() bool {
	if a.actual != nil && !reflect.ValueOf(a.actual).IsNil() {
		return true
	}
	err := fmt.Errorf("IsNotNil assertion failed. \n Actual: \n\t\t %v", a.actual)
	currentT.Error(err)
	return false
}

// IsNil asserts that actual value is nil
func (a *AnyAssertions) IsNil() bool {
	if a.actual == nil || reflect.ValueOf(a.actual).IsNil() {
		return true
	}
	err := fmt.Errorf("IsNil assertion failed. \n Actual: \n\t\t %v", a.actual)
	currentT.Error(err)
	return false
}

// HasLen asserts that actual value has an expected length
func (a *AnyAssertions) HasLen(expected int) bool {
	length := reflect.ValueOf(a.actual).Len()
	if expected == length {
		return true
	}
	err := fmt.Errorf("HasLen assertion failed. \n Expected: \n\t\t %d \n Actual: \n\t\t %d", expected, length)
	currentT.Error(err)
	return false
}

// Panics asserts that actual value panics whenever called
func (a *AnyAssertions) Panics() (panicked bool) {
	mustBeFunction(a.actual)
	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()
	reflect.ValueOf(a.actual).Call([]reflect.Value{})
	if !panicked {
		err := errors.New("Panics assertion failed. \n Given function didn't panic")
		currentT.Error(err)
	}
	return
}

// EventuallyEquals asserts that, within 30 seconds, the actual function will return same value as expected value
func (a *AnyAssertions) EventuallyEquals(expected interface{}) bool {
	mustBeFunction(a.actual)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		values := reflect.ValueOf(a.actual).Call([]reflect.Value{})
		if reflect.DeepEqual(expected, values[0].Interface()) {
			return true
		}
		select {
		case <-ctx.Done():
			err := fmt.Errorf("EventuallyEquals assertion failed. \n Expected: \n\t\t %s \n Actual: \n\t\t %s", describe(expected), describe(a.actual))
			currentT.Error(err)
			return false
		case <-ticker.C:
		}
	}
}

// WithinTime asserts that actual value is between a range of other time value
func (a *AnyAssertions) WithinTime(other time.Time, diff time.Duration) bool {
	var t time.Time

	if actual, ok := a.actual.(*time.Time); ok {
		if actual == nil {
			panic("Value is nil")
		}

		t = *actual
	} else if actual, ok := a.actual.(time.Time); ok {
		t = actual
	} else {
		panic("Value is not a time")
	}

	upperBound := other.Add(diff)
	lowerBound := other.Add(diff * -1)
	if t.After(lowerBound) && t.Before(upperBound) {
		return true
	}

	err := fmt.Errorf("WithinTime assertion failed. \n Range: \n\t\t %s ~ %s \n Actual: \n\t\t %s", lowerBound, upperBound, t)
	currentT.Error(err)
	return false
}
