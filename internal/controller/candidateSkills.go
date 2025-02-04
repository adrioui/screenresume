package controller

import (
	"screenresume/internal/models"
	"screenresume/internal/services"

	"github.com/go-fuego/fuego"
)

type CandidateSkillsResources struct {
	CandidateSkillsService services.CandidateSkillsService
}

func (rs CandidateSkillsResources) Routes(s *fuego.Server) {
	CandidateSkillsGroup := fuego.Group(s, "/candidate-skills")

	fuego.Get(CandidateSkillsGroup, "/", rs.getAllCandidateSkills)
	fuego.Post(CandidateSkillsGroup, "/", rs.postCandidateSkills)

	fuego.Get(CandidateSkillsGroup, "/{id}", rs.getCandidateSkills)
	fuego.Put(CandidateSkillsGroup, "/{id}", rs.putCandidateSkills)
	fuego.Delete(CandidateSkillsGroup, "/{id}", rs.deleteCandidateSkills)
}

func (rs CandidateSkillsResources) getAllCandidateSkills(c fuego.ContextNoBody) ([]models.CandidateSkills, error) {
	return rs.CandidateSkillsService.GetAllCandidateSkills(c.Context())
}

func (rs CandidateSkillsResources) postCandidateSkills(c fuego.ContextWithBody[models.CandidateSkillsCreate]) (models.CandidateSkills, error) {
	body, err := c.Body()
	if err != nil {
		return models.CandidateSkills{}, err
	}

	new, err := rs.CandidateSkillsService.CreateCandidateSkills(c.Context(), body)
	if err != nil {
		return models.CandidateSkills{}, err
	}

	return new, nil
}

func (rs CandidateSkillsResources) getCandidateSkills(c fuego.ContextNoBody) (models.CandidateSkills, error) {
	jobRoleID := c.PathParam("candidate_id")
	skillID := c.PathParam("skill_id")

	return rs.CandidateSkillsService.GetCandidateSkills(c.Context(), jobRoleID, skillID)
}

func (rs CandidateSkillsResources) putCandidateSkills(c fuego.ContextWithBody[models.CandidateSkillsUpdate]) (models.CandidateSkills, error) {
	jobRoleID := c.PathParam("candidate_id")
	skillID := c.PathParam("skill_id")

	body, err := c.Body()
	if err != nil {
		return models.CandidateSkills{}, err
	}

	new, err := rs.CandidateSkillsService.UpdateCandidateSkills(c.Context(), jobRoleID, skillID, body)
	if err != nil {
		return models.CandidateSkills{}, err
	}

	return new, nil
}

func (rs CandidateSkillsResources) deleteCandidateSkills(c fuego.ContextNoBody) (any, error) {
	return rs.CandidateSkillsService.DeleteCandidateSkills(c.Context(), c.PathParam("candidate_id"), c.PathParam("skill_id"))
}
