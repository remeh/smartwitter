package suggest

// SetKeywords stores the given keywords for the user.
// It also resets the last run on these, forcing a near-to-be executed
// crawling.
func SetKeywords(userUid uuid.UUID, keywords []string, position int) error {
	if _, err := storage.DB().Exec(`
		
	`); err != nil {
		return err
	}
}
