package components

import (
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
)

func NewRouter(corsOptions cors.Options) *Router {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	m := &Router{
		router:          httprouter.New(),
		cors:            cors.New(corsOptions),
		logger:          logger,
		pathPrefix:      "",
		routeNamePrefix: "",
	}
	m.initNotFoundHandler()
	m.initNotAllowedHandler()
	m.initPanicHandler()
	return m
}

type MiddlewareFunc func(writer http.ResponseWriter, request *http.Request, params httprouter.Params, routeName string) bool
type Router struct {
	router          *httprouter.Router
	middlewares     []MiddlewareFunc
	cors            *cors.Cors
	logger          *logrus.Logger
	pathPrefix      string
	routeNamePrefix string
}

func (m *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	m.cors.ServeHTTP(writer, request, m.router.ServeHTTP)
}
func (m *Router) GET(path string, handle httprouter.Handle, routeName string) {
	m.handle(http.MethodGet, path, handle, routeName)
}
func (m *Router) HEAD(path string, handle httprouter.Handle, routeName string) {
	m.handle(http.MethodHead, path, handle, routeName)
}
func (m *Router) OPTIONS(path string, handle httprouter.Handle, routeName string) {
	m.handle(http.MethodOptions, path, handle, routeName)
}
func (m *Router) POST(path string, handle httprouter.Handle, routeName string) {
	m.handle(http.MethodPost, path, handle, routeName)
}
func (m *Router) PUT(path string, handle httprouter.Handle, routeName string) {
	m.handle(http.MethodPut, path, handle, routeName)
}
func (m *Router) PATCH(path string, handle httprouter.Handle, routeName string) {
	m.handle(http.MethodPatch, path, handle, routeName)
}
func (m *Router) DELETE(path string, handle httprouter.Handle, routeName string) {
	m.handle(http.MethodDelete, path, handle, routeName)
}
func (m *Router) Group(pathPrefix string, routeNamePrefix string, handleGroup func(router *Router)) *Router {
	newRouter := &Router{
		router:          m.router,
		cors:            m.cors,
		logger:          m.logger,
		middlewares:     m.middlewares,
		pathPrefix:      m.pathPrefix + pathPrefix,
		routeNamePrefix: m.routeNamePrefix + routeNamePrefix,
	}
	handleGroup(newRouter)
	return newRouter
}
func (m *Router) AddMiddleware(middlewares ...MiddlewareFunc) {
	m.middlewares = append(m.middlewares, middlewares...)
}
func (m *Router) SetMiddleware(middlewares ...MiddlewareFunc) {
	m.middlewares = middlewares
}
func (m *Router) handle(method string, path string, handle httprouter.Handle, routeName string) {
	m.router.Handle(method, m.pathPrefix+path, func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
		for _, middleware := range m.middlewares {
			if !middleware(w, r, p, m.routeNamePrefix+routeName) {
				return
			}
		}
		handle(w, r, p)
	})
}
func (m *Router) initNotFoundHandler() {
	m.router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.WithFields(logrus.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Warn("router not found")

		helpers_http.SendResponseJson(w, http.StatusNotFound, nil)
	})
}
func (m *Router) initNotAllowedHandler() {
	m.router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.WithFields(logrus.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Warn("router method not allowed")

		helpers_http.SendResponseJson(w, http.StatusMethodNotAllowed, nil)
	})
}
func (m *Router) initPanicHandler() {
	m.router.PanicHandler = func(w http.ResponseWriter, r *http.Request, p any) {
		m.logger.WithFields(logrus.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Error(p)

		helpers_http.SendResponseJson(w, http.StatusInternalServerError, nil)
	}
}
