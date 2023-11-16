package sanitizer

import "github.com/microcosm-cc/bluemonday"

var XSS = bluemonday.StrictPolicy()
