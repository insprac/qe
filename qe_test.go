package qe

import (
	"strings"
	"testing"
)

func assertEquals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatalf("Equal assertion failed: %v != %v", a, b)
	}
}

func assertNilError(t *testing.T, err error) {
	if err != nil {
		t.Fatalf("Error wasn't nil: %v", err)
	}
}

func assertErrorContains(t *testing.T, err error, term string) {
	if err == nil {
		t.Fatalf("Error was nil, but should've contained '%s'", term)
	} else {
		if !strings.Contains(err.Error(), term) {
			t.Fatalf("Error '%s' didn't contain term '%s'", err.Error(), term)
		}
	}
}

func TestMarshal(t *testing.T) {
	tests := map[string]interface{}{
		"str=test": struct {
			Str string `q:"str"`
		}{"test"},

		"bool=true": struct {
			Bool bool `q:"bool"`
		}{true},

		"int=123&float=12.345": struct {
			Int   uint8   `q:"int"`
			Float float32 `q:"float"`
		}{123, 12.345},

		"esc=escaping+test%3Dtrue": struct {
			Esc string `q:"esc"`
		}{"escaping test=true"},

		"list=1%2C2%2C3": struct {
			List []uint8 `q:"list"`
		}{[]uint8{1, 2, 3}},

		"list=true%2Cfalse%2Ctrue": struct {
			List []bool `q:"list"`
		}{[]bool{true, false, true}},

		"list=a%2Cb%2Cc": struct {
			List []string `q:"list"`
		}{[]string{"a", "b", "c"}},

		"": struct {
			String string   `q:"ignored"`
			Int    uint8    `q:"int"`
			Empty  []string `q:"empty"`
		}{"", 0, []string{}},
	}

	for expected, data := range tests {
		t.Run(expected, func(t *testing.T) {
			encoded, err := Marshal(data)
			assertNilError(t, err)
			assertEquals(t, encoded, expected)
		})
	}
}

func TestMarshalRequired(t *testing.T) {
	type Require struct {
		Req []uint8 `q:"req" required:"true"`
	}

	valid, err := Marshal(Require{[]uint8{1, 2}})

	assertNilError(t, err)
	assertEquals(t, valid, "req=1%2C2")

	invalid, err := Marshal(Require{nil})

	assertErrorContains(t, err, "required")
	assertEquals(t, invalid, "")
}
