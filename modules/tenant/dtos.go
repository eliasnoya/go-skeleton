package tenant

import (
	"encoding/json"
	"tuples/helpers"
)

type Settings struct {
	Theme Theme `json:"theme" binding:"required"`
}

type Theme struct {
	Logo    string      `json:"logo" binding:"required"`
	Default string      `json:"default" binding:"required"`
	Light   ColorScheme `json:"light" binding:"required"`
	Dark    ColorScheme `json:"dark" binding:"required"`
}

type ColorScheme struct {
	Background string `json:"background" binding:"required"`
	Primary    string `json:"primary" binding:"required"`
	Secondary  string `json:"secondary" binding:"required"`
}

func (cs *ColorScheme) AsJson() string {
	res, err := json.Marshal(cs)
	if err != nil {
		return "{}"
	}
	return string(res)
}

// from json String get ColorSchema DTO
func NewColorSchemeFromJson(jsonStr string) ColorScheme {
	var cs ColorScheme
	json.Unmarshal([]byte(jsonStr), &cs)
	return cs
}

func DefaultColorScheme(typo string) ColorScheme {
	if typo == "light" {
		return ColorScheme{
			Background: "#f5f5f5",
			Primary:    helpers.TuplesPrimary,
			Secondary:  helpers.TuplesSecondary,
		}
	}
	// DARK
	return ColorScheme{
		Background: "#121212",
		Primary:    helpers.TuplesSecondary,
		Secondary:  helpers.TuplesTertiary,
	}
}

type SetupStepRequest struct {
	Step int `json:"step" validate:"exists"`
}

type SetupThemeRequest struct {
	NextStep     int         `json:"next_step" binding:"required"`
	Logo         string      `json:"logo" binding:"required"`
	DefaultTheme string      `json:"default_theme" binding:"required"`
	Light        ColorScheme `json:"light" binding:"required"`
	Dark         ColorScheme `json:"dark" binding:"required"`
}
