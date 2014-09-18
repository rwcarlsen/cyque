package queue

/*
#cgo CFLAGS: -I $SRCDIR/cctools/work_queue/src
#cgo CFLAGS: -I $SRCDIR/cctools/dttools/src
#cgo CFLAGS: -I $SRCDIR/cctools/chirp/src
#cgo LDFLAGS: $SRCDIR/cctools/work_queue/src/libwork_queue.a
#cgo LDFLAGS: $SRCDIR/cctools/dttools/src/libdttools.a
#cgo LDFLAGS: $SRCDIR/cctools/chirp/src/libchirp.a
#cgo LDFLAGS: -lm

#include <stdlib.h>
#include <errno.h>
#include <string.h>

#include "work_queue.h"
#include "list.h"

typedef struct work_queue work_queue;
typedef struct work_queue_task work_queue_task;
*/
import "C"
import (
	"fmt"
	"runtime"
	"time"
	"unsafe"
)

type Queue struct {
	q     *C.work_queue
	tasks map[int]*Task
}

func NewQueue(port int) (*Queue, error) {
	q, err := C.work_queue_create(C.int(port))
	if q == nil {
		return nil, err
	}

	queue := &Queue{q, make(map[int]*Task)}
	runtime.SetFinalizer(queue, freequeue)
	return queue, nil
}

func (q *Queue) Submit(t *Task) (taskid int) {
	id := int(C.work_queue_submit(q.q, t.t))
	q.tasks[id] = t
	return id
}

func (q *Queue) KillWorkers() (n int) {
	return int(C.work_queue_shut_down_workers(q.q, 0))
}

func (q *Queue) Empty() bool { return C.work_queue_empty(q.q) == 0 }

func (q *Queue) CancelAllTasks() {
	list := C.work_queue_cancel_all_tasks(q.q)
	C.list_delete(list)
	q.tasks = map[int]*Task{}
}

func (q *Queue) CancelTask(id int) {
	C.work_queue_cancel_by_taskid(q.q, C.int(id))
	delete(q.tasks, id)
}

func (q *Queue) Wait() *Task {
	t := C.work_queue_wait(q.q, C.WORK_QUEUE_WAITFORTASK)
	if t == nil {
		return nil
	}

	task := q.tasks[int(t.taskid)]
	delete(q.tasks, int(t.taskid))
	task.update()
	return task
}

func freequeue(q *Queue) { C.work_queue_delete(q.q) }

type Task struct {
	Id            int
	Command       string
	Stdout        string
	ReturnStatus  int
	Result        int
	Host          string
	Hostname      string
	Finished      time.Time
	Committed     time.Time
	Submitted     time.Time
	NSubmissions  int
	TotalExecTime time.Duration
	t             *C.work_queue_task
}

func (t *Task) update() {
	t.Command = C.GoString(t.t.command_line)
	t.Stdout = C.GoString(t.t.output)
	t.ReturnStatus = int(t.t.return_status)
	t.Result = int(t.t.result)
	t.Host = C.GoString(t.t.host)
	t.Hostname = C.GoString(t.t.hostname)
	t.Finished = convTime(t.t.time_task_finish)
	t.Submitted = convTime(t.t.time_task_submit)
	t.Committed = convTime(t.t.time_committed)
	t.NSubmissions = int(t.t.total_submissions)
	t.TotalExecTime = time.Duration(int(t.t.total_cmd_execution_time) * 1000)
}

func convTime(tm C.timestamp_t) time.Time {
	return time.Unix(int64(tm)/1000000, (int64(tm)%1000000)*1000)
}

func NewTask(cmd string) (*Task, error) {
	ccmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(ccmd))

	t := C.work_queue_task_create(ccmd)

	if t == nil {
		return nil, fmt.Errorf("failed to create task")
	}

	task := &Task{t: t}
	runtime.SetFinalizer(task, freetask)
	return task, nil
}

func freetask(t *Task) { C.work_queue_task_delete(t.t) }

func (t *Task) AddInfile(localpath, remotepath string, cache bool) error {
	return t.addfile(localpath, remotepath, cache, false)
}

func (t *Task) AddOutfile(localpath, remotepath string, cache bool) error {
	return t.addfile(localpath, remotepath, cache, true)
}

func (t *Task) addfile(local, remote string, cache bool, outfile bool) error {
	clocal := C.CString(local)
	cremote := C.CString(remote)
	defer C.free(unsafe.Pointer(clocal))
	defer C.free(unsafe.Pointer(cremote))

	ftype := C.WORK_QUEUE_INPUT
	if outfile {
		ftype = C.WORK_QUEUE_OUTPUT
	}

	wantcache := C.WORK_QUEUE_NOCACHE
	if cache {
		wantcache = C.WORK_QUEUE_CACHE
	}

	status, err := C.work_queue_task_specify_file(t.t, clocal, cremote, C.int(ftype), C.int(wantcache))
	if int(status) == 0 {
		return err
	}
	return nil
}
