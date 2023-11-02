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

	r.New().SetMiddleware(middlewares.AuthHandle, middlewareAcl.Handle).
		Group("/keep", "keep.", func(r *gorillamux_router.Router) {

			transaksi := keep_handlers_restful.NewTransaksiRestfulHandler(dependency_injection.InitKeepTransaksiServices())
			r.Group("/transaksi", "transaksi.", func(r *gorillamux_router.Router) {
				r.GET("", transaksi.Get, "get")
				r.GET("/trash", transaksi.GetTrashed, "trash")
				r.POST("", transaksi.Insert, "insert")
				r.PATCH("/{transaksiId:[0-9]+}", transaksi.Update, "update")
				r.PATCH("/{transaksiId:[0-9]+}/trash/restore", transaksi.RestoreTrashedById, "trash")
				r.DELETE("/{transaksiId:[0-9]+}", transaksi.DeleteById, "delete")
				r.DELETE("/{transaksiId:[0-9]+}/trash", transaksi.DeleteTrashedById, "trash") // Dangerous
			})

			lokasi := keep_handlers_restful.NewLokasiRestfulHandler(dependency_injection.InitKeepLokasiServices())
			r.Group("/lokasi", "lokasi.", func(r *gorillamux_router.Router) {
				r.GET("", lokasi.Get, "get")
			})

			barang := keep_handlers_restful.NewBarangRestfulHandler(dependency_injection.InitKeepBarangServices())
			r.Group("/barang", "barang.", func(r *gorillamux_router.Router) {
				r.GET("", barang.Get, "get")
			})

			pos := keep_handlers_restful.NewPosRestfulHandler(dependency_injection.InitKeepPosServices())
			r.Group("/pos", "pos.", func(r *gorillamux_router.Router) {
				r.GET("", pos.Get, "get")
				r.GET("/trash", pos.GetTrashed, "trash")
				r.POST("", pos.Insert, "insert")
				r.PATCH("/order", pos.UpdateUrutan, "update")
				r.PATCH("/visibility", pos.UpdateVisivility, "update")
				r.PATCH("/{posId:[0-9]+}", pos.Update, "update")
				r.PATCH("/{posId:[0-9]+}/trash/restore", pos.RestoreTrashedById, "trash")
				r.DELETE("/{posId:[0-9]+}", pos.DeleteById, "delete")
				//r.DELETE("/{posId:[0-9]+}/trash", pos.DeleteTrashedById, "trash") // Dangerous
			})

			kantong := keep_handlers_restful.NewKantongRestfulHandler(dependency_injection.InitKeepKantongServices())
			r.Group("/kantong", "kantong.", func(r *gorillamux_router.Router) {
				r.GET("", kantong.Get, "get")
				r.GET("/trash", kantong.GetTrashed, "trash")
				r.POST("", kantong.Insert, "insert")
				r.PATCH("/order", kantong.UpdateUrutan, "update")
				r.PATCH("/visibility", kantong.UpdateVisivility, "update")
				r.PATCH("/{kantongId:[0-9]+}", kantong.Update, "update")
				r.PATCH("/{kantongId:[0-9]+}/trash/restore", kantong.RestoreTrashedById, "trash")
				r.DELETE("/{kantongId:[0-9]+}", kantong.DeleteById, "delete")
				//r.DELETE("/{kantongId:[0-9]+}/trash", kantong.DeleteTrashedById, "trash") // Dangerous

				kantongHistory := keep_handlers_restful.NewKantongHistoryRestfulHandler(dependency_injection.InitKeepKantongHistoryServices())
				r.Group("/{kantongId:[0-9]+}/history", "history.", func(r *gorillamux_router.Router) {
					r.GET("", kantongHistory.Get, "get")
					r.POST("", kantongHistory.InsertAndUpdateSaldoKantong, "insert")
					r.PATCH("/{kantongHistoryId:[0-9]+}", kantongHistory.Update, "update")
					r.DELETE("/{kantongHistoryId:[0-9]+}", kantongHistory.DeleteById, "delete")
				})
			})
		})
}

func RegisterKeepListeners(
	posService keep_service_interfaces.IPosServices,
	kantongServices keep_service_interfaces.IKantongServices,
	kantongHistoryServices keep_service_interfaces.IKantongHistoryServices,
	lokasiServices keep_service_interfaces.ILokasiServices,
	barangServices keep_service_interfaces.IBarangServices,
) {
	d := helpers_events.GetDispatcher()

	pos := keep_handlers_events.NewPosEventListenerHandler(posService)
	kantong := keep_handlers_events.NewKantongEventListenerHandler(kantongServices, kantongHistoryServices)
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
		lokasi.TransaksiUpdated,
		barang.TransaksiUpdated,
	)
	_ = d.Register(keep_events.TransaksiSoftDeleted,
		pos.TransaksiSoftDeleted,
		lokasi.TransaksiSoftDeleted,
		barang.TransaksiSoftDeleted,
	)
	_ = d.Register(keep_events.TransaksiRestored,
		pos.TransaksiRestored,
		lokasi.TransaksiRestored,
		barang.TransaksiRestored,
	)
	_ = d.Register(keep_events.KantongHistoryCreated,
		kantong.KantongHistoryCreated,
	)
}
