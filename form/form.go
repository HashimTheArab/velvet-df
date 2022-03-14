package form

import "github.com/df-mc/dragonfly/server/player/form"

type nopSubmit struct{}

var NopSubmit nopSubmit

func (nopSubmit) Submit(form.Submitter, form.Button) {}
