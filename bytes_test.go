package pflag

import (
	"encoding/base64"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBytesHex(t *testing.T) {
	newFlag := func(bytesHex *[]byte) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.BytesHexVar(bytesHex, "bytes", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "Some bytes in HEX")
		f.BytesHexVarP(bytesHex, "bytes2", "B", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "Some bytes in HEX")
		return f
	}

	testCases := []struct {
		input    string
		success  bool
		expected string
	}{
		/// Positive cases
		{"", true, ""}, // Is empty string OK ?
		{"01", true, "01"},
		{"0101", true, "0101"},
		{"1234567890abcdef", true, "1234567890ABCDEF"},
		{"1234567890ABCDEF", true, "1234567890ABCDEF"},

		// Negative cases
		{"0", false, ""},   // Short string
		{"000", false, ""}, /// Odd-length string
		{"qq", false, ""},  /// non-hex character
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull

	for i := range testCases {
		var bytesHex []byte
		f := newFlag(&bytesHex)

		tc := &testCases[i]

		// --bytes
		args := []string{
			fmt.Sprintf("--bytes=%s", tc.input),
			fmt.Sprintf("-B  %s", tc.input),
			fmt.Sprintf("--bytes2=%s", tc.input),
		}

		for _, arg := range args {
			err := f.Parse([]string{arg})

			if !tc.success {
				require.Errorf(t, err,
					"expected failure while processing %q", tc.input,
				)

				continue
			}

			require.NoErrorf(t, err, "expected success, got %q", err)

			bytesHex, err := f.GetBytesHex("bytes")
			require.NoErrorf(t, err,
				"got error trying to fetch the 'bytes' flag: %v", err,
			)

			require.Equalf(t, tc.expected, fmt.Sprintf("%X", bytesHex),
				"expected %q, got '%X'", tc.expected, bytesHex,
			)

		}
	}
}

func TestBytesBase64(t *testing.T) {
	newFlag := func(bytesBase64 *[]byte) *FlagSet {
		f := NewFlagSet("test", ContinueOnError)
		f.BytesBase64Var(bytesBase64, "bytes", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "Some bytes in Base64")
		f.BytesBase64VarP(bytesBase64, "bytes2", "B", []byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0}, "Some bytes in Base64")
		return f
	}

	testCases := []struct {
		input    string
		success  bool
		expected string
	}{
		/// Positive cases
		{"", true, ""}, // Is empty string OK ?
		{"AQ==", true, "AQ=="},

		// Negative cases
		{"AQ", false, ""}, // Padding removed
		{"ï", false, ""},  // non-base64 characters
	}

	devnull, _ := os.Open(os.DevNull)
	os.Stderr = devnull

	for i := range testCases {
		var bytesBase64 []byte
		f := newFlag(&bytesBase64)
		tc := &testCases[i]

		args := []string{
			fmt.Sprintf("--bytes=%s", tc.input),
			fmt.Sprintf("-B  %s", tc.input),
			fmt.Sprintf("--bytes2=%s", tc.input),
		}

		for _, arg := range args {
			err := f.Parse([]string{arg})
			if !tc.success {
				require.Errorf(t, err,
					"expected failure while processing %q", tc.input,
				)

				continue
			}

			require.NoErrorf(t, err, "expected success, got %q", err)

			bytesBase64, err := f.GetBytesBase64("bytes")
			require.NoErrorf(t, err,
				"got error trying to fetch the 'bytes' flag: %v", err,
			)
			require.Equalf(t, tc.expected, base64.StdEncoding.EncodeToString(bytesBase64),
				"expected %q, got '%X'", tc.expected, bytesBase64,
			)
		}
	}
}
