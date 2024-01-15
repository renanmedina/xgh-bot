package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/renanmedina/xgh-bot/gohorse"
	"github.com/renanmedina/xgh-bot/slack"
)

type CommandRequest struct {
	CommandType   string `json:"command" form:"command" binding:"required"`
	RequestedType string `json:"text" form:"text" binding:"required"`
}

func (c CommandRequest) isXGH() bool {
	return c.CommandType == "/xgh"
}

func (c CommandRequest) isRandomAxiomType() bool {
	return c.RequestedType == "random"
}

func (c CommandRequest) axiomNumber() (int, error) {
	number, err := strconv.Atoi(c.RequestedType)
	if err != nil {
		return 0, errors.New("Invalid axiom number parameter")
	}

	return number, nil
}

func SlackBotHandler(c *gin.Context) {
	var request CommandRequest
	c.Bind(&request)

	if !request.isXGH() {
		c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Command %s not available", request.CommandType)})
		return
	}

	service := gohorse.NewAxiomsService()

	if request.isRandomAxiomType() {
		axiom := service.PickOneRandom()
		c.JSON(http.StatusOK, slack.NewSlackSimpleResponse(slack.EPHEMERAL, axiom.ToQuote()))
		return
	}

	axiom_number, err := request.axiomNumber()

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	axiom, err := service.GetByNumber(axiom_number)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, slack.NewSlackSimpleResponse(slack.EPHEMERAL, axiom.ToQuote()))
}
