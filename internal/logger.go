/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package internal

import (
	"io"
	"os"
	"strings"

	"github.com/buildpacks/libcnb/poet"
	"github.com/heroku/color"
)

// TODO: Remove once TTY support is in place
func init() {
	color.Enabled()
}

// Logger logs message to a writer.
type Logger struct {
	poet.Logger
}

// Option is a function for configuring a Logger instance.
type Option func(logger Logger) Logger

// WithDebug configures the debug Writer.
func WithDebug(writer io.Writer) Option {
	return func(logger Logger) Logger {
		logger.Logger = poet.WithDebug(writer)(logger.Logger)
		return logger
	}
}

// NewLoggerWithOptions create a new instance of Logger.  It configures the Logger with options.
func NewLoggerWithOptions(writer io.Writer, options ...Option) Logger {
	l := Logger{
		Logger: poet.NewLogger(writer),
	}

	for _, option := range options {
		l = option(l)
	}

	return l
}

// NewLogger creates a new instance of Logger.  It configures debug logging if $BP_DEBUG is set.
func NewLogger(writer io.Writer) Logger {
	var options []Option

	// check for presence and value of log level environment variable
	options = LogLevel(options, writer)

	return NewLoggerWithOptions(writer, options...)
}

func LogLevel(options []Option, writer io.Writer) []Option {

	// Check for older log level env variable
	_, dbSet := os.LookupEnv("BP_DEBUG")

	// Then check for common buildpack log level env variable - if either are set to DEBUG/true, enable Debug Writer
	if level, ok := os.LookupEnv("BP_LOG_LEVEL"); (ok && strings.ToLower(level) == "debug") || dbSet {

		options = append(options, WithDebug(writer))
	}
	return options
}
