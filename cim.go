package main

import (
	"fmt"
	"os"

	"charm.land/bubbles/v2/viewport"
	tea "charm.land/bubbletea/v2"
	"charm.land/glamour/v2"
	"charm.land/glamour/v2/styles"
	"charm.land/lipgloss/v2"
)

const content = `
# CIM - Vim IMpaired

## Cim version 0.0.1

by Connor Jones

Cim is open source and freely distributable

Lorem ipsum dolor sit amet consectetur adipiscing elit. Quisque faucibus ex sapien vitae pellentesque sem placerat. In id cursus mi pretium tellus duis convallis. Tempus leo eu aenean sed diam urna tempor. Pulvinar vivamus fringilla lacus nec metus bibendum egestas. Iaculis massa nisl malesuada lacinia integer nunc posuere. Ut hendrerit semper vel class aptent taciti sociosqu. Ad litora torquent per conubia nostra inceptos himenaeos.


iasdasd Lorem ipsum dolor sit amet consectetur adipiscing elit. Quisque faucibus ex sapien vitae pellentesque sem placerat. In id cursus mi pretium tellus duis convallis. Tempus leo eu aenean sed diam urna tempor. Pulvinar vivamus fringilla lacus nec metus bibendum egestas. Iaculis massa nisl malesuada lacinia integer nunc posuere. Ut hendrerit semper vel class aptent taciti sociosqu. Ad litora torquent per conubia nostra inceptos himenaeos.
`

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

/*-------------------------------------------*/
/*The model will store the applications state*/
/*-------------------------------------------*/
type model struct {
	viewport viewport.Model
}

/*-------------------------------------------*/
/*Define apps initial state*/
/*-------------------------------------------*/
func initModel(isDark bool) (*model, error) {
	const (
		width  = 78
		height = 20
	)

	vp := viewport.New()
	vp.SetWidth(width)
	vp.SetHeight(height)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)

	// We need to adjust the width of the glamour render from our main width
	// to account for a few things:
	//
	//  * The viewport border width
	//  * The viewport padding
	//  * The viewport margins
	//  * The gutter glamour applies to the left side of the content
	//
	const glamourGutter = 3
	glamourRenderWidth := width - vp.Style.GetHorizontalFrameSize() - glamourGutter

	style := styles.DarkStyleConfig
	if !isDark {
		style = styles.LightStyleConfig
	}
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStyles(style),
		glamour.WithWordWrap(glamourRenderWidth),
	)
	if err != nil {
		return nil, err
	}

	str, err := renderer.Render(content)
	if err != nil {
		return nil, err
	}

	vp.SetContent(str)

	return &model{
		viewport: vp,
	}, nil
}

func (m model) Init() tea.Cmd {
	return nil
}

/*-------------------------------------------*/
/*Handle State Changes in Update*/
/*-------------------------------------------*/
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "q", "ctrl+c", "esc":
			return m, tea.Quit
		default:
			var cmd tea.Cmd
			m.viewport, cmd = m.viewport.Update(msg)
			return m, cmd
		}
	default:
		return m, nil
	}
}

/*-------------------------------------------*/
/*Render the UI*/
/*-------------------------------------------*/
func (m model) View() tea.View {
	return tea.NewView(m.viewport.View() + m.helpView())
}

func (m model) helpView() string {
	return helpStyle("\n  j/k: Navigate • q: Quit\n")
}

/*-------------------------------------------*/
/*Let It Rip*/
/*-------------------------------------------*/
func main() {
	hasDarkBg := lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
	example, err := initModel(hasDarkBg)
	if err != nil {
		fmt.Println("Could not initialize Bubble Tea model:", err)
		os.Exit(1)
	}

	if _, err := tea.NewProgram(example).Run(); err != nil {
		fmt.Println("Bummer, there's been an error:", err)
		os.Exit(1)
	}
}
