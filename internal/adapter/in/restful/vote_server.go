package restful

import (
	"net/http"

	"github.com/CorrectRoadH/Likit/internal/application/domain"
	"github.com/CorrectRoadH/Likit/internal/port/in"
	"github.com/labstack/echo/v4"
)

type VoteServer struct {
	voteUseCase in.VoteUseCase
}

func NewVoteServer(voteUseCase in.VoteUseCase) *VoteServer {
	s := &VoteServer{
		voteUseCase: voteUseCase,
	}
	return s
}

type VoteUsersEvent struct {
	List []string `json:"userList"`
}

type VoteEvent struct {
	BusinessId string `json:"businessId"`
	MessageId  string `json:"messageId"`
	UserId     string `json:"userId"`
}

type ErrorResponse struct {
	Status  int    `json:"status"`
	Message string `json:"msg"`
}

type VoteResponse struct {
	Status int   `json:"status"`
	Count  int64 `json:"count"`
}

func (v *VoteServer) Vote(c echo.Context) error {
	var event VoteEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	val, err := v.voteUseCase.Vote(c.Request().Context(), event.BusinessId, event.MessageId, event.UserId)
	if err != nil {
		if err == domain.ErrBusinessNotExist {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})

		}

		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// need to define a struct for response
	return c.JSON(http.StatusOK, VoteResponse{
		Status: http.StatusOK,
		Count:  int64(val),
	})
}

func (v *VoteServer) UnVote(c echo.Context) error {
	var event VoteEvent
	if err := c.Bind(&event); err != nil {
		return c.JSON(http.StatusBadRequest, err)
	}

	val, err := v.voteUseCase.UnVote(c.Request().Context(), event.BusinessId, event.MessageId, event.UserId)

	if err != nil {
		if err == domain.ErrBusinessNotExist {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})

		}

		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// need to define a struct for response
	return c.JSON(http.StatusOK, VoteResponse{
		Status: http.StatusOK,
		Count:  int64(val),
	})
}

func (v *VoteServer) Count(c echo.Context) error {
	businessId := c.Param("businessId")
	messageId := c.Param("messageId")

	count, err := v.voteUseCase.Count(c.Request().Context(), businessId, messageId)
	if err != nil {
		if err == domain.ErrBusinessNotExist {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})

		}

		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	// need to define a struct for response
	return c.JSON(http.StatusOK, count)
}

func (v *VoteServer) ListUser(c echo.Context) error {
	businessId := c.Param("businessId")
	messageId := c.Param("messageId")

	voteUsers, err := v.voteUseCase.VotedUsers(c.Request().Context(), businessId, messageId)
	if err != nil {
		if err == domain.ErrBusinessNotExist {
			return c.JSON(http.StatusBadRequest, ErrorResponse{
				Status:  http.StatusBadRequest,
				Message: err.Error(),
			})

		}

		return c.JSON(http.StatusInternalServerError, ErrorResponse{
			Status:  http.StatusInternalServerError,
			Message: err.Error(),
		})
	}

	result := VoteUsersEvent{
		List: voteUsers,
	}

	return c.JSON(http.StatusOK, result)
}

func (v *VoteServer) register(g *echo.Group) error {
	g.GET("/count/:businessId/:messageId", v.Count)
	g.GET("/list/:businessId/:messageId", v.ListUser)

	g.POST("/vote", v.Vote)
	g.POST("/unvote", v.UnVote)

	return nil
}
