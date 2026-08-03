package main

import (
	"os/exec"
)

func extract_background_args(args []string) (is_background bool, background_args []string) {
	is_background = false
	for _, arg := range args {
		if arg == "&" {
			is_background = true
		} else {
			background_args = append(background_args, arg)
		}
	}
	return is_background, background_args
}

var _background_jobs map[int]*exec.Cmd = make(map[int]*exec.Cmd)

func queue_job(prog *exec.Cmd) (job_id int, pid int, err error) {
	err = prog.Start()
	if (err != nil) {
		return -1, -1, err
	}
	job_id = 1
	var exists bool
	for true {
		_, exists = _background_jobs[job_id]
		if (!exists) {
			break
		}
		job_id ++
	}
	_background_jobs[job_id] = prog
	return job_id, prog.Process.Pid, nil
}