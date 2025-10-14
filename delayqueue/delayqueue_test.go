package delayqueue

import (
	"context"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

func TestDelayQueue(t *testing.T) {
	// 创建延迟队列实例
	queue := NewDelayQueue("test_queue", &redis.Options{
		Addr: "localhost:637" +
			"9",
		DB: 0,
	})

	ctx := context.Background()

	// 添加延迟任务
	task := map[string]interface{}{
		"id":      1,
		"message": "test message",
	}
	err := queue.Push(ctx, task, 5*time.Second)
	if err != nil {
		log.Fatalf("Push task failed: %v", err)
	}

	// 获取队列长度
	size, err := queue.Size(ctx)
	if err != nil {
		log.Fatalf("Get queue size failed: %v", err)
	}
	fmt.Printf("Queue size: %d\n", size)

	// 非阻塞获取任务
	task1, err := queue.Pop(ctx, false, 0)
	if err != nil {
		log.Fatalf("Pop task failed: %v", err)
	}
	if task1 != nil {
		fmt.Printf("Got task immediately: %+v\n", task1)
	} else {
		fmt.Println("No task available")
	}

	// 阻塞等待任务
	fmt.Println("Waiting for task...")
	task2, err := queue.Pop(ctx, true, 10*time.Second)
	if err != nil {
		log.Fatalf("Pop task failed: %v", err)
	}
	if task2 != nil {
		fmt.Printf("Got task after waiting: %+v\n", task2)
	} else {
		fmt.Println("Timeout waiting for task")
	}
}
