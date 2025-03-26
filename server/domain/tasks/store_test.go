package tasks

import (
	"testing"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zemzale/ubiquitest/storage"
)

func TestStore(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		giveTask  Task
		prepareDB func(t *testing.T, db *sqlx.DB)
		wantErr   bool
	}{
		{
			name: "store task",
			prepareDB: func(t *testing.T, db *sqlx.DB) {
				t.Helper()
				_, err := db.Exec("INSERT INTO users (username) VALUES (?)", "user")
				require.NoError(t, err, "failed to insert user")
			},
			giveTask: Task{
				ID:        uuid.MustParse("a3afc3d5-9717-40d8-9e66-2c0b9c2b6a53"),
				Title:     "Create a new task",
				CreatedBy: 1,
			},
		},
		{
			name:      "fail to store without user",
			prepareDB: func(t *testing.T, db *sqlx.DB) { t.Helper() },
			giveTask: Task{
				ID:    uuid.MustParse("a3afc3d5-9717-40d8-9e66-2c0b9c2b6a53"),
				Title: "Create a new task",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := sqlx.Open("sqlite3", ":memory:")
			require.NoError(t, err, "failed to open database")
			t.Cleanup(func() {
				db.Close()
			})

			require.NoError(t, storage.CreateDB(db))
			tt.prepareDB(t, db)

			taskRepo := storage.NewTaskRepository(db)

			action := NewStore(NewUpdateParentCost(NewFindAllParents(taskRepo), taskRepo), taskRepo, storage.NewUserRepository(db))
			if tt.wantErr {
				assert.Error(t, action.Run(tt.giveTask), "expected error")
				return
			}
			assert.NoError(t, action.Run(tt.giveTask), "failed to run task")

			var id string
			err = db.Get(&id, "SELECT id FROM tasks WHERE id = ?", tt.giveTask.ID.String())
			require.NoError(t, err, "failed to get task id from DB")
			assert.Equal(t, tt.giveTask.ID.String(), id, "task id does not match expected")
		})
	}
}
