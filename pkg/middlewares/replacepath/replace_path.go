package replacepath

import (
	"context"
	"net/http"
	"net/url"

	"github.com/opentracing/opentracing-go/ext"
	"github.com/txdultra/traefik/v2/pkg/config/dynamic"
	"github.com/txdultra/traefik/v2/pkg/log"
	"github.com/txdultra/traefik/v2/pkg/middlewares"
	"github.com/txdultra/traefik/v2/pkg/tracing"
)

const (
	// ReplacedPathHeader is the default header to set the old path to.
	ReplacedPathHeader = "X-Replaced-Path"
	typeName           = "ReplacePath"
)

// ReplacePath is a middleware used to replace the path of a URL request.
type replacePath struct {
	next http.Handler
	path string
	name string
}

// New creates a new replace path middleware.
func New(ctx context.Context, next http.Handler, config dynamic.ReplacePath, name string) (http.Handler, error) {
	log.FromContext(middlewares.GetLoggerCtx(ctx, name, typeName)).Debug("Creating middleware")

	return &replacePath{
		next: next,
		path: config.Path,
		name: name,
	}, nil
}

func (r *replacePath) GetTracingInformation() (string, ext.SpanKindEnum) {
	return r.name, tracing.SpanKindNoneEnum
}

func (r *replacePath) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	currentPath := req.URL.RawPath
	if currentPath == "" {
		currentPath = req.URL.EscapedPath()
	}

	req.Header.Add(ReplacedPathHeader, currentPath)
	req.URL.RawPath = r.path

	var err error
	req.URL.Path, err = url.PathUnescape(req.URL.RawPath)
	if err != nil {
		log.FromContext(middlewares.GetLoggerCtx(context.Background(), r.name, typeName)).Error(err)
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	req.RequestURI = req.URL.RequestURI()

	r.next.ServeHTTP(rw, req)
}
