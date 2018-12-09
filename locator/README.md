# Service Locator

Implements a very simple service locator that does allows for modular initialization with a deterministic init order.

A module has a simple 4 phase lifecycle: Init, Start, Reload and Stop. You can enable or disable a feature in the config.
This hooks into the watching infrastructure, so you can also enable or disable modules by just editing config or changing a remote value.

Name | Description
-----|------------
Init | Called on initial creation of the module
Start | Called when the module is started, or enabled at runtime
Reload | Called when the config has changed and the module needs to reconfigure itself
Stop | Called when the module is stopped, or disabled at runtime

Each module is identified by a unique name, this defaults to its package name,

## Usage

To use it, a package that serves as a module needs to export a method or variable that implements the Module interface.

```go
package orders

import "github.com/casualjim/go-appserver/locator"

var Module = locator.MakeModule(
  locator.Init(func(app locator.Application) error {
    orders := new(ordersService)
    locator.Set("ordersService", orders)
    orders.app = app
    return nil
  }),
  locator.Reload(func(app locator.Application) error {
    // you can reconfigure the services that belong to this module here
    return nil
  })
)

type Order struct {
  ID      int64
  Product int64
}

type odersService struct {
  app locator.Application
}

func (o *ordersService) Create(o *Order) error {
  var db OrdersStore
  o.app.Get("ordersDb", &db)
  return db.Save(o)
}
```

In the main package you would then write a main function that could look like this:

```go
func main() {
  app := locator.New()
  app.RegisterModule(orders.Module)

  if err := app.Init(); err != nil {
    app.Logger().Fatalln(err)
  }

  app.Logger().Infoln("application initialized, starting...")

  if err := app.Start(); err != nil {
    app.Logger().Fatalln(err)
  }

  app.Logger().Infoln("application initialized, starting...")
  // do a blocking operation here, like run a http server

  if err := app.Stop(); err != nil {
    app.Logger().Fatalln(err)
  }
}
```
