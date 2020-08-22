package urlquery

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMarshalScalar(t *testing.T) {
	require := require.New(t)
	{
		ret, err := Marshal("test ?with &stuff", "string")
		require.Nil(err)
		require.Equal("string=test+%3Fwith+%26stuff", ret)
	}
	{
		ret, err := Marshal(true, "bool")
		require.Nil(err)
		require.Equal("bool=true", ret)
	}
	{
		ret, err := Marshal(42, "int")
		require.Nil(err)
		require.Equal("int=42", ret)

	}
	{
		ret, err := Marshal(3.14, "float")
		require.Nil(err)
		require.Equal("float=3.14", ret)
	}
}

func TestMarshalArray(t *testing.T) {
	require := require.New(t)

	arr := []int{5, 16, 49, 50}
	ret, err := Marshal(arr, "y")
	require.Nil(err)
	require.Equal("y=5&y=16&y=49&y=50", ret)
}

func TestMarshalArrayOfArrays(t *testing.T) {
	require := require.New(t)

	arr := [][]string{}
	a := []string{"this", "is", "a", "test"}
	arr = append(arr, a)
	a = []string{"another", "test", "whatever"}
	arr = append(arr, a)

	ret, err := Marshal(arr, "myArr")
	require.Nil(err)
	require.Equal("myArr.0=this&myArr.0=is&myArr.0=a&myArr.0=test&myArr.1=another&myArr.1=test&myArr.1=whatever", ret)
}

func TestMarshalMap(t *testing.T) {
	require := require.New(t)

	m := map[string]string{
		"name": "John",
		"Age":  "oLD",
	}

	ret, err := Marshal(m, "")
	require.Nil(err)
	is := ret == "name=John&Age=oLD" || ret == "Age=oLD&name=John"
	require.Equal(is, true)
}

func TestMarshalWrongMap(t *testing.T) {
	require := require.New(t)

	m := map[int]string{
		17: "something",
	}

	ret, err := Marshal(m, "")
	require.NotNil(err)
	require.Equal("map key must be string", err.Error())
	require.Equal("", ret)
}

func TestMarshalStruct(t *testing.T) {
	require := require.New(t)

	s := &struct {
		Host     string
		Port     int
		Valid    bool
		notValid string
	}{
		Host:     "localhost",
		Port:     8080,
		Valid:    true,
		notValid: "qwerty",
	}

	ret, err := Marshal(s, "")
	require.Nil(err)
	require.Equal("Host=localhost&Port=8080&Valid=true", ret)
}

func TestMarshalStructCustomKeys(t *testing.T) {
	require := require.New(t)

	s := &struct {
		Host     string `url:"host"`
		Port     int    `url:"port"`
		Valid    bool   `url:"valid"`
		notValid string
	}{
		Host:     "localhost",
		Port:     8080,
		Valid:    true,
		notValid: "qwerty",
	}

	ret, err := Marshal(s, "")
	require.Nil(err)
	require.Equal("host=localhost&port=8080&valid=true", ret)
}

func TestMarshalNestedStruct(t *testing.T) {
	require := require.New(t)

	type rec struct {
		ID    int
		Value string
	}
	s := &struct {
		Host  string
		Port  int
		Valid bool
		Recs  []rec
	}{
		Host:  "localhost",
		Port:  8080,
		Valid: true,
		Recs: []rec{
			rec{5, "some value"},
			rec{17, "$ymbl%s;/"},
		},
	}

	ret, err := Marshal(s, "")
	require.Nil(err)
	require.Equal("Host=localhost&Port=8080&Valid=true&Recs.0.ID=5&Recs.0.Value=some+value&Recs.1.ID=17&Recs.1.Value=%24ymbl%25s%3B%2F", ret)
}

func TestMarshalNil(t *testing.T) {
	require := require.New(t)
	ret, err := Marshal(nil, "")
	require.Nil(err)
	require.Equal("", ret)
}
