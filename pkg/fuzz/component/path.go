package component

import (
	"context"
	"net/url"
	"strconv"
	"strings"

	"github.com/projectdiscovery/nuclei/v3/pkg/fuzz/dataformat"
	"github.com/projectdiscovery/retryablehttp-go"
)

// Path is a component for a request Path
type Path struct {
	value *Value

	originalPath string
	req          *retryablehttp.Request
}

var _ Component = &Path{}

// NewPath creates a new URL component
func NewPath() *Path {
	return &Path{}
}

// Name returns the name of the component
func (q *Path) Name() string {
	return RequestPathComponent
}

// Parse parses the component and returns the
// parsed component
func (q *Path) Parse(req *retryablehttp.Request) (bool, error) {
	q.req = req
	// Split the escaped path so that an encoded slash such as a%2Fb stays
	// part of its segment instead of starting a new one.
	q.originalPath = req.EscapedPath()
	q.value = NewValue("")

	split := strings.Split(q.originalPath, "/")
	values := make(map[string]interface{})
	for i, segment := range split {
		if segment == "" && i == 0 {
			// Skip the first empty segment from leading "/"
			continue
		}
		if segment == "" {
			// Skip any other empty segments
			continue
		}
		// Use 1-based indexing and store individual segments
		key := strconv.Itoa(len(values) + 1)
		if decoded, err := url.PathUnescape(segment); err == nil {
			segment = decoded
		}
		values[key] = segment
	}
	q.value.SetParsed(dataformat.KVMap(values), "")
	return true, nil
}

// Iterate iterates through the component
func (q *Path) Iterate(callback func(key string, value interface{}) error) (err error) {
	q.value.parsed.Iterate(func(key string, value any) bool {
		if errx := callback(key, value); errx != nil {
			err = errx
			return false
		}
		return true
	})
	return
}

// SetValue sets a value in the component
// for a key
func (q *Path) SetValue(key string, value string) error {
	if !q.value.SetParsedValue(key, value) {
		return ErrSetValue
	}
	return nil
}

// Delete deletes a key from the component
func (q *Path) Delete(key string) error {
	if !q.value.Delete(key) {
		return ErrKeyNotFound
	}
	return nil
}

// Rebuild returns a new request with the
// component rebuilt
func (q *Path) Rebuild() (*retryablehttp.Request, error) {
	// Use the snapshot of the original path to avoid mutation from Clone/UpdateRelPath
	originalSplit := strings.Split(q.originalPath, "/")

	// Create a new slice to hold the rebuilt segments
	rebuiltSegments := make([]string, 0, len(originalSplit))

	// Add the first empty segment (from leading "/")
	if len(originalSplit) > 0 && originalSplit[0] == "" {
		rebuiltSegments = append(rebuiltSegments, "")
	}

	// Process each segment
	segmentIndex := 1 // 1-based indexing for our stored values
	for i := 1; i < len(originalSplit); i++ {
		originalSegment := originalSplit[i]
		if originalSegment == "" {
			// Skip empty segments
			continue
		}

		// Check if we have a replacement for this segment
		key := strconv.Itoa(segmentIndex)
		newValue, exists := q.value.parsed.Map.GetOrDefault(key, "").(string)
		if decoded, err := url.PathUnescape(originalSegment); err == nil && newValue == decoded {
			// Unchanged segments keep their original escaping.
			exists = false
		}
		if exists && newValue != "" {
			rebuiltSegments = append(rebuiltSegments, escapePathSegment(newValue))
		} else {
			rebuiltSegments = append(rebuiltSegments, originalSegment)
		}
		segmentIndex++
	}

	// Join the segments back into a path. Every segment is escaped at this
	// point, so a "?" or "#" in the original path or in a payload stays in the
	// path instead of starting the query or the fragment.
	rebuiltPath := strings.Join(rebuiltSegments, "/")
	decodedPath, err := url.PathUnescape(rebuiltPath)
	if err != nil {
		decodedPath = rebuiltPath
	}

	// Clone the request and set the path directly. The raw path keeps the
	// escaping chosen above, including %XX sequences already in a payload, so
	// they are sent as they are rather than decoded or encoded twice.
	cloned := q.req.Clone(context.Background())
	cloned.Path = decodedPath
	cloned.RawPath = rebuiltPath
	return cloned, nil
}

// Clones current state to a new component
func (q *Path) Clone() Component {
	return &Path{
		value:        q.value.Clone(),
		originalPath: q.originalPath,
		req:          q.req.Clone(context.Background()),
	}
}

// escapePathSegment escapes a value for use as one path segment. It keeps
// "/" so that traversal payloads work, and it keeps valid %XX sequences so
// that payloads which are already encoded are sent unchanged. Any other
// character that is not allowed in a path, including "?", "#" and a "%" that
// does not start an escape, is percent-encoded.
func escapePathSegment(value string) string {
	const upperhex = "0123456789ABCDEF"
	var b strings.Builder
	b.Grow(len(value))
	for i := 0; i < len(value); i++ {
		c := value[i]
		switch {
		case c == '%' && i+2 < len(value) && isHex(value[i+1]) && isHex(value[i+2]):
			b.WriteByte(c)
		case isPathByte(c):
			b.WriteByte(c)
		default:
			b.WriteByte('%')
			b.WriteByte(upperhex[c>>4])
			b.WriteByte(upperhex[c&15])
		}
	}
	return b.String()
}

func isHex(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}

// isPathByte reports whether c can appear unescaped in a URL path
// (RFC 3986 pchar plus "/").
func isPathByte(c byte) bool {
	if ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z') || ('0' <= c && c <= '9') {
		return true
	}
	return strings.IndexByte("-._~!$&'()*+,;=:@/", c) >= 0
}
