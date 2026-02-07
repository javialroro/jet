package main

import (
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type tickMsg time.Time

type Time struct {
	ID int64
	Name string
	Minutes uint
	totalSeconds uint
	elapsedSeconds uint
}

const (
	welcomeView                = 0
	StudyTimeSelectorView uint = 1
	BreakTimeSelectorView uint = 2
	StudyView             uint = 3
	BreakView             uint = 4
)

var (
	debugStudy = Time{5,"Debug", 1, 60, 0}
	shortStudy = Time{0, "Short Session", 25, 1500, 0}
	longStudy  = Time{1, "Long Session", 45, 2700, 0}

	debugBreak = Time{5,"Debug", 1, 60, 0}
	shortBreak = Time{0, "Short Break", 5, 300, 0}
	longBreak  = Time{1, "Long Break", 15, 900, 0}
)

var studyTimes = []Time{
	debugStudy,
	shortStudy,
	longStudy,
}

var breakTimes = []Time{
	debugBreak,
	shortBreak,
	longBreak,
}

type model struct {
	studyProgress     progress.Model
	breakProgress     progress.Model
	selectedStudyTime int
	selectedBreakTime int
	breakTimes        []Time
	studyTimes        []Time
	state             uint
	currStudyTime     Time
	currBreakTime     Time
}

func NewModel() model {
	return model{
		breakTimes:        breakTimes,
		studyTimes:        studyTimes,
		selectedStudyTime: 0,
		selectedBreakTime: 0,
		state:             welcomeView,
		studyProgress:     progress.New(progress.WithDefaultGradient()),
		breakProgress:     progress.New(progress.WithDefaultGradient()),
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var (
		cmds []tea.Cmd
		cmd  tea.Cmd
	)
	cmds = append(cmds, cmd)

	cmds = append(cmds, cmd)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		key := msg.String()
		switch m.state {
		case welcomeView:
			switch key {
			case "q":
				return m, tea.Quit
			case "enter":
				m.state = StudyTimeSelectorView

			}
		case StudyTimeSelectorView:
			switch key {
			case "up", "k":
				if m.selectedStudyTime > 0 {
					m.selectedStudyTime--
				}
			case "down", "j":
				if m.selectedStudyTime < len(studyTimes)-1 {
					m.selectedStudyTime++
				}
			case "enter":
				m.currStudyTime = studyTimes[m.selectedStudyTime]
				m.state = BreakTimeSelectorView

			case "esc":
				m.state = welcomeView
			case "q":
				return m, tea.Quit
			}

		case BreakTimeSelectorView:
			switch key {
			case "up", "k":
				if m.selectedBreakTime > 0 {
					m.selectedBreakTime--
				}
			case "down", "j":
				if m.selectedBreakTime < len(breakTimes)-1 {
					m.selectedBreakTime++
				}
			case "enter":
				m.currBreakTime = breakTimes[m.selectedBreakTime]
				m.state = StudyView
				return m, tickCmd()

			case "esc":
				m.state = welcomeView
			case "q":
				return m, tea.Quit
			}

		case StudyView:
			switch key {

			case "esc":
				m.state = StudyTimeSelectorView
			case "q":
				return m, tea.Quit
			}
		case BreakView:
			switch key {
			case "esc":
				m.state = welcomeView
			}

		}
	case tickMsg:
		if m.state != StudyView && m.state != BreakView {
			return m, nil
		}
		if m.state == StudyView {
			if m.currStudyTime.elapsedSeconds > m.currStudyTime.totalSeconds{
				m.state = BreakView
				m.breakProgress.SetPercent(0)
				m.currStudyTime.elapsedSeconds = 0
				return m, tickCmd()
			}
			m.currStudyTime.elapsedSeconds++
			percent := float64(m.currStudyTime.elapsedSeconds) / float64(m.currStudyTime.totalSeconds)

			cmd := m.studyProgress.SetPercent(percent)
			return m, tea.Batch(tickCmd(), cmd)
		} else {
			if m.currBreakTime.elapsedSeconds > m.currBreakTime.totalSeconds{
				m.state = StudyView
				m.studyProgress.SetPercent(0)
				m.currBreakTime.elapsedSeconds = 0
				return m, tickCmd()
			}
			m.currBreakTime.elapsedSeconds++
			percent := float64(m.currBreakTime.elapsedSeconds) / float64(m.currBreakTime.totalSeconds)

			cmd := m.breakProgress.SetPercent(percent)
			return m, tea.Batch(tickCmd(), cmd)

		}
	case progress.FrameMsg:
		if m.state != StudyView && m.state != BreakView {
			return m, nil
		}

		if m.state == StudyView {
			pm, cmd := m.studyProgress.Update(msg)
			m.studyProgress = pm.(progress.Model)
			return m, cmd
		}

		pm, cmd := m.breakProgress.Update(msg)
		m.breakProgress = pm.(progress.Model)
		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*1, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
