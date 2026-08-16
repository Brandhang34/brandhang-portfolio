package web

import "unicode"

// contactFieldClass is shared by every field on the contact form; the class
// string was previously repeated verbatim on each input.
const contactFieldClass = "font-primary block p-2.5 w-full text-sm rounded-lg border shadow-sm " +
	"bg-gray-700 border-gray-600 placeholder-gray-400 text-white " +
	"focus:ring-blue-500 focus:border-blue-500"

// initial returns the first letter of s, uppercased, for use in the placeholder
// tile shown when a portfolio entry has no image.
func initial(s string) string {
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			return string(unicode.ToUpper(r))
		}
	}
	return "?"
}
