package simpleTaskPool

import (
	"context"
	"fmt"
	"time"
)

//TaskFuncType task实体
type TaskFuncType func(ctx context.Context) (result interface{}, error error)

//Task 任务的数据结构
type Task struct {
	Name string
	Ctx  context.Context //可用于保存参数
	Func TaskFuncType
}

// TaskExecResult 任务执行返回值的数据结构
type TaskExecResult struct {
	Name   string
	Finish bool
	Result interface{}
	Err    error
}

//TaskPool 任务池的数据结构
type TaskPool struct {
	TaskList   []*Task
	Timeout    time.Duration
	ResultInfo map[string]*TaskExecResult
}

//NewTaskPool 创建任务池
func NewTaskPool(timeout time.Duration) *TaskPool {
	tp := &TaskPool{
		TaskList:   []*Task{},
		Timeout:    timeout,
		ResultInfo: make(map[string]*TaskExecResult),
	}
	return tp
}

//AddTask 向任务池中添加任务
func (tp *TaskPool) AddTask(task *Task) error {

	var err error

	if task == nil {
		err = fmt.Errorf("task is nil")
		return err
	}

	if task.Func == nil {
		err = fmt.Errorf("task func is nil")
		return err
	}
	if task.Ctx == nil {
		err = fmt.Errorf("task ctx is nil")
		return err
	}

	if task.Name == "" {
		err = fmt.Errorf("task name is nil")
		return err
	}

	for k, _ := range tp.ResultInfo {
		if k == task.Name {
			err = fmt.Errorf("repeat task name %v", task.Name)
			return err
		}
	}

	tp.TaskList = append(tp.TaskList, task)
	tp.ResultInfo[task.Name] = &TaskExecResult{}

	return nil
}

//Run 运行
func (tp *TaskPool) Run() error {
	leftTaskNum := len(tp.TaskList)
	if leftTaskNum == 0 {
		return nil
	}

	resultChan := make(chan *TaskExecResult, len(tp.TaskList))

	for _, task := range tp.TaskList {
		//并发开启任务
		go func(t *Task) {
			result, err := t.Func(t.Ctx)
			taskExecResult := &TaskExecResult{
				Name:   t.Name,
				Result: result,
				Err:    err,
				Finish: true,
			}

			resultChan <- taskExecResult
		}(task)
	}

	//timeout
	timeout := time.After(tp.Timeout)
	for leftTaskNum > 0 {
		select {
		case taskExecResult := <-resultChan:
			tp.ResultInfo[taskExecResult.Name] = taskExecResult
			leftTaskNum--
		case <-timeout:
			err := fmt.Errorf("timeout leftTaskNum:%v", leftTaskNum)
			return err
		}
	}
	return nil
}
