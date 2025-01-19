package locale

import (
	"fmt"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// Localizer stores the localizer for a given language.
type Localizer struct {
	ID      string
	printer *message.Printer
}

// Initialize a slice which holds the initialized Localizer types for
// each of our supported locales.
var locales = []Localizer{
	{
		ID:      "en-US",
		printer: message.NewPrinter(language.MustParse("en-US")),
	},
	{
		ID:      "ja-JP",
		printer: message.NewPrinter(language.MustParse("ja-JP")),
	},
	{
		ID:      "tr-TR",
		printer: message.NewPrinter(language.MustParse("tr-TR")),
	},
}

// Get returns the Localizer for the given ID.
func Get(id string) Localizer {
	for _, locale := range locales {
		if id == locale.ID {
			return locale
		}
	}
	// If we can't find our exact locale, see if one exists with the prefix...
	for _, locale := range locales {
		if locale.ID[:2] == id[:2] {
			return locale
		}
	}
	// Finally, just fall back to en-US.
	fmt.Printf("%s locale unsupported, falling back to en-US\n", id)
	return locales[0]
}

// T translates, yo.
func (l Localizer) T(key message.Reference, args ...interface{}) string {
	return l.printer.Sprintf(key, args...)
}
