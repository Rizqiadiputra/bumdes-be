package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type UserLogHandler struct {
	logService service.UserLogService
}

func NewUserLogHandler(logService service.UserLogService) *UserLogHandler {
	return &UserLogHandler{logService: logService}
}

// ListLogs godoc
// @Summary List Aktivitas User
// @Description Menampilkan riwayat aktivitas user: membuka halaman (view), create, update, dan delete
// @Tags User Logs
// @Produce json
// @Security BearerAuth
// @Param user_id query string false "Filter berdasarkan user id (UUID)"
// @Param action query string false "Filter berdasarkan aksi (view, create, update, delete)"
// @Param module query string false "Filter berdasarkan module (misal: accounts, roles, me)"
// @Param page query int false "Halaman (default 1)"
// @Param limit query int false "Jumlah data per halaman (default 20, maksimal 100)"
// @Success 200 {object} utils.SuccessResponse{data=dto.PaginatedUserLogResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /logs [get]
func (h *UserLogHandler) ListLogs(c *gin.Context) {
	var filter repository.UserLogFilter

	if v := c.Query("user_id"); v != "" {
		if id, err := uuid.Parse(v); err == nil {
			filter.UserID = id
		}
	}
	filter.Action = c.Query("action")
	filter.Module = c.Query("module")
	filter.Page, _ = strconv.Atoi(c.DefaultQuery("page", "1"))
	filter.Limit, _ = strconv.Atoi(c.DefaultQuery("limit", "20"))

	result, err := h.logService.ListLogs(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar log", result)
}
