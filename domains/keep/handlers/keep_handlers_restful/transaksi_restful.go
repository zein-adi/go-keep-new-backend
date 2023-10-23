package keep_handlers_restful

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"net/http"
	"strconv"
	"time"
)

func NewTransaksiRestfulHandler(transaksiService keep_service_interfaces.ITransaksiServices) *TransaksiRestfulHandler {
	return &TransaksiRestfulHandler{
		service: transaksiService,
	}
}

type TransaksiRestfulHandler struct {
	service keep_service_interfaces.ITransaksiServices
}

func (x *TransaksiRestfulHandler) Get(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	q := r.URL.Query()
	request := keep_request.NewGetTransaksi()
	request.PosId = q.Get("posId")
	request.KantongId = q.Get("kantongId")
	request.JenisTanggal = q.Get("jenisTanggal")
	request.Jenis = q.Get("jenis")
	tanggalString := q.Get("tanggal")
	if tanggalString != "" {
		tanggalInt, err := strconv.ParseInt(tanggalString, 10, 64)
		if err != nil {
			e := helpers_error.NewValidationErrors("tanggal", "type", "integer")
			h.SendErrorResponse(w, 400, e.Error())
			return
		}
		request.Tanggal = tanggalInt
	}
	waktuAwalString := q.Get("waktuAwal")
	if waktuAwalString != "" {
		waktuAwalInt, err := strconv.ParseInt(waktuAwalString, 10, 64)
		if err != nil {
			e := helpers_error.NewValidationErrors("tanggal", "type", "integer")
			h.SendErrorResponse(w, 400, e.Error())
			return
		}
		request.WaktuAwal = waktuAwalInt
	}

	models := x.service.Get(ctx, request)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}
func (x *TransaksiRestfulHandler) Insert(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.TransaksiInputUpdate{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	model, err := x.service.Insert(ctx, input)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, model)
}
func (x *TransaksiRestfulHandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	input := &keep_request.TransaksiInputUpdate{}
	if !h.ReadRequest(w, r, h.NewDefaultFormRequest(input)) {
		return
	}

	vars := mux.Vars(r)
	input.Id = vars["transaksiId"]
	model, err := x.service.Update(ctx, input)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else if errors.Is(err, helpers_error.ValidationError) {
			h.SendErrorResponse(w, http.StatusBadRequest, errors.Unwrap(err).Error())
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}

	h.SendSingleResponse(w, http.StatusOK, model)
}
func (x *TransaksiRestfulHandler) DeleteById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["transaksiId"]
	affected, err := x.service.DeleteById(ctx, id)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}
	h.SendSingleResponse(w, http.StatusOK, affected)
}
func (x *TransaksiRestfulHandler) GetTrashed(w http.ResponseWriter, _ *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	models := x.service.GetTrashed(ctx)
	h.SendMultiResponse(w, http.StatusOK, models, len(models))
}
func (x *TransaksiRestfulHandler) RestoreTrashedById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["transaksiId"]
	affected, err := x.service.RestoreTrashedById(ctx, id)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}
	h.SendSingleResponse(w, http.StatusOK, affected)
}
func (x *TransaksiRestfulHandler) DeleteTrashedById(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	vars := mux.Vars(r)
	id := vars["transaksiId"]
	affected, err := x.service.DeleteTrashedById(ctx, id)
	if err != nil {
		if errors.Is(err, helpers_error.EntryNotFoundError) {
			h.SendErrorResponse(w, http.StatusNotFound, "")
		} else {
			h.SendErrorResponse(w, http.StatusInternalServerError, "")
		}
		return
	}
	h.SendSingleResponse(w, http.StatusOK, affected)
}
