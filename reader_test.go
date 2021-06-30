package bass_test

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vito/bass"
)

type ReaderExample struct {
	Source string
	Result bass.Value
}

func TestReader(t *testing.T) {
	for _, example := range []ReaderExample{
		{
			Source: "_",
			Result: bass.Ignore{},
		},
		{
			Source: "null",
			Result: bass.Null{},
		},
		{
			Source: "false",
			Result: bass.Bool(false),
		},
		{
			Source: "true",
			Result: bass.Bool(true),
		},
		{
			Source: "42",
			Result: bass.Int(42),
		},

		{
			Source: "hello",
			Result: bass.Symbol("hello"),
		},

		{
			Source: `"hello world"`,
			Result: bass.String("hello world"),
		},

		{
			Source: `"hello \"\n\\\t\a\f\r\b\v"`,
			Result: bass.String("hello \"\n\\\t\a\a\r\b\v"),
		},

		{
			Source: `[]`,
			Result: bass.Empty{},
		},
		{
			Source: `[1 true "three"]`,
			Result: bass.InertPair{
				A: bass.Int(1),
				D: bass.InertPair{
					A: bass.Bool(true),
					D: bass.InertPair{
						A: bass.String("three"),
						D: bass.Empty{},
					},
				},
			},
		},

		{
			Source: `()`,
			Result: bass.Empty{},
		},
		{
			Source: `(foo . bar)`,
			Result: bass.Pair{
				A: bass.Symbol("foo"),
				D: bass.Symbol("bar"),
			},
		},
		{
			Source: `(foo 1 . bar)`,
			Result: bass.Pair{
				A: bass.Symbol("foo"),
				D: bass.Pair{
					A: bass.Int(1),
					D: bass.Symbol("bar"),
				},
			},
		},
		{
			Source: `(foo 1 true "three")`,
			Result: bass.Pair{
				A: bass.Symbol("foo"),
				D: bass.NewList(
					bass.Int(1),
					bass.Bool(true),
					bass.String("three"),
				),
			},
		},
		{
			Source: `(foo 1 (two "three"))`,
			Result: bass.Pair{
				A: bass.Symbol("foo"),
				D: bass.NewList(
					bass.Int(1),
					bass.Pair{
						A: bass.Symbol("two"),
						D: bass.NewList(bass.String("three")),
					},
				),
			},
		},

		// TODO: add tests covering syntax that Bass does *not* support:
		//
		// * syntax-quote
		// * unquote
		//
		// these tests should be a little defensive because we rely on
		// spy16/slurp's Reader, which has a few default macros - a PR upstream to
		// remove these and make them options would be nice.
	} {
		example.Run(t)
	}
}

func TestReaderComments(t *testing.T) {
	for _, example := range []ReaderExample{
		{
			Source: `; hello!
_`,
			Result: bass.Annotated{
				Comment: "hello!",
				Value:   bass.Ignore{},
			},
		},
		{
			Source: `;;; hello!
_`,
			Result: bass.Annotated{
				Comment: "hello!",
				Value:   bass.Ignore{},
			},
		},
		{
			Source: `;; ; hello!
_`,
			Result: bass.Annotated{
				Comment: "; hello!",
				Value:   bass.Ignore{},
			},
		},
		{
			Source: `;;;   hello!
_`,
			Result: bass.Annotated{
				Comment: "hello!",
				Value:   bass.Ignore{},
			},
		},
		{
			Source: `; hello!
; multiline!
_`,
			Result: bass.Annotated{
				Comment: "hello! multiline!",
				Value:   bass.Ignore{},
			},
		},
		{
			Source: `; hello!
;
; multi paragraph!
_`,
			Result: bass.Annotated{
				Comment: "hello!\n\nmulti paragraph!",
				Value:   bass.Ignore{},
			},
		},
		{
			Source: `123 ; hello!`,
			Result: bass.Annotated{
				Comment: "hello!",
				Value:   bass.Int(123),
			},
		},
		{
			Source: `; outer
[123 ; hello!
 456 ; goodbye!

 ; inner
 foo
]
`,
			Result: bass.Annotated{
				Comment: "outer",
				Value: bass.NewInertList(
					bass.Annotated{
						Comment: "hello!",
						Value:   bass.Int(123),
					},
					bass.Annotated{
						Comment: "goodbye!",
						Value:   bass.Int(456),
					},
					bass.Annotated{
						Comment: "inner",
						Value:   bass.Symbol("foo"),
					},
				),
			},
		},
	} {
		example.Run(t)
	}
}

func (example ReaderExample) Run(t *testing.T) {
	t.Run(example.Source, func(t *testing.T) {
		reader := bass.NewReader(bytes.NewBufferString(example.Source))

		form, err := reader.Next()
		require.NoError(t, err)
		require.True(t, form.Equal(example.Result), "%s != %s", form, example.Result)
	})
}
