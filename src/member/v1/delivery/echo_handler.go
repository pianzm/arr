package delivery

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/pianzm/arr/helper"
	"github.com/pianzm/arr/src/member/v1/model"
	"github.com/pianzm/arr/src/member/v1/usecase"
)

type HTTPHandler struct {
	MemberUsecase usecase.MemberUsecase
}

func NewHandler(uc usecase.MemberUsecase) *HTTPHandler {
	return &HTTPHandler{
		MemberUsecase: uc,
	}
}

func (h *HTTPHandler) Mount(group *echo.Group) {
	group.GET("", h.GetMembers)
	group.GET("/status/:reqId", h.getStatus)
	group.GET("/download", h.download)
}

func (h *HTTPHandler) GetMembers(c echo.Context) error {
	params := model.QueryParameters{}
	if err := c.Bind(&params); err != nil {
		return helper.NewJSONResponse(http.StatusBadRequest, "invalid query parameters").JSON(c)
	}
	reqID := helper.GetRequestID(c)
	response := model.StatusRequest{
		RequestID: reqID,
	}
	qParam := model.QueueStatus{
		RequestID: reqID,
		Parameter: params,
	}
	if err := h.MemberUsecase.Publish(c.Request().Context(), &qParam); err != nil {
		return helper.NewJSONResponse(http.StatusBadRequest, "failed publish message to redis").JSON(c)
	}

	return helper.NewJSONResponse(http.StatusAccepted, "task accepted", response).JSON(c)
}

func (h *HTTPHandler) getStatus(c echo.Context) error {
	reqID := c.Param("reqId")
	result, err := h.MemberUsecase.GetStatus(c.Request().Context(), reqID)
	if err != nil {
		return helper.NewJSONResponse(http.StatusBadRequest, "failed get key from redis").JSON(c)
	}
	message := "task found"
	if result.Completed {
		message = "task completed"
	}
	return helper.NewJSONResponse(http.StatusOK, message, result).JSON(c)
}

func (h *HTTPHandler) download(c echo.Context) error {
	return nil
}
