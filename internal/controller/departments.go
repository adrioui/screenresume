package controller

import (
	"screenresume/internal/models"
	"screenresume/internal/services"

	"github.com/go-fuego/fuego"
)

type DepartmentsResources struct {
	DepartmentsService services.DepartmentsService
}

func (rs DepartmentsResources) Routes(s *fuego.Server) {
	departmentsGroup := fuego.Group(s, "/departments")

	fuego.Get(departmentsGroup, "/", rs.getAllDepartments)
	fuego.Post(departmentsGroup, "/", rs.postDepartments)

	fuego.Get(departmentsGroup, "/{id}", rs.getDepartments)
	fuego.Put(departmentsGroup, "/{id}", rs.putDepartments)
	fuego.Delete(departmentsGroup, "/{id}", rs.deleteDepartments)
}

func (rs DepartmentsResources) getAllDepartments(c fuego.ContextNoBody) ([]models.Departments, error) {
	return rs.DepartmentsService.GetAllDepartments(c.Context())
}

func (rs DepartmentsResources) postDepartments(c fuego.ContextWithBody[models.DepartmentsCreate]) (models.Departments, error) {
	body, err := c.Body()
	if err != nil {
		return models.Departments{}, err
	}

	new, err := rs.DepartmentsService.CreateDepartments(c.Context(), body)
	if err != nil {
		return models.Departments{}, err
	}

	return new, nil
}

func (rs DepartmentsResources) getDepartments(c fuego.ContextNoBody) (models.Departments, error) {
	id := c.PathParam("id")

	return rs.DepartmentsService.GetDepartments(c.Context(), id)
}

func (rs DepartmentsResources) putDepartments(c fuego.ContextWithBody[models.DepartmentsUpdate]) (models.Departments, error) {
	id := c.PathParam("id")

	body, err := c.Body()
	if err != nil {
		return models.Departments{}, err
	}

	new, err := rs.DepartmentsService.UpdateDepartments(c.Context(), id, body)
	if err != nil {
		return models.Departments{}, err
	}

	return new, nil
}

func (rs DepartmentsResources) deleteDepartments(c fuego.ContextNoBody) (any, error) {
	return rs.DepartmentsService.DeleteDepartments(c.Context(), c.PathParam("id"))
}
