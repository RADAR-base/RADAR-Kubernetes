package wizard

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/huh"
)

// backKeyMap returns a huh KeyMap that binds Esc (in addition to Ctrl+C) to the
// Quit action and advertises it as "esc back" in the form's help footer. The
// step loop in runWizard treats huh.ErrUserAborted as "go back one step" rather
// than as a cancellation.
func backKeyMap() *huh.KeyMap {
	km := huh.NewDefaultKeyMap()
	km.Quit = key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("esc", "back"),
	)
	return km
}

// runForm runs a huh.Form with the back-aware keymap. Use this in place of
// `form.Run()` everywhere in the wizard so the user can press Esc at any prompt
// to return to the previous step.
func runForm(form *huh.Form) error {
	return form.WithKeyMap(backKeyMap()).Run()
}
