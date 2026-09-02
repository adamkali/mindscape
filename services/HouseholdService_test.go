package services

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

// <method var=services.MockHouseholdService.CreateHousehold>

// Run_HouseholdCreate_AddsOwnerAsAdmin verifies a new household seeds the owner
// as an admin member (owner-only setup).
func Run_HouseholdCreate_AddsOwnerAsAdmin(t *testing.T, service *MockHouseholdService) {
	household, err := service.CreateHousehold("Family", MockHouseholdOwnerID)

	assert.NoError(t, err)
	assert.NotNil(t, household)
	assert.Equal(t, "Family", household.Name)
	assert.Equal(t, MockHouseholdOwnerID, household.OwnerID)

	member, err := service.GetMember(household.ID, MockHouseholdOwnerID)
	assert.NoError(t, err)
	assert.Equal(t, "admin", member.Role)
}

func Run_HouseholdCreate_EmptyName(t *testing.T, service *MockHouseholdService) {
	household, err := service.CreateHousehold("", MockHouseholdOwnerID)
	assert.Error(t, err)
	assert.Nil(t, household)
}

var HouseholdCreateTestMap = map[string]func(*testing.T, *MockHouseholdService){
	"AddsOwnerAsAdmin": Run_HouseholdCreate_AddsOwnerAsAdmin,
	"EmptyName":        Run_HouseholdCreate_EmptyName,
}

func Test_MockHouseholdService_Create(t *testing.T) {
	for name, testFunc := range HouseholdCreateTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockHouseholdService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockHouseholdService.AcceptInvite>

// Run_Accept_Valid accepts a pending, unexpired invite and adds the member.
func Run_Accept_Valid(t *testing.T, service *MockHouseholdService) {
	household, err := service.AcceptInvite("valid-code", MockHouseholdMemberID)

	assert.NoError(t, err)
	assert.NotNil(t, household)
	assert.Equal(t, MockHouseholdID, household.ID)

	member, err := service.GetMember(MockHouseholdID, MockHouseholdMemberID)
	assert.NoError(t, err)
	assert.Equal(t, "member", member.Role)
}

// Run_Accept_Expired rejects an expired invite (AC: invite expired).
func Run_Accept_Expired(t *testing.T, service *MockHouseholdService) {
	household, err := service.AcceptInvite("expired-code", MockHouseholdMemberID)

	assert.Error(t, err)
	assert.Nil(t, household)
	assert.Contains(t, err.Error(), "expired")
}

// Run_Accept_DoubleAccept rejects an already-accepted invite (AC: accepted twice).
func Run_Accept_DoubleAccept(t *testing.T, service *MockHouseholdService) {
	household, err := service.AcceptInvite("accepted-code", MockHouseholdMemberID)

	assert.Error(t, err)
	assert.Nil(t, household)
	assert.Contains(t, err.Error(), "already accepted")
}

// Run_Accept_WrongUser rejects an invite that belongs to another user.
func Run_Accept_WrongUser(t *testing.T, service *MockHouseholdService) {
	household, err := service.AcceptInvite("valid-code", MockHouseholdOutsiderID)

	assert.Error(t, err)
	assert.Nil(t, household)
	assert.Contains(t, err.Error(), "does not belong")
}

var HouseholdAcceptTestMap = map[string]func(*testing.T, *MockHouseholdService){
	"Valid":        Run_Accept_Valid,
	"Expired":      Run_Accept_Expired,
	"DoubleAccept": Run_Accept_DoubleAccept,
	"WrongUser":    Run_Accept_WrongUser,
}

func Test_MockHouseholdService_AcceptInvite(t *testing.T) {
	for name, testFunc := range HouseholdAcceptTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockHouseholdService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockHouseholdService.RemoveMember>

func Run_RemoveMember_Valid(t *testing.T, service *MockHouseholdService) {
	// Seed the member first by accepting their invite.
	_, err := service.AcceptInvite("valid-code", MockHouseholdMemberID)
	assert.NoError(t, err)

	err = service.RemoveMember(MockHouseholdID, MockHouseholdMemberID)
	assert.NoError(t, err)

	_, err = service.GetMember(MockHouseholdID, MockHouseholdMemberID)
	assert.Error(t, err) // member no longer present
}

var HouseholdRemoveTestMap = map[string]func(*testing.T, *MockHouseholdService){
	"Valid": Run_RemoveMember_Valid,
}

func Test_MockHouseholdService_RemoveMember(t *testing.T) {
	for name, testFunc := range HouseholdRemoveTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockHouseholdService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>

// <method var=services.MockHouseholdService.RejectInvite>

func Run_Reject_Valid(t *testing.T, service *MockHouseholdService) {
	err := service.RejectInvite("valid-code", MockHouseholdMemberID)
	assert.NoError(t, err)

	// Rejected invite is gone — accepting it now fails.
	_, err = service.AcceptInvite("valid-code", MockHouseholdMemberID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

var HouseholdRejectTestMap = map[string]func(*testing.T, *MockHouseholdService){
	"Valid": Run_Reject_Valid,
}

func Test_MockHouseholdService_RejectInvite(t *testing.T) {
	for name, testFunc := range HouseholdRejectTestMap {
		t.Run(name, func(t *testing.T) {
			service := CreateMockHouseholdService(context.Background(), nil)
			service.Reset()
			testFunc(t, service)
		})
	}
}

// </method>
