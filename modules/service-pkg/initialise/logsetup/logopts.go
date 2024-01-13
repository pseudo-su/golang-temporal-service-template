package logsetup

import "log/slog"

type LogMode string

const (
	LogModeText     LogMode = "text"
	LogModeJson     LogMode = "json"
	LogModeDisabled LogMode = "disabled"
)

type loggerOpts struct {
	mode  LogMode
	level slog.Level
}

type LoggerOpt func(opts *loggerOpts)

func WithLevel(in slog.Level) LoggerOpt {
	return func(opts *loggerOpts) {
		opts.level = in
	}
}

func WithLevelStr(s string) LoggerOpt {
	return func(opts *loggerOpts) {
		var level slog.Level
		_ = level.UnmarshalText([]byte(s))
		opts.level = level
	}
}

func WithMode(mode LogMode) LoggerOpt {
	return func(opts *loggerOpts) {
		opts.mode = mode
	}
}

func WithModeStr(s string) LoggerOpt {
	return func(opts *loggerOpts) {
		opts.mode = modeFromStr(s)
	}
}

func modeFromStr(s string) LogMode {
	if s == string(LogModeJson) {
		return LogModeJson
	}
	if s == string(LogModeDisabled) {
		return LogModeDisabled
	}
	return LogModeText
}
