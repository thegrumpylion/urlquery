package urlquery

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUnmarshalArray(t *testing.T) {
	require := require.New(t)

	arr := []int{}
	err := UnmarshalName(&arr, "y", "y=5&y=16&y=49&y=50")
	require.Nil(err)
	require.Equal([]int{5, 16, 49, 50}, arr)
}

func TestUnmarshalNestedStruct(t *testing.T) {
	require := require.New(t)

	type rec struct {
		ID    int
		Value string
	}
	s := &struct {
		Env   map[string]map[string]string
		Host  string
		Port  int
		Valid bool
		Recs  []rec
	}{}

	err := Unmarshal(s, "Env.Ena.Ena=enaena&Env.Ena.Dio=enadio&Env.Dio.Ena=Dioena&Env.Dio.Dio=diodio&Host=localhost&Port=8080&Valid=true&Recs.0.ID=5&Recs.0.Value=some+value&Recs.1.ID=17&Recs.1.Value=%24ymbl%25s%3B%2F")

	require.Nil(err)

	fmt.Println(s.Env)
	fmt.Println(s.Recs)

}
