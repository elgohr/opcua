// Copyright 2020 Converter Systems LLC. All rights reserved.

package opcua_test

import (
	"bytes"
	"reflect"
	"testing"
	"time"

	ua "github.com/awcullen/opcua"
	uuid "github.com/google/uuid"
	"github.com/pascaldekloe/goe/verify"
)

func TestBoolean(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "Boolean",
			In:   true,
			Bytes: []byte{
				0x01,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestInt32(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "Int32",
			In:   int32(1_000_000_000),
			Bytes: []byte{
				0x00, 0xCA, 0x9A, 0x3B,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestFloat(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "Single",
			In:   float32(-6.5),
			Bytes: []byte{
				0x00, 0x00, 0xD0, 0xC0,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestString(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "String",
			In:   "水Boy",
			Bytes: []byte{
				0x06, 0x00, 0x00, 0x00, 0xE6, 0xB0, 0xB4, 0x42, 0x6F, 0x79,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestTime(t *testing.T) {
	t1, _ := time.Parse(time.RFC3339, "2020-07-04T12:00:00Z")
	cases := []encoderTestCase{
		{
			Name:  "DateTime",
			In:    t1,
			Bytes: []byte{0x00, 0xa0, 0xa5, 0xa4, 0xfa, 0x51, 0xd6, 0x01},
		},
	}
	runEncoderTest(t, cases)
}

func TestGUID(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "spec",
			In:   uuid.MustParse("72962B91-FA75-4AE6-8D28-B404DC7DAF63"),
			Bytes: []byte{
				// data1 (inverse order)
				0x91, 0x2b, 0x96, 0x72,
				// data2 (inverse order)
				0x75, 0xfa,
				// data3 (inverse order)
				0xe6, 0x4a,
				// data4 (same order)
				0x8d, 0x28, 0xb4, 0x04, 0xdc, 0x7d, 0xaf, 0x63,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestNodeID(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "TwoByte",
			In:   ua.NewNodeIDNumeric(0, 255),
			Bytes: []byte{
				// mask
				0x00,
				// id
				0xff,
			},
		},
		{
			Name: "FourByte",
			In:   ua.NewNodeIDNumeric(2, 65535),
			Bytes: []byte{
				// mask
				0x01,
				// namespace
				0x02,
				// id
				0xff, 0xff,
			},
		},
		{
			Name: "Numeric",
			In:   ua.NewNodeIDNumeric(10, 4294967295),
			Bytes: []byte{
				// mask
				0x02,
				// namespace
				0x0a, 0x00,
				// id
				0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			Name: "String",
			In:   ua.NewNodeIDString(2, "bar"),
			Bytes: []byte{
				// mask
				0x03,
				// namespace
				0x02, 0x00,
				// value
				0x03, 0x00, 0x00, 0x00, // len
				0x62, 0x61, 0x72, // char
			},
		},
		{
			Name: "Guid",
			In:   ua.NewNodeIDGUID(2, uuid.MustParse("AAAABBBB-CCDD-EEFF-0102-0123456789AB")),
			Bytes: []byte{
				// mask
				0x04,
				// namespace
				0x02, 0x00,
				// value
				// data1 (inverse order)
				0xbb, 0xbb, 0xaa, 0xaa,
				// data2 (inverse order)
				0xdd, 0xcc,
				// data3 (inverse order)
				0xff, 0xee,
				// data4 (same order)
				0x01, 0x02, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab,
			},
		},
		{
			Name: "Opaque",
			In:   ua.NewNodeIDOpaque(2, ua.ByteString("\x00\x10\x20\x30\x40\x50\x60\x70")),
			Bytes: []byte{
				// mask
				0x05,
				// namespace
				0x02, 0x00,
				// value
				0x08, 0x00, 0x00, 0x00, // len
				0x00, 0x10, 0x20, 0x30, 0x40, 0x50, 0x60, 0x70, // bytes
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestQualifiedName(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "has-both",
			In:   ua.QualifiedName{NamespaceIndex: 2, Name: "bar"},
			Bytes: []byte{
				0x02, 0x00,
				// name: "bar"
				0x03, 0x00, 0x00, 0x00, 0x62, 0x61, 0x72,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestLocalizedText(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name:  "nothing",
			In:    ua.LocalizedText{},
			Bytes: []byte{0x00},
		},
		{
			Name: "has-locale",
			In:   ua.LocalizedText{Locale: "foo"},
			Bytes: []byte{
				0x01,
				0x03, 0x00, 0x00, 0x00, 0x66, 0x6f, 0x6f,
			},
		},
		{
			Name: "has-text",
			In:   ua.LocalizedText{Text: "bar"},
			Bytes: []byte{
				0x02,
				0x03, 0x00, 0x00, 0x00, 0x62, 0x61, 0x72,
			},
		},
		{
			Name: "has-both",
			In:   ua.LocalizedText{Text: "bar", Locale: "foo"},
			Bytes: []byte{
				0x03,
				0x03, 0x00, 0x00, 0x00, 0x66, 0x6f, 0x6f,
				// second String: "bar"
				0x03, 0x00, 0x00, 0x00, 0x62, 0x61, 0x72,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestDataValue(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "value only",
			In:   ua.NewDataValueFloat(float32(2.50025), 0, time.Time{}, 0, time.Time{}, 0),
			Bytes: []byte{
				// EncodingMask
				0x01,
				// Value
				0x0a,                   // type
				0x19, 0x04, 0x20, 0x40, // value
			},
		},
		{
			Name: "value, source timestamp, server timestamp",
			In: ua.NewDataValueFloat(float32(2.50017), 0,
				time.Date(2018, time.September, 17, 14, 28, 29, 112000000, time.UTC), 0,
				time.Date(2018, time.September, 17, 14, 28, 29, 112000000, time.UTC), 0),
			Bytes: []byte{
				// EncodingMask
				0x0d,
				// Value
				0x0a,                   // type
				0xc9, 0x02, 0x20, 0x40, // value
				// SourceTimestamp
				0x80, 0x3b, 0xe8, 0xb3, 0x92, 0x4e, 0xd4, 0x01,
				// SeverTimestamp
				0x80, 0x3b, 0xe8, 0xb3, 0x92, 0x4e, 0xd4, 0x01,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestEnum(t *testing.T) {
	cases := []encoderTestCase{
		{
			Name: "Enum",
			In:   ua.MessageSecurityModeSignAndEncrypt,
			Bytes: []byte{
				// int32
				0x03, 0x00, 0x00, 0x00,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestStruct(t *testing.T) {
	t0, _ := time.Parse(time.RFC3339, "1601-01-01T12:00:00Z")
	nodesToRead := make([]*ua.ReadValueID, 1)
	for index := 0; index < len(nodesToRead); index++ {
		nodesToRead[index] = &ua.ReadValueID{
			AttributeID: ua.AttributeIDValue,
			NodeID:      ua.NewNodeIDNumeric(0, 255),
		}
	}
	cases := []encoderTestCase{
		{
			Name: "ReadRequest",
			In:   &ua.ReadRequest{RequestHeader: ua.RequestHeader{Timestamp: t0}, NodesToRead: nodesToRead},
			Bytes: []byte{
				0x00, 0x00, 0x00, 0xe0, 0x34, 0x95, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, // max age
				0x00, 0x00, 0x00, 0x00, // timestamps
				0x01, 0x00, 0x00, 0x00, // len
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
			},
		},
		{
			Name: "CreateSessionRequest",
			In:   &ua.CreateSessionRequest{RequestHeader: ua.RequestHeader{Timestamp: t0}, ClientDescription: &ua.ApplicationDescription{}},
			Bytes: []byte{
				0x00, 0x00, 0x00, 0xe0, 0x34, 0x95, 0x64, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
				0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestSlice(t *testing.T) {
	nodesToRead := make([]*ua.ReadValueID, 10)
	for index := 0; index < len(nodesToRead); index++ {
		nodesToRead[index] = &ua.ReadValueID{
			AttributeID: ua.AttributeIDValue,
			NodeID:      ua.NewNodeIDNumeric(0, 255),
		}
	}
	cases := []encoderTestCase{
		{
			Name: "ReadValueID",
			In:   nodesToRead,
			Bytes: []byte{
				// int32
				0x0a, 0x00, 0x00, 0x00,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
				0x00, 0xff, 0x0d, 0x00, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff, 0x00, 0x00, 0xff, 0xff, 0xff, 0xff,
			},
		},
	}
	runEncoderTest(t, cases)
}

func TestSliceVariant(t *testing.T) {
	variants := []*ua.Variant{
		ua.NewVariantString("foo"),
		ua.NewVariantUInt16(255),
	}
	cases := []encoderTestCase{
		{
			Name: "ReadValueID",
			In:   variants,
			Bytes: []byte{
				0x02, 0x00, 0x00, 0x00, // len
				0x0c, 0x03, 0x00, 0x00, 0x00, 0x66, 0x6f, 0x6f, // foo
				0x05, 0xff, 0x00, // 255
			},
		},
	}
	runEncoderTest(t, cases)
}

type encoderTestCase struct {
	Name  string
	In    interface{}
	Bytes []byte
}

func runEncoderTest(t *testing.T, cases []encoderTestCase) {
	t.Helper()
	for _, c := range cases {
		t.Run(c.Name, func(t *testing.T) {
			t.Run("encode", func(t *testing.T) {
				bs := make([]byte, 0, 200)
				buf := bytes.NewBuffer(bs)
				enc := ua.NewBinaryEncoder(buf, ua.NewEncodingContext())
				if err := enc.Encode(c.In); err != nil {
					t.Fatal(err)
				}
				// t.Logf("% #x\n", buf.Bytes())
				// t.Logf("% #x\n", c.Bytes)
				verify.Values(t, "", buf.Bytes(), c.Bytes)
			})
			t.Run("decode", func(t *testing.T) {
				buf := bytes.NewBuffer(c.Bytes)
				dec := ua.NewBinaryDecoder(buf, ua.NewEncodingContext())
				out := reflect.New(reflect.TypeOf(c.In)).Interface()
				if err := dec.Decode(out); err != nil {
					t.Fatal(err)
				}
				out = reflect.ValueOf(out).Elem().Interface()
				// t.Logf("%+v\n", c.In)
				// t.Logf("%+v\n", out)
				verify.Values(t, "", out, c.In)
			})
		})
	}
}
