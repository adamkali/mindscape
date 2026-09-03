package requests

import "github.com/google/uuid"

type ShareFolderWithHouseholdMemeberRequest struct {
	FolderID uuid.UUID `json:"folder_id"`
	HouseholdID uuid.UUID `json:"household_id"`
	MemberID uuid.UUID `json:"member_id"`

} // @Name ShareFolderWithHouseholdMemeberRequest
