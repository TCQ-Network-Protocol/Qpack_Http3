//	Project: QPACK HTTP3
//	Author: Trần Nguyên Hiền (c)
//	Major: Electronic And Communication Engineering
//	Email: trannguyenhien29085@gmail.com
//	Date: 2/3/2026
//	GPL-3.0 Licence
//
// ----------------------------------------------------------------
package qpack

import (
	"bytes"
	"fmt"
	"io"
	"math/rand/v2"
	"testing"

	"github.com/stretchr/testify/require"
)

func randomString(l int) string {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	s := make([]byte, l)
	for i := range s {
		s[i] = charset[rand.IntN(len(charset))]
	}
	return string(s)
}

func getEncoder() (*Encoder, *bytes.Buffer) {
	output := &bytes.Buffer{}
	return NewEncoder(output), output
}

func TestEncodeDecode(t *testing.T) {
	hfs := []HeaderField{
		{Name: "foo", Value: "bar"},
		{Name: "lorem", Value: "ipsum"},
		{Name: randomString(15), Value: randomString(20)},
	}
	encoder, output := getEncoder()
	for _, hf := range hfs {
		require.NoError(t, encoder.WriteField(hf))
	}
	headerFields := decodeAll(t, NewDecoder().Decode(output.Bytes()))
	require.Equal(t, hfs, headerFields)
}

// replace one character by a random character at a random position
func replaceRandomCharacter(s string) string {
	pos := rand.IntN(len(s))
	new := s[:pos]
	for {
		if c := randomString(1); c != string(s[pos]) {
			new += c
			break
		}
	}
	new += s[pos+1:]
	return new
}

func check(t *testing.T, encoded []byte, hf HeaderField) {
	t.Helper()

	headerFields := decodeAll(t, NewDecoder().Decode(encoded))
	require.Len(t, headerFields, 1)
	require.Equal(t, hf, headerFields[0])
}

func TestStaticTableForFieldNamesWithoutValues(t *testing.T) {
	for i := range 10 {
		t.Run(fmt.Sprintf("run %d", i), func(t *testing.T) {
			testStaticTableForFieldNamesWithoutValues(t)
		})
	}
}

func testStaticTableForFieldNamesWithoutValues(t *testing.T) {
	var hf HeaderField
	for {
		if entry := staticTableEntries[rand.IntN(len(staticTableEntries))]; len(entry.Value) == 0 {
			hf = HeaderField{Name: entry.Name}
			break
		}
	}
	encoder, output := getEncoder()
	require.NoError(t, encoder.WriteField(hf))
	encodedLen := output.Len()
	check(t, output.Bytes(), hf)
	encoder, output = getEncoder()
	oldName := hf.Name
	hf.Name = replaceRandomCharacter(hf.Name)
	require.NoError(t, encoder.WriteField(hf))
	t.Logf("Encoding field name:\n\t%s: %d bytes\n\t%s: %d bytes\n", oldName, encodedLen, hf.Name, output.Len())
	require.Greater(t, output.Len(), encodedLen)
}

func TestStaticTableForFieldNamesWithCustomValues(t *testing.T) {
	for i := range 10 {
		t.Run(fmt.Sprintf("run %d", i), func(t *testing.T) {
			testStaticTableForFieldNamesWithCustomValues(t)
		})
	}
}

func testStaticTableForFieldNamesWithCustomValues(t *testing.T) {
	var hf HeaderField
	for {
		if entry := staticTableEntries[rand.IntN(len(staticTableEntries))]; len(entry.Value) == 0 {
			hf = HeaderField{
				Name:  entry.Name,
				Value: randomString(5),
			}
			break
		}
	}
	encoder, output := getEncoder()
	require.NoError(t, encoder.WriteField(hf))
	encodedLen := output.Len()
	check(t, output.Bytes(), hf)
	encoder, output = getEncoder()
	oldName := hf.Name
	hf.Name = replaceRandomCharacter(hf.Name)
	require.NoError(t, encoder.WriteField(hf))
	t.Logf("Encoding field name:\n\t%s: %d bytes\n\t%s: %d bytes", oldName, encodedLen, hf.Name, output.Len())
	require.Greater(t, output.Len(), encodedLen)
}

func TestStaticTableForFieldNamesWithValues(t *testing.T) {
	for i := range 10 {
		t.Run(fmt.Sprintf("run %d", i), func(t *testing.T) {
			testStaticTableForFieldNamesWithValues(t)
		})
	}
}

func testStaticTableForFieldNamesWithValues(t *testing.T) {
	var hf HeaderField
	for {
		// Only use values with at least 2 characters.
		// This makes sure that Huffman encoding doesn't compress them as much as encoding it using the static table would.
		if entry := staticTableEntries[rand.IntN(len(staticTableEntries))]; len(entry.Value) > 1 {
			hf = HeaderField{
				Name:  entry.Name,
				Value: randomString(20),
			}
			break
		}
	}
	encoder, output := getEncoder()
	require.NoError(t, encoder.WriteField(hf))
	encodedLen := output.Len()
	check(t, output.Bytes(), hf)
	encoder, output = getEncoder()
	oldName := hf.Name
	hf.Name = replaceRandomCharacter(hf.Name)
	require.NoError(t, encoder.WriteField(hf))
	t.Logf("Encoding field name:\n\t%s: %d bytes\n\t%s: %d bytes", oldName, encodedLen, hf.Name, output.Len())
	require.Greater(t, output.Len(), encodedLen)
}

func TestStaticTableForFieldValues(t *testing.T) {
	for i := range 10 {
		t.Run(fmt.Sprintf("run %d", i), func(t *testing.T) {
			testStaticTableForFieldValues(t)
		})
	}
}

func testStaticTableForFieldValues(t *testing.T) {
	var hf HeaderField
	for {
		// Only use values with at least 2 characters.
		// This makes sure that Huffman encoding doesn't compress them as much as encoding it using the static table would.
		if entry := staticTableEntries[rand.IntN(len(staticTableEntries))]; len(entry.Value) > 1 {
			hf = HeaderField{
				Name:  entry.Name,
				Value: entry.Value,
			}
			break
		}
	}
	encoder, output := getEncoder()
	require.NoError(t, encoder.WriteField(hf))
	encodedLen := output.Len()
	check(t, output.Bytes(), hf)
	encoder, output = getEncoder()
	oldValue := hf.Value
	hf.Value = replaceRandomCharacter(hf.Value)
	require.NoError(t, encoder.WriteField(hf))
	t.Logf(
		"Encoding field value:\n\t%s: %s -> %d bytes\n\t%s: %s -> %d bytes",
		hf.Name, oldValue, encodedLen,
		hf.Name, hf.Value, output.Len(),
	)
	require.Greater(t, output.Len(), encodedLen)
}

func BenchmarkRoundTrip(b *testing.B) {
	b.Run("typical HTTP request", func(b *testing.B) {
		fields := []HeaderField{
			{Name: ":method", Value: "GET"},
			{Name: ":scheme", Value: "https"},
			{Name: ":path", Value: "/"},
			{Name: "accept", Value: "*/*"},
			{Name: "accept-encoding", Value: "gzip, deflate, br"},
			{Name: "user-agent", Value: "benchmark-client/1.0"},
		}
		benchmarkRoundTrip(b, fields)
	})

	b.Run("typical HTTP response", func(b *testing.B) {
		fields := []HeaderField{
			{Name: ":status", Value: "200"},
			{Name: "content-type", Value: "application/json"},
			{Name: "content-length", Value: "1234"},
			{Name: "cache-control", Value: "no-cache"},
			{Name: "vary", Value: "accept-encoding"},
			{Name: "server", Value: "qpack-test-server"},
			{Name: "x-request-id", Value: "req-abcdef-123456"},
		}
		benchmarkRoundTrip(b, fields)
	})
}

func benchmarkRoundTrip(b *testing.B, fields []HeaderField) {
	b.ReportAllocs()

	output := &bytes.Buffer{}
	encoder := NewEncoder(output)
	decoder := NewDecoder()
	for b.Loop() {
		output.Reset()
		_ = encoder.Close()
		for _, hf := range fields {
			if err := encoder.WriteField(hf); err != nil {
				b.Fatalf("encode error: %v", err)
			}
		}
		decodeFn := decoder.Decode(output.Bytes())
		count := 0
		for {
			_, err := decodeFn()
			if err == io.EOF {
				break
			}
			if err != nil {
				b.Fatalf("decode error: %v", err)
			}
			count++
		}
		if count != len(fields) {
			b.Fatalf("expected %d fields, got %d", len(fields), count)
		}
	}
}
