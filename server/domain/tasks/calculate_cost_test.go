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
			},
			want: []Task{
				{
					ID:        parentID,
					Title:     "Create a new task",
					CreatedBy: 1,
					Cost:      20,
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
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			action := NewCalculateCost()
			assert.Equal(t, tt.want, action.Run(tt.give))
		})
	}
}
