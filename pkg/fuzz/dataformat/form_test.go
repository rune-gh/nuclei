package dataformat

import (
	"net/url"
	"reflect"
	"testing"
)

func TestFormDecodeEncode_DuplicateParameters(t *testing.T) {
	form := NewForm()
	decoded, err := form.Decode("foo=a&foo=b&foo=c")
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := form.Encode(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if encoded != "foo=a&foo=b&foo=c" {
		t.Fatalf("unexpected form encoding: %q", encoded)
	}
}

func TestFormDecodeEncode_DoesNotMergePrefixParameterNames(t *testing.T) {
	form := NewForm()
	decoded, err := form.Decode("foo=a&foobar=b&foobar=c")
	if err != nil {
		t.Fatal(err)
	}

	encoded, err := form.Encode(decoded)
	if err != nil {
		t.Fatal(err)
	}
	if encoded != "foo=a&foobar=b&foobar=c" {
		t.Fatalf("unexpected form encoding: %q", encoded)
	}
}

func TestFormDecodeEncode_KeepsParametersNamedLikeNormalizedKeys(t *testing.T) {
	form := NewForm()
	for _, input := range []string{
		"address=home&address_1=apt5",
		"q=1&q_2=z",
		"item_1=x&item=y",
		"foo=a&foo=b&foo_7=c",
	} {
		decoded, err := form.Decode(input)
		if err != nil {
			t.Fatal(err)
		}

		encoded, err := form.Encode(decoded)
		if err != nil {
			t.Fatal(err)
		}
		assertSameFormValues(t, input, encoded)
	}
}

func TestFormEncode_DuplicateParameterAliasAfterClone(t *testing.T) {
	form := NewForm()
	decoded, err := form.Decode("foo=a&foo=b&foo_2=c")
	if err != nil {
		t.Fatal(err)
	}

	cloned := decoded.Clone()
	cloned.Set("foo_1", "x")
	encoded, err := form.Encode(cloned)
	if err != nil {
		t.Fatal(err)
	}
	assertSameFormValues(t, "foo=x&foo=b&foo_2=c", encoded)
}

func assertSameFormValues(t *testing.T, expected, actual string) {
	t.Helper()
	expectedValues, err := url.ParseQuery(expected)
	if err != nil {
		t.Fatal(err)
	}
	actualValues, err := url.ParseQuery(actual)
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(expectedValues, actualValues) {
		t.Fatalf("unexpected form encoding: got %q, want the parameters of %q", actual, expected)
	}
}
