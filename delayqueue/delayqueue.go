package delayqueue

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// DelayedTask 延迟任务结构
type DelayedTask struct {
	Data       interface{} `json:"data"`        // 任务数据
	CreateTime int64       `json:"create_time"` // 创建时间
}

// DelayQueue Redis延迟队列
type DelayQueue struct {
	client    *redis.Client
	queueName string
}

// NewDelayQueue 创建延迟队列实例
func NewDelayQueue(queueName string, options *redis.Options) *DelayQueue {
	if options == nil {
		options = &redis.Options{
			Addr:     "localhost:6379",
			Password: "", // no password set
			DB:       0,  // use default DB
		}
	}

	return &DelayQueue{
		client:    redis.NewClient(options),
		queueName: queueName,
	}
}

// Push 添加延迟任务
func (dq *DelayQueue) Push(ctx context.Context, data interface{}, delay time.Duration) error {
	task := DelayedTask{
		Data:       data,
		CreateTime: time.Now().Unix(),
	}

	// 序列化任务数据
	taskBytes, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("marshal task failed: %w", err)
	}

	// 计算执行时间
	executeTime := time.Now().Add(delay)

	// 添加到有序集合
	err = dq.client.ZAdd(ctx, dq.queueName, redis.Z{
		Score:  float64(executeTime.Unix()),
		Member: string(taskBytes),
	}).Err()

	if err != nil {
		return fmt.Errorf("add task to redis failed: %w", err)
	}

	return nil
}

// Pop 获取到期的任务
func (dq *DelayQueue) Pop(ctx context.Context, block bool, timeout time.Duration) (*DelayedTask, error) {
	for {
		// 获取当前时间之前的任务
		result, err := dq.client.ZRangeByScore(ctx, dq.queueName, &redis.ZRangeBy{
			Min:    "0",
			Max:    fmt.Sprintf("%d", time.Now().Unix()),
			Offset: 0,
			Count:  1,
		}).Result()

		if err != nil {
			return nil, fmt.Errorf("get task from redis failed: %w", err)
		}

		if len(result) > 0 {
			// 尝试从队列中移除该任务
			removed, err := dq.client.ZRem(ctx, dq.queueName, result[0]).Result()
			if err != nil {
				return nil, fmt.Errorf("remove task from redis failed: %w", err)
			}

			if removed > 0 {
				var task DelayedTask
				err = json.Unmarshal([]byte(result[0]), &task)
				if err != nil {
					return nil, fmt.Errorf("unmarshal task failed: %w", err)
				}
				return &task, nil
			}
		}

		if !block {
			return nil, nil
		}

		// 如果设置了超时时间且已超时，则返回
		if timeout > 0 {
			select {
			case <-time.After(time.Second):
				timeout -= time.Second
				if timeout <= 0 {
					return nil, nil
				}
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		} else {
			// 无超时时间，持续等待
			select {
			case <-time.After(time.Second):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
}

// Size 获取队列长度
func (dq *DelayQueue) Size(ctx context.Context) (int64, error) {
	return dq.client.ZCard(ctx, dq.queueName).Result()
}

// Clear 清空队列
func (dq *DelayQueue) Clear(ctx context.Context) error {
	return dq.client.Del(ctx, dq.queueName).Err()
}

// GetDelayedTasks 获取所有延迟任务
func (dq *DelayQueue) GetDelayedTasks(ctx context.Context) ([]map[string]interface{}, error) {
	tasks, err := dq.client.ZRangeWithScores(ctx, dq.queueName, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("get tasks from redis failed: %w", err)
	}

	var result []map[string]interface{}
	for _, task := range tasks {
		var taskData DelayedTask
		err = json.Unmarshal([]byte(task.Member.(string)), &taskData)
		if err != nil {
			log.Printf("unmarshal task failed: %v", err)
			continue
		}

		result = append(result, map[string]interface{}{
			"data":         taskData.Data,
			"create_time":  taskData.CreateTime,
			"execute_time": int64(task.Score),
		})
	}

	return result, nil
}
