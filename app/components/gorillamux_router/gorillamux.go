package gorillamux_router

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
)

func NewRouter(corsOptions cors.Options) *Router {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})

	m := &Router{
		router:          mux.NewRouter(),
		cors:            cors.New(corsOptions),
		logger:          logger,
		pathPrefix:      "",
		routeNamePrefix: "",
	}
	m.initNotFoundHandler()
	m.initNotAllowedHandler()

	return m
}

type MiddlewareFunc func(writer http.ResponseWriter, request *http.Request, routeName string) bool
type HttpHandler func(writer http.ResponseWriter, request *http.Request)
type Router struct {
	router          *mux.Router
	middlewares     []MiddlewareFunc
	cors            *cors.Cors
	logger          *logrus.Logger
	pathPrefix      string
	routeNamePrefix string
}

func (m *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	defer func() {
		r := recover()
		if r != nil {
			m.panicHandler(writer, request)
		}
	}()
	m.cors.ServeHTTP(writer, request, m.router.ServeHTTP)
}
func (m *Router) GET(path string, handle HttpHandler, routeName string) {
	m.handle(http.MethodGet, path, handle, routeName)
}
func (m *Router) HEAD(path string, handle HttpHandler, routeName string) {
	m.handle(http.MethodHead, path, handle, routeName)
}
func (m *Router) OPTIONS(path string, handle HttpHandler, routeName string) {
	m.handle(http.MethodOptions, path, handle, routeName)
}
func (m *Router) POST(path string, handle HttpHandler, routeName string) {
	m.handle(http.MethodPost, path, handle, routeName)
}
func (m *Router) PUT(path string, handle HttpHandler, routeName string) {
	m.handle(http.MethodPut, path, handle, routeName)
}
func (m *Router) PATCH(path string, handle HttpHandler, routeName string) {
	m.handle(http.MethodPatch, path, handle, routeName)
}
func (m *Router) DELETE(path string, handle HttpHandler, routeName string) {
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
func (m *Router) handle(method string, path string, handle HttpHandler, routeName string) {
	fmt.Printf("%-10s\t%-30s\t%-40s\n", method, m.routeNamePrefix+routeName, m.pathPrefix+path)
	m.router.HandleFunc(m.pathPrefix+path, func(w http.ResponseWriter, r *http.Request) {
		for _, middleware := range m.middlewares {
			if !middleware(w, r, m.routeNamePrefix+routeName) {
				return
			}
		}
		handle(w, r)
	}).Methods(method)
}
func (m *Router) initNotFoundHandler() {
	m.router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.WithFields(logrus.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Warn("router not found")

		helpers_http.SendResponseJson(w, http.StatusNotFound, nil)
	})
}
func (m *Router) initNotAllowedHandler() {
	m.router.MethodNotAllowedHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		m.logger.WithFields(logrus.Fields{
			"method": r.Method,
			"uri":    r.RequestURI,
		}).Warn("router method not allowed")

		helpers_http.SendResponseJson(w, http.StatusMethodNotAllowed, nil)
	})
}
func (m *Router) panicHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	m.logger.WithFields(logrus.Fields{
		"method": r.Method,
		"uri":    r.RequestURI,
	}).Error(vars)

	helpers_http.SendResponseJson(w, http.StatusInternalServerError, nil)
}
