package bass_test

// import (
// 	"testing"

// 	"github.com/stretchr/testify/require"
// 	"github.com/vito/bass"
// )

// func TestPathDecode(t *testing.T) {
// 	var foo string
// 	err := bass.Path("foo").Decode(&foo)
// 	require.NoError(t, err)
// 	require.Equal(t, foo, "foo")

// 	err = bass.Path("bar").Decode(&foo)
// 	require.NoError(t, err)
// 	require.Equal(t, foo, "bar")

// 	var str bass.Path
// 	err = bass.Path("bar").Decode(&str)
// 	require.NoError(t, err)
// 	require.Equal(t, str, bass.Path("bar"))
// }

// func TestPathEqual(t *testing.T) {
// 	require.True(t, bass.Path("hello").Equal(bass.Path("hello")))
// 	require.True(t, bass.Path("").Equal(bass.Path("")))
// 	require.False(t, bass.Path("hello").Equal(bass.Path("")))
// 	require.True(t, bass.Path("hello").Equal(wrappedValue{bass.Path("hello")}))
// 	require.False(t, bass.Path("hello").Equal(wrappedValue{bass.Path("")}))
// }
