package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components/gorillamux_router"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
)

func injectKeepRoutes(r *gorillamux_router.Router) {
	middlewareAcl := dependency_injection.InitAclMiddleware()

	r.Group("/keep", "keep.", func(r *gorillamux_router.Router) {

		pos := dependency_injection.InitKeepPosRestful()
		r.Group("/posts", "pos.", func(r *gorillamux_router.Router) {
			r.GET("", pos.Get, "get")
			r.GET("/trash", pos.GetTrashed, "trash")
			r.POST("", pos.Insert, "insert")
			r.PATCH("/{posId:[0-9]+}", pos.Update, "update")
			r.PATCH("/{posId:[0-9]+}/trash", pos.RestoreTrashedById, "trash")
			r.DELETE("/{posId:[0-9]+}", pos.DeleteById, "delete")
			r.DELETE("/{posId:[0-9]+}/trash", pos.DeleteTrashedById, "trash")
		})

		kantong := dependency_injection.InitKeepKantongRestful()
		r.Group("/kantong", "kantong.", func(r *gorillamux_router.Router) {
			r.GET("", kantong.Get, "get")
			r.GET("/trash", kantong.GetTrashed, "trash")
			r.POST("", kantong.Insert, "insert")
			r.PATCH("/{kantongId:[0-9]+}", kantong.Update, "update")
			r.PATCH("/{kantongId:[0-9]+}/trash", kantong.RestoreTrashedById, "trash")
			r.DELETE("/{kantongId:[0-9]+}", kantong.DeleteById, "delete")
			r.DELETE("/{kantongId:[0-9]+}/trash", kantong.DeleteTrashedById, "trash")

			kantongHistory := dependency_injection.InitKeepKantongHistoryRestful()
			r.Group("/{kantongId:[0-9]+}/history", "history.", func(r *gorillamux_router.Router) {
				r.GET("", kantongHistory.Get, "get")
				r.POST("", kantongHistory.Insert, "insert")
				r.PATCH("/{kantongHistoryId:[0-9]+}", kantongHistory.Update, "update")
				r.DELETE("/{kantongHistoryId:[0-9]+}", kantongHistory.DeleteById, "delete")
			})
		})
	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}
