package responses

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type CoolifyWidgetApplication struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Fqdn        string `json:"fqdn"`
	Status      string `json:"status"`
	Redirect    string `json:"redirect"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
} // @name CoolifyWidgetApplication

func newCoolifyWidgetService(
	service clients.CoolifyService,
) *CoolifyWidgetService {
	return &CoolifyWidgetService{
		ID:          service.ID,
		UUID:        service.UUID,
		Name:        service.Name,
		Description: service.Description,
		ServiceType: service.ServiceType,
		CreatedAt:   service.CreatedAt,
		UpdatedAt:   service.UpdatedAt,
		DeletedAt:   service.DeletedAt,
	}
}

type CoolifyWidgetService struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid"`
	Name        string `json:"name"`
	Description string `json:"description"`
	ServiceType string `json:"service_type"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at"`
} // @name CoolifyWidgetService

func newCoolifyWidgetApplication(
	app clients.CoolifyApplication,
) *CoolifyWidgetApplication {
	return &CoolifyWidgetApplication{
		ID:          app.ID,
		Description: app.Description,
		UUID:        app.UUID,
		Name:        app.Name,
		Fqdn:        app.Fqdn,
		Status:      app.Status,
		Redirect:    app.Redirect,
		CreatedAt:   app.CreatedAt,
		UpdatedAt:   app.UpdatedAt,
		DeletedAt:   app.DeletedAt,
	}
}

type CoolifyWidgetData struct {
	Applications []CoolifyWidgetApplication `json:"applications"`
	Services     []CoolifyWidgetService     `json:"services"`
} // @name CoolifyWidgetData


type CoolifyWidgetResponse struct {
	Data    *CoolifyWidgetData `json:"data"`
	Message string             `json:"message"`
	Success bool               `json:"success"`
} // @name CoolifyWidgetResponse

func NewCoolifyWidgetResponse() *CoolifyWidgetResponse {
	return &CoolifyWidgetResponse{
		Data:    &CoolifyWidgetData{},
		Success: true,
		Message: "Ok",
	}
}

func (w *CoolifyWidgetResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *CoolifyWidgetResponse) Successful(
	ctx echo.Context,
	applications []clients.CoolifyApplication,
	services []clients.CoolifyService,
) error {
	apps := make([]CoolifyWidgetApplication, len(applications))
	sers := make([]CoolifyWidgetService, len(services))
	for i, val := range applications {
		apps[i] = *newCoolifyWidgetApplication(val)
	}
	for i, val := range services {
		sers[i] = *newCoolifyWidgetService(val)
	}
	w.Success = true
	w.Data = &CoolifyWidgetData{
		Applications: apps,
		Services:     sers,
	} 
	return ctx.JSON(200, w)
}
