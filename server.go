package appserver

import (
	"net/http"

	"github.com/casualjim/go-appserver/locator"
	"github.com/casualjim/go-appserver/middleware"
	httpd "github.com/casualjim/go-httpd"
	"github.com/go-chi/chi"
	"github.com/heptiolabs/healthcheck"
	"go.opencensus.io/plugin/ochttp"
	"go.opencensus.io/plugin/ochttp/propagation/b3"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
	"go.opencensus.io/zpages"
)

type Option func(*defaultServer)

// WithAdminAPI registers a new api handler for the admin routes
//noinspection GoUnusedExportedFunction
func WithAdminAPI(handler chi.Router, listeners ...httpd.ServerListener) Option {
	if len(listeners) == 0 {
		listeners = []httpd.ServerListener{
			&httpd.HTTPFlags{
				Prefix: "metrics",
				Port:   10239,
			},
		}
	}
	return func(srv *defaultServer) {
		srv.adminApp = handler
		srv.adminOpts = []httpd.Option{
			httpd.WithAdminListeners(listeners[0], listeners[1:]...),
			httpd.HandlesAdminWith(handler),
		}
	}
}

// WithHTTPDOption configures the httpd server.
//noinspection GoUnusedExportedFunction
func WithHTTPDOption(opts ...httpd.Option) Option {
	return func(srv *defaultServer) {
		srv.httpdOpts = append(srv.httpdOpts, opts...)
	}
}

// IsPublic lets the server know that it's going to be the first hop
//noinspection GoUnusedExportedFunction
func IsPublic() Option {
	return func(server *defaultServer) {
		server.isPublic = true
	}
}

// WithTraceExporter enable opencensus trace exporting
//noinspection GoUnusedExportedFunction
func WithTraceExporter(exp trace.Exporter) Option {
	return func(server *defaultServer) {
		server.tracer = exp
	}
}

// WithMetricsExporter enable opencensus metrics exporter
//noinspection GoUnusedExportedFunction
func WithMetricsExporter(exp view.Exporter) Option {
	return func(server *defaultServer) {
		server.metrics = exp
	}
}

// New creates a new app server with default config
func New(log middleware.Logger, opts ...Option) Server {
	adminApp := chi.NewRouter()
	health := healthcheck.NewHandler()

	adminApp.Use(middleware.NoCache)
	adminApp.Mount("/debug", middleware.Profiler())
	adminApp.Get("/healthz", health.LiveEndpoint)
	adminApp.Get("/readyz", health.ReadyEndpoint)
	adminApp.Get("/version", VersionHandler(log, NewVersionInfo()))

	app := chi.NewRouter()
	app.Use(
		middleware.ProxyHeaders,
		middleware.Recover(log),
	)

	srv := &defaultServer{
		Application: locator.New(),
		app:         app,
		adminApp:    adminApp,
		Handler:     health,
		adminOpts: []httpd.Option{
			httpd.WithAdminListeners(&httpd.HTTPFlags{
				Prefix: "metrics",
				Port:   10239,
			}),
			httpd.HandlesAdminWith(adminApp),
		},
		httpdOpts: []httpd.Option{
			httpd.LogsWith(log),
		},
	}
	for _, configure := range opts {
		configure(srv)
	}

	if srv.metrics != nil {
		if pe, ok := srv.metrics.(http.Handler); ok {
			srv.adminApp.Mount("/metrics", pe)
		}
		view.RegisterExporter(srv.metrics)
	}

	if srv.tracer != nil {
		trace.RegisterExporter(srv.tracer)
		srv.app.Use(
			func(next http.Handler) http.Handler {
				return &ochttp.Handler{
					Handler:          next,
					Propagation:      &b3.HTTPFormat{},
					IsPublicEndpoint: srv.isPublic,
				}
			},
			func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					route := chi.RouteContext(r.Context())
					ochttp.WithRouteTag(next, route.RoutePath).ServeHTTP(w, r)
				})
			},
		)

		mux := http.NewServeMux()
		zpages.Handle(mux, "/")
		adminApp.Mount("/rpcz", mux)
		adminApp.Mount("/tracez", mux)
		adminApp.Mount("/public/*", mux)
	}
	return srv
}

type Server interface {
	locator.Application

	// App returns the application for the http server
	App() chi.Router

	// AddLivenessCheck adds a check that indicates that this instance of the
	// application should be destroyed or restarted. A failed liveness check
	// indicates that this instance is unhealthy, not some upstream dependency.
	// Every liveness check is also included as a readiness check.
	AddLivenessCheck(name string, check healthcheck.Check)

	// AddReadinessCheck adds a check that indicates that this instance of the
	// application is currently unable to serve requests because of an upstream
	// or some transient failure. If a readiness check fails, this instance
	// should no longer receiver requests, but should not be restarted or
	// destroyed.
	AddReadinessCheck(name string, check healthcheck.Check)
}

type defaultServer struct {
	locator.Application
	healthcheck.Handler

	server    httpd.Server
	adminOpts []httpd.Option
	httpdOpts []httpd.Option
	app       chi.Router
	adminApp  chi.Router
	isPublic  bool

	tracer  trace.Exporter
	metrics view.Exporter
}

// App returns the application for the web server
func (d *defaultServer) App() chi.Router {
	return d.app
}

// Init the application and its modules with the config.
func (d *defaultServer) Init() error {
	err := d.Application.Init()
	if err != nil {
		return err
	}

	srvOpts := d.httpdOpts[:]
	srvOpts = append(srvOpts, d.adminOpts...) // force admin config
	srvOpts = append(srvOpts,
		httpd.HandlesRequestsWith(d.app), // force handler config
	)
	d.server = httpd.New(srvOpts...)
	return nil
}

// Start the application an its enabled modules
func (d *defaultServer) Start() error {
	if err := d.Application.Start(); err != nil {
		return err
	}

	if err := d.server.Listen(); err != nil {
		return err
	}

	return d.server.Serve()
}

// Stop the application an its enabled modules
func (d *defaultServer) Stop() error {
	if err := d.server.Shutdown(); err != nil {
		return err
	}
	return d.Application.Stop()
}
