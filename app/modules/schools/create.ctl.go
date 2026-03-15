package schools

import (
	"eduflow/app/utils"
	"eduflow/app/utils/base"

	"github.com/gin-gonic/gin"
)

type CreateRequest struct {
	Name        string  `json:"name" binding:"required"`
	LogoURL     *string `json:"logo_url"`
	ThemeColor  *string `json:"theme_color"`
	Address     *string `json:"address"`
	Description *string `json:"description"`
}

func (c *Controller) Create(ctx *gin.Context) {
	span, log := utils.LogSpanFromGin(ctx)
	defer span.End()

	var req CreateRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		base.BadRequest(ctx, "invalid-request", nil)
		return
	}

	school, err := c.svc.Create(ctx.Request.Context(), req.Name, req.LogoURL, req.ThemeColor, req.Address, req.Description)
	if err != nil {
		c.handleServiceError(ctx, log, err, "school-create-failed")
		return
	}

	base.Success(ctx, school, "success")
}

func (c *Controller) CreateSchoolController(ctx *gin.Context) {
	c.Create(ctx)
}

func (c *Controller) SchoolsCreate(ctx *gin.Context) {
	c.Create(ctx)
}
