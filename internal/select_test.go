package internal

import (
	"testing"

	tea "charm.land/bubbletea/v2"
)

func TestSelectModelEnterChoosesSelectedItem(t *testing.T) {
	m := newSelectModel(map[string]string{"b.yaml": "/b.yaml", "a.yaml": "/a.yaml"})

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(selectModel)

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(selectModel)

	if m.chosen != "a.yaml" { // sorted first
		t.Errorf("expected chosen %q, got %q", "a.yaml", m.chosen)
	}
	if cmd == nil {
		t.Error("expected a quit command after selecting")
	}
}

func TestSelectModelQuitCancelsWithoutChoosing(t *testing.T) {
	m := newSelectModel(map[string]string{"a.yaml": "/a.yaml"})

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(selectModel)

	updated, _ = m.Update(tea.KeyPressMsg{Text: "q", Code: 'q'})
	m = updated.(selectModel)

	if m.chosen != "" {
		t.Errorf("expected no choice after quitting, got %q", m.chosen)
	}
}

func TestSelectModelEnterWhileFilteringDoesNotChoose(t *testing.T) {
	m := newSelectModel(map[string]string{"a.yaml": "/a.yaml", "b.yaml": "/b.yaml"})

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 24})
	m = updated.(selectModel)

	// Enter filtering mode ("/"), then press enter to confirm the (empty) filter.
	updated, _ = m.Update(tea.KeyPressMsg{Text: "/", Code: '/'})
	m = updated.(selectModel)

	updated, _ = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(selectModel)

	if m.chosen != "" {
		t.Errorf("expected no choice from confirming a filter, got %q", m.chosen)
	}
}
