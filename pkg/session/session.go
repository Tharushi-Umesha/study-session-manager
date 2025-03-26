package session

import (
	"fmt"
	"time"
)

// Subject represents a study subject
type Subject struct {
	Name        string
	Description string
}

// Session represents a study session
type Session struct {
	ID          int
	Subject     Subject
	StartTime   time.Time
	EndTime     time.Time
	Duration    time.Duration
	Notes       string
	Completed   bool
}

// Manager handles study sessions
type Manager struct {
	sessions []Session
	nextID   int
}

// NewManager creates a new session manager
func NewManager() *Manager {
	return &Manager{
		sessions: make([]Session, 0),
		nextID:   1,
	}
}

// StartSession begins a new study session
func (m *Manager) StartSession(subject Subject) (Session, error) {
	session := Session{
		ID:        m.nextID,
		Subject:   subject,
		StartTime: time.Now(),
		Completed: false,
	}
	
	m.sessions = append(m.sessions, session)
	m.nextID++
	
	return session, nil
}

// EndSession completes a study session
func (m *Manager) EndSession(id int, notes string) error {
	for i, s := range m.sessions {
		if s.ID == id && !s.Completed {
			// Update the session
			m.sessions[i].EndTime = time.Now()
			m.sessions[i].Duration = m.sessions[i].EndTime.Sub(s.StartTime)
			m.sessions[i].Notes = notes
			m.sessions[i].Completed = true
			
			return nil
		}
	}
	
	return fmt.Errorf("session with ID %d not found or already completed", id)
}

// GetActiveSessions returns all active (not completed) sessions
func (m *Manager) GetActiveSessions() []Session {
	var active []Session
	
	for _, s := range m.sessions {
		if !s.Completed {
			active = append(active, s)
		}
	}
	
	return active
}

// GetCompletedSessions returns all completed sessions
func (m *Manager) GetCompletedSessions() []Session {
	var completed []Session
	
	for _, s := range m.sessions {
		if s.Completed {
			completed = append(completed, s)
		}
	}
	
	return completed
}

// GetAllSessions returns all sessions
func (m *Manager) GetAllSessions() []Session {
	// Return a copy to prevent modification
	sessions := make([]Session, len(m.sessions))
	copy(sessions, m.sessions)
	return sessions
}

// GetSessionByID finds a session by its ID
func (m *Manager) GetSessionByID(id int) (Session, error) {
	for _, s := range m.sessions {
		if s.ID == id {
			return s, nil
		}
	}
	
	return Session{}, fmt.Errorf("session with ID %d not found", id)
}

// GetTotalStudyTime calculates total time spent studying
func (m *Manager) GetTotalStudyTime() time.Duration {
	var total time.Duration
	
	for _, s := range m.sessions {
		if s.Completed {
			total += s.Duration
		}
	}
	
	return total
}

// GetSubjectStudyTime calculates time spent on a specific subject
func (m *Manager) GetSubjectStudyTime(subjectName string) time.Duration {
	var total time.Duration
	
	for _, s := range m.sessions {
		if s.Completed && s.Subject.Name == subjectName {
			total += s.Duration
		}
	}
	
	return total
}