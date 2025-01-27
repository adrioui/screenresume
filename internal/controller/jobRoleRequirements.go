package controller

import (
	"screenresume/internal/models"
	"screenresume/internal/services"

	"github.com/go-fuego/fuego"
)

type JobRoleRequirementsResources struct {
	JobRoleRequirementsService services.JobRoleRequirementsService
}

func (rs JobRoleRequirementsResources) Routes(s *fuego.Server) {
	JobRoleRequirementsGroup := fuego.Group(s, "/job_role_requirements")

	fuego.Get(JobRoleRequirementsGroup, "/", rs.getAllJobRoleRequirements)
	fuego.Post(JobRoleRequirementsGroup, "/", rs.postJobRoleRequirements)

	fuego.Get(JobRoleRequirementsGroup, "/{id}", rs.getJobRoleRequirements)
	fuego.Put(JobRoleRequirementsGroup, "/{id}", rs.putJobRoleRequirements)
	fuego.Delete(JobRoleRequirementsGroup, "/{id}", rs.deleteJobRoleRequirements)
}

func (rs JobRoleRequirementsResources) getAllJobRoleRequirements(c fuego.ContextNoBody) ([]models.JobRoleRequirements, error) {
	return rs.JobRoleRequirementsService.GetAllJobRoleRequirements(c.Context())
}

func (rs JobRoleRequirementsResources) postJobRoleRequirements(c fuego.ContextWithBody[models.JobRoleRequirementsCreate]) (models.JobRoleRequirements, error) {
	body, err := c.Body()
	if err != nil {
		return models.JobRoleRequirements{}, err
	}

	new, err := rs.JobRoleRequirementsService.CreateJobRoleRequirements(c.Context(), body)
	if err != nil {
		return models.JobRoleRequirements{}, err
	}

	return new, nil
}

func (rs JobRoleRequirementsResources) getJobRoleRequirements(c fuego.ContextNoBody) (models.JobRoleRequirements, error) {
	jobRoleID := c.PathParam("job_role_id")
	skillID := c.PathParam("skill_id")

	return rs.JobRoleRequirementsService.GetJobRoleRequirements(c.Context(), jobRoleID, skillID)
}

func (rs JobRoleRequirementsResources) putJobRoleRequirements(c fuego.ContextWithBody[models.JobRoleRequirementsUpdate]) (models.JobRoleRequirements, error) {
	jobRoleID := c.PathParam("job_role_id")
	skillID := c.PathParam("skill_id")

	body, err := c.Body()
	if err != nil {
		return models.JobRoleRequirements{}, err
	}

	new, err := rs.JobRoleRequirementsService.UpdateJobRoleRequirements(c.Context(), jobRoleID, skillID, body)
	if err != nil {
		return models.JobRoleRequirements{}, err
	}

	return new, nil
}

func (rs JobRoleRequirementsResources) deleteJobRoleRequirements(c fuego.ContextNoBody) (any, error) {
	return rs.JobRoleRequirementsService.DeleteJobRoleRequirements(c.Context(), c.PathParam("job_role_id"), c.PathParam("skill_id"))
}
