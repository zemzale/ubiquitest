package tasks

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCalculateCost(t *testing.T) {
	t.Parallel()

	var (
		parentID = uuid.MustParse("a3afc3d5-9717-40d8-9e66-2c0b9c2b6a51")
		childID1 = uuid.MustParse("a3afc3d5-9717-40d8-9e66-2c0b9c2b6a52")
		childID2 = uuid.MustParse("a3afc3d5-9717-40d8-9e66-2c0b9c2b6a53")
		childID3 = uuid.MustParse("a3afc3d5-9717-40d8-9e66-2c0b9c2b6a54")
	)

	tests := []struct {
		name string
		give []Task
		want []Task
	}{
		{
			name: "calculate cost",
			give: []Task{
				{
					ID:        parentID,
					Title:     "Create a new task",
					CreatedBy: 1,
				},
				{
					ID:        childID1,
					Title:     "Create a new task",
					CreatedBy: 1,
					ParentID:  parentID,
					Cost:      10,
				},
				{
					ID:        childID2,
					Title:     "Create a new task",
					CreatedBy: 1,
					ParentID:  parentID,
					Cost:      10,
				},
				{
					ID:        childID3,
					Title:     "Create a new task",
					CreatedBy: 1,
					ParentID:  childID2,
					Cost:      33,
				},
			},
			want: []Task{
				{
					ID:        parentID,
					Title:     "Create a new task",
					CreatedBy: 1,
					Cost:      10 + 10 + 33,
				},
				{
					ID:        childID1,
					Title:     "Create a new task",
					CreatedBy: 1,
					ParentID:  parentID,
					Cost:      10,
				},
				{
					ID:        childID2,
					Title:     "Create a new task",
					CreatedBy: 1,
					ParentID:  parentID,
					Cost:      10 + 33,
				},
				{
					ID:        childID3,
					Title:     "Create a new task",
					CreatedBy: 1,
					ParentID:  childID2,
					Cost:      33,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := NewCalculateCost()
			assert.ElementsMatch(t, tt.want, action.Run(tt.give))
		})
	}
}
