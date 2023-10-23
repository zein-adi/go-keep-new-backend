package routes

import (
	"github.com/zein-adi/go-keep-new-backend/app/components/gorillamux_router"
	"github.com/zein-adi/go-keep-new-backend/app/middlewares"
	"github.com/zein-adi/go-keep-new-backend/dependency_injection"
	"github.com/zein-adi/go-keep-new-backend/domains/auth/handlers/auth_handlers_local"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/handlers/keep_handlers_events"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/handlers/keep_handlers_restful"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_events"
)

func injectKeepRoutes(r *gorillamux_router.Router) {

	middlewareAcl := middlewares.NewMiddlewareAcl(auth_handlers_local.NewRoleLocalHandler(dependency_injection.InitUserRoleServices()))

	r.Group("/keep", "keep.", func(r *gorillamux_router.Router) {

		pos := keep_handlers_restful.NewPosRestfulHandler(dependency_injection.InitKeepPosServices())
		r.Group("/posts", "pos.", func(r *gorillamux_router.Router) {
			r.GET("", pos.Get, "get")
			r.GET("/trash", pos.GetTrashed, "trash")
			r.POST("", pos.Insert, "insert")
			r.PATCH("/{posId:[0-9]+}", pos.Update, "update")
			r.PATCH("/{posId:[0-9]+}/trash", pos.RestoreTrashedById, "trash")
			r.DELETE("/{posId:[0-9]+}", pos.DeleteById, "delete")
			//r.DELETE("/{posId:[0-9]+}/trash", pos.DeleteTrashedById, "trash") // Dangerous
		})

		kantong := keep_handlers_restful.NewKantongRestfulHandler(dependency_injection.InitKeepKantongServices())
		r.Group("/kantong", "kantong.", func(r *gorillamux_router.Router) {
			r.GET("", kantong.Get, "get")
			r.GET("/trash", kantong.GetTrashed, "trash")
			r.POST("", kantong.Insert, "insert")
			r.PATCH("/{kantongId:[0-9]+}", kantong.Update, "update")
			r.PATCH("/{kantongId:[0-9]+}/trash", kantong.RestoreTrashedById, "trash")
			r.DELETE("/{kantongId:[0-9]+}", kantong.DeleteById, "delete")
			//r.DELETE("/{kantongId:[0-9]+}/trash", kantong.DeleteTrashedById, "trash") // Dangerous

			kantongHistory := keep_handlers_restful.NewKantongHistoryRestfulHandler(dependency_injection.InitKeepKantongHistoryServices())
			r.Group("/{kantongId:[0-9]+}/history", "history.", func(r *gorillamux_router.Router) {
				r.GET("", kantongHistory.Get, "get")
				r.POST("", kantongHistory.Insert, "insert")
				r.PATCH("/{kantongHistoryId:[0-9]+}", kantongHistory.Update, "update")
				r.DELETE("/{kantongHistoryId:[0-9]+}", kantongHistory.DeleteById, "delete")
			})
		})
	}).SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle)
}

func RegisterKeepListeners(
	posService keep_service_interfaces.IPosServices,
	kantongServices keep_service_interfaces.IKantongServices,
	lokasiServices keep_service_interfaces.ILokasiServices,
	barangServices keep_service_interfaces.IBarangServices,
) {
	d := helpers_events.GetDispatcher()
	pos := keep_handlers_events.NewPosEventListenerHandler(posService)
	kantong := keep_handlers_events.NewKantongEventListenerHandler(kantongServices)
	lokasi := keep_handlers_events.NewLokasiEventListenerHandler(lokasiServices)
	barang := keep_handlers_events.NewBarangEventListenerHandler(barangServices)
	_ = d.Register(keep_events.TransaksiCreated,
		pos.TransaksiCreated,
		kantong.TransaksiCreated,
		lokasi.TransaksiCreated,
		barang.TransaksiCreated,
	)
	_ = d.Register(keep_events.TransaksiUpdated,
		pos.TransaksiUpdated,
		kantong.TransaksiUpdated,
		lokasi.TransaksiUpdated,
		barang.TransaksiUpdated,
	)
	_ = d.Register(keep_events.TransaksiSoftDeleted,
		pos.TransaksiSoftDelete,
		kantong.TransaksiSoftDelete,
		lokasi.TransaksiSoftDelete,
		barang.TransaksiSoftDelete,
	)
	_ = d.Register(keep_events.TransaksiRestored,
		pos.TransaksiRestore,
		kantong.TransaksiRestore,
		lokasi.TransaksiRestore,
		barang.TransaksiRestore,
	)
}
