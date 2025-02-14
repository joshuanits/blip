package main

import (
	"fmt"
	"os"
  "strings"
  "syscall"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
  "github.com/twpayne/go-shell"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
  tags []string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc + " " + i.FormatTags() }
func (i item) FilterValue() string { return i.title + strings.Join(i.tags, " ") }
func (i item) FormatTags() string {
  if len(i.tags) == 0 {
    return ""
  }

  return "(tags: " + strings.Join(i.tags, ", ") + ")"
}

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

func runCommand(cmd string) {
  shell, ok := shell.CurrentUserShell()
  if !ok { shell = "sh" }

  if err := syscall.Exec("/usr/bin/sh", []string{shell, "-c", cmd}, os.Environ()); err != nil {
    fmt.Println(err)
    os.Exit(1)
	}

}

func main() {
	items := []list.Item{
    item{title: "Raspberry Pi’s", desc: "I have ’em all over my house", tags: []string{"computer", "electronics"}},
    item{title: "Warm light", desc: "Like around 2700 Kelvin", tags: []string{"electronics"}},
		item{title: "The vernal equinox", desc: "The autumnal equinox is pretty good too"},
		item{title: "Gaffer’s tape", desc: "Basically sticky fabric"},
		item{title: "Terrycloth", desc: "In other words, towel fabric"},
	}

	m := model{list: list.New(items, list.NewDefaultDelegate(), 0, 0)}
	m.list.Title = "My Fave Things"

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
