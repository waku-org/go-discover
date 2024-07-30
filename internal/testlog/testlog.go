// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package testlog provides a log handler for unit tests.
package testlog

import (
	"context"
	"log/slog"
	"os"
	"sync"
	"testing"

	"github.com/ethereum/go-ethereum/log"
)

type logger struct {
	l      log.Logger
	mu     *sync.Mutex
	prefix string
}

// Logger returns a logger which logs to the unit test log of t.
func Logger(t *testing.T, level slog.Level, optionalPrefix ...string) log.Logger {
	h := log.NewTerminalHandlerWithSource(os.Stderr, level, false)
	prefix := ""
	if len(optionalPrefix) > 0 {
		prefix = optionalPrefix[0]
	}
	return &logger{l: log.NewLoggerWithOpts(h, &log.LoggerOptions{SkipCallers: 1}), mu: new(sync.Mutex), prefix: prefix}
}

func (l *logger) With(ctx ...interface{}) log.Logger {
	return &logger{l: l.l.With(ctx...), mu: l.mu, prefix: l.prefix}
}

func (l *logger) New(ctx ...interface{}) log.Logger {
	return l.With(ctx...)
}

func (l *logger) Log(level slog.Level, msg string, ctx ...interface{}) {
	l.Write(level, msg, ctx...)
}

func (l *logger) Trace(msg string, ctx ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Write(log.LevelTrace, msg, ctx...)
}

func (l *logger) Debug(msg string, ctx ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Write(log.LevelDebug, msg, ctx...)
}

func (l *logger) Info(msg string, ctx ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Write(log.LevelInfo, msg, ctx...)
}

func (l *logger) Warn(msg string, ctx ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Write(log.LevelWarn, msg, ctx...)
}

func (l *logger) Error(msg string, ctx ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Write(log.LevelError, msg, ctx...)
}

func (l *logger) Crit(msg string, ctx ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.Write(log.LevelCrit, msg, ctx...)
	os.Exit(1)
}

func (l *logger) Write(level slog.Level, msg string, attrs ...any) {
	if l.prefix != "" {
		msg = l.prefix + " " + msg
	}
	l.l.Write(level, msg, attrs...)
}

func (l *logger) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

func (l *logger) Handler() slog.Handler {
	return l.l.Handler()
}
