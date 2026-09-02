package requests

import "github.com/google/uuid"

type CreateInviteRequest struct {
	InvitedUserID uuid.UUID `json:"invited_user_id"`
} // @name CreateInviteRequest
