package structx

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

// TimeWheel time wheel struct
type TimeWheel struct {
	interval       time.Duration
	ticker         *time.Ticker
	slots          []*list.List
	currentPos     int
	slotNum        int
	addTaskChannel chan *task
	stopChannel    chan bool
	taskRecord     *sync.Map
}

// Job callback function
type Job func(TaskData)

// TaskData callback params
type TaskData map[interface{}]interface{}

// task struct
type task struct {
	interval time.Duration
	times    int //-1:no limit >=1:run times
	circle   int
	key      interface{}
	job      Job
	taskData TaskData
}

// New create a empty time wheel
func New(interval time.Duration, slotNum int) *TimeWheel {
	if interval <= 0 || slotNum <= 0 {
		return nil
	}
	tw := &TimeWheel{
		interval:       interval,
		slots:          make([]*list.List, slotNum),
		currentPos:     0,
		slotNum:        slotNum,
		addTaskChannel: make(chan *task),
		stopChannel:    make(chan bool),
		taskRecord:     &sync.Map{},
	}

	tw.init()

	return tw
}

// Start start the time wheel
func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

// Stop stop the time wheel
func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case t := <-tw.addTaskChannel:
			tw.addTask(t)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

// AddTask add new task to the time wheel
func (tw *TimeWheel) AddTask(interval time.Duration, times int, key interface{}, data TaskData, job Job) error {
	if interval <= 0 || key == nil || job == nil || times < -1 || times == 0 {
		return errors.New("illegal task params")
	}

	_, ok := tw.taskRecord.Load(key)
	if ok {
		return errors.New("duplicate task key")
	}

	tw.addTaskChannel <- &task{interval: interval, times: times, key: key, taskData: data, job: job}
	return nil
}

// RemoveTask remove the task from time wheel
func (tw *TimeWheel) RemoveTask(key interface{}) error {
	if key == nil {
		return nil
	}

	value, ok := tw.taskRecord.Load(key)

	if !ok {
		return errors.New("task not exists, please check you task key")
	} else {
		// lazy remove task
		t := value.(*task)
		t.times = 0
		tw.taskRecord.Delete(t.key)
	}
	return nil
}

// UpdateTask update task times and data
func (tw *TimeWheel) UpdateTask(key interface{}, interval time.Duration, taskData TaskData) error {
	if key == nil {
		return errors.New("illegal key, please try again")
	}

	value, ok := tw.taskRecord.Load(key)

	if !ok {
		return errors.New("task not exists, please check you task key")
	}
	t := value.(*task)
	t.taskData = taskData
	t.interval = interval
	return nil
}

// time wheel initialize
func (tw *TimeWheel) init() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

//
func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAddRunTask(l)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

// add task
func (tw *TimeWheel) addTask(t *task) {
	if t.times == 0 {
		return
	}

	pos, circle := tw.getPositionAndCircle(t.interval)
	t.circle = circle

	tw.slots[pos].PushBack(t)

	//record the task
	tw.taskRecord.Store(t.key, t)
}

// scan task list and run the task
func (tw *TimeWheel) scanAddRunTask(l *list.List) {

	if l == nil {
		return
	}

	for item := l.Front(); item != nil; {
		t := item.Value.(*task)

		if t.times == 0 {
			next := item.Next()
			l.Remove(item)
			tw.taskRecord.Delete(t.key)
			item = next
			continue
		}

		if t.circle > 0 {
			t.circle--
			item = item.Next()
			continue
		}

		go t.job(t.taskData)
		next := item.Next()
		l.Remove(item)
		item = next

		if t.times == 1 {
			t.times = 0
			tw.taskRecord.Delete(t.key)
		} else {
			if t.times > 0 {
				t.times--
			}
			tw.addTask(t)
		}
	}
}

// get the task position
func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (int, int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle := delaySeconds / intervalSeconds / tw.slotNum
	pos := (tw.currentPos + delaySeconds/intervalSeconds) % tw.slotNum
	return circle, pos
}
