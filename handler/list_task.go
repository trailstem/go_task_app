package handler

import (
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/trailstem/go_task_app/entity"
	"github.com/trailstem/go_task_app/store"
	"github.com/trailstem/go_task_app/testutil"
)

type ListTask struct {
	// Store *store.TaskStore
	DB   *sqlx.DB
	Repo store.Repository
}

type task struct {
	ID     entity.TaskID     `json:"id"`
	Title  string            `json:"title"`
	Status entity.TaskStatus `json:"status"`
}

func (lt *ListTask) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	tasks, err := lt.Repo.ListTasks(ctx, lt.DB, 1)

	if err != nil {
		testutil.RespondJSON(ctx, w, &testutil.ErrResponse{
			Message: err.Error(),
		}, http.StatusInternalServerError)
		return
	}

	rsp := []task{}

	for _, t := range tasks {

		rsp = append(rsp, task{
			ID:     t.ID,
			Title:  t.Title,
			Status: t.Status,
		})
		testutil.RespondJSON(ctx, w, rsp, http.StatusOK)
	}
}
