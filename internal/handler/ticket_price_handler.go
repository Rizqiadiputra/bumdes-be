package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/middleware"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type TicketPriceHandler struct {
	ticketPriceService service.TicketPriceService
	userLogService     service.UserLogService
}

func NewTicketPriceHandler(ticketPriceService service.TicketPriceService, userLogService service.UserLogService) *TicketPriceHandler {
	return &TicketPriceHandler{ticketPriceService: ticketPriceService, userLogService: userLogService}
}

// ListTicketPrices godoc
// @Summary List Tarif Tiket
// @Description Menampilkan daftar tarif tiket, bisa difilter dengan pencarian type_name
// @Tags Ticket Prices
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan type_name"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.TicketPriceResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/ticket-prices [get]
func (h *TicketPriceHandler) ListTicketPrices(c *gin.Context) {
	filter := repository.TicketPriceFilter{Search: c.Query("search")}

	result, err := h.ticketPriceService.ListTicketPrices(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar tarif tiket", result)
}

// GetTicketPrice godoc
// @Summary Detail Tarif Tiket
// @Description Menampilkan detail satu tarif tiket
// @Tags Ticket Prices
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket Price ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.TicketPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /master-data/ticket-prices/{id} [get]
func (h *TicketPriceHandler) GetTicketPrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif tiket tidak valid")
		return
	}

	result, err := h.ticketPriceService.GetTicketPrice(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "tarif tiket tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail tarif tiket", result)
}

// CreateTicketPrice godoc
// @Summary Buat Tarif Tiket
// @Description Membuat tarif tiket baru
// @Tags Ticket Prices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTicketPriceRequest true "Data tarif tiket baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.TicketPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/ticket-prices [post]
func (h *TicketPriceHandler) CreateTicketPrice(c *gin.Context) {
	var req dto.CreateTicketPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.ticketPriceService.CreateTicketPrice(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "ticket-prices",
			fmt.Sprintf("membuat tarif tiket %s", result.TypeName), req)
	}

	utils.Success(c, http.StatusCreated, "tarif tiket berhasil dibuat", result)
}

// UpdateTicketPrice godoc
// @Summary Ubah Tarif Tiket
// @Description Mengubah data tarif tiket
// @Tags Ticket Prices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket Price ID (UUID)"
// @Param request body dto.UpdateTicketPriceRequest true "Data tarif tiket yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.TicketPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/ticket-prices/{id} [put]
func (h *TicketPriceHandler) UpdateTicketPrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif tiket tidak valid")
		return
	}

	var req dto.UpdateTicketPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.ticketPriceService.UpdateTicketPrice(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tarif tiket tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "ticket-prices",
			fmt.Sprintf("mengubah tarif tiket %s", newResp.TypeName), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "tarif tiket berhasil diperbarui", newResp)
}

// DeleteTicketPrice godoc
// @Summary Hapus Tarif Tiket
// @Description Menghapus tarif tiket
// @Tags Ticket Prices
// @Produce json
// @Security BearerAuth
// @Param id path string true "Ticket Price ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/ticket-prices/{id} [delete]
func (h *TicketPriceHandler) DeleteTicketPrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif tiket tidak valid")
		return
	}

	if err := h.ticketPriceService.DeleteTicketPrice(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tarif tiket tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "ticket-prices",
			fmt.Sprintf("menghapus tarif tiket id %s", id))
	}

	utils.Success(c, http.StatusOK, "tarif tiket berhasil dihapus", nil)
}
