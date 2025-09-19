package services

import (
	"fmt"
	
	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
)

// MockNoteService provides a simple mock implementation of INoteService for testing
type MockNoteService struct {
	// Test behavior controls
	ShouldFailGetByFolder bool
	GetByFolderErrorMessage string
}

func NewMockNoteService() *MockNoteService {
	return &MockNoteService{}
}

func (m *MockNoteService) GetAll() ([]repository.Note, error) {
	return []repository.Note{}, nil
}

func (m *MockNoteService) Get(id uuid.UUID) (*repository.Note, error) {
	return nil, nil
}

func (m *MockNoteService) GetMostRecent(user_id uuid.UUID) (*repository.Note, error) {
	return nil, nil
}

func (m *MockNoteService) GetMostRecents(user_id uuid.UUID) ([]repository.Note, error) {
	return []repository.Note{}, nil
}

func (m *MockNoteService) GetByUser(user_id uuid.UUID) ([]repository.Note, error) {
	return []repository.Note{}, nil
}

func (m *MockNoteService) GetByFolder(folder_id uuid.UUID) ([]repository.Note, error) {
	if m.ShouldFailGetByFolder {
		return nil, fmt.Errorf("%s", m.GetByFolderErrorMessage)
	}
	return []repository.Note{}, nil
}

func (m *MockNoteService) GetByDateRange(params *repository.FindNotesByUserIDDateTimeRangeParams) ([]repository.Note, error) {
	return []repository.Note{}, nil
}

func (m *MockNoteService) Create(params *repository.CreateNoteParams) (*repository.Note, error) {
	return nil, nil
}

func (m *MockNoteService) Update(id uuid.UUID, params *repository.UpdateNoteParams) (*repository.Note, error) {
	return nil, nil
}

func (m *MockNoteService) Delete(id uuid.UUID) error {
	return nil
}

func (m *MockNoteService) Move(params *repository.MoveNoteParams) (*repository.Note, error) {
	return nil, nil
}