package controller

import (
	"screenresume/internal/models"
	"screenresume/internal/services"

	"github.com/go-fuego/fuego"
)

type SkillsResources struct {
	SkillsService services.SkillsService
}

func (rs SkillsResources) Routes(s *fuego.Server) {
	skillsGroup := fuego.Group(s, "/skills")

	fuego.Get(skillsGroup, "/", rs.getAllSkills)
	fuego.Post(skillsGroup, "/", rs.postSkills)

	fuego.Get(skillsGroup, "/{id}", rs.getSkills)
	fuego.Put(skillsGroup, "/{id}", rs.putSkills)
	fuego.Delete(skillsGroup, "/{id}", rs.deleteSkills)
}

func (rs SkillsResources) getAllSkills(c fuego.ContextNoBody) ([]models.Skills, error) {
	return rs.SkillsService.GetAllSkills(c.Context())
}

func (rs SkillsResources) postSkills(c fuego.ContextWithBody[models.SkillsCreate]) (models.Skills, error) {
	body, err := c.Body()
	if err != nil {
		return models.Skills{}, err
	}

	new, err := rs.SkillsService.CreateSkills(c.Context(), body)
	if err != nil {
		return models.Skills{}, err
	}

	return new, nil
}

func (rs SkillsResources) getSkills(c fuego.ContextNoBody) (models.Skills, error) {
	id := c.PathParam("id")

	return rs.SkillsService.GetSkills(c.Context(), id)
}

func (rs SkillsResources) putSkills(c fuego.ContextWithBody[models.SkillsUpdate]) (models.Skills, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.Skills{}, err
	}

	new, err := rs.SkillsService.UpdateSkills(c.Context(), id, body)
	if err != nil {
		return models.Skills{}, err
	}

	return new, nil
}

func (rs SkillsResources) deleteSkills(c fuego.ContextNoBody) (any, error) {
	return rs.SkillsService.DeleteSkills(c.Context(), c.PathParam("id"))
}
