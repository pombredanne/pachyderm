package clock

import (
	"testing"

	"github.com/pachyderm/pachyderm/src/client/pkg/require"
	"github.com/pachyderm/pachyderm/src/server/pfs/persist"
)

func TestNewBranchClocks(t *testing.T) {
	b := NewBranchClocks("master")

	require.Equal(t, 1, len(b))           // Only one branch
	require.Equal(t, 1, len(b[0].Clocks)) // Only one commit
	require.Equal(t, "master", b[0].Clocks[0].Branch)
	require.Equal(t, 0, b[0].Clocks[0].Clock)
}

func TestNewChild(t *testing.T) {
	b := NewBranchClocks("master")
	b2 := NewChild(b[0])

	require.Equal(t, 1, len(b2.Clocks))

	b1 := &persist.BranchClock{
		Clocks: []*persist.Clock{
			{
				Branch: "master",
				Clock:  0,
			},
			{
				Branch: "foo",
				Clock:  0,
			},
		},
	}

	expected := persist.BranchClocks{
		{
			Clocks: []*persist.Clock{
				{
					Branch: "master",
					Clock:  0,
				},
				{
					Branch: "foo",
					Clock:  1,
				},
			},
		},
	}
	b2 = NewChild(b1)
	require.Equal(t, expected, b2)
}

func TestNewChildOfBranchClocks(t *testing.T) {
	input := persist.BranchClocks{
		{
			Clocks: []*persist.Clock{
				{
					Branch: "foo",
					Clock:  1,
				},
				{
					Branch: "bar",
					Clock:  1,
				},
			},
		},
		{
			Clocks: []*persist.Clock{
				{
					Branch: "master",
					Clock:  1,
				},
			},
		},
	}
	expected := persist.BranchClocks{{
		Clocks: []*persist.Clock{
			{
				Branch: "master",
				Clock:  2,
			},
		}}}
	actual, err := NewChildOfBranchClocks(input, "master")
	require.NoError(t, err)
	require.Equal(t, expected, actual)

	actual, err = NewChildOfBranchClocks(input, "abba")
	require.YesError(t, err)
	require.Nil(t, actual)
}

func TestAddClock(t *testing.T) {
	// AddClock adds a BranchClock to a BranchClocks.
	// Returns an error if the BranchClock already exists in the BranchClocks
	input := persist.BranchClocks{
		{
			Clocks: []*persist.Clock{
				{
					Branch: "foo",
					Clock:  1,
				},
				{
					Branch: "bar",
					Clock:  1,
				},
			},
		},
		{
			Clocks: []*persist.Clock{
				{
					Branch: "master",
					Clock:  1,
				},
			},
		},
	}
	newClock := &persist.BranchClock{
		Clocks: []*persist.Clock{
			{
				Branch: "bar",
				Clock:  4,
			},
		},
	}
	b, err := AddClock(input, newClock)
	require.NoError(t, err)
	expected := append(input, newClock)
	require.Equal(t, expected, b)

	b2, err := AddClock(b, newClock)
	require.YesError(t, err)
	require.Nil(t, b2)
}