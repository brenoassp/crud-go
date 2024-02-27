package clientsctrl

import (
	"encoding/json"
	"strconv"

	"github.com/brenoassp/crud-go/domain"
	"github.com/brenoassp/crud-go/domain/clients"
	"github.com/gofiber/fiber/v2"
)

type Controller struct {
	clientsService *clients.Service
}

func NewController(clientsService *clients.Service) *Controller {
	return &Controller{
		clientsService: clientsService,
	}
}

func (c *Controller) GetClients(ctx *fiber.Ctx) error {
	var page int = 1
	var pageSize int = 10
	var err error
	m := ctx.Queries()
	if m["page"] != "" {
		page, err = strconv.Atoi(m["page"])
		if err != nil {
			return err
		}
	}

	if m["size"] != "" {
		pageSize, err = strconv.Atoi(m["size"])
		if err != nil {
			return err
		}
	}

	response, err := c.clientsService.GetClients(ctx.Context(), page, pageSize)
	if err != nil {
		return err
	}

	return ctx.JSON(response)
}

func (c *Controller) CreateClient(ctx *fiber.Ctx) error {
	var client domain.CreateClientEvent
	err := json.Unmarshal(ctx.Body(), &client)
	if err != nil {
		return err
	}

	err = c.clientsService.CreateClientEvent(ctx.Context(), client)
	if err != nil {
		return err
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (c *Controller) DeleteClient(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return err
	}

	err = c.clientsService.DeleteClient(ctx.Context(), id)
	if err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (c *Controller) UpdateClient(ctx *fiber.Ctx) error {
	id, err := strconv.Atoi(ctx.Params("id"))
	if err != nil {
		return err
	}

	var client domain.Client
	err = json.Unmarshal(ctx.Body(), &client)
	if err != nil {
		return err
	}
	client.ID = id

	err = c.clientsService.UpdateClient(ctx.Context(), client)
	if err != nil {
		return err
	}
	return ctx.SendStatus(fiber.StatusOK)
}
