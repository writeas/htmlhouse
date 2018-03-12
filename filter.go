package htmlhouse

func passesPublicFilter(app *app, html string) bool {
	spam := app.cfg.BlacklistReg.MatchString(html)
	return !spam
}
