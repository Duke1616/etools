package logger

import (
	"log/slog"
)

type SLogger struct {
	l *slog.Logger
}

func NewSLogger(l *slog.Logger) *SLogger {
	return &SLogger{
		l: l,
	}
}

func (s *SLogger) Debug(msg string, args ...Field) {
	s.l.Debug(msg, s.toArgs(args))
}

func (s *SLogger) Info(msg string, args ...Field) {
	s.l.Info(msg, s.toArgs(args))
}

func (s *SLogger) Warn(msg string, args ...Field) {
	s.l.Warn(msg, s.toArgs(args))
}

func (s *SLogger) Error(msg string, args ...Field) {
	s.l.Error(msg, s.toArgs(args))
}

func (s *SLogger) toArgs(args []Field) []slog.Attr {
	res := make([]slog.Attr, 0, len(args))
	for _, arg := range args {
		res = append(res, slog.Any(arg.Key, arg.Val))
	}
	return res
}
