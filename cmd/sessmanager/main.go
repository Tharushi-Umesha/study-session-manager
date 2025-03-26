package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Tharushi-Umesha/study-session-manager/pkg/session"
)

func main() {
	fmt.Println("Welcome to Study Session Manager!")
	fmt.Println("=================================")
	
	manager := session.NewManager()
	reader := bufio.NewReader(os.Stdin)
	
	for {
		fmt.Println("\nOptions:")
		fmt.Println("1. Start a new study session")
		fmt.Println("2. End a study session")
		fmt.Println("3. View active sessions")
		fmt.Println("4. View completed sessions")
		fmt.Println("5. View study statistics")
		fmt.Println("6. Exit")
		fmt.Print("\nEnter your choice (1-6): ")
		
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		
		switch choice {
		case "1":
			startStudySession(manager, reader)
		case "2":
			endStudySession(manager, reader)
		case "3":
			viewActiveSessions(manager)
		case "4":
			viewCompletedSessions(manager)
		case "5":
			viewStatistics(manager)
		case "6":
			fmt.Println("Thank you for using Study Session Manager. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}

func startStudySession(manager *session.Manager, reader *bufio.Reader) {
	fmt.Print("\nEnter subject name: ")
	name, _ := reader.ReadString('\n')
	name = strings.TrimSpace(name)
	
	fmt.Print("Enter subject description: ")
	desc, _ := reader.ReadString('\n')
	desc = strings.TrimSpace(desc)
	
	subject := session.Subject{
		Name:        name,
		Description: desc,
	}
	
	sess, err := manager.StartSession(subject)
	if err != nil {
		fmt.Printf("Error starting session: %v\n", err)
		return
	}
	
	fmt.Printf("Started study session #%d for %s at %s\n", 
		sess.ID, 
		sess.Subject.Name, 
		sess.StartTime.Format("15:04:05"))
}

func endStudySession(manager *session.Manager, reader *bufio.Reader) {
	// Show active sessions first
	activeSessions := manager.GetActiveSessions()
	if len(activeSessions) == 0 {
		fmt.Println("No active sessions to end.")
		return
	}
	
	fmt.Println("\nActive sessions:")
	for _, s := range activeSessions {
		fmt.Printf("#%d - %s (started at %s)\n", 
			s.ID, 
			s.Subject.Name, 
			s.StartTime.Format("15:04:05"))
	}
	
	fmt.Print("\nEnter session ID to end: ")
	idStr, _ := reader.ReadString('\n')
	idStr = strings.TrimSpace(idStr)
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("Invalid ID. Please enter a number.")
		return
	}
	
	fmt.Print("Enter session notes: ")
	notes, _ := reader.ReadString('\n')
	notes = strings.TrimSpace(notes)
	
	err = manager.EndSession(id, notes)
	if err != nil {
		fmt.Printf("Error ending session: %v\n", err)
		return
	}
	
	session, _ := manager.GetSessionByID(id)
	duration := formatDuration(session.Duration)
	
	fmt.Printf("Ended study session #%d for %s. Duration: %s\n", 
		id, 
		session.Subject.Name,
		duration)
}

func viewActiveSessions(manager *session.Manager) {
	sessions := manager.GetActiveSessions()
	
	if len(sessions) == 0 {
		fmt.Println("\nNo active study sessions.")
		return
	}
	
	fmt.Println("\nActive study sessions:")
	fmt.Println("---------------------")
	for _, s := range sessions {
		duration := time.Since(s.StartTime)
		fmt.Printf("#%d - %s\n", s.ID, s.Subject.Name)
		fmt.Printf("    Started: %s\n", s.StartTime.Format("15:04:05"))
		fmt.Printf("    Running for: %s\n", formatDuration(duration))
		fmt.Println()
	}
}

func viewCompletedSessions(manager *session.Manager) {
	sessions := manager.GetCompletedSessions()
	
	if len(sessions) == 0 {
		fmt.Println("\nNo completed study sessions.")
		return
	}
	
	fmt.Println("\nCompleted study sessions:")
	fmt.Println("------------------------")
	for _, s := range sessions {
		fmt.Printf("#%d - %s\n", s.ID, s.Subject.Name)
		fmt.Printf("    Duration: %s\n", formatDuration(s.Duration))
		fmt.Printf("    Notes: %s\n", s.Notes)
		fmt.Println()
	}
}

func viewStatistics(manager *session.Manager) {
	totalTime := manager.GetTotalStudyTime()
	
	fmt.Println("\nStudy Statistics:")
	fmt.Println("----------------")
	fmt.Printf("Total study time: %s\n", formatDuration(totalTime))
	
	// Get unique subjects
	sessions := manager.GetAllSessions()
	subjectMap := make(map[string]bool)
	
	for _, s := range sessions {
		subjectMap[s.Subject.Name] = true
	}
	
	fmt.Println("\nTime by subject:")
	for subject := range subjectMap {
		subjectTime := manager.GetSubjectStudyTime(subject)
		if subjectTime > 0 {
			fmt.Printf("- %s: %s\n", subject, formatDuration(subjectTime))
		}
	}
}

func formatDuration(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second
	
	if h > 0 {
		return fmt.Sprintf("%dh %dm %ds", h, m, s)
	}
	if m > 0 {
		return fmt.Sprintf("%dm %ds", m, s)
	}
	return fmt.Sprintf("%ds", s)
}