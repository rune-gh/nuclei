package dataformat

import (
	"errors"
	"io"
	"strings"

	"github.com/projectdiscovery/nuclei/v3/pkg/utils/json"
)

// JSON is a JSON encoder
//
// For now JSON only supports objects as the root data type
// and not arrays
//
// TODO: Support arrays + other JSON oddities by
// adding more attributes to the map[string]interface{}
type JSON struct{}

var (
	_ DataFormat = &JSON{}
)

// NewJSON returns a new JSON encoder
func NewJSON() *JSON {
	return &JSON{}
}

// IsType returns true if the data is JSON encoded
func (j *JSON) IsType(data string) bool {
	return strings.HasPrefix(data, "{") && strings.HasSuffix(data, "}")
}

// Encode encodes the data into JSON format
func (j *JSON) Encode(data KV) (string, error) {
	encoded, err := json.Marshal(data.Map)
	return string(encoded), err
}

// Decode decodes the data from JSON format
func (j *JSON) Decode(data string) (KV, error) {
	var decoded map[string]interface{}
	// Numbers are kept as json.Number so that the fields that are not fuzzed
	// are sent back as they came. A float64 would round integers above 2^53,
	// such as 64-bit IDs.
	decoder := json.NewDecoder(strings.NewReader(data))
	decoder.UseNumber()
	if err := decoder.Decode(&decoded); err != nil {
		return KVMap(decoded), err
	}
	if err := decoder.Decode(new(interface{})); !errors.Is(err, io.EOF) {
		return KVMap(decoded), errors.New("invalid json: unexpected data after the top-level value")
	}
	return KVMap(decoded), nil
}

// Name returns the name of the encoder
func (j *JSON) Name() string {
	return JSONDataFormat
}
