package controller

import (
	"screenresume/internal/models"
	"screenresume/internal/services"

	"github.com/go-fuego/fuego"
)

type CandidatesResources struct {
	CandidatesService services.CandidatesService
}

func (rs CandidatesResources) Routes(s *fuego.Server) {
	candidatesGroup := fuego.Group(s, "/candidates")

	fuego.Get(candidatesGroup, "/", rs.getAllCandidates)
	fuego.Post(candidatesGroup, "/", rs.postCandidates)

	fuego.Get(candidatesGroup, "/{id}", rs.getCandidates)
	fuego.Put(candidatesGroup, "/{id}", rs.putCandidates)
	fuego.Delete(candidatesGroup, "/{id}", rs.deleteCandidates)
}

func (rs CandidatesResources) getAllCandidates(c fuego.ContextNoBody) ([]models.Candidates, error) {
	return rs.CandidatesService.GetAllCandidates(c.Context())
}

func (rs CandidatesResources) postCandidates(c fuego.ContextWithBody[models.CandidatesCreate]) (models.Candidates, error) {
	body, err := c.Body()
	if err != nil {
		return models.Candidates{}, err
	}

	new, err := rs.CandidatesService.CreateCandidates(c.Context(), body)
	if err != nil {
		return models.Candidates{}, err
	}

	return new, nil
}

func (rs CandidatesResources) getCandidates(c fuego.ContextNoBody) (models.Candidates, error) {
	id := c.PathParam("id")

	return rs.CandidatesService.GetCandidates(c.Context(), id)
}

func (rs CandidatesResources) putCandidates(c fuego.ContextWithBody[models.CandidatesUpdate]) (models.Candidates, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.Candidates{}, err
	}

	new, err := rs.CandidatesService.UpdateCandidates(c.Context(), id, body)
	if err != nil {
		return models.Candidates{}, err
	}

	return new, nil
}

func (rs CandidatesResources) deleteCandidates(c fuego.ContextNoBody) (any, error) {
	return rs.CandidatesService.DeleteCandidates(c.Context(), c.PathParam("id"))
}
