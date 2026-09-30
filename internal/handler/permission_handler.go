package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type PermissionHandler struct {
	permissionService service.PermissionService
}

func NewPermissionHandler(permissionService service.PermissionService) *PermissionHandler {
	return &PermissionHandler{permissionService: permissionService}
}

// ListPermissions godoc
// @Summary List Permission
// @Description Menampilkan seluruh permission yang tersedia di sistem
// @Tags Role Management
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse{data=[]dto.PermissionResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /permissions [get]
func (h *PermissionHandler) ListPermissions(c *gin.Context) {
	result, err := h.permissionService.ListPermissions()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar permission", result)
}
