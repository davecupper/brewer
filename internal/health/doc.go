// Package health provides utilities for waiting until a service becomes
// healthy before proceeding with dependent service startup.
//
// A health check can be configured per service in brewer.yaml using either
// an HTTP endpoint or a TCP address. The checker will poll at a fixed
// interval until the check passes or the configured timeout is exceeded.
//
// Supported check types:
//
//	http  — performs an HTTP GET and expects a 2xx response
//	tcp   — attempts a TCP dial and expects a successful connection
//
// Example configuration:
//
//	services:
//	  - name: api
//	    command: ./api-server
//	    health_check:
//	      type: http
//	      target: http://localhost:8080/health
//	      interval: 2s
//	      timeout: 30s
package health
