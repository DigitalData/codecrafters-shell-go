package main

import (
	"os"
	"os/exec"
	"slices"
	"strings"
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
	go prog.Wait()
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

	var job_id_index, job_id int
	var job *BackgroundJob
	var delete_job_ids []int
	var delete_job_id_indexes []int
	for job_id_index, job_id = range _job_ids {
		job, _ = _background_jobs[job_id]
		var raw string = job.raw

		symbol := " "
		switch job_id {
		case current_id:
			symbol = "+"
		case last_id:
			symbol = "-"
		}
		
		var prog *exec.Cmd = job.cmd
		var pstate *os.ProcessState = prog.ProcessState
		status := "Running"
		if (pstate != nil) {
			status = "Done"
			raw = strings.TrimSuffix(raw, " &")
			delete_job_ids = append(delete_job_ids, job_id)
			delete_job_id_indexes = append(delete_job_id_indexes, job_id_index)
		}
		outputs.outf("[%d]%s  %17s %s\n", job_id, symbol, status, raw)
	}

	for didx := range len(delete_job_ids) {
		delete_job_id := delete_job_ids[didx]
		delete_job_id_idx := delete_job_id_indexes[didx]
		delete(_background_jobs, delete_job_id)
		_job_ids = slices.Delete(_job_ids, delete_job_id_idx, delete_job_id_idx + 1)
	}
}