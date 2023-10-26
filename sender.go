package mslibpak

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/microsoft/ApplicationInsights-Go/appinsights"
	"github.com/paketo-buildpacks/libpak/bard"
)

type sender interface {
	send(context.Context, Metric)
	shutdown(context.Context) error
}

type Metric struct {
	Name       string
	Properties map[string]string
}

type aisender struct {
	c appinsights.TelemetryClient
}

func newAISender() (sender, error) {
	ikey := os.Getenv(EnvAIConnectionStringKey)
	if len(ikey) == 0 {
		return nil, fmt.Errorf("Instrumentation key is empty")
	}
	return &aisender{c: appinsights.NewTelemetryClient(ikey)}, nil
}

func (a aisender) send(_ context.Context, m Metric) {
	event := appinsights.NewEventTelemetry(m.Name)
	event.Properties = m.Properties

	go a.c.Track(event)
}

func (a aisender) shutdown(ctx context.Context) error {
	select {
	case <-a.c.Channel().Close(10 * time.Second):
		// If we got here, then all telemetry was submitted
		// successfully, and we can proceed to exiting.
	case <-time.After(30 * time.Second):
	case <-ctx.Done():
	}
	return nil
}

type consolesender struct {
	*bard.Logger
}

func newConsoleSender(logger *bard.Logger) (sender, error) {
	return consolesender{logger}, nil
}

func (c consolesender) send(ctx context.Context, m Metric) {
	c.Logger.Header(fmt.Sprintf("%s %s", m.Name, m.Properties))
}

func (c consolesender) shutdown(ctx context.Context) error {
	// do nothing
	return nil
}

type noopSender struct {
}

func newNoopSender() (sender, error) {
	return noopSender{}, nil
}

func (n noopSender) send(ctx context.Context, m Metric) {
}

func (n noopSender) shutdown(ctx context.Context) error {
	return nil
}
