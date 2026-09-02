// Package household_handlers implements the IHandler pattern (mirror of
// folder_handlers/) for the Households epic. Handlers are consolidated into one
// file to avoid ~10 near-identical boilerplate files; each endpoint is still a
// distinct IHandler with its own Handle()/JSON(), sharing the common
// SetCode/SetError/Code/Data/Error plumbing via baseHandler.
package household_handlers

import (
	"errors"
	"time"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/adamkali/mindscape/models/handlers"
	"github.com/adamkali/mindscape/models/requests"
	"github.com/adamkali/mindscape/models/responses"
	"github.com/adamkali/mindscape/services"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// baseHandler supplies the shared IHandler plumbing (everything except
// Handle/JSON, which each handler defines).
type baseHandler struct {
	ctx  echo.Context
	code int
	err  error
}

func (h *baseHandler) SetCode(code int) handlers.IHandler   { h.code = code; return h }
func (h *baseHandler) SetError(err error) handlers.IHandler { h.err = err; return h }
func (h *baseHandler) Code() int                            { return h.code }
func (h *baseHandler) Error() error                         { return h.err }
func (h *baseHandler) Data() any                            { return nil }

// Handle/JSON are no-ops so *baseHandler satisfies IHandler; each concrete
// handler overrides them. Lock/requireMember are always called with the
// concrete handler so JSON() dispatches to the right response type.
func (h *baseHandler) Handle() handlers.IHandler { return h }
func (h *baseHandler) JSON() error               { return nil }

// authUser validates the JWT and returns the authenticated user's ID.
func authUser(ctx echo.Context, auth services.IAuthService) (uuid.UUID, error) {
	token := ctx.Get("user").(*jwt.Token)
	if err := auth.CheckToken(token.Raw); err != nil {
		return uuid.Nil, err
	}
	claims := token.Claims.(*services.CustomJwt)
	return claims.UserId, nil
}

// requireMember returns 404 (not 403) when the requester is not a member, so the
// existence of a household is never leaked to outsiders.
func requireMember(h handlers.IHandler, hs services.IHouseholdService, householdID, userID uuid.UUID) (*repository.HouseholdMember, handlers.IHandler) {
	member, err := hs.GetMember(householdID, userID)
	if err != nil {
		return nil, handlers.Lock(h, 404, errors.New("household not found"))
	}
	return member, nil
}

// ── CreateHousehold ───────────────────────────────────────────────────────────

type CreateHouseholdHandler struct {
	baseHandler
	data             *repository.Household
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewCreateHouseholdHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *CreateHouseholdHandler {
	return &CreateHouseholdHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *CreateHouseholdHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	req := new(requests.CreateHouseholdRequest)
	if err = h.ctx.Bind(req); err != nil {
		return handlers.Lock(h, 400, err)
	}
	if req.Name == "" {
		return handlers.Lock(h, 400, errors.New("household name is required"))
	}
	if h.data, err = h.HouseholdService.CreateHousehold(req.Name, userID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *CreateHouseholdHandler) JSON() error {
	if h.err != nil {
		return responses.NewHouseholdResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewHouseholdResponse().Successful(h.ctx, *h.data)
}

// ── GetMyHouseholds ────────────────────────────────────────────────────────────

type GetMyHouseholdsHandler struct {
	baseHandler
	data             []repository.Household
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewGetMyHouseholdsHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *GetMyHouseholdsHandler {
	return &GetMyHouseholdsHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *GetMyHouseholdsHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	if h.data, err = h.HouseholdService.GetByUser(userID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *GetMyHouseholdsHandler) JSON() error {
	if h.err != nil {
		return responses.NewHouseholdsResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewHouseholdsResponse().Successful(h.ctx, h.data)
}

// ── GetHouseholdMembers ────────────────────────────────────────────────────────

type GetMembersHandler struct {
	baseHandler
	data             []repository.GetMembersByHouseholdIDRow
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewGetMembersHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *GetMembersHandler {
	return &GetMembersHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *GetMembersHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	householdID, err := uuid.Parse(h.ctx.Param("household_id"))
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	if _, locked := requireMember(h, h.HouseholdService, householdID, userID); locked != nil {
		return locked
	}
	if h.data, err = h.HouseholdService.GetMembers(householdID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *GetMembersHandler) JSON() error {
	if h.err != nil {
		return responses.NewMembersResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewMembersResponse().Successful(h.ctx, h.data)
}

// ── GetMyInvites (needed for accept flow) ──────────────────────────────────────

type GetMyInvitesHandler struct {
	baseHandler
	data             []repository.GetInvitesByUserIDRow
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewGetMyInvitesHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *GetMyInvitesHandler {
	return &GetMyInvitesHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *GetMyInvitesHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	if h.data, err = h.HouseholdService.GetInvites(userID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *GetMyInvitesHandler) JSON() error {
	if h.err != nil {
		return responses.NewInvitesResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewInvitesResponse().Successful(h.ctx, h.data)
}

// ── CreateInvite (admin only) ──────────────────────────────────────────────────

type CreateInviteHandler struct {
	baseHandler
	data             *string
	ttlHours         int
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewCreateInviteHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService, ttlHours int) *CreateInviteHandler {
	return &CreateInviteHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth, ttlHours: ttlHours}
}

func (h *CreateInviteHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	householdID, err := uuid.Parse(h.ctx.Param("household_id"))
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	member, locked := requireMember(h, h.HouseholdService, householdID, userID)
	if locked != nil {
		return locked
	}
	if member.Role != "admin" {
		return handlers.Lock(h, 403, errors.New("only admins can invite members"))
	}
	req := new(requests.CreateInviteRequest)
	if err = h.ctx.Bind(req); err != nil {
		return handlers.Lock(h, 400, err)
	}
	if req.InvitedUserID == uuid.Nil {
		return handlers.Lock(h, 400, errors.New("invited_user_id is required"))
	}
	ttl := time.Duration(h.ttlHours) * time.Hour
	invite, err := h.HouseholdService.CreateInvite(householdID, req.InvitedUserID, userID, ttl)
	if err != nil {
		return handlers.Lock(h, 500, err)
	}
	h.data = &invite.Code
	return h
}

func (h *CreateInviteHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewStringResponse().Successful(h.ctx, *h.data)
}

// ── AcceptInvite ───────────────────────────────────────────────────────────────

type AcceptInviteHandler struct {
	baseHandler
	data             *repository.Household
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewAcceptInviteHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *AcceptInviteHandler {
	return &AcceptInviteHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *AcceptInviteHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	code := h.ctx.Param("code")
	if code == "" {
		return handlers.Lock(h, 400, errors.New("invite code is required"))
	}
	if h.data, err = h.HouseholdService.AcceptInvite(code, userID); err != nil {
		return handlers.Lock(h, 400, err)
	}
	return h
}

func (h *AcceptInviteHandler) JSON() error {
	if h.err != nil {
		return responses.NewHouseholdResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewHouseholdResponse().Successful(h.ctx, *h.data)
}

// ── RejectInvite ───────────────────────────────────────────────────────────────

type RejectInviteHandler struct {
	baseHandler
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewRejectInviteHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *RejectInviteHandler {
	return &RejectInviteHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *RejectInviteHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	code := h.ctx.Param("code")
	if code == "" {
		return handlers.Lock(h, 400, errors.New("invite code is required"))
	}
	if err = h.HouseholdService.RejectInvite(code, userID); err != nil {
		return handlers.Lock(h, 400, err)
	}
	return h
}

func (h *RejectInviteHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewStringResponse().Successful(h.ctx, "invite rejected")
}

// ── RemoveMember (admin only; cannot remove the owner) ─────────────────────────

type RemoveMemberHandler struct {
	baseHandler
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewRemoveMemberHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *RemoveMemberHandler {
	return &RemoveMemberHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *RemoveMemberHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	householdID, err := uuid.Parse(h.ctx.Param("household_id"))
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	targetID, err := uuid.Parse(h.ctx.Param("user_id"))
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	member, locked := requireMember(h, h.HouseholdService, householdID, userID)
	if locked != nil {
		return locked
	}
	if member.Role != "admin" {
		return handlers.Lock(h, 403, errors.New("only admins can remove members"))
	}
	household, err := h.HouseholdService.Get(householdID)
	if err != nil {
		return handlers.Lock(h, 404, errors.New("household not found"))
	}
	if household.OwnerID == targetID {
		return handlers.Lock(h, 403, errors.New("cannot remove the household owner"))
	}
	if err = h.HouseholdService.RemoveMember(householdID, targetID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *RemoveMemberHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewStringResponse().Successful(h.ctx, "member removed")
}

// ── LeaveHousehold (owner cannot leave — must delete) ──────────────────────────

type LeaveHandler struct {
	baseHandler
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewLeaveHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *LeaveHandler {
	return &LeaveHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *LeaveHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	householdID, err := uuid.Parse(h.ctx.Param("household_id"))
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	if _, locked := requireMember(h, h.HouseholdService, householdID, userID); locked != nil {
		return locked
	}
	household, err := h.HouseholdService.Get(householdID)
	if err != nil {
		return handlers.Lock(h, 404, errors.New("household not found"))
	}
	if household.OwnerID == userID {
		return handlers.Lock(h, 403, errors.New("owner cannot leave; delete the household instead"))
	}
	if err = h.HouseholdService.RemoveMember(householdID, userID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *LeaveHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewStringResponse().Successful(h.ctx, "left household")
}

// ── DeleteHousehold (owner only) ───────────────────────────────────────────────

type DeleteHandler struct {
	baseHandler
	AuthService      services.IAuthService
	HouseholdService services.IHouseholdService
}

func NewDeleteHandler(ctx echo.Context, hs services.IHouseholdService, auth services.IAuthService) *DeleteHandler {
	return &DeleteHandler{baseHandler: baseHandler{ctx: ctx, code: 200}, HouseholdService: hs, AuthService: auth}
}

func (h *DeleteHandler) Handle() handlers.IHandler {
	userID, err := authUser(h.ctx, h.AuthService)
	if err != nil {
		return handlers.Lock(h, 401, err)
	}
	householdID, err := uuid.Parse(h.ctx.Param("household_id"))
	if err != nil {
		return handlers.Lock(h, 400, err)
	}
	if _, locked := requireMember(h, h.HouseholdService, householdID, userID); locked != nil {
		return locked
	}
	household, err := h.HouseholdService.Get(householdID)
	if err != nil {
		return handlers.Lock(h, 404, errors.New("household not found"))
	}
	if household.OwnerID != userID {
		return handlers.Lock(h, 403, errors.New("only the owner can delete the household"))
	}
	if err = h.HouseholdService.Delete(householdID); err != nil {
		return handlers.Lock(h, 500, err)
	}
	return h
}

func (h *DeleteHandler) JSON() error {
	if h.err != nil {
		return responses.NewStringResponse().Fail(h.ctx, h.code, h.err)
	}
	return responses.NewStringResponse().Successful(h.ctx, "household deleted")
}
