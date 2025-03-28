package status_transition_rule

import (
	"Group03-EX-StudentManagementAppBE/internal/handlers/base"
	"Group03-EX-StudentManagementAppBE/internal/services"
	"Group03-EX-StudentManagementAppBE/middleware"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	base.Handler
}

func NewHandler(service *services.Service) *Handler {
	return &Handler{
		Handler: base.NewHandler(service),
	}
}


func (h *Handler) RegisterRoutes(rg *gin.RouterGroup) {
	statusTransitionGroup := rg.Group("/v1/status-transition-rules")
	{
		statusTransitionGroup.PATCH("",middleware.UserAuthentication, h.OnOffTransition)
	}
}