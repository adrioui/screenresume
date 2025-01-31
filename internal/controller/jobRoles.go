package controller

import (
	"screenresume/internal/models"
	"screenresume/internal/services"

	"github.com/go-fuego/fuego"
)

type JobRolesResources struct {
	JobRolesService services.JobRolesService
}

func (rs JobRolesResources) Routes(s *fuego.Server) {
	JobRolesGroup := fuego.Group(s, "/job-roles")

	fuego.Get(JobRolesGroup, "/", rs.getAllJobRoles)
	fuego.Post(JobRolesGroup, "/", rs.postJobRoles)

	fuego.Get(JobRolesGroup, "/{id}", rs.getJobRoles)
	fuego.Put(JobRolesGroup, "/{id}", rs.putJobRoles)
	fuego.Delete(JobRolesGroup, "/{id}", rs.deleteJobRoles)
}

func (rs JobRolesResources) getAllJobRoles(c fuego.ContextNoBody) ([]models.JobRoles, error) {
	return rs.JobRolesService.GetAllJobRoles(c.Context())
}

func (rs JobRolesResources) postJobRoles(c fuego.ContextWithBody[models.JobRolesCreate]) (models.JobRoles, error) {
	body, err := c.Body()
	if err != nil {
		return models.JobRoles{}, err
	}

	new, err := rs.JobRolesService.CreateJobRoles(c.Context(), body)
	if err != nil {
		return models.JobRoles{}, err
	}

	return new, nil
}

func (rs JobRolesResources) getJobRoles(c fuego.ContextNoBody) (models.JobRoles, error) {
	id := c.PathParam("id")

	return rs.JobRolesService.GetJobRoles(c.Context(), id)
}

func (rs JobRolesResources) putJobRoles(c fuego.ContextWithBody[models.JobRolesUpdate]) (models.JobRoles, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.JobRoles{}, err
	}

	new, err := rs.JobRolesService.UpdateJobRoles(c.Context(), id, body)
	if err != nil {
		return models.JobRoles{}, err
	}

	return new, nil
}

func (rs JobRolesResources) deleteJobRoles(c fuego.ContextNoBody) (any, error) {
	return rs.JobRolesService.DeleteJobRoles(c.Context(), c.PathParam("id"))
}
