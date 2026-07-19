package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/term"
)


func get_completion_parts(line string) (parts []string) {
	simple_line := strings.ReplaceAll(line, "\t", " ")
	return strings.Split(simple_line, " ")
}

func match_commands(partial string, matches []string) (new_matches []string, exact bool) {
	commands := []string{CMD_ECHO, CMD_EXIT}
	for _, command := range commands {
		if command == partial {
			return []string{command}, true
		} else if strings.HasPrefix(command, partial) {
			matches = append(matches, command)
		}
	}
	return matches, false
}

func match_path_programs(partial string, matches []string) (new_matches []string, exact bool) {
	var raw_path string = os.Getenv("PATH")
	raw_path = strings.ReplaceAll(raw_path, ";", ":")
	var paths []string = strings.Split(raw_path, ":")
	for _, path := range paths {
		files, err := os.ReadDir(path)
		if err != nil {
			continue
		}

		for _, file := range files {
			raw_filename := file.Name()
			filename := strings.TrimSuffix(raw_filename, filepath.Ext(file.Name()))
			if filename == partial {
				return []string{partial}, true
			} else if raw_filename == partial {
				return []string{partial}, true
			} else if strings.HasPrefix(raw_filename, partial) {
				matches = append(matches, raw_filename)
			}
		}
	}

	slices.Sort(matches)
	return matches, false
}

func match_dir(partial string, matches []string) (new_matches []string, exact bool) {
	var err error

	dir := path.Dir(partial)
	var wd string
	wd, err = os.Getwd()
	if (err != nil) {
		panic(err)
	}
	wd_fs := os.DirFS(wd)
	files, err := fs.ReadDir(wd_fs, dir)

	for  _, file := range files {
		filename := file.Name()
		if (dir != ".") {
			filename = dir + string(os.PathSeparator) + filename
		}
		if (file.IsDir()) {
			filename += string(os.PathSeparator)
		}
		if filename == partial {
			return []string{partial}, true
		} else if (strings.HasPrefix(filename, partial)) {
			matches = append(matches, filename)
		}
	}
	slices.Sort(matches)
	return matches, false
}

func match_completion_script(line string, parts []string) (matches []string) {
	num_parts := len(parts)
	if (num_parts == 0) {
		return matches
	}
	program := parts[0]
	completion_script, exists := _completions[program]
	if (!exists) {
		return matches
	}

	partial := ""
	if (num_parts > 1) {
		partial = parts[num_parts - 1]
	}
	previous_word := ""
	if (num_parts > 1) {
		previous_word = parts[num_parts - 2]
	}


	prog := exec.Command(completion_script, program, partial, previous_word)
	prog.Env = append(prog.Env, fmt.Sprintf("COMP_LINE=%s", line))
	prog.Env = append(prog.Env, fmt.Sprintf("COMP_POINT=%d", len(line)))
	out, err := prog.CombinedOutput()
	if (err != nil) {
		log.Fatal(err)
	}
	if (len(out) == 0) {
		fmt.Print("\a")
	} else {
		match_output := strings.TrimRight(string(out), "\r\n")
		matches = strings.Split(match_output, "\n")
	}
	return matches
}

func match_autocomplete(partial string, is_arg bool) (matches []string, exact bool) {
	if (is_arg) {
		return match_dir(partial, matches)
	}

	matches, exact = match_commands(partial, matches)
	if len(matches) > 0 {
		return matches, exact
	}

	return match_path_programs(partial, matches)
}

func print_matches(matches []string) {
	entry_size := 0
	for _, match := range matches {
		entry_size = max(entry_size, len(match))
	}

	entry_size += 3
	slices.Sort(matches)
	text := ""
	width, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	max_cols := int(width / entry_size)
	num_cols := 0
	for _, match := range matches {
		text += fmt.Sprintf("%-*s", entry_size, match)
		num_cols += 1
		if num_cols >= max_cols {
			text += "\n"
			num_cols = 0
		}
	}

	fmt.Printf("\r\n%s\r\n", text)
}

func handle_matches(line string, partial string, matches []string, double_tab bool) (new_line string) {
	partial_suffix := ""

	if len(matches) == 1 {
		partial_suffix = strings.TrimPrefix(matches[0], partial)
		len_suffix := len(partial_suffix)
		if (len_suffix > 0 && partial_suffix[len_suffix - 1] != os.PathSeparator) {
			partial_suffix += " "
		} 
	} else {
		lcp := partial
		for {
			lcp_index := len(lcp)
			first_match := matches[0]
			if len(first_match) <= lcp_index {
				break
			}
			lcp_char := first_match[lcp_index]
			break_early := false
			for _, match := range matches[1:] {
				if len(match) <= lcp_index {
					break_early = true
					break
				}
				if lcp_char != match[lcp_index] {
					break_early = true
					break
				}
			}
			if break_early {
				break
			}
			lcp += string(lcp_char)
		}

		partial_suffix = strings.TrimPrefix(lcp, partial)
		if lcp == partial && double_tab {
			print_matches(matches)
			fmt.Printf("$ %s", line)
		} else if !double_tab {
			fmt.Print("\a")
		}
	}
	fmt.Print(partial_suffix)
	return line + partial_suffix
}

func handle_autocomplete(line string, double_tab bool) (new_line string, new_double_tab bool) {
	var matches []string
	var exact_match bool	
	var parts []string
	parts = get_completion_parts(line)
	num_parts := len(parts)
	is_arg := num_parts > 1
	partial := parts[num_parts - 1]

	matches = match_completion_script(line, parts)
	if (len(matches) == 0) {
		matches, exact_match = match_autocomplete(partial, is_arg)
	}
	if exact_match {
		fmt.Print(" ")
		line += " "
	} else if len(matches) == 0 {
		fmt.Print("\a")
	} else {
		line = handle_matches(line, partial, matches, double_tab)
	}
	return line, !double_tab
}
