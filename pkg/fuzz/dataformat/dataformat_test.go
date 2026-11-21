package dataformat

import (
	"testing"
)

func TestDataformatDecodeEncode_JSON(t *testing.T) {
	obj := `{"foo":"bar"}`

	decoded, err := Decode(obj)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.DataFormat != "json" {
		t.Fatal("unexpected data format")
	}
	if decoded.Data.Get("foo") != "bar" {
		t.Fatal("unexpected data")
	}

	encoded, err := Encode(decoded.Data, decoded.DataFormat)
	if err != nil {
		t.Fatal(err)
	}
	if encoded != obj {
		t.Fatal("unexpected data")
	}
}

func TestDataformatDecodeEncode_XML(t *testing.T) {
	obj := `<foo attr="baz">bar</foo>`

	decoded, err := Decode(obj)
	if err != nil {
		t.Fatal(err)
	}
	if decoded.DataFormat != "xml" {
		t.Fatal("unexpected data format")
	}
	fooValue := decoded.Data.Get("foo")
	if fooValue == nil {
		t.Fatal("key 'foo' not found")
	}
	fooMap, ok := fooValue.(map[string]interface{})
	if !ok {
		t.Fatal("type assertion to map[string]interface{} failed")
	}
	if fooMap["#text"] != "bar" {
		t.Fatal("unexpected data for '#text'")
	}
	if fooMap["-attr"] != "baz" {
		t.Fatal("unexpected data for '-attr'")
	}

	encoded, err := Encode(decoded.Data, decoded.DataFormat)
	if err != nil {
		t.Fatal(err)
	}
	if encoded != obj {
		t.Fatal("unexpected data")
	}
}

func TestDataformatDecodeEncode_JSONKeepsLargeIntegers(t *testing.T) {
	obj := `{"id":9007199254740993,"price":1.10,"snowflake":1795243533176012345}`

	decoded, err := Decode(obj)
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := Encode(decoded.Data, decoded.DataFormat)
	if err != nil {
		t.Fatal(err)
	}
	if encoded != obj {
		t.Fatalf("unexpected data: %s", encoded)
	}
}

func TestDataformatDecode_JSONRejectsTrailingData(t *testing.T) {
	if _, err := NewJSON().Decode(`{"foo":"bar"} {"foo":"baz"}`); err == nil {
		t.Fatal("expected an error for data after the JSON object")
	}
}
