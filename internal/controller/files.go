package controller

import (
	"screenresume/internal/models"
	"screenresume/internal/services"

	"github.com/go-fuego/fuego"
)

type FilesResources struct {
	FilesService services.FilesService
}

func (rs FilesResources) Routes(s *fuego.Server) {
	filesGroup := fuego.Group(s, "/files")

	fuego.Get(filesGroup, "/", rs.getAllFiles)
	fuego.Post(filesGroup, "/", rs.postFiles)

	fuego.Get(filesGroup, "/{id}", rs.getFiles)
	fuego.Put(filesGroup, "/{id}", rs.putFiles)
	fuego.Delete(filesGroup, "/{id}", rs.deleteFiles)
}

func (rs FilesResources) getAllFiles(c fuego.ContextNoBody) ([]models.Files, error) {
	return rs.FilesService.GetAllFiles(c.Context())
}

func (rs FilesResources) postFiles(c fuego.ContextWithBody[models.FilesCreate]) (models.Files, error) {
	body, err := c.Body()
	if err != nil {
		return models.Files{}, err
	}

	new, err := rs.FilesService.CreateFiles(c.Context(), body)
	if err != nil {
		return models.Files{}, err
	}

	return new, nil
}

func (rs FilesResources) getFiles(c fuego.ContextNoBody) (models.Files, error) {
	id := c.PathParam("id")

	return rs.FilesService.GetFiles(c.Context(), id)
}

func (rs FilesResources) putFiles(c fuego.ContextWithBody[models.FilesUpdate]) (models.Files, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.Files{}, err
	}

	new, err := rs.FilesService.UpdateFiles(c.Context(), id, body)
	if err != nil {
		return models.Files{}, err
	}

	return new, nil
}

func (rs FilesResources) deleteFiles(c fuego.ContextNoBody) (any, error) {
	return rs.FilesService.DeleteFiles(c.Context(), c.PathParam("id"))
}
