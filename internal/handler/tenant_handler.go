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

const tenantSearchDefaultLimit = 10

type TenantHandler struct {
	tenantService       service.TenantService
	tenantTypeService   service.TenantTypeService
	kiosLocationService service.KiosLocationService
	userLogService      service.UserLogService
}

func NewTenantHandler(
	tenantService service.TenantService,
	tenantTypeService service.TenantTypeService,
	kiosLocationService service.KiosLocationService,
	userLogService service.UserLogService,
) *TenantHandler {
	return &TenantHandler{
		tenantService:       tenantService,
		tenantTypeService:   tenantTypeService,
		kiosLocationService: kiosLocationService,
		userLogService:      userLogService,
	}
}

// SearchTenantTypes godoc
// @Summary Cari Tenant Type (untuk form Tenant)
// @Description Mencari tenant type berdasarkan kolom type, dipakai untuk pengisian tenant_type_id di form tenant
// @Tags Tenants
// @Produce json
// @Security BearerAuth
// @Param q query string false "Kata kunci pencarian (kolom type)"
// @Param limit query int false "Jumlah maksimal hasil (default 10)"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.TenantTypeResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /tenants/tenant-types [get]
func (h *TenantHandler) SearchTenantTypes(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(tenantSearchDefaultLimit)))
	if err != nil || limit <= 0 {
		limit = tenantSearchDefaultLimit
	}

	filter := repository.TenantTypeFilter{Search: c.Query("q"), Limit: limit}

	result, err := h.tenantTypeService.ListTenantTypes(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar tipe tenant", result)
}

// SearchKiosLocations godoc
// @Summary Cari Kios Location (untuk form Tenant)
// @Description Mencari lokasi kios berdasarkan kolom kios, dipakai untuk pengisian kios_location_id di form tenant
// @Tags Tenants
// @Produce json
// @Security BearerAuth
// @Param q query string false "Kata kunci pencarian (kolom kios)"
// @Param limit query int false "Jumlah maksimal hasil (default 10)"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.KiosLocationResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /tenants/kios_locations [get]
func (h *TenantHandler) SearchKiosLocations(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", strconv.Itoa(tenantSearchDefaultLimit)))
	if err != nil || limit <= 0 {
		limit = tenantSearchDefaultLimit
	}

	filter := repository.KiosLocationFilter{Search: c.Query("q"), Limit: limit}

	result, err := h.kiosLocationService.ListKiosLocations(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar lokasi kios", result)
}

// ListTenants godoc
// @Summary List Tenant
// @Description Menampilkan daftar tenant (penyewa kios) dengan pagination, bisa difilter dengan pencarian nama
// @Tags Tenants
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan nama tenant"
// @Param page query int false "Halaman (default 1)"
// @Param limit query int false "Jumlah data per halaman (default 20, maksimal 100)"
// @Success 200 {object} utils.SuccessResponse{data=dto.PaginatedTenantResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /tenants [get]
func (h *TenantHandler) ListTenants(c *gin.Context) {
	filter := repository.TenantFilter{Search: c.Query("search")}
	filter.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	filter.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.tenantService.ListTenants(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar tenant", result)
}

// GetTenant godoc
// @Summary Detail Tenant
// @Description Menampilkan detail satu tenant
// @Tags Tenants
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.TenantResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /tenants/{id} [get]
func (h *TenantHandler) GetTenant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tenant tidak valid")
		return
	}

	result, err := h.tenantService.GetTenant(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "tenant tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail tenant", result)
}

// CreateTenant godoc
// @Summary Buat Tenant
// @Description Membuat tenant baru. Field type otomatis diambil dari tenant_type_id, kios_name otomatis diambil dari kios_location_id.
// @Tags Tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTenantRequest true "Data tenant baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.TenantResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /tenants [post]
func (h *TenantHandler) CreateTenant(c *gin.Context) {
	var req dto.CreateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.tenantService.CreateTenant(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrTenantTypeNotFound), errors.Is(err, service.ErrKiosLocationNotFound):
			utils.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrInvalidStartPaymentDate), errors.Is(err, service.ErrInvalidEndPeriodeDate):
			utils.Error(c, http.StatusBadRequest, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "tenants",
			fmt.Sprintf("membuat tenant %s", result.Name), req)
	}

	utils.Success(c, http.StatusCreated, "tenant berhasil dibuat", result)
}

// UpdateTenant godoc
// @Summary Ubah Tenant
// @Description Mengubah data tenant. Field type dan kios_name otomatis dihitung ulang berdasarkan tenant_type_id/kios_location_id yang baru.
// @Tags Tenants
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID (UUID)"
// @Param request body dto.UpdateTenantRequest true "Data tenant yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.TenantResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /tenants/{id} [put]
func (h *TenantHandler) UpdateTenant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tenant tidak valid")
		return
	}

	var req dto.UpdateTenantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.tenantService.UpdateTenant(id, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.Error(c, http.StatusNotFound, "tenant tidak ditemukan")
		case errors.Is(err, service.ErrTenantTypeNotFound), errors.Is(err, service.ErrKiosLocationNotFound):
			utils.Error(c, http.StatusNotFound, err.Error())
		case errors.Is(err, service.ErrInvalidStartPaymentDate), errors.Is(err, service.ErrInvalidEndPeriodeDate):
			utils.Error(c, http.StatusBadRequest, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "tenants",
			fmt.Sprintf("mengubah tenant %s", newResp.Name), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "tenant berhasil diperbarui", newResp)
}

// DeleteTenant godoc
// @Summary Hapus Tenant
// @Description Menghapus tenant
// @Tags Tenants
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /tenants/{id} [delete]
func (h *TenantHandler) DeleteTenant(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tenant tidak valid")
		return
	}

	if err := h.tenantService.DeleteTenant(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tenant tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "tenants",
			fmt.Sprintf("menghapus tenant id %s", id))
	}

	utils.Success(c, http.StatusOK, "tenant berhasil dihapus", nil)
}
