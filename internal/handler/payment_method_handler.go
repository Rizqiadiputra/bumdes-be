package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/middleware"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type PaymentMethodHandler struct {
	paymentMethodService service.PaymentMethodService
	userLogService       service.UserLogService
}

func NewPaymentMethodHandler(paymentMethodService service.PaymentMethodService, userLogService service.UserLogService) *PaymentMethodHandler {
	return &PaymentMethodHandler{paymentMethodService: paymentMethodService, userLogService: userLogService}
}

// ListPaymentMethods godoc
// @Summary List Payment Methods
// @Description Menampilkan daftar metode pembayaran, bisa difilter dengan pencarian method
// @Tags Payment Methods
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan method"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.PaymentMethodResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/payment-methods [get]
func (h *PaymentMethodHandler) ListPaymentMethods(c *gin.Context) {
	filter := repository.PaymentMethodFilter{Search: c.Query("search")}

	result, err := h.paymentMethodService.ListPaymentMethods(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar metode pembayaran", result)
}

// GetPaymentMethod godoc
// @Summary Detail Payment Method
// @Description Menampilkan detail satu metode pembayaran
// @Tags Payment Methods
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment Method ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.PaymentMethodResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /master-data/payment-methods/{id} [get]
func (h *PaymentMethodHandler) GetPaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id metode pembayaran tidak valid")
		return
	}

	result, err := h.paymentMethodService.GetPaymentMethod(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "metode pembayaran tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail metode pembayaran", result)
}

// CreatePaymentMethod godoc
// @Summary Buat Payment Method
// @Description Membuat metode pembayaran baru
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreatePaymentMethodRequest true "Data metode pembayaran baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.PaymentMethodResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/payment-methods [post]
func (h *PaymentMethodHandler) CreatePaymentMethod(c *gin.Context) {
	var req dto.CreatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.paymentMethodService.CreatePaymentMethod(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "payment-methods",
			fmt.Sprintf("membuat metode pembayaran %s", result.Method), req)
	}

	utils.Success(c, http.StatusCreated, "metode pembayaran berhasil dibuat", result)
}

// UpdatePaymentMethod godoc
// @Summary Ubah Payment Method
// @Description Mengubah data metode pembayaran
// @Tags Payment Methods
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment Method ID"
// @Param request body dto.UpdatePaymentMethodRequest true "Data metode pembayaran yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.PaymentMethodResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/payment-methods/{id} [put]
func (h *PaymentMethodHandler) UpdatePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id metode pembayaran tidak valid")
		return
	}

	var req dto.UpdatePaymentMethodRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.paymentMethodService.UpdatePaymentMethod(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "metode pembayaran tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "payment-methods",
			fmt.Sprintf("mengubah metode pembayaran %s", newResp.Method), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "metode pembayaran berhasil diperbarui", newResp)
}

// DeletePaymentMethod godoc
// @Summary Hapus Payment Method
// @Description Menghapus metode pembayaran
// @Tags Payment Methods
// @Produce json
// @Security BearerAuth
// @Param id path int true "Payment Method ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/payment-methods/{id} [delete]
func (h *PaymentMethodHandler) DeletePaymentMethod(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id metode pembayaran tidak valid")
		return
	}

	if err := h.paymentMethodService.DeletePaymentMethod(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "metode pembayaran tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "payment-methods",
			fmt.Sprintf("menghapus metode pembayaran id %d", id))
	}

	utils.Success(c, http.StatusOK, "metode pembayaran berhasil dihapus", nil)
}
