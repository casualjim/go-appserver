# App Server

A small wrapper around [httpd](https://github.com/casualjim/go-httpd) and [chi](https://github.com/go-chi/chi).
Comes with a simple service locator included.

## Features

* Fast router configurable with middlewares.
* Service Locator with application lifecycle hooks.
* Default middleware stack (can be replaced or augmented):
  * Recover from panics
  * Compression
  * Read Load Balancer headers for correct remote address and scheme
  * OpenCensus tracing
  * Request logging with zap

## Usage
