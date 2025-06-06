package header

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/platform-inf/go-ratchet/errlist"
	"github.com/platform-inf/go-ratchet/keys"
)

func TestHeaderSuccessEncodeAndDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		header Header
		bytes  []byte
	}{
		{
			"full header",
			Header{
				PublicKey:                         keys.Public{Bytes: []byte{0x01, 0x02, 0x03, 0x04, 0x05}},
				PreviousSendingChainMessagesCount: 123,
				MessageNumber:                     321,
			},
			[]byte{
				0x41, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x7b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x01, 0x02, 0x03, 0x04, 0x05,
			},
		},
		{
			"zero header",
			Header{},
			[]byte{
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			bytes := test.header.Encode()
			if !slices.Equal(bytes, test.bytes) {
				t.Fatalf("%+v.Encode(): expected %v but got %v", test.header, test.bytes, bytes)
			}

			header, err := Decode(bytes)
			if err != nil {
				t.Fatalf("Decode(%v): expected no error but got %v", bytes, err)
			}

			if !reflect.DeepEqual(header, test.header) {
				t.Fatalf("Decode(%v): expected %+v but got %+v", bytes, test.header, header)
			}
		})
	}
}

func TestHeaderDecode(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name          string
		bytes         []byte
		errorCategory error
		errorString   string
	}{
		{
			"not enough bytes",
			[]byte{
				0x12, 0x00, 0x00, 0x00, 0x22, 0x00, 0x00, 0x0F,
				0x55, 0x00, 0x00, 0x00, 0x77, 0x00, 0x0B,
			},
			errlist.ErrInvalidValue,
			"invalid value: not enough bytes",
		},
		{"nil bytes slice", nil, errlist.ErrInvalidValue, "invalid value: not enough bytes"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			if _, err := Decode(test.bytes); !errors.Is(err, test.errorCategory) || err.Error() != test.errorString {
				t.Fatalf("Decode(%v) expected error %q but got %v", test.bytes, test.errorString, err)
			}
		})
	}
}
