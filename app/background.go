package main

import (
	"os/exec"
	"slices"
)

type BackgroundJob struct {
	raw string
	cmd *exec.Cmd
}

var _background_jobs map[int]*BackgroundJob = make(map[int]*BackgroundJob)
var _job_ids []int

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


func queue_job(raw string, prog *exec.Cmd) (job_id int, pid int, err error) {
	err = prog.Start()
	if (err != nil) {
		return -1, -1, err
	}
	job_id = 1
	if (len(_job_ids) > 0) {
		var exists bool
		max_job_id := slices.Max(_job_ids)
		for job_id <= max_job_id {
			_, exists = _background_jobs[job_id]
			if (!exists) {
				break
			}
			job_id ++
		}
	}
	_background_jobs[job_id] = &BackgroundJob{raw, prog}
	_job_ids = append(_job_ids, job_id)
	slices.Sort(_job_ids)
	return job_id, prog.Process.Pid, nil
}

func handle_jobs(raw_line string, cmd string, cmd_args []string, has_args bool, outputs *Outputs) {
	var num_jobs = len(_job_ids)

	if (num_jobs == 0) {
		return
	}

	current_id := _job_ids[num_jobs - 1]
	last_id := -1
	if (num_jobs > 1) {
		last_id = _job_ids[num_jobs - 2]
	}

	var job_id int
	var job *BackgroundJob
	for _, job_id = range _job_ids {
		job, _ = _background_jobs[job_id]
		symbol := " "
		if (job_id == current_id) {
			symbol = "+"
		} else if (job_id == last_id) {
			symbol = "-"
		}
		// var prog *exec.Cmd = job.cmd
		// var pstate *os.ProcessState = prog.ProcessState
		status := "Running"
		// if (pstate.Success()) {
		// 	status = "Success"
		// } else if (pstate.Exited()) {
		// 	status = "Exited"
		// }
		outputs.outf("[%d]%s  %17s %s\n", job_id, symbol, status, job.raw)
	}
}