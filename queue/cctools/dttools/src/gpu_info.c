/*
Copyright (C) 2013- The University of Notre Dame
This software is distributed under the GNU General Public License.
See the file COPYING for details.
*/

#include "gpu_info.h"

#include <stddef.h>
#include <stdio.h>
#include <stdlib.h>
#include <sys/types.h>
#include <sys/wait.h>
#include <unistd.h>

#define GPU_AUTODETECT "cctools_gpu_autodetect"

int gpu_info_get()
{
	pid_t pid;
	int pipefd[2];
	pipe(pipefd);

	switch (pid = fork()) {
	case -1:
		return 0;
	case 0:
		close(pipefd[0]);
		dup2(pipefd[1], fileno(stdout));
		char *args[] = {GPU_AUTODETECT, NULL};
		if(!access(GPU_AUTODETECT, R_OK|X_OK)){
			execv(GPU_AUTODETECT, args);
		} else {
			execvp(GPU_AUTODETECT, args);
		}
		return 0;
	default:
		close(pipefd[1]);
		int status = 0;
		int gpu_count = 0;
		char buffer[10]; /* Enough characters to hold a decimal representation of a 32 bit int. */
		if(read(pipefd[0], buffer, 10)){
			waitpid(pid, &status, 0);
			gpu_count = atoi(buffer);
		}

		close(pipefd[0]);
		return gpu_count;
	}
}

/* vim: set noexpandtab tabstop=4: */
