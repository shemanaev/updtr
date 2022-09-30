package container

import (
	"html/template"

	"github.com/containrrr/watchtower/pkg/container"
)

type Info struct {
	ID          string
	Name        string
	ImageName   string
	State       State               // TODO: add message for Error state
	Changelog   template.HTML       `json:"-"`
	LatestImage string              `json:"-"`
	Container   container.Container `json:"-"`
}

//go:generate go run github.com/dmarkham/enumer -type=State -json
type State int

const (
	Stale State = iota
	RequestedUpdate
	Updating
	Error
	Fresh
)
