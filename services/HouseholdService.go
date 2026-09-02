package services

import (
	"context"
	"errors"
	"time"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type HouseholdService struct {
	ctx  context.Context
	pool *pgxpool.Pool
}

func CreateHouseholdService(
	ctx context.Context,
	pool *pgxpool.Pool,
) IHouseholdService {
	return &HouseholdService{
		ctx:  ctx,
		pool: pool,
	}
}

// CreateHousehold creates a household and adds the owner as an admin member in
// a single transaction.
func (s HouseholdService) CreateHousehold(name string, ownerID uuid.UUID) (*repository.Household, error) {
	if name == "" {
		return nil, errors.New("household name cannot be empty")
	}
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)

	household, err := repo.CreateHousehold(s.ctx, repository.CreateHouseholdParams{
		Name:    name,
		OwnerID: ownerID,
	})
	if err != nil {
		return nil, err
	}
	if _, err = repo.AddMember(s.ctx, repository.AddMemberParams{
		HouseholdID: household.ID,
		UserID:      ownerID,
		Role:        "admin",
	}); err != nil {
		return nil, err
	}
	tx.Commit(s.ctx)
	return &household, nil
}

func (s HouseholdService) GetByUser(userID uuid.UUID) ([]repository.Household, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []repository.Household{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	households, err := repo.GetHouseholdsByUserID(s.ctx, userID)
	if err != nil {
		return []repository.Household{}, err
	}
	tx.Commit(s.ctx)
	return households, nil
}

func (s HouseholdService) Get(id uuid.UUID) (*repository.Household, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	household, err := repo.GetHouseholdByID(s.ctx, id)
	if err != nil {
		return nil, err
	}
	tx.Commit(s.ctx)
	return &household, nil
}

func (s HouseholdService) Delete(id uuid.UUID) error {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	if err := repo.DeleteHousehold(s.ctx, id); err != nil {
		return err
	}
	tx.Commit(s.ctx)
	return nil
}

func (s HouseholdService) GetMembers(householdID uuid.UUID) ([]repository.GetMembersByHouseholdIDRow, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []repository.GetMembersByHouseholdIDRow{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	members, err := repo.GetMembersByHouseholdID(s.ctx, householdID)
	if err != nil {
		return []repository.GetMembersByHouseholdIDRow{}, err
	}
	tx.Commit(s.ctx)
	return members, nil
}

func (s HouseholdService) GetMember(householdID uuid.UUID, userID uuid.UUID) (*repository.HouseholdMember, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	member, err := repo.GetMember(s.ctx, repository.GetMemberParams{
		HouseholdID: householdID,
		UserID:      userID,
	})
	if err != nil {
		return nil, err
	}
	tx.Commit(s.ctx)
	return &member, nil
}

func (s HouseholdService) RemoveMember(householdID uuid.UUID, userID uuid.UUID) error {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	if err := repo.RemoveMember(s.ctx, repository.RemoveMemberParams{
		HouseholdID: householdID,
		UserID:      userID,
	}); err != nil {
		return err
	}
	tx.Commit(s.ctx)
	return nil
}

func (s HouseholdService) CreateInvite(householdID uuid.UUID, invitedUserID uuid.UUID, createdBy uuid.UUID, ttl time.Duration) (*repository.HouseholdInvite, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	expires := time.Now().Add(ttl)
	code := uuid.NewString()
	invite, err := repo.CreateInvite(s.ctx, repository.CreateInviteParams{
		HouseholdID:   householdID,
		InvitedUserID: invitedUserID,
		Code:          code,
		ExpiresAt:     &expires,
		CreatedBy:     createdBy,
	})
	if err != nil {
		return nil, err
	}
	tx.Commit(s.ctx)
	return &invite, nil
}

func (s HouseholdService) GetInvites(userID uuid.UUID) ([]repository.GetInvitesByUserIDRow, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return []repository.GetInvitesByUserIDRow{}, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	invites, err := repo.GetInvitesByUserID(s.ctx, userID)
	if err != nil {
		return []repository.GetInvitesByUserIDRow{}, err
	}
	tx.Commit(s.ctx)
	return invites, nil
}

// AcceptInvite validates the invite (ownership, single-use, expiry), marks it
// accepted, and adds the user as a member — all in one transaction.
func (s HouseholdService) AcceptInvite(code string, userID uuid.UUID) (*repository.Household, error) {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)

	invite, err := repo.GetInviteByCode(s.ctx, code)
	if err != nil {
		return nil, errors.New("invite not found")
	}
	if invite.InvitedUserID != userID {
		return nil, errors.New("invite does not belong to this user")
	}
	if invite.AcceptedAt.Valid {
		return nil, errors.New("invite already accepted")
	}
	if invite.ExpiresAt != nil && invite.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("invite expired")
	}

	if err = repo.MarkInviteAccepted(s.ctx, invite.ID); err != nil {
		return nil, err
	}
	if _, err = repo.AddMember(s.ctx, repository.AddMemberParams{
		HouseholdID: invite.HouseholdID,
		UserID:      userID,
		Role:        "member",
	}); err != nil {
		return nil, err
	}
	household, err := repo.GetHouseholdByID(s.ctx, invite.HouseholdID)
	if err != nil {
		return nil, err
	}
	tx.Commit(s.ctx)
	return &household, nil
}

// RejectInvite deletes the invite (single-use; rejection is permanent).
func (s HouseholdService) RejectInvite(code string, userID uuid.UUID) error {
	tx, err := s.pool.Begin(s.ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(s.ctx)
	repo := repository.New(tx)
	invite, err := repo.GetInviteByCode(s.ctx, code)
	if err != nil {
		return errors.New("invite not found")
	}
	if invite.InvitedUserID != userID {
		return errors.New("invite does not belong to this user")
	}
	if err = repo.DeleteInvite(s.ctx, invite.ID); err != nil {
		return err
	}
	tx.Commit(s.ctx)
	return nil
}
