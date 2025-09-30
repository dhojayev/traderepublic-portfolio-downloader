package app

import (
	"os"
	"path/filepath"

	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal"
	"github.com/dhojayev/traderepublic-portfolio-downloader/v2/internal/bus"
)

type Faker struct {
	eventBus *bus.EventBus
}

func NewFaker(eventBus *bus.EventBus) *Faker {
	return &Faker{
		eventBus: eventBus,
	}
}

func (f *Faker) Run() error {
	dir := internal.ResponseBaseDir + "/timelineDetailV2"
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		contents, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			return err
		}

		f.eventBus.Publish(bus.NewEvent(bus.TopicTimelineDetailsV2, filenameWithoutExt(entry.Name()), bus.EventNameTimelineDetailV2Received, contents))
	}

	return nil
}

func filenameWithoutExt(path string) string {
	return filepath.Base(path[:len(path)-len(filepath.Ext(path))])
}
