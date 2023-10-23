package keep_handlers_restful

import (
	"context"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"time"
)

func NewLokasiRestfulHandler(lokasiService keep_service_interfaces.ILokasiServices) *LokasiRestfulHandler {
	return &LokasiRestfulHandler{
		service: lokasiService,
	}
}

type LokasiRestfulHandler struct {
	service keep_service_interfaces.ILokasiServices
}

func (x *LokasiRestfulHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	q := r.URL.Query()
	search := ""
	if q.Has("search") {
		search = q.Get("search")
	}

	models := x.service.Get(ctx, search)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}
