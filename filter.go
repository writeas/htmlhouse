package htmlhouse

func passesPublicFilter(app *app, html string) bool {
	if app.cfg.BlacklistTerms == "" {
		return true
	}

	spam := app.cfg.BlacklistReg.MatchString(html)
	return !spam
}
