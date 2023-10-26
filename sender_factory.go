package mslibpak

import (
	"os"

	"github.com/paketo-buildpacks/libpak/bard"
)

type senderFactory interface {
	getSender() (sender, error)
}

type senderFactoryImpl struct {
	*bard.Logger
}

func newSenderFactory(logger *bard.Logger) senderFactory {
	return &senderFactoryImpl{logger}
}

func (sf senderFactoryImpl) getSender() (sender, error) {
	t := os.Getenv(MonitoringTools)
	switch t {
	case TelemetryOutputConsole:
		return newConsoleSender(sf.Logger)
	case TelemetryOutputAi:
		return newAISender()
	default:
		return newConsoleSender(sf.Logger)
	}
}
