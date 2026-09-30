package handler

import (
	"errors"
	"fmt"
	"net/http"
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

const attractionDateLayout = "2006-01-02"

type AttractionHandler struct {
	attractionService service.AttractionService
	userLogService    service.UserLogService
}

func NewAttractionHandler(attractionService service.AttractionService, userLogService service.UserLogService) *AttractionHandler {
	return &AttractionHandler{attractionService: attractionService, userLogService: userLogService}
}

// ListAttractions godoc
// @Summary List Transaksi Wahana
// @Description Menampilkan daftar transaksi wahana. Default menampilkan transaksi hari ini, bisa difilter dengan date, attraction_name, dan payment_method (berdasarkan payment_method_name).
// @Tags Attraction
// @Produce json
// @Security BearerAuth
// @Param date query string false "Filter tanggal, format YYYY-MM-DD (default: hari ini)"
// @Param attraction_name query string false "Filter berdasarkan nama wahana"
// @Param payment_method query string false "Filter berdasarkan payment_method_name"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.AttractionResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /attractions [get]
func (h *AttractionHandler) ListAttractions(c *gin.Context) {
	filterDate := time.Now()
	if v := c.Query("date"); v != "" {
		parsed, err := time.Parse(attractionDateLayout, v)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "format date tidak valid, gunakan YYYY-MM-DD")
			return
		}
		filterDate = parsed
	}

	filter := repository.AttractionFilter{
		Date:              filterDate,
		AttractionName:    c.Query("attraction_name"),
		PaymentMethodName: c.Query("payment_method"),
	}

	result, err := h.attractionService.ListAttractions(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar transaksi wahana", result)
}

// GetAttraction godoc
// @Summary Detail Transaksi Wahana
// @Description Menampilkan detail satu transaksi wahana
// @Tags Attraction
// @Produce json
// @Security BearerAuth
// @Param id path string true "Attraction ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.AttractionResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /attractions/{id} [get]
func (h *AttractionHandler) GetAttraction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id transaksi wahana tidak valid")
		return
	}

	result, err := h.attractionService.GetAttraction(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "transaksi wahana tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail transaksi wahana", result)
}

// CreateAttraction godoc
// @Summary Buat Transaksi Wahana
// @Description Membuat transaksi wahana baru. Field attraction_name dan price otomatis diambil dari tarif wahana (attraction_price_id), payment_method_name otomatis diambil dari payment_id. user_id otomatis diisi dari user yang sedang login, date otomatis diisi tanggal hari ini.
// @Tags Attraction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAttractionRequest true "Data transaksi wahana baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.AttractionResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /attractions [post]
func (h *AttractionHandler) CreateAttraction(c *gin.Context) {
	var req dto.CreateAttractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "token tidak valid")
		return
	}

	result, err := h.attractionService.CreateAttraction(req, actorID.(uuid.UUID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrAttractionPriceNotFound), errors.Is(err, service.ErrPaymentMethodNotFound):
			utils.Error(c, http.StatusNotFound, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	meta := utils.GetRequestMeta(c)
	_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "attractions",
		fmt.Sprintf("membuat transaksi wahana %s", result.AttractionName), req)

	utils.Success(c, http.StatusCreated, "transaksi wahana berhasil dibuat", result)
}

// UpdateAttraction godoc
// @Summary Ubah Transaksi Wahana
// @Description Mengubah data transaksi wahana. Field attraction_name, price, dan payment_method_name otomatis dihitung ulang berdasarkan attraction_price_id/payment_id yang baru.
// @Tags Attraction
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Attraction ID (UUID)"
// @Param request body dto.UpdateAttractionRequest true "Data transaksi wahana yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.AttractionResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /attractions/{id} [put]
func (h *AttractionHandler) UpdateAttraction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id transaksi wahana tidak valid")
		return
	}

	var req dto.UpdateAttractionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.attractionService.UpdateAttraction(id, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.Error(c, http.StatusNotFound, "transaksi wahana tidak ditemukan")
		case errors.Is(err, service.ErrAttractionPriceNotFound), errors.Is(err, service.ErrPaymentMethodNotFound):
			utils.Error(c, http.StatusNotFound, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "attractions",
			fmt.Sprintf("mengubah transaksi wahana %s", newResp.AttractionName), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "transaksi wahana berhasil diperbarui", newResp)
}

// DeleteAttraction godoc
// @Summary Hapus Transaksi Wahana
// @Description Menghapus transaksi wahana
// @Tags Attraction
// @Produce json
// @Security BearerAuth
// @Param id path string true "Attraction ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /attractions/{id} [delete]
func (h *AttractionHandler) DeleteAttraction(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id transaksi wahana tidak valid")
		return
	}

	if err := h.attractionService.DeleteAttraction(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "transaksi wahana tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "attractions",
			fmt.Sprintf("menghapus transaksi wahana id %s", id))
	}

	utils.Success(c, http.StatusOK, "transaksi wahana berhasil dihapus", nil)
}
