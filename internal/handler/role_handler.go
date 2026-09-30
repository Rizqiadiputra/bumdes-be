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
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type RoleHandler struct {
	roleService    service.RoleService
	userLogService service.UserLogService
}

func NewRoleHandler(roleService service.RoleService, userLogService service.UserLogService) *RoleHandler {
	return &RoleHandler{roleService: roleService, userLogService: userLogService}
}

// ListRoles godoc
// @Summary List Role
// @Description Menampilkan daftar role beserta permission yang bisa diakses oleh masing-masing role
// @Tags Role Management
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse{data=[]dto.RoleResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /roles [get]
func (h *RoleHandler) ListRoles(c *gin.Context) {
	result, err := h.roleService.ListRoles()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar role", result)
}

// GetRole godoc
// @Summary Detail Role
// @Description Menampilkan detail satu role beserta permission yang dimiliki
// @Tags Role Management
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.RoleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /roles/{id} [get]
func (h *RoleHandler) GetRole(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id role tidak valid")
		return
	}

	result, err := h.roleService.GetRole(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "role tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail role", result)
}

// UpdatePermissions godoc
// @Summary Ubah Permission Role
// @Description Mengganti seluruh permission yang dimiliki sebuah role (assign/revoke)
// @Tags Role Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Role ID"
// @Param request body dto.UpdateRolePermissionsRequest true "Daftar ID permission yang dimiliki role"
// @Success 200 {object} utils.SuccessResponse{data=dto.RoleResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /roles/{id}/permissions [put]
func (h *RoleHandler) UpdatePermissions(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id role tidak valid")
		return
	}

	var req dto.UpdateRolePermissionsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.roleService.UpdatePermissions(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "role tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "roles",
			fmt.Sprintf("mengubah permission role %s", newResp.Name), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "permission role berhasil diperbarui", newResp)
}
