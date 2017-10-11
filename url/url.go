package url

// Url base url structure
type URL struct {
	url string
}

// FromString generate a new Url structure from a sample url string
func FromString(url string) *URL {
	u := new(URL)
	u.url = url
	return u
}
