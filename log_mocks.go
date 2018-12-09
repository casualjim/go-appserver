// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package appserver

import (
	"sync"
	"github.com/go-logr/logr"
)

var (
	lockLoggerMockEnabled    sync.RWMutex
	lockLoggerMockError      sync.RWMutex
	lockLoggerMockInfo       sync.RWMutex
	lockLoggerMockV          sync.RWMutex
	lockLoggerMockWithName   sync.RWMutex
	lockLoggerMockWithValues sync.RWMutex
)

// LoggerMock is a mock implementation of Logger.
//
//     func TestSomethingThatUsesLogger(t *testing.T) {
//
//         // make and configure a mocked Logger
//         mockedLogger := &LoggerMock{
//             EnabledFunc: func() bool {
// 	               panic("mock out the Enabled method")
//             },
//             ErrorFunc: func(err error, msg string, keysAndValues ...interface{})  {
// 	               panic("mock out the Error method")
//             },
//             InfoFunc: func(msg string, keysAndValues ...interface{})  {
// 	               panic("mock out the Info method")
//             },
//             VFunc: func(level int) logr.InfoLogger {
// 	               panic("mock out the V method")
//             },
//             WithNameFunc: func(name string) Logger {
// 	               panic("mock out the WithName method")
//             },
//             WithValuesFunc: func(keysAndValues ...interface{}) Logger {
// 	               panic("mock out the WithValues method")
//             },
//         }
//
//         // use mockedLogger in code that requires Logger
//         // and then make assertions.
//
//     }
type LoggerMock struct {
	// EnabledFunc mocks the Enabled method.
	EnabledFunc func() bool

	// ErrorFunc mocks the Error method.
	ErrorFunc func(err error, msg string, keysAndValues ...interface{})

	// InfoFunc mocks the Info method.
	InfoFunc func(msg string, keysAndValues ...interface{})

	// VFunc mocks the V method.
	VFunc func(level int) logr.InfoLogger

	// WithNameFunc mocks the WithName method.
	WithNameFunc func(name string) logr.Logger

	// WithValuesFunc mocks the WithValues method.
	WithValuesFunc func(keysAndValues ...interface{}) logr.Logger

	// calls tracks calls to the methods.
	calls struct {
		// Enabled holds details about calls to the Enabled method.
		Enabled []struct {
		}
		// Error holds details about calls to the Error method.
		Error []struct {
			// Err is the err argument value.
			Err error
			// Msg is the msg argument value.
			Msg string
			// KeysAndValues is the keysAndValues argument value.
			KeysAndValues []interface{}
		}
		// Info holds details about calls to the Info method.
		Info []struct {
			// Msg is the msg argument value.
			Msg string
			// KeysAndValues is the keysAndValues argument value.
			KeysAndValues []interface{}
		}
		// V holds details about calls to the V method.
		V []struct {
			// Level is the level argument value.
			Level int
		}
		// WithName holds details about calls to the WithName method.
		WithName []struct {
			// Name is the name argument value.
			Name string
		}
		// WithValues holds details about calls to the WithValues method.
		WithValues []struct {
			// KeysAndValues is the keysAndValues argument value.
			KeysAndValues []interface{}
		}
	}
}

// Enabled calls EnabledFunc.
func (mock *LoggerMock) Enabled() bool {
	if mock.EnabledFunc == nil {
		panic("LoggerMock.EnabledFunc: method is nil but Logger.Enabled was just called")
	}
	callInfo := struct {
	}{}
	lockLoggerMockEnabled.Lock()
	mock.calls.Enabled = append(mock.calls.Enabled, callInfo)
	lockLoggerMockEnabled.Unlock()
	return mock.EnabledFunc()
}

// EnabledCalls gets all the calls that were made to Enabled.
// Check the length with:
//     len(mockedLogger.EnabledCalls())
func (mock *LoggerMock) EnabledCalls() []struct {
} {
	var calls []struct {
	}
	lockLoggerMockEnabled.RLock()
	calls = mock.calls.Enabled
	lockLoggerMockEnabled.RUnlock()
	return calls
}

// Error calls ErrorFunc.
func (mock *LoggerMock) Error(err error, msg string, keysAndValues ...interface{}) {
	if mock.ErrorFunc == nil {
		panic("LoggerMock.ErrorFunc: method is nil but Logger.Error was just called")
	}
	callInfo := struct {
		Err           error
		Msg           string
		KeysAndValues []interface{}
	}{
		Err:           err,
		Msg:           msg,
		KeysAndValues: keysAndValues,
	}
	lockLoggerMockError.Lock()
	mock.calls.Error = append(mock.calls.Error, callInfo)
	lockLoggerMockError.Unlock()
	mock.ErrorFunc(err, msg, keysAndValues...)
}

// ErrorCalls gets all the calls that were made to Error.
// Check the length with:
//     len(mockedLogger.ErrorCalls())
func (mock *LoggerMock) ErrorCalls() []struct {
	Err           error
	Msg           string
	KeysAndValues []interface{}
} {
	var calls []struct {
		Err           error
		Msg           string
		KeysAndValues []interface{}
	}
	lockLoggerMockError.RLock()
	calls = mock.calls.Error
	lockLoggerMockError.RUnlock()
	return calls
}

// Info calls InfoFunc.
func (mock *LoggerMock) Info(msg string, keysAndValues ...interface{}) {
	if mock.InfoFunc == nil {
		panic("LoggerMock.InfoFunc: method is nil but Logger.Info was just called")
	}
	callInfo := struct {
		Msg           string
		KeysAndValues []interface{}
	}{
		Msg:           msg,
		KeysAndValues: keysAndValues,
	}
	lockLoggerMockInfo.Lock()
	mock.calls.Info = append(mock.calls.Info, callInfo)
	lockLoggerMockInfo.Unlock()
	mock.InfoFunc(msg, keysAndValues...)
}

// InfoCalls gets all the calls that were made to Info.
// Check the length with:
//     len(mockedLogger.InfoCalls())
func (mock *LoggerMock) InfoCalls() []struct {
	Msg           string
	KeysAndValues []interface{}
} {
	var calls []struct {
		Msg           string
		KeysAndValues []interface{}
	}
	lockLoggerMockInfo.RLock()
	calls = mock.calls.Info
	lockLoggerMockInfo.RUnlock()
	return calls
}

// V calls VFunc.
func (mock *LoggerMock) V(level int) logr.InfoLogger {
	if mock.VFunc == nil {
		panic("LoggerMock.VFunc: method is nil but Logger.V was just called")
	}
	callInfo := struct {
		Level int
	}{
		Level: level,
	}
	lockLoggerMockV.Lock()
	mock.calls.V = append(mock.calls.V, callInfo)
	lockLoggerMockV.Unlock()
	return mock.VFunc(level)
}

// VCalls gets all the calls that were made to V.
// Check the length with:
//     len(mockedLogger.VCalls())
func (mock *LoggerMock) VCalls() []struct {
	Level int
} {
	var calls []struct {
		Level int
	}
	lockLoggerMockV.RLock()
	calls = mock.calls.V
	lockLoggerMockV.RUnlock()
	return calls
}

// WithName calls WithNameFunc.
func (mock *LoggerMock) WithName(name string) logr.Logger {
	if mock.WithNameFunc == nil {
		panic("LoggerMock.WithNameFunc: method is nil but Logger.WithName was just called")
	}
	callInfo := struct {
		Name string
	}{
		Name: name,
	}
	lockLoggerMockWithName.Lock()
	mock.calls.WithName = append(mock.calls.WithName, callInfo)
	lockLoggerMockWithName.Unlock()
	return mock.WithNameFunc(name)
}

// WithNameCalls gets all the calls that were made to WithName.
// Check the length with:
//     len(mockedLogger.WithNameCalls())
func (mock *LoggerMock) WithNameCalls() []struct {
	Name string
} {
	var calls []struct {
		Name string
	}
	lockLoggerMockWithName.RLock()
	calls = mock.calls.WithName
	lockLoggerMockWithName.RUnlock()
	return calls
}

// WithValues calls WithValuesFunc.
func (mock *LoggerMock) WithValues(keysAndValues ...interface{}) logr.Logger {
	if mock.WithValuesFunc == nil {
		panic("LoggerMock.WithValuesFunc: method is nil but Logger.WithValues was just called")
	}
	callInfo := struct {
		KeysAndValues []interface{}
	}{
		KeysAndValues: keysAndValues,
	}
	lockLoggerMockWithValues.Lock()
	mock.calls.WithValues = append(mock.calls.WithValues, callInfo)
	lockLoggerMockWithValues.Unlock()
	return mock.WithValuesFunc(keysAndValues...)
}

// WithValuesCalls gets all the calls that were made to WithValues.
// Check the length with:
//     len(mockedLogger.WithValuesCalls())
func (mock *LoggerMock) WithValuesCalls() []struct {
	KeysAndValues []interface{}
} {
	var calls []struct {
		KeysAndValues []interface{}
	}
	lockLoggerMockWithValues.RLock()
	calls = mock.calls.WithValues
	lockLoggerMockWithValues.RUnlock()
	return calls
}
