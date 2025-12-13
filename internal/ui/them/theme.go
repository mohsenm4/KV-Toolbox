package them

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

const (
	ThemeDark   = "dark"
	ThemeLight  = "light"
	ThemeCustom = "custom"
)

func GetThemeKey(app fyne.App) string {
	t := app.Settings().Theme()
	currentBG := t.Color(theme.ColorNameBackground, app.Settings().ThemeVariant())
	darkBG := theme.DarkTheme().Color(theme.ColorNameBackground, app.Settings().ThemeVariant())
	lightBG := theme.LightTheme().Color(theme.ColorNameBackground, app.Settings().ThemeVariant())

	switch {
	case currentBG == darkBG:
		return ThemeDark
	case currentBG == lightBG:
		return ThemeLight
	default:
		return ThemeCustom
	}
}

func SetThemeByKey(app fyne.App, mytheme string) {
	if mytheme == ThemeDark {
		app.Settings().SetTheme(theme.DarkTheme())
	} else if mytheme == ThemeLight {
		app.Settings().SetTheme(theme.LightTheme())
	}
}
