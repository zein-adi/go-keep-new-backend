package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
)

func injectKeepRoutes(r *components.Router) {
	middlewareAcl := dependency_injection.InitAclMiddleware()
	pos := dependency_injection.InitKeepPosRestful()

	r.Group("/keep", "keep.", func(r *components.Router) {

		r.GET("/posts", pos.Get, "pos.get")
		r.GET("/posts/trash", pos.GetTrashed, "pos.trash")
		r.POST("/posts", pos.Insert, "pos.insert")
		r.PATCH("/posts/:posId", pos.Update, "pos.update")
		r.PATCH("/posts/:posId/trash", pos.RestoreTrashedById, "pos.trash")
		r.DELETE("/posts/:posId", pos.DeleteById, "pos.delete")
		r.DELETE("/posts/:posId/trash", pos.DeleteTrashedById, "pos.trash")

	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}
