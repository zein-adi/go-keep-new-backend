package keep_handlers_restful

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"time"
)

func NewBarangRestfulHandler(barangServices keep_service_interfaces.IBarangServices) *BarangRestfulHandler {
	return &BarangRestfulHandler{
		service: barangServices,
	}
}

type BarangRestfulHandler struct {
	service keep_service_interfaces.IBarangServices
}

func (x *BarangRestfulHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	q := r.URL.Query()
	search := q.Get("search")
	lokasi := q.Get("lokasi")

	models := x.service.Get(ctx, search, lokasi)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}
