package controller

import (
	"screenresume/internal/models"
	"screenresume/internal/services"

	"github.com/go-fuego/fuego"
)

type ApplicationResources struct {
	ApplicationService services.ApplicationService
}

func (rs ApplicationResources) Routes(s *fuego.Server) {
	applicationsGroup := fuego.Group(s, "/applications")

	fuego.Get(applicationsGroup, "/", rs.getAllApplication)
	fuego.Post(applicationsGroup, "/", rs.postApplication)

	fuego.Get(applicationsGroup, "/{id}", rs.getApplication)
	fuego.Put(applicationsGroup, "/{id}", rs.putApplication)
	fuego.Delete(applicationsGroup, "/{id}", rs.deleteApplication)
}

func (rs ApplicationResources) getAllApplication(c fuego.ContextNoBody) ([]models.Application, error) {
	return rs.ApplicationService.GetAllApplication(c.Context())
}

func (rs ApplicationResources) postApplication(c fuego.ContextWithBody[models.ApplicationCreate]) (models.Application, error) {
	body, err := c.Body()
	if err != nil {
		return models.Application{}, err
	}

	new, err := rs.ApplicationService.CreateApplication(c.Context(), body)
	if err != nil {
		return models.Application{}, err
	}

	return new, nil
}

func (rs ApplicationResources) getApplication(c fuego.ContextNoBody) (models.Application, error) {
	id := c.PathParam("id")

	return rs.ApplicationService.GetApplication(c.Context(), id)
}

func (rs ApplicationResources) putApplication(c fuego.ContextWithBody[models.ApplicationUpdate]) (models.Application, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.Application{}, err
	}

	new, err := rs.ApplicationService.UpdateApplication(c.Context(), id, body)
	if err != nil {
		return models.Application{}, err
	}

	return new, nil
}

func (rs ApplicationResources) deleteApplication(c fuego.ContextNoBody) (any, error) {
	return rs.ApplicationService.DeleteApplication(c.Context(), c.PathParam("id"))
}
