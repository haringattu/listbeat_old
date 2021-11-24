package beater

import (
	"fmt"
	"time"
	"io/ioutil"
	"path/filepath"

	"github.com/elastic/beats/v7/libbeat/beat"
	"github.com/elastic/beats/v7/libbeat/common"
	"github.com/elastic/beats/v7/libbeat/logp"

	"github.com/haringattu/listbeat/config"
)

// listbeat configuration.
type listbeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
}

// New creates an instance of listbeat.
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &listbeat{
		done:   make(chan struct{}),
		config: c,
	}
	return bt, nil
}

// Run starts listbeat.
func (bt *listbeat) Run(b *beat.Beat) error {
	logp.Info("listbeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}

	ticker := time.NewTicker(bt.config.Period)
	counter := 1
	for {
		bt.listDir(bt.config.Path, b.Info.Name) // call listDir
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		/*event := beat.Event{
			Timestamp: time.Now(),
			Fields: common.MapStr{
				"type":    b.Info.Name,
				"counter": counter,
			},
		}
		bt.client.Publish(event)*/
		logp.Info("Event sent")
		counter++
	}
}

// Stop stops listbeat.
func (bt *listbeat) Stop() {
	bt.client.Close()
	close(bt.done)
}

func (bt *listbeat) listDir(dirFile string, beatname string) {
	files, _ := ioutil.ReadDir(dirFile)
	for _, f := range files {
		t := f.ModTime()
		path := filepath.Join(dirFile, f.Name())
		event := beat.Event {
			Fields: common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       beatname,
				"modtime":    common.Time(t),
				"filename":   f.Name(),
				"path":       path,
				"directory":  f.IsDir(),
				"filesize":   f.Size(),
			},
		}
		bt.client.Publish(event)
		if f.IsDir() {
			bt.listDir(path, beatname)
		}
	}
}
