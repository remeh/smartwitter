package follow

import (
	"context"
	"net/url"
	"time"

	"github.com/remeh/anaconda"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/twitter"
)

func Follow(ctx context.Context) {
	for {
		after := time.After(time.Second * 3)

		// ----------------------

		select {
		case <-after:
			log.Debug("Follow is starting.")
			if err := run(ctx); err != nil {
				log.Error("while running Follow:", err)
			}
			log.Debug("Follow is ending.")
		case <-ctx.Done():
			log.Debug("Follow canceled.")
			return
		}
	}
}

func run(ctx context.Context) error {

	// TODO(remy): build the slice of ids to follow
	// ----------------------

	v := url.Values{
		"follow": []string{"4927618108", "24744541", "25073877"},
	}

	stream := twitter.GetApi().PublicStreamFilter(v)

	for {
		select {
		case s := <-stream.C:
			t := s.(anaconda.Tweet)
			log.Debug(t.User.Name, ": ", t.Text)
		case <-ctx.Done():
			stream.Stop()
			return nil
		}
	}

	return nil
}
