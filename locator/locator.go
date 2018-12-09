package locator

import (
	"github.com/hashicorp/go-multierror"
	"reflect"
	"sync"
)

type errString string

func (e errString) Error() string {
	return string(e)
}

const (
	// ErrModuleUnknown returned when no module can be found for the specified key
	ErrModuleUnknown errString = "unknown module"
	ErrOnlyPointerAllowed errString = "app can only load values into pointers"
)

// Application is an application level context package
// It can be used as a kind of dependency injection container
type Application interface {
	// RegisterModule with the application context
	RegisterModule(...Module) error

	// LocateService for the specified key, thread-safe
	LocateService(interface{}, interface{}) error

	// SetService for the specified key, this should be safe across multiple threads
	SetService(interface{}, interface{}) error

	// Init the application and its modules with the config.
	Init() error

	// Start the application and its enabled modules
	Start() error

	// Stop the application and its enabled modules
	Stop() error

	// Reload the application and its enabled modules
	Reload() error
}

// LifecycleCallback function definition
type LifecycleCallback interface {
	Call(Application) error
}

// Init is an initializer for an initialization function
type Init func(Application) error

// Call implements the callback interface
func (fn Init) Call(app Application) error {
	return fn(app)
}

// Start is an initializer for a start function
type Start func(Application) error

// Call implements the callback interface
func (fn Start) Call(app Application) error {
	return fn(app)
}

// Stop is an initializer for a stop function
type Stop func(Application) error

// Call implements the callback interface
func (fn Stop) Call(app Application) error {
	return fn(app)
}

// Reload is an initalizater for a reload function
type Reload func(Application) error

// Call implements the callback interface
func (fn Reload) Call(app Application) error {
	return fn(app)
}

// A Module is a component that has a specific lifecycle
type Module interface {
	Init(Application) error
	Start(Application) error
	Stop(Application) error
	Reload(Application) error
}

// MakeModule by passing the callback functions.
// You can pass multiple callback functions of the same type if you want
func MakeModule(callbacks ...LifecycleCallback) Module {
	var (
		init   []Init
		start  []Start
		reload []Reload
		stop   []Stop
	)

	for _, callback := range callbacks {
		switch cb := callback.(type) {
		case Init:
			init = append(init, cb)
		case Start:
			start = append(start, cb)
		case Stop:
			stop = append(stop, cb)
		case Reload:
			reload = append(reload, cb)
		}
	}

	return &dynamicModule{
		init:   init,
		start:  start,
		reload: reload,
		stop:   stop,
	}
}

type dynamicModule struct {
	init   []Init
	start  []Start
	stop   []Stop
	reload []Reload
}

func (d *dynamicModule) Init(app Application) error {
	for _, cb := range d.init {
		if err := cb.Call(app); err != nil {
			return err
		}
	}
	return nil
}

func (d *dynamicModule) Start(app Application) error {
	for _, cb := range d.start {
		if err := cb.Call(app); err != nil {
			return err
		}
	}
	return nil
}

func (d *dynamicModule) Stop(app Application) error {
	for _, cb := range d.stop {
		if err := cb.Call(app); err != nil {
			return err
		}
	}
	return nil
}

func (d *dynamicModule) Reload(app Application) error {
	for _, cb := range d.reload {
		if err := cb.Call(app); err != nil {
			return err
		}
	}
	return nil
}


// New creates an application context
func New() Application {
	return &defaultApplication{}
}

type defaultApplication struct {
	modules []Module

	registry sync.Map
}

func (d *defaultApplication) RegisterModule(modules ...Module) error {
	d.modules = append(d.modules, modules...)
	return nil
}

func (d *defaultApplication) LocateService(key interface{}, mod interface{}) error {
	rv := reflect.ValueOf(mod)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return ErrOnlyPointerAllowed
	}

	val, ok := d.registry.Load(key)
	if !ok {
		return ErrModuleUnknown
	}

	vv := reflect.Indirect(reflect.ValueOf(val))
	reflect.Indirect(rv).Set(vv)
	return nil
}

func (d *defaultApplication) SetService(key interface{}, module interface{}) error {
	d.registry.Store(key, module)
	return nil
}

func (d *defaultApplication) Init() error {
	for _, mod := range d.modules {
		if err := mod.Init(d); err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultApplication) Start() error {
	for _, mod := range d.modules {
		if err := mod.Start(d); err != nil {
			return err
		}
	}
	return nil
}

func (d *defaultApplication) Stop() error {
	// We stop in reverse order to respect dependency chains
	var result *multierror.Error
	for i := len(d.modules); i > 0; i-- {
		if err := d.modules[i-1].Stop(d); err != nil {
			result = multierror.Append(result, err)
		}
	}
	return result.ErrorOrNil()
}

func (d *defaultApplication) Reload() error {
	for _, mod := range d.modules {
		if err := mod.Reload(d); err != nil {
			return err
		}
	}
	return nil
}
