package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type BillingHistoryHandler struct {
	billingService service.BillingService
}

func NewBillingHistoryHandler(billingService service.BillingService) *BillingHistoryHandler {
	return &BillingHistoryHandler{billingService: billingService}
}

// ListBillingHistories godoc
// @Summary List Billing History
// @Description Menampilkan daftar tagihan billing tenant dengan pagination, bisa difilter dengan tenant_id dan status
// @Tags Billing Histories
// @Produce json
// @Security BearerAuth
// @Param tenant_id query string false "Filter berdasarkan Tenant ID (UUID)"
// @Param status query string false "Filter berdasarkan status (Unpaid, Paid, Partially Paid)"
// @Param page query int false "Halaman (default 1)"
// @Param limit query int false "Jumlah data per halaman (default 20, maksimal 100)"
// @Success 200 {object} utils.SuccessResponse{data=dto.PaginatedBillingHistoryResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /billing-histories [get]
func (h *BillingHistoryHandler) ListBillingHistories(c *gin.Context) {
	filter := repository.BillingHistoryFilter{Status: c.Query("status")}

	if tenantIDStr := c.Query("tenant_id"); tenantIDStr != "" {
		tenantID, err := uuid.Parse(tenantIDStr)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "tenant_id tidak valid")
			return
		}
		filter.TenantID = &tenantID
	}

	filter.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	filter.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.billingService.ListBillingHistories(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar billing history", result)
}

// GetBillingHistory godoc
// @Summary Detail Billing History
// @Description Menampilkan detail satu billing history
// @Tags Billing Histories
// @Produce json
// @Security BearerAuth
// @Param id path string true "Billing History ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.BillingHistoryResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /billing-histories/{id} [get]
func (h *BillingHistoryHandler) GetBillingHistory(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id billing history tidak valid")
		return
	}

	result, err := h.billingService.GetBillingHistory(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "billing history tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail billing history", result)
}
