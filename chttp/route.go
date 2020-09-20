package chttp

import (
	"net/http"
	"regexp"
	"sort"
	"strings"
)

// Route represents a single path that the http server is accepting requests on.
// The route can be configured with middleware functions.
// Additionally, it can be limited to accept requests on specific http methods.
type Route struct {
	MiddlewareFuncs []MiddlewareFunc
	Path            string
	Methods         []string
	Handler         http.Handler
}

func sortRoutes(routes []Route) {
	const matcherPlaceholder = "{{matcher}}"

	var re = regexp.MustCompile(`(?U)(\{.*\})`)

	sort.Slice(routes, func(i, j int) bool {
		aPath := re.ReplaceAllString(routes[i].Path, matcherPlaceholder)
		bPath := re.ReplaceAllString(routes[j].Path, matcherPlaceholder)

		aParts := strings.Split(aPath, "/")
		bParts := strings.Split(bPath, "/")

		if aPath == "/" {
			aParts = nil
		}

		if bPath == "/" {
			bParts = nil
		}

		if len(aParts) != len(bParts) {
			return len(aParts) > len(bParts)
		}

		for i, aPart := range aParts {
			bPart := bParts[i]

			if aPart == matcherPlaceholder {
				return false
			}

			if bPart == matcherPlaceholder {
				return true
			}
		}

		return false
	})
}
