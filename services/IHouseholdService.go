package services


import (
	"time"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)

type IHouseholdService interface {
	// Households
	CreateHousehold(name string, ownerID uuid.UUID) (*repository.Household, error)
	GetByUser(userID uuid.UUID) ([]repository.Household, error)
	Get(id uuid.UUID) (*repository.Household, error)
	Delete(id uuid.UUID) error

	// Members
	GetMembers(householdID uuid.UUID) ([]repository.GetMembersByHouseholdIDRow, error)
	GetMember(householdID uuid.UUID, userID uuid.UUID) (*repository.HouseholdMember, error)
	RemoveMember(householdID uuid.UUID, userID uuid.UUID) error

	// Invites
	CreateInvite(householdID uuid.UUID, invitedUserID uuid.UUID, createdBy uuid.UUID, ttl time.Duration) (*repository.HouseholdInvite, error)
	GetInvites(userID uuid.UUID) ([]repository.GetInvitesByUserIDRow, error)
	AcceptInvite(code string, userID uuid.UUID) (*repository.Household, error)
	RejectInvite(code string, userID uuid.UUID) error
}
