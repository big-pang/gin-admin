package router

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

var routerNames = make(map[string]string)

type NamedRoutes struct {
	routes map[string]string
	engine *gin.Engine
}

func NewNamedRoutes(engine *gin.Engine) *NamedRoutes {
	return &NamedRoutes{
		routes: routerNames,
		engine: engine,
	}
}

// AddRoute adds a new route with a given name and method, path and handlers.
func (nr *NamedRoutes) AddRoute(name, method, path string, handlers ...gin.HandlerFunc) {
	nr.engine.Handle(method, path, handlers...)
	nr.routes[name] = path
}

func (nr *NamedRoutes) URLFor(name string, params map[string]string) (string, error) {
	path, exists := nr.routes[name]
	if !exists {
		return "", fmt.Errorf("route %s not found", name)
	}
	url := path
	for key, value := range params {
		url = strings.Replace(url, fmt.Sprintf(":%s", key), value, -1)
	}
	return url, nil
}
