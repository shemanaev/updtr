package config

import (
	"time"

	env "github.com/caitlinelfring/go-env-default"
)

type Settings struct {
	Debug             bool
	Trace             bool
	RemoveVolumes     bool
	IncludeRestarting bool
	IncludeStopped    bool
	ReviveStopped     bool
	NoPull            bool
	WarnOnHeadFailure string
	PollInterval      int
	Schedule          string
	EnableLabel       bool
	Scope             string
	StopTimeout       time.Duration
}

func LoadConfig() Settings {
	var s Settings

	day := int((time.Hour * 24).Seconds())
	timeout, _ := time.ParseDuration("10s")

	s.Debug = env.GetBoolDefault("WATCHTOWER_DEBUG", false)
	s.Trace = env.GetBoolDefault("WATCHTOWER_TRACE", false)
	s.RemoveVolumes = env.GetBoolDefault("WATCHTOWER_REMOVE_VOLUMES", false)
	s.IncludeRestarting = env.GetBoolDefault("WATCHTOWER_INCLUDE_RESTARTING", false)
	s.IncludeStopped = env.GetBoolDefault("WATCHTOWER_INCLUDE_STOPPED", false)
	s.ReviveStopped = env.GetBoolDefault("WATCHTOWER_REVIVE_STOPPED", false)
	s.NoPull = env.GetBoolDefault("WATCHTOWER_NO_PULL", false)
	s.WarnOnHeadFailure = env.GetDefault("WATCHTOWER_WARN_ON_HEAD_FAILURE", "auto")
	s.PollInterval = env.GetIntDefault("WATCHTOWER_POLL_INTERVAL", day)
	s.Schedule = env.GetDefault("WATCHTOWER_SCHEDULE", "")
	s.EnableLabel = env.GetBoolDefault("WATCHTOWER_LABEL_ENABLE", false)
	s.Scope = env.GetDefault("WATCHTOWER_SCOPE", "")
	s.StopTimeout = env.GetDurationDefault("WATCHTOWER_TIMEOUT", timeout)

	return s
}
