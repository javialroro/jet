package main

import (
	"strconv"
	"strings"
	"github.com/charmbracelet/lipgloss"
	"math"
)

var (
	appNameStyle     = lipgloss.NewStyle().Background(lipgloss.Color("99")).Padding(0, 1)
	faintStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("255")).Faint(true)
	enumeratorStyler = lipgloss.NewStyle().Foreground(lipgloss.Color("99")).MarginRight(1)
)

const (
	padding  = 2
)

func (m model) View() string {
	s := appNameStyle.Render(("JET")) + "\n\n"

	if m.state == welcomeView {
		s += "Jet!\n\nUse your time efficiently\n"
		s += faintStyle.Render("s - configure study time, b - configure break time, q - quit")
	}
	if m.state == StudyView {
		remaining := m.currStudyTime.totalSeconds - m.currStudyTime.elapsedSeconds
		minutesLeft := uint(math.Ceil(float64(remaining) / 60))

		pad := strings.Repeat(" ", padding)
        return "\n" + 
		"Studying..." + faintStyle.Render(uintToString(minutesLeft)) + " minutes left\n\n"+ 
        pad + m.studyProgress.View() + "\n\n" +
		pad + faintStyle.Render("esc - cancel and go back")
	}
	if m.state == BreakView {
		remaining := m.currBreakTime.totalSeconds - m.currBreakTime.elapsedSeconds
		minutesLeft := uint(math.Ceil(float64(remaining) / 60))

		pad := strings.Repeat(" ", padding)
		return "\n" + 
		    "Break time! Rest for: "+ faintStyle.Render(uintToString(minutesLeft)) + 
			" minutes\n\n"+
            pad + m.breakProgress.View() + "\n\n" +
            pad + faintStyle.Render("esc - cancel and go back")
	}
	if m.state == StudyTimeSelectorView {
		for i, n := range m.studyTimes {
			prefix := ""
			if i == m.selectedStudyTime {
				prefix = ">"
			}
			s += enumeratorStyler.Render(prefix) + n.Name + " | " + 
			faintStyle.Render(uintToString(n.Minutes)) + " minutes" +"\n\n"
		}
		s += faintStyle.Render(" esc - go back, enter - set")
	}
	if m.state == BreakTimeSelectorView {
		for i, n := range m.breakTimes {
			prefix := ""
			if i == m.selectedBreakTime {
				prefix = ">"
			}
			s += enumeratorStyler.Render(prefix) + n.Name + " | " + 
			faintStyle.Render(uintToString(n.Minutes)) + " minutes" +"\n\n"
		}
		s += faintStyle.Render(" esc - go back, enter - set")
	}

	return s
}

func uintToString(num uint) string{
	return strconv.FormatUint(uint64(num), 10)
}