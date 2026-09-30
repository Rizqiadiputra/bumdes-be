package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/middleware"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

const ticketingDateLayout = "2006-01-02"

type TicketingHandler struct {
	ticketingService service.TicketingService
	userLogService   service.UserLogService
}

func NewTicketingHandler(ticketingService service.TicketingService, userLogService service.UserLogService) *TicketingHandler {
	return &TicketingHandler{ticketingService: ticketingService, userLogService: userLogService}
}

// ListTicketings godoc
// @Summary List Transaksi Tiket
// @Description Menampilkan daftar transaksi tiket. Default menampilkan transaksi hari ini, bisa difilter dengan date, user_id, jenis (type_name), dan payment (payment_id).
// @Tags Ticketing
// @Produce json
// @Security BearerAuth
// @Param date query string false "Filter tanggal, format YYYY-MM-DD (default: hari ini)"
// @Param user_id query string false "Filter berdasarkan user id (UUID)"
// @Param jenis query string false "Filter berdasarkan jenis tiket (type_name)"
// @Param payment query int false "Filter berdasarkan payment_id"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.TicketingResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /ticketing [get]
func (h *TicketingHandler) ListTicketings(c *gin.Context) {
	filterDate := time.Now()
	if v := c.Query("date"); v != "" {
		parsed, err := time.Parse(ticketingDateLayout, v)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "format date tidak valid, gunakan YYYY-MM-DD")
			return
		}
		filterDate = parsed
	}

	filter := repository.TicketingFilter{
		Date:     filterDate,
		TypeName: c.Query("jenis"),
	}

	if v := c.Query("user_id"); v != "" {
		id, err := uuid.Parse(v)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "user_id tidak valid")
			return
		}
		filter.UserID = id
	}

	if v := c.Query("payment"); v != "" {
		id, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "payment tidak valid")
			return
		}
		filter.PaymentID = uint(id)
	}

	result, err := h.ticketingService.ListTicketings(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar transaksi tiket", result)
}

// GetTicketing godoc
// @Summary Detail Transaksi Tiket
// @Description Menampilkan detail satu transaksi tiket
// @Tags Ticketing
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticketing ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.TicketingResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /ticketing/{id} [get]
func (h *TicketingHandler) GetTicketing(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id transaksi tiket tidak valid")
		return
	}

	result, err := h.ticketingService.GetTicketing(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "transaksi tiket tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail transaksi tiket", result)
}

// CreateTicketing godoc
// @Summary Buat Transaksi Tiket
// @Description Membuat transaksi tiket baru. Field type_name dan price otomatis diambil dari tarif tiket (ticket_price_id), payment_method_name otomatis diambil dari payment_id. user_id otomatis diisi dari user yang sedang login, date otomatis diisi tanggal hari ini.
// @Tags Ticketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTicketingRequest true "Data transaksi tiket baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.TicketingResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /ticketing [post]
func (h *TicketingHandler) CreateTicketing(c *gin.Context) {
	var req dto.CreateTicketingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "token tidak valid")
		return
	}

	result, err := h.ticketingService.CreateTicketing(req, actorID.(uuid.UUID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTicketPriceNotFound), errors.Is(err, service.ErrPaymentMethodNotFound):
			utils.Error(c, http.StatusNotFound, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	meta := utils.GetRequestMeta(c)
	_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "ticketing",
		fmt.Sprintf("membuat transaksi tiket %s", result.TypeName), req)

	utils.Success(c, http.StatusCreated, "transaksi tiket berhasil dibuat", result)
}

// UpdateTicketing godoc
// @Summary Ubah Transaksi Tiket
// @Description Mengubah data transaksi tiket. Field type_name, price, dan payment_method_name otomatis dihitung ulang berdasarkan ticket_price_id/payment_id yang baru.
// @Tags Ticketing
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticketing ID (UUID)"
// @Param request body dto.UpdateTicketingRequest true "Data transaksi tiket yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.TicketingResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /ticketing/{id} [put]
func (h *TicketingHandler) UpdateTicketing(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id transaksi tiket tidak valid")
		return
	}

	var req dto.UpdateTicketingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.ticketingService.UpdateTicketing(id, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.Error(c, http.StatusNotFound, "transaksi tiket tidak ditemukan")
		case errors.Is(err, service.ErrTicketPriceNotFound), errors.Is(err, service.ErrPaymentMethodNotFound):
			utils.Error(c, http.StatusNotFound, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "ticketing",
			fmt.Sprintf("mengubah transaksi tiket %s", newResp.TypeName), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "transaksi tiket berhasil diperbarui", newResp)
}

// DeleteTicketing godoc
// @Summary Hapus Transaksi Tiket
// @Description Menghapus transaksi tiket
// @Tags Ticketing
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticketing ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /ticketing/{id} [delete]
func (h *TicketingHandler) DeleteTicketing(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id transaksi tiket tidak valid")
		return
	}

	if err := h.ticketingService.DeleteTicketing(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "transaksi tiket tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "ticketing",
			fmt.Sprintf("menghapus transaksi tiket id %s", id))
	}

	utils.Success(c, http.StatusOK, "transaksi tiket berhasil dihapus", nil)
}
