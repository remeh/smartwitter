package suggest

import (
	"fmt"

	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/uuid"

	"github.com/lib/pq"
)

// SetKeywords stores the given keywords for the user.
// It also resets the last run on these, forcing a near-to-be executed
// crawling.
func SetKeywords(userUid uuid.UUID, keywords []string, position int) error {
	if position < 0 {
		return fmt.Errorf("SetKeywords: position <1")
	}

	if _, err := storage.DB().Exec(`
		UPDATE "twitter_keywords_watcher"
		SET
			"keywords" = $1
		WHERE
			"user_uid" = $2
			AND
			"position" = $3
	`, pq.Array(keywords), userUid, position); err != nil {
		return err
	}

	return nil
}
