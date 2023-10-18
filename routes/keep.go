package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
)

func injectKeepRoutes(r *components.Router) {
	middlewareAcl := dependency_injection.InitAclMiddleware()

	r.Group("/keep", "keep.", func(r *components.Router) {

		pos := dependency_injection.InitKeepPosRestful()
		r.Group("/posts", "pos.", func(r *components.Router) {
			r.GET("/", pos.Get, "get")
			r.GET("/trash/", pos.GetTrashed, "trash")
			r.POST("/", pos.Insert, "insert")
			r.PATCH("/:posId/", pos.Update, "update")
			r.PATCH("/:posId/trash/", pos.RestoreTrashedById, "trash")
			r.DELETE("/:posId/", pos.DeleteById, "delete")
			r.DELETE("/:posId/trash/", pos.DeleteTrashedById, "trash")
		})

		kantong := dependency_injection.InitKeepKantongRestful()
		r.Group("/kantong", "kantong.", func(r *components.Router) {
			r.GET("/", kantong.Get, "get")
			r.GET("/trash", kantong.GetTrashed, "trash")
			r.POST("/", kantong.Insert, "insert")
			r.PATCH("/:kantongId/", kantong.Update, "update")
			r.PATCH("/:kantongId/trash/", kantong.RestoreTrashedById, "trash")
			r.DELETE("/:kantongId/", kantong.DeleteById, "delete")
			r.DELETE("/:kantongId/trash/", kantong.DeleteTrashedById, "trash")

			kantongHistory := dependency_injection.InitKeepKantongHistoryRestful()
			r.Group("/:kantongId/history", "history.", func(r *components.Router) {
				r.GET("/", kantongHistory.Get, "get")
				r.POST("/", kantongHistory.Insert, "insert")
				r.PATCH("/:kantongHistoryId/", kantongHistory.Update, "update")
				r.DELETE("/:kantongHistoryId/", kantongHistory.DeleteById, "delete")
			})
		})
	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}
