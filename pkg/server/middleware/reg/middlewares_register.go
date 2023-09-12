package reg

import (
	"context"
	"net/http"
)

type NewFunc func(ctx context.Context, next http.Handler, cfg any, name string) (http.Handler, error)

var thirdMiddlewares = make(map[string]NewFunc)

func Register(typeName string, fn NewFunc) {
	if _, ok := thirdMiddlewares[typeName]; !ok {
		thirdMiddlewares[typeName] = fn
	}
}

func NewMiddlewareFunc(typeName string) NewFunc {
	return thirdMiddlewares[typeName]
}
