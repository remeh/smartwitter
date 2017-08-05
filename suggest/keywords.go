package suggest

import (
	"fmt"
	"time"

	"github.com/remeh/smartwitter/log"
	"github.com/remeh/smartwitter/storage"
	"github.com/remeh/uuid"

	"github.com/lib/pq"
)

type Keywords []keywords
type keywords struct {
	Uid          uuid.UUID `json:"-"`
	UserUid      uuid.UUID `json:"-"`
	Keywords     []string  `json:"keywords"`
	Label        string    `json:"label"`
	Position     int       `json:"position"`
	CreationTime time.Time `json:"creation_time"`
	LastUpdate   time.Time `json:"last_update"`
}

func (k Keywords) Len() int               { return len(k) }
func (k Keywords) Less(i int, j int) bool { return k[i].Position < k[j].Position }
func (k Keywords) Swap(i int, j int)      { k[i], k[j] = k[j], k[j] }

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

// GetKeywords returns users watched keywords.
func GetKeywords(userUid uuid.UUID) (Keywords, error) {
	rv := make(Keywords, 0)

	if userUid.IsNil() {
		log.Warning("GetKeywords called with a nil userUid")
		return rv, nil
	}

	rows, err := storage.DB().Query(`
		SELECT "uid", "label", "keywords", "position", "creation_time", "last_update"
		FROM
			"twitter_keywords_watcher"
		WHERE
			"user_uid" = $1
		ORDER BY "position"
	`, userUid)
	if err != nil {
		return rv, log.Err("GetKeywords", err)
	}

	defer rows.Close()

	for rows.Next() {
		k := keywords{}
		if err := rows.Scan(
			&k.Uid,
			&k.Label,
			pq.Array(&k.Keywords),
			&k.Position,
			&k.CreationTime,
			&k.LastUpdate,
		); err != nil {
			return rv, log.Err("GetKeywords", err)
		}
		k.UserUid = userUid
		rv = append(rv, k)
	}
	return rv, nil
}
