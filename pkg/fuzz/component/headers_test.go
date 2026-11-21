package component

import (
	"net/http"
	"testing"

	"github.com/projectdiscovery/retryablehttp-go"
	"github.com/stretchr/testify/require"
)

func TestHeaderComponent(t *testing.T) {
	req, err := retryablehttp.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	req.Header.Set("User-Agent", "test-agent")

	header := NewHeader()
	_, err = header.Parse(req)
	if err != nil {
		t.Fatal(err)
	}

	var keys []string
	var values []string
	_ = header.Iterate(func(key string, value interface{}) error {
		keys = append(keys, key)
		switch v := value.(type) {
		case string:
			values = append(values, v)
		case []string:
			values = append(values, v...)
		}
		return nil
	})

	require.Equal(t, []string{"User-Agent"}, keys, "unexpected keys")
	require.Equal(t, []string{"test-agent"}, values, "unexpected values")

	err = header.SetValue("User-Agent", "new-agent")
	if err != nil {
		t.Fatal(err)
	}

	rebuilt, err := header.Rebuild()
	if err != nil {
		t.Fatal(err)
	}

	require.Equal(t, "new-agent", rebuilt.Header.Get("User-Agent"), "unexpected header value")
}

func TestHeaderComponentMultiValueRebuildDoesNotDuplicate(t *testing.T) {
	req, err := retryablehttp.NewRequest(http.MethodGet, "https://example.com", nil)
	require.NoError(t, err)
	req.Header.Add("X-Forwarded-For", "10.0.0.1")
	req.Header.Add("X-Forwarded-For", "10.0.0.2")

	header := NewHeader()
	_, err = header.Parse(req)
	require.NoError(t, err)

	rebuilt, err := header.Rebuild()
	require.NoError(t, err)
	require.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, rebuilt.Header.Values("X-Forwarded-For"), "an unchanged header should be sent as it was")

	require.NoError(t, header.SetValue("X-Forwarded-For", "PAYLOAD"))
	rebuilt, err = header.Rebuild()
	require.NoError(t, err)
	values := rebuilt.Header.Values("X-Forwarded-For")
	require.Len(t, values, 2, "the fuzzed header should keep its number of values: %v", values)
	require.Contains(t, values, "PAYLOAD")
	require.Equal(t, []string{"10.0.0.1", "10.0.0.2"}, req.Header.Values("X-Forwarded-For"), "the original request should not change")
}
