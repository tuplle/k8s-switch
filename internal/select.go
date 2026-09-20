package internal

import (
	"maps"
	"slices"

	"charm.land/bubbles/v2/key"
	"charm.land/bubbles/v2/list"
	tea "charm.land/bubbletea/v2"
)

// keyEnter/keyQuit are shown in the list's help footer in place of bubbles/list's own
// defaults: as of bubbles v2.2.1, list.DefaultKeyMap's Quit binding is "v" with help
// text "select" (it quits without selecting anything — an upstream default we don't
// want to expose), so we disable it in newSelectModel and own cancel-key handling
// ourselves below.
var (
	keyEnter = key.NewBinding(key.WithKeys("enter"), key.WithHelp("enter", "select"))
	keyQuit  = key.NewBinding(key.WithKeys("q", "esc", "ctrl+c"), key.WithHelp("q/esc", "quit"))
)

// configItem adapts a kubeconfig file for display in a bubbles/list.
type configItem struct {
	name        string
	description string
}

func (i configItem) Title() string       { return i.name }
func (i configItem) Description() string { return i.description }
func (i configItem) FilterValue() string { return i.name }

// describeConfig returns a "<context> · <server>" summary of the kubeconfig at path,
// read via ReadKubeconfigSummary, for display in the picker. It falls back to path
// itself if the file can't be read/parsed or has no usable cluster/context info, so a
// single bad file doesn't break browsing the rest of the list.
func describeConfig(path string) string {
	summary, err := ReadKubeconfigSummary(path)
	if err != nil || (summary.ContextName == "" && summary.ServerAddress == "") {
		return path
	}
	switch {
	case summary.ContextName == "":
		return summary.ServerAddress
	case summary.ServerAddress == "":
		return summary.ContextName
	default:
		return summary.ContextName + "  ·  " + summary.ServerAddress
	}
}

// selectModel is the Bubble Tea model backing the interactive config picker.
type selectModel struct {
	list   list.Model
	chosen string
}

func newSelectModel(configs map[string]string) selectModel {
	names := slices.Sorted(maps.Keys(configs))
	items := make([]list.Item, len(names))
	for i, name := range names {
		items[i] = configItem{name: name, description: describeConfig(configs[name])}
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select config"
	// SetEnabled(false) alone isn't enough: list.Model.SetSize -> updateKeybindings
	// re-enables KeyMap.Quit on every resize unless this persistent flag is set too.
	l.DisableQuitKeybindings()
	l.AdditionalShortHelpKeys = func() []key.Binding { return []key.Binding{keyEnter, keyQuit} }
	l.AdditionalFullHelpKeys = func() []key.Binding { return []key.Binding{keyEnter, keyQuit} }

	return selectModel{list: l}
}

func (m selectModel) Init() tea.Cmd { return nil }

func (m selectModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.list.SetSize(msg.Width, msg.Height)
		return m, nil
	case tea.KeyPressMsg:
		filtering := m.list.FilterState() == list.Filtering
		switch msg.String() {
		case "ctrl+c":
			// Always cancels immediately, even mid-filter.
			return m, tea.Quit
		case "enter":
			if !filtering {
				if item, ok := m.list.SelectedItem().(configItem); ok {
					m.chosen = item.name
				}
				return m, tea.Quit
			}
		case "q", "esc":
			// While filtering, let these fall through to the list itself
			// ("q" types into the filter query, "esc" clears the filter).
			if !filtering {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m selectModel) View() tea.View {
	return tea.NewView(m.list.View())
}

// SelectConfig runs the interactive picker over configs (name -> absolute path)
// and returns the chosen name. ok is false if the user cancelled (q/esc/ctrl+c).
func SelectConfig(configs map[string]string) (name string, ok bool, err error) {
	finalModel, err := tea.NewProgram(newSelectModel(configs)).Run()
	if err != nil {
		return "", false, err
	}
	m := finalModel.(selectModel)
	return m.chosen, m.chosen != "", nil
}
