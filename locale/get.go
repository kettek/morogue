package locale

import "github.com/cubiest/jibberjabber"

// Locale returns the user's locale.
func Locale() string {
	userLocale, err := jibberjabber.DetectIETF()
	if err != nil {
		return "en-US" // Just default what we know to support 100% if we can't detect.
	}
	return userLocale
}
