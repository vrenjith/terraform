package variables

import (
	"flag"
	"fmt"
	"reflect"
	"testing"

	"github.com/davecgh/go-spew/spew"
)

func TestFlag_impl(t *testing.T) {
	var _ flag.Value = new(Flag)
}

func TestFlag(t *testing.T) {
	cases := []struct {
		Input  interface{}
		Output map[string]interface{}
		Error  bool
	}{
		{
			"key=value",
			map[string]interface{}{"key": "value"},
			false,
		},

		{
			"key=",
			map[string]interface{}{"key": ""},
			false,
		},

		{
			"key=foo=bar",
			map[string]interface{}{"key": "foo=bar"},
			false,
		},

		{
			"key=false",
			map[string]interface{}{"key": "false"},
			false,
		},

		{
			"map.key=foo",
			map[string]interface{}{"map.key": "foo"},
			false,
		},

		{
			"key",
			nil,
			true,
		},

		{
			`key=["hello", "world"]`,
			map[string]interface{}{"key": []interface{}{"hello", "world"}},
			false,
		},

		{
			`key={"hello" = "world", "foo" = "bar"}`,
			map[string]interface{}{
				"key": map[string]interface{}{
					"hello": "world",
					"foo":   "bar",
				},
			},
			false,
		},

		{
			`key={"hello" = "world", "foo" = "bar"}\nkey2="invalid"`,
			nil,
			true,
		},

		{
			"key=/path",
			map[string]interface{}{"key": "/path"},
			false,
		},

		{
			"key=1234.dkr.ecr.us-east-1.amazonaws.com/proj:abcdef",
			map[string]interface{}{"key": "1234.dkr.ecr.us-east-1.amazonaws.com/proj:abcdef"},
			false,
		},

		// simple values that can parse as numbers should remain strings
		{
			"key=1",
			map[string]interface{}{
				"key": "1",
			},
			false,
		},
		{
			"key=1.0",
			map[string]interface{}{
				"key": "1.0",
			},
			false,
		},
		{
			"key=0x10",
			map[string]interface{}{
				"key": "0x10",
			},
			false,
		},

		// Test setting multiple times
		{
			[]string{
				"foo=bar",
				"bar=baz",
			},
			map[string]interface{}{
				"foo": "bar",
				"bar": "baz",
			},
			false,
		},

		// Test map merging
		{
			[]string{
				`foo={ foo = "bar" }`,
				`foo={ bar = "baz" }`,
			},
			map[string]interface{}{
				"foo": map[string]interface{}{
					"foo": "bar",
					"bar": "baz",
				},
			},
			false,
		},
	}

	for i, tc := range cases {
		t.Run(fmt.Sprintf("%d-%s", i, tc.Input), func(t *testing.T) {
			var input []string
			switch v := tc.Input.(type) {
			case string:
				input = []string{v}
			case []string:
				input = v
			default:
				t.Fatalf("bad input type: %T", tc.Input)
			}

			f := new(Flag)
			for i, single := range input {
				err := f.Set(single)

				// Only check for expected errors on the final input
				expected := tc.Error && i == len(input)-1
				if err != nil != expected {
					t.Fatalf("bad error. Input: %#v\n\nError: %s", single, err)
				}
			}

			actual := map[string]interface{}(*f)
			if !reflect.DeepEqual(actual, tc.Output) {
				t.Fatalf("bad:\nexpected: %s\n\n     got: %s\n", spew.Sdump(tc.Output), spew.Sdump(actual))
			}
		})
	}
}
