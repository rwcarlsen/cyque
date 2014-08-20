#!/usr/bin/env python

# Copyright (c) 2010- The University of Notre Dame.
# This software is distributed under the GNU General Public License.
# See the file COPYING for details.

# This program is a very simple example of how to use Work Queue.
# It accepts a list of files on the command line.
# Each file is compressed with gzip and returned to the user.

import work_queue as wq

import os
import sys

# Main program
if __name__ == '__main__':
  port = wq.WORK_QUEUE_DEFAULT_PORT

  if len(sys.argv) < 2:
    print "work_queue_example <file1> [file2] [file3] ..."
    print "Each file given on the command line will be compressed using a remote worker."
    sys.exit(1)

  # We create the tasks queue using the default port. If this port is already
  # been used by another program, you can try setting port = 0 to use an
  # available port.
  try:
      q = wq.WorkQueue(port)
  except:
      print "Instantiation of Work Queue failed!" 
      sys.exit(1)

  print "listening on port %d..." % q.port

  # We create and dispatch a task for each filename given in the argument list 
  for i in range(1, len(sys.argv)):
      infile = "%s" % sys.argv[i] 
      outfile = "%s.metric" % sys.argv[i]

      # Note that we write ./gzip here, to guarantee that the gzip version we
      # are using is the one being sent to the workers.
      command = "tar -xzf cyclus-cde.tar.gz && ./cde-package/cde-exec /home/r/cyc/install/bin/cyclus %s && cycdriver -metric cyclus.sqlite > %s && rm cyclus.sqlite;" % (infile, outfile)
      
      t = wq.Task(command)

      t.specify_file("cyclus-cde.tar.gz", type=wq.WORK_QUEUE_INPUT, flags=wq.WORK_QUEUE_CACHE)
      t.specify_file("cycdriver", type=wq.WORK_QUEUE_INPUT, flags=wq.WORK_QUEUE_CACHE)

      # files to be compressed are different across all tasks, so we do not
      # cache them. This is, of course, application specific. Sometimes you may
      # want to cache an output file if is the input of a later task.
      t.specify_file(infile, type=wq.WORK_QUEUE_INPUT, flags=wq.WORK_QUEUE_NOCACHE)
      t.specify_file(outfile, type=wq.WORK_QUEUE_OUTPUT, flags=wq.WORK_QUEUE_NOCACHE)
      
      # Once all files has been specified, we are ready to submit the task to the queue.
      taskid = q.submit(t)
      print "submitted task (id# %d): %s" % (taskid, t.command)

  print "waiting for tasks to complete..."
  while not q.empty():
      t = q.wait(5)
      if t:
          print "task (id# %d) complete: %s (return code %d)" % (t.id, t.command, t.return_status)
          if t.return_status != 0:
            # The task failed. Error handling (e.g., resubmit with new parameters, examine logs, etc.) here
            None
      #task object will be garbage collected by Python automatically when it goes out of scope

  print "all tasks complete!"

  #work queue object will be garbage collected by Python automatically when it goes out of scope
  sys.exit(0)
