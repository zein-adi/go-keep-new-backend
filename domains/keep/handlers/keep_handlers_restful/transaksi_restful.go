package keep_handlers_restful

import (
	"context"
	"errors"
	"github.com/gorilla/mux"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_request"
	"github.com/zein-adi/go-keep-new-backend/domains/keep/core/keep_service_interfaces"
	"github.com/zein-adi/go-keep-new-backend/helpers/helpers_error"
	h "github.com/zein-adi/go-keep-new-backend/helpers/helpers_http"
	"github.com/zein-adi/go-keep-new-backend/helpers/validator"
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
	v := validator.New()
	err := v.ValidateMap(map[string]any{
		"skip":      q.Get("skip"),
		"take":      q.Get("take"),
		"tanggal":   q.Get("tanggal"),
		"waktuAwal": q.Get("waktuAwal"),
	}, map[string]any{
		"skip":      "omitempty,number",
		"take":      "omitempty,number",
		"tanggal":   "omitempty,number",
		"waktuAwal": "omitempty,number",
	})
	if err != nil {
		h.SendErrorResponse(w, 400, err.Error())
		return
	}

	request := keep_request.NewGetTransaksi()
	request.Skip, err = strconv.Atoi(q.Get("skip"))
	helpers_error.PanicIfError(err)
	request.Take, err = strconv.Atoi(q.Get("take"))
	helpers_error.PanicIfError(err)
	request.Tanggal, err = strconv.ParseInt(q.Get("tanggal"), 10, 64)
	helpers_error.PanicIfError(err)
	request.WaktuAwal, err = strconv.ParseInt(q.Get("waktuAwal"), 10, 64)
	helpers_error.PanicIfError(err)
	request.Search = q.Get("search")
	request.PosId = q.Get("posId")
	request.KantongId = q.Get("kantongId")
	request.JenisTanggal = q.Get("jenisTanggal")
	request.Jenis = q.Get("jenis")

	err = v.ValidateStruct(request)
	if err != nil {
		h.SendErrorResponse(w, 400, err.Error())
		return
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
