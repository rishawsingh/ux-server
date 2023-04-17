package group

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/remotestate/golang/internal"
	"github.com/remotestate/golang/models"
	"github.com/remotestate/golang/services"
	"net/http"
)

type Controller struct {
	logger            *internal.Logger
	groupService      services.Group
	locationService   services.Location
	attachmentService services.Attachment
	tx                *internal.Transactor
	tracer            *internal.Tracer
}

func NewController(logger *internal.Logger,
	groupService services.Group,
	tx *internal.Transactor,
	tracer *internal.Tracer) *Controller {
	return &Controller{
		logger:       logger,
		groupService: groupService,
		tx:           tx,
		tracer:       tracer,
	}
}

func (c *Controller) CreateGroup(ctx *gin.Context) {
	groupDetails := models.GroupDetails{}
	if parseErr := ctx.ShouldBindJSON(&groupDetails); parseErr != nil {
		ctx.AbortWithError(http.StatusBadRequest, errors.Errorf("%s, stack: %+v",
			"failed to parse request body.", parseErr))
		return
	}

	userID := fmt.Sprintf("%v", ctx.MustGet("userID"))

	groupID, groupErr := c.groupService.CreateGroup(groupDetails, userID, nil)
	if groupErr != nil {
		_ = ctx.AbortWithError(http.StatusInternalServerError, errors.Errorf("%s, stack: %+v",
			"failed to create group", groupErr))
		return
	}

	ctx.JSON(http.StatusOK, groupID)
}
