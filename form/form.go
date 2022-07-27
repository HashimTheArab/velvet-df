package form

import "github.com/df-mc/dragonfly/server/player/form"

// nopSubmit is a nop submitter for a form that does nothing extra when submitted.
type nopSubmit struct{}

// NopSubmit ...
var NopSubmit nopSubmit

// Submit ...
func (nopSubmit) Submit(form.Submitter, form.Button) {}
