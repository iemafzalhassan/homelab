package registry

import (
	"net/http"
	"net/http/httputil"
	"net/url"

	"go.uber.org/zap"
)

type Proxy struct {
	logger *zap.SugaredLogger
	target *url.URL
	proxy  *httputil.ReverseProxy
}

func NewProxy(logger *zap.SugaredLogger, targetURL string) (*Proxy, error) {
	target, err := url.Parse(targetURL)
	if err != nil {
		return nil, err
	}

	proxy := httputil.NewSingleHostReverseProxy(target)

	// Modify the Director to ensure Host headers are preserved correctly
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		// Maintain the original host so the backend registry knows who the client is
		req.Host = req.URL.Host

		// If we are passing through the authorization, we could inject a different token here
		// But currently we've validated the token and the backend registry is configured
		// to trust us (auth disabled on the backend, or we forward a system token).
		// For now, we just pass the request through.
	}

	return &Proxy{
		logger: logger,
		target: target,
		proxy:  proxy,
	}, nil
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The Docker client uses /v2/ to ping the registry.
	// If it hits the proxy, we forward it to the backend.
	p.logger.Infow("proxying request", "method", r.Method, "path", r.URL.Path)

	// Set standard proxy headers
	r.Header.Set("X-Forwarded-Host", r.Host)
	r.Header.Set("X-Forwarded-Proto", "http") // Change to https in production
	if r.TLS != nil {
		r.Header.Set("X-Forwarded-Proto", "https")
	}

	p.proxy.ServeHTTP(w, r)
}

// IsPublicEndpoint returns true if the endpoint doesn't require authentication.
// For example, /v2/ ping endpoint might be open, but we still want to challenge it
// so the docker client knows how to authenticate.
func IsPublicEndpoint(path string, method string) bool {
	// V2 ping is generally expected to return 401 with WWW-Authenticate header
	// if auth is required. So we don't treat it as public.
	// We might treat specific catalog or public repo pulls as public later.
	return false
}
