package router

import (
	"github.com/kataras/iris"
)

type task struct {
	app *iris.Application
	c   common
}

func NewRouterTask(app *iris.Application) *task {
	return &task{
		app: app,
		c:   common{},
	}
}

func (t *task) AddTask() {
	//ToDo
	t.app.Get("/task", t.root)
}

func (t *task) root(ctx iris.Context) {
	//state := taskService.NewTaskStateSnap()
	//var body string
	//todo
}
