package encoder

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawDecoder_BytePointerReadsAll(t *testing.T) {
	data := strings.Repeat("a", 100000)

	var result []byte
	err := RawDecoderFactory()(strings.NewReader(data))(&result)

	require.NoError(t, err)
	assert.Equal(t, data, string(result))
}

func TestRawDecoder_Writer(t *testing.T) {
	result := &bytes.Buffer{}
	err := RawDecoderFactory()(strings.NewReader("abc"))(result)

	require.NoError(t, err)
	assert.Equal(t, "abc", result.String())
}

func TestRawDecoder_UnsupportedValue(t *testing.T) {
	err := RawDecoderFactory()(strings.NewReader("abc"))(make([]byte, 3))
	assert.Error(t, err)
}

func TestRawEncoder(t *testing.T) {
	s := "abc"
	b := []byte("abc")
	for _, v := range []interface{}{s, &s, b, &b, strings.NewReader("abc")} {
		w := &bytes.Buffer{}
		require.NoError(t, RawEncoderFactory()(w)(v))
		assert.Equal(t, "abc", w.String())
	}
}
