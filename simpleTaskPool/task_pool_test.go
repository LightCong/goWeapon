package simpleTaskPool

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"
)

func taskfunc1(ctx context.Context) (interface{}, error) {
	fmt.Println("hello word !")
	return nil, nil
}

func taskfunc2(ctx context.Context) (interface{}, error) {
	val := ctx.Value("name")
	name, ok := val.(string)
	if !ok {
		err := fmt.Errorf("invalid param %v", val)
		log.Println(err.Error())
		return nil, err
	}
	fmt.Println("hello ", name)
	return nil, nil
}

func taskfunc3(ctx context.Context) (interface{}, error) {
	time.Sleep(2 * time.Second)
	return nil, nil
}

func TestTaskPool_Run(t *testing.T) {
	taskPool := NewTaskPool(1 * time.Second)
	task1 := Task{
		Name: "task1",
		Ctx:  context.Background(),
		Func: taskfunc1,
	}

	ctx := context.Background()
	nameCtx := context.WithValue(ctx, "name", "congxiaoliang")

	task2 := Task{
		Name: "task2",
		Ctx:  nameCtx,
		Func: taskfunc2,
	}

	task3 := Task{
		Name: "task3",
		Ctx:  context.Background(),
		Func: taskfunc3,
	}

	taskPool.AddTask(&task1)
	taskPool.AddTask(&task2)
	taskPool.AddTask(&task3)

	err := taskPool.Run()
	fmt.Println(err)

	for k, v := range taskPool.ResultInfo {
		fmt.Println(k, v.Name, v.Finish, v.Result, v.Err)
	}
}
