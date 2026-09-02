package services

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/adamkali/mindscape/db/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// MockNoteService provides an in-memory implementation of INoteService for testing
type MockNoteService struct {
	ctx           context.Context
	pool          *pgxpool.Pool
	notes         map[uuid.UUID]*repository.Note
	notesByUser   map[uuid.UUID][]*repository.Note
	notesByFolder map[uuid.UUID][]*repository.Note
	mutex         sync.RWMutex

	// Test behavior controls
	ShouldFailGetAll      bool
	ShouldFailGet         bool
	ShouldFailGetByFolder bool
	ShouldFailGetByUser   bool
	ShouldFailCreate      bool
	ShouldFailUpdate      bool
	ShouldFailMove        bool
	ShouldFailDelete      bool

	GetAllErrorMessage      string
	GetErrorMessage         string
	GetByFolderErrorMessage string
	GetByUserErrorMessage   string
	CreateErrorMessage      string
	UpdateErrorMessage      string
	MoveErrorMessage        string
	DeleteErrorMessage      string

	// Test data tracking
	GetAllCallCount      int
	GetCallCount         int
	GetByFolderCallCount int
	GetByUserCallCount   int
	CreateCallCount      int
	UpdateCallCount      int
	MoveCallCount        int
	DeleteCallCount      int

	LastGetID         uuid.UUID
	LastGetByFolderID uuid.UUID
	LastGetByUserID   uuid.UUID
	LastCreateParams  *repository.CreateNoteParams
	LastUpdateID      uuid.UUID
	LastUpdateParams  *repository.UpdateNoteParams
	LastMoveParams    *repository.MoveNoteParams
	LastDeleteID      uuid.UUID
}

// NewMockNoteService creates a bare MockNoteService (no seeded data). Retained for
// callers that only need the default empty behaviour (e.g. folder handler tests).
func NewMockNoteService() *MockNoteService {
	return &MockNoteService{}
}

// CreateMockNoteService creates a new MockNoteService with seeded test data
func CreateMockNoteService(ctx context.Context, pool *pgxpool.Pool) *MockNoteService {
	service := &MockNoteService{
		ctx:           ctx,
		pool:          pool,
		notes:         make(map[uuid.UUID]*repository.Note),
		notesByUser:   make(map[uuid.UUID][]*repository.Note),
		notesByFolder: make(map[uuid.UUID][]*repository.Note),

		GetAllErrorMessage:      "Mock GetAll failure",
		GetErrorMessage:         "Mock Get failure",
		GetByFolderErrorMessage: "Mock GetByFolder failure",
		GetByUserErrorMessage:   "Mock GetByUser failure",
		CreateErrorMessage:      "Mock Create failure",
		UpdateErrorMessage:      "Mock Update failure",
		MoveErrorMessage:        "Mock Move failure",
		DeleteErrorMessage:      "Mock Delete failure",
	}

	service.seedTestData()
	return service
}

// seedTestData populates the mock with test notes
func (m *MockNoteService) seedTestData() {
	now := time.Now()
	testUserID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	adminUserID := uuid.MustParse("22222222-2222-2222-2222-222222222222")
	testFolderID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	adminFolderID := uuid.MustParse("cccccccc-cccc-cccc-cccc-cccccccccccc")

	recentNote := &repository.Note{
		ID:              uuid.MustParse("dddddddd-dddd-dddd-dddd-dddddddddddd"),
		UserID:          testUserID,
		FolderID:        testFolderID,
		Name:            "Recent Note",
		Description:     noteStringPtr("Most recent note"),
		Content:         "# Recent Note\n\nSome markdown content.",
		CreatedDatetime: &now,
		UpdatedDatetime: &now,
	}

	olderTime := now.Add(-1 * time.Hour)
	olderNote := &repository.Note{
		ID:              uuid.MustParse("eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"),
		UserID:          testUserID,
		FolderID:        testFolderID,
		Name:            "Older Note",
		Description:     noteStringPtr("Older note for testing"),
		Content:         "- list item\n- another",
		CreatedDatetime: &olderTime,
		UpdatedDatetime: &olderTime,
	}

	adminNote := &repository.Note{
		ID:              uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
		UserID:          adminUserID,
		FolderID:        adminFolderID,
		Name:            "Admin Note",
		Description:     noteStringPtr("Admin's note"),
		Content:         "admin only",
		CreatedDatetime: &now,
		UpdatedDatetime: &now,
	}

	m.notes[recentNote.ID] = recentNote
	m.notes[olderNote.ID] = olderNote
	m.notes[adminNote.ID] = adminNote

	m.notesByUser[testUserID] = []*repository.Note{recentNote, olderNote}
	m.notesByUser[adminUserID] = []*repository.Note{adminNote}

	m.notesByFolder[testFolderID] = []*repository.Note{recentNote, olderNote}
	m.notesByFolder[adminFolderID] = []*repository.Note{adminNote}
}

// Helper function for string pointers
func noteStringPtr(s string) *string {
	return &s
}

// Reset clears call tracking and resets test data
func (m *MockNoteService) Reset() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	m.GetAllCallCount = 0
	m.GetCallCount = 0
	m.GetByFolderCallCount = 0
	m.GetByUserCallCount = 0
	m.CreateCallCount = 0
	m.UpdateCallCount = 0
	m.MoveCallCount = 0
	m.DeleteCallCount = 0

	m.LastGetID = uuid.Nil
	m.LastGetByFolderID = uuid.Nil
	m.LastGetByUserID = uuid.Nil
	m.LastCreateParams = nil
	m.LastUpdateID = uuid.Nil
	m.LastUpdateParams = nil
	m.LastMoveParams = nil
	m.LastDeleteID = uuid.Nil

	m.ShouldFailGetAll = false
	m.ShouldFailGet = false
	m.ShouldFailGetByFolder = false
	m.ShouldFailGetByUser = false
	m.ShouldFailCreate = false
	m.ShouldFailUpdate = false
	m.ShouldFailMove = false
	m.ShouldFailDelete = false

	m.notes = make(map[uuid.UUID]*repository.Note)
	m.notesByUser = make(map[uuid.UUID][]*repository.Note)
	m.notesByFolder = make(map[uuid.UUID][]*repository.Note)
	m.seedTestData()
}

// INoteService implementation

func (m *MockNoteService) GetAll() ([]repository.Note, error) {
	m.mutex.Lock()
	m.GetAllCallCount++
	m.mutex.Unlock()

	if m.ShouldFailGetAll {
		return nil, errors.New(m.GetAllErrorMessage)
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	notes := make([]repository.Note, 0)
	for _, note := range m.notes {
		notes = append(notes, *note)
	}
	return notes, nil
}

func (m *MockNoteService) Get(id uuid.UUID) (*repository.Note, error) {
	m.mutex.Lock()
	m.GetCallCount++
	m.LastGetID = id
	m.mutex.Unlock()

	if m.ShouldFailGet {
		return nil, errors.New(m.GetErrorMessage)
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	note, exists := m.notes[id]
	if !exists {
		return nil, errors.New("note not found")
	}
	return note, nil
}

func (m *MockNoteService) GetMostRecent(user_id uuid.UUID) (*repository.Note, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	var mostRecent *repository.Note
	if notes, exists := m.notesByUser[user_id]; exists {
		for _, note := range notes {
			if mostRecent == nil ||
				(note.UpdatedDatetime != nil && mostRecent.UpdatedDatetime != nil &&
					note.UpdatedDatetime.After(*mostRecent.UpdatedDatetime)) {
				mostRecent = note
			}
		}
	}
	if mostRecent == nil {
		return nil, errors.New("no notes found for user")
	}
	return mostRecent, nil
}

func (m *MockNoteService) GetMostRecents(user_id uuid.UUID) ([]repository.Note, error) {
	return m.GetByUser(user_id)
}

func (m *MockNoteService) GetByUser(user_id uuid.UUID) ([]repository.Note, error) {
	m.mutex.Lock()
	m.GetByUserCallCount++
	m.LastGetByUserID = user_id
	m.mutex.Unlock()

	if m.ShouldFailGetByUser {
		return nil, errors.New(m.GetByUserErrorMessage)
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	userNotes := make([]repository.Note, 0)
	if notes, exists := m.notesByUser[user_id]; exists {
		for _, note := range notes {
			userNotes = append(userNotes, *note)
		}
	}
	return userNotes, nil
}

func (m *MockNoteService) GetByFolder(folder_id uuid.UUID) ([]repository.Note, error) {
	m.mutex.Lock()
	m.GetByFolderCallCount++
	m.LastGetByFolderID = folder_id
	m.mutex.Unlock()

	if m.ShouldFailGetByFolder {
		return nil, fmt.Errorf("%s", m.GetByFolderErrorMessage)
	}

	m.mutex.RLock()
	defer m.mutex.RUnlock()

	folderNotes := make([]repository.Note, 0)
	if notes, exists := m.notesByFolder[folder_id]; exists {
		for _, note := range notes {
			folderNotes = append(folderNotes, *note)
		}
	}
	return folderNotes, nil
}

func (m *MockNoteService) GetByDateRange(params *repository.FindNotesByUserIDDateTimeRangeParams) ([]repository.Note, error) {
	if params == nil {
		return []repository.Note{}, nil
	}
	return m.GetByUser(params.UserID)
}

func (m *MockNoteService) Create(params *repository.CreateNoteParams) (*repository.Note, error) {
	m.mutex.Lock()
	m.CreateCallCount++
	m.LastCreateParams = params
	m.mutex.Unlock()

	if m.ShouldFailCreate {
		return nil, errors.New(m.CreateErrorMessage)
	}

	if params == nil {
		return nil, errors.New("create parameters cannot be nil")
	}
	if params.Name == "" {
		return nil, errors.New("note name cannot be empty")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.notes == nil {
		m.notes = make(map[uuid.UUID]*repository.Note)
		m.notesByUser = make(map[uuid.UUID][]*repository.Note)
		m.notesByFolder = make(map[uuid.UUID][]*repository.Note)
	}

	now := time.Now()
	newNote := &repository.Note{
		ID:              uuid.New(),
		UserID:          params.UserID,
		FolderID:        params.FolderID,
		Name:            params.Name,
		Description:     params.Description,
		Content:         params.Content,
		CreatedDatetime: &now,
		UpdatedDatetime: &now,
	}

	m.notes[newNote.ID] = newNote
	m.notesByUser[params.UserID] = append(m.notesByUser[params.UserID], newNote)
	m.notesByFolder[params.FolderID] = append(m.notesByFolder[params.FolderID], newNote)

	return newNote, nil
}

func (m *MockNoteService) Update(id uuid.UUID, params *repository.UpdateNoteParams) (*repository.Note, error) {
	m.mutex.Lock()
	m.UpdateCallCount++
	m.LastUpdateID = id
	m.LastUpdateParams = params
	m.mutex.Unlock()

	if m.ShouldFailUpdate {
		return nil, errors.New(m.UpdateErrorMessage)
	}
	if params == nil {
		return nil, errors.New("update parameters cannot be nil")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	note, exists := m.notes[id]
	if !exists {
		return nil, errors.New("note not found")
	}
	if params.Name == "" {
		return nil, errors.New("note name cannot be empty")
	}

	now := time.Now()
	note.Name = params.Name
	note.Description = params.Description
	note.Content = params.Content
	note.UpdatedDatetime = &now

	return note, nil
}

func (m *MockNoteService) Delete(id uuid.UUID) error {
	m.mutex.Lock()
	m.DeleteCallCount++
	m.LastDeleteID = id
	m.mutex.Unlock()

	if m.ShouldFailDelete {
		return errors.New(m.DeleteErrorMessage)
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	note, exists := m.notes[id]
	if !exists {
		return errors.New("note not found")
	}

	delete(m.notes, id)

	if userNotes, exists := m.notesByUser[note.UserID]; exists {
		for i, n := range userNotes {
			if n.ID == id {
				m.notesByUser[note.UserID] = append(userNotes[:i], userNotes[i+1:]...)
				break
			}
		}
	}
	if folderNotes, exists := m.notesByFolder[note.FolderID]; exists {
		for i, n := range folderNotes {
			if n.ID == id {
				m.notesByFolder[note.FolderID] = append(folderNotes[:i], folderNotes[i+1:]...)
				break
			}
		}
	}

	return nil
}

func (m *MockNoteService) Move(params *repository.MoveNoteParams) (*repository.Note, error) {
	m.mutex.Lock()
	m.MoveCallCount++
	m.LastMoveParams = params
	m.mutex.Unlock()

	if m.ShouldFailMove {
		return nil, errors.New(m.MoveErrorMessage)
	}
	if params == nil {
		return nil, errors.New("move parameters cannot be nil")
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	note, exists := m.notes[params.ID]
	if !exists {
		return nil, errors.New("note not found")
	}

	if note.FolderID != params.FolderID {
		if folderNotes, exists := m.notesByFolder[note.FolderID]; exists {
			for i, n := range folderNotes {
				if n.ID == params.ID {
					m.notesByFolder[note.FolderID] = append(folderNotes[:i], folderNotes[i+1:]...)
					break
				}
			}
		}
		m.notesByFolder[params.FolderID] = append(m.notesByFolder[params.FolderID], note)

		note.FolderID = params.FolderID
		now := time.Now()
		note.UpdatedDatetime = &now
	}

	return note, nil
}
