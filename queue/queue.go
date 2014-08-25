package queue

/*

#cgo CFLAGS: -I "/home/robert/git/cyque/queue/work_queue"
#cgo LDFLAGS: "libwork_queue.a"

#include "work_queue.h"

#include <stdlib.h>
#include <errno.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"runtime"
	"unsafe"
)

type Queue struct {
	q     *C.work_queue
	tasks map[uintptr]*Task
}

func NewQueue(port int) (*Queue, error) {
	q := C.work_queue_create(C.int(port))
	if q == nil {
		return nil, fmt.Errorf("couldn't listen on port %v: %v", port, C.GoString(C.strerror(errno)))
	}

	queue := &Queue{q}
	runtime.SetFinalizer(queue, freequeue)
	return queue, nil
}

func (q *Queue) Submit(t *Task) (taskid int) {
	q.tasks[uintptr(t.t)] = t
	return int(C.work_queue_submit(q.q, t.t))
}

func (q *Queue) Empty() bool { return C.work_queue_empty(q.q) == 0 }

func (q *Queue) Wait(secs int) *Task {
	if secs < 0 {
		secs = int(C.WORK_QUEUE_WAITFORTASK)
	}

	t := C.work_queue_wait(q.q, C.int(secs))
	if t == nil {
		return nil
	}

	task := q.tasks[uintptr(t)]
	delete(q.tasks, uintptr(t))
	return task
}

func freequeue(q *Queue) { C.work_queue_delete(q.q) }

type Task struct {
	t *C.work_queue_task
}

func NewTask(cmd string) (*Task, error) {
	ccmd := C.CString(cmd)
	defer C.free(unsafe.Pointer(ccmd))

	t := C.work_queue_task_create(ccmd)

	if t == nil {
		return nil, fmt.Errorf("failed to create task")
	}

	task := &Task{t}
	runtime.SetFinalizer(task, freetask)
	return task
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

	status := C.work_queue_task_specify_file(t.t, clocal, cremote, ftype, wantcache)
	if int(status) == 0 {
		return fmt.Errorf("failed to add file %v to task: %v", local, C.GoString(strerror(errno)))
	}
	return nil
}
