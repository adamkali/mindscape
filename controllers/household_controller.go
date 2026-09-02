package controllers

import (
	"github.com/adamkali/mindscape/cmd/configuration"
	"github.com/adamkali/mindscape/models/handlers/household_handlers"
	"github.com/adamkali/mindscape/services"
	"github.com/labstack/echo/v4"
)

const defaultInviteTTLHours = 168 // 7 days

type HouseholdController struct {
	Name             string
	Config           *configuration.Configuration
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
	ValidatorService *services.ValidatorService
}

func (c HouseholdController) ControllerName() string {
	return c.Name
}

func BuildHouseholdController(p *services.Registrar) HouseholdController {
	return HouseholdController{
		Name:             "/households",
		Config:           p.Config,
		AuthService:      p.AuthService,
		HouseholdService: p.HouseholdService,
		ValidatorService: p.ValidatorService,
	}
}

func (c HouseholdController) inviteTTLHours() int {
	if c.Config != nil && c.Config.Households.InviteTTLHours > 0 {
		return c.Config.Households.InviteTTLHours
	}
	return defaultInviteTTLHours
}

// @Summary     Create a new Household
// @Description Create a Household; the caller becomes its owner and admin.
// @ID          CreateHousehold
// @Tags        Households
// @Accept      json
// @Produce     json
// @Param       CreateHouseholdRequest body requests.CreateHouseholdRequest true "Create Household Request"
// @Security    BearerAuth
// @Success     200 {object} HouseholdResponse
// @Failure     400 {object} HouseholdResponse
// @Failure     401 {object} HouseholdResponse
// @Failure     500 {object} HouseholdResponse
// @Router      /households [post]
func (c HouseholdController) CreateHousehold(e echo.Context) error {
	return household_handlers.NewCreateHouseholdHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

// @Summary     Get My Households
// @Description Get all households the caller is a member of.
// @ID          GetMyHouseholds
// @Tags        Households
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} HouseholdsResponse
// @Failure     401 {object} HouseholdsResponse
// @Failure     500 {object} HouseholdsResponse
// @Router      /households [get]
func (c HouseholdController) GetMyHouseholds(e echo.Context) error {
	return household_handlers.NewGetMyHouseholdsHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

// @Summary     Get My Pending Invites
// @Description Get the caller's pending, unexpired household invites.
// @ID          GetMyInvites
// @Tags        Households
// @Produce     json
// @Security    BearerAuth
// @Success     200 {object} InvitesResponse
// @Failure     401 {object} InvitesResponse
// @Failure     500 {object} InvitesResponse
// @Router      /households/invites [get]
func (c HouseholdController) GetMyInvites(e echo.Context) error {
	return household_handlers.NewGetMyInvitesHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

// @Summary     Get Household Members
// @Description Get a household's members. Returns 404 to non-members.
// @ID          GetHouseholdMembers
// @Tags        Households
// @Produce     json
// @Security    BearerAuth
// @Param       household_id path string true "Household ID"
// @Success     200 {object} MembersResponse
// @Failure     401 {object} MembersResponse
// @Failure     404 {object} MembersResponse
// @Router      /households/{household_id}/members [get]
func (c HouseholdController) GetHouseholdMembers(e echo.Context) error {
	return household_handlers.NewGetMembersHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

// @Summary     Create an Invite
// @Description Invite a user to a household (admin only).
// @ID          CreateInvite
// @Tags        Households
// @Accept      json
// @Produce     json
// @Security    BearerAuth
// @Param       household_id path string true "Household ID"
// @Param       CreateInviteRequest body requests.CreateInviteRequest true "Create Invite Request"
// @Success     200 {object} StringResponse
// @Failure     400 {object} StringResponse
// @Failure     401 {object} StringResponse
// @Failure     403 {object} StringResponse
// @Failure     404 {object} StringResponse
// @Router      /households/{household_id}/invites [post]
func (c HouseholdController) CreateInvite(e echo.Context) error {
	return household_handlers.NewCreateInviteHandler(e, c.HouseholdService, c.AuthService, c.inviteTTLHours()).Handle().JSON()
}

// @Summary     Accept an Invite
// @Description Accept a pending invite by its code.
// @ID          AcceptInvite
// @Tags        Households
// @Produce     json
// @Security    BearerAuth
// @Param       code path string true "Invite Code"
// @Success     200 {object} HouseholdResponse
// @Failure     400 {object} HouseholdResponse
// @Failure     401 {object} HouseholdResponse
// @Router      /households/invites/{code}/accept [post]
func (c HouseholdController) AcceptInvite(e echo.Context) error {
	return household_handlers.NewAcceptInviteHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

// @Summary     Reject an Invite
// @Description Reject (and delete) a pending invite by its code.
// @ID          RejectInvite
// @Tags        Households
// @Produce     json
// @Security    BearerAuth
// @Param       code path string true "Invite Code"
// @Success     200 {object} StringResponse
// @Failure     400 {object} StringResponse
// @Failure     401 {object} StringResponse
// @Router      /households/invites/{code}/reject [post]
func (c HouseholdController) RejectInvite(e echo.Context) error {
	return household_handlers.NewRejectInviteHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

// @Summary     Remove a Member
// @Description Remove a member from a household (admin only; not the owner).
// @ID          RemoveMember
// @Tags        Households
// @Produce     json
// @Security    BearerAuth
// @Param       household_id path string true "Household ID"
// @Param       user_id path string true "User ID"
// @Success     200 {object} StringResponse
// @Failure     401 {object} StringResponse
// @Failure     403 {object} StringResponse
// @Failure     404 {object} StringResponse
// @Router      /households/{household_id}/members/{user_id} [delete]
func (c HouseholdController) RemoveMember(e echo.Context) error {
	return household_handlers.NewRemoveMemberHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

// @Summary     Leave a Household
// @Description Leave a household. The owner cannot leave — they must delete it.
// @ID          LeaveHousehold
// @Tags        Households
// @Produce     json
// @Security    BearerAuth
// @Param       household_id path string true "Household ID"
// @Success     200 {object} StringResponse
// @Failure     401 {object} StringResponse
// @Failure     403 {object} StringResponse
// @Failure     404 {object} StringResponse
// @Router      /households/{household_id}/leave [post]
func (c HouseholdController) LeaveHousehold(e echo.Context) error {
	return household_handlers.NewLeaveHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

// @Summary     Delete a Household
// @Description Delete a household (owner only). Cascades members and invites.
// @ID          DeleteHousehold
// @Tags        Households
// @Produce     json
// @Security    BearerAuth
// @Param       household_id path string true "Household ID"
// @Success     200 {object} StringResponse
// @Failure     401 {object} StringResponse
// @Failure     403 {object} StringResponse
// @Failure     404 {object} StringResponse
// @Router      /households/{household_id} [delete]
func (c HouseholdController) DeleteHousehold(e echo.Context) error {
	return household_handlers.NewDeleteHandler(e, c.HouseholdService, c.AuthService).Handle().JSON()
}

func (c HouseholdController) Attatch(e *echo.Echo, middlewares ...echo.MiddlewareFunc) {
	api := e.Group("/api" + c.Name)
	api.POST("", c.CreateHousehold, middlewares...)
	api.GET("", c.GetMyHouseholds, middlewares...)
	api.GET("/invites", c.GetMyInvites, middlewares...)
	api.POST("/invites/:code/accept", c.AcceptInvite, middlewares...)
	api.POST("/invites/:code/reject", c.RejectInvite, middlewares...)
	api.GET("/:household_id/members", c.GetHouseholdMembers, middlewares...)
	api.POST("/:household_id/invites", c.CreateInvite, middlewares...)
	api.DELETE("/:household_id/members/:user_id", c.RemoveMember, middlewares...)
	api.POST("/:household_id/leave", c.LeaveHousehold, middlewares...)
	api.DELETE("/:household_id", c.DeleteHousehold, middlewares...)
}
