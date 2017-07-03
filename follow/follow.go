package follow

import (
	"context"
	"log"
	"net/url"
	"time"

	"github.com/ChimeraCoder/anaconda"
	"github.com/remeh/smartwitter/twitter"
)

func Follow(ctx context.Context) {
	for {
		after := time.After(time.Second * 3)

		// ----------------------

		select {
		case <-after:
			log.Println("debug: Follow is starting.")
			if err := run(ctx); err != nil {
				log.Println("error: while running Follow:", err)
			}
			log.Println("debug: Follow is ending.")
		case <-ctx.Done():
			log.Println("debug: Follow canceled.")
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
			log.Println(t.User.Name, ": ", t.Text)
		case <-ctx.Done():
			stream.Interrupt() // NOTE(remy): don't know if better to interrupt or to End() here
			return nil
		}
	}

	return nil
}
