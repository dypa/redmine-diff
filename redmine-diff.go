package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"sort"
)

const redmineUrl = "http://redmine/issue/"

func main() {
	if len(os.Args) != 4 {
		fmt.Printf("usage: %s <repo-path> <first-branch> <second-branch>\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}

	gitpath := os.Args[1]
	_, err := os.Open(gitpath + "/.git")
	if err != nil {
		fmt.Println(".git folder not found")
		os.Exit(2)
	}
	branch1 := os.Args[2]
	branch2 := os.Args[3]

	command := fmt.Sprintf("log --pretty=oneline %s..%s", branch1, branch2)

	os.Chdir(gitpath)
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/C", "git.exe " + command)
	} else {
		cmd = exec.Command("sh", "-c", "git " + command)
	}
	output, _ := cmd.CombinedOutput()

	pattern := regexp.MustCompile(`#(\d+)`)
	if matches := pattern.FindAllStringSubmatch(string(output), -1); matches != nil {
		var issues []int
		for _, issue := range matches {
			integer, _ := strconv.Atoi(issue[1])
			issues = append(issues, integer)
		}
		
		issues = removeDuplicates(issues)
		sort.Sort(sort.IntSlice(issues))
		
		for _, issue := range issues {
			fmt.Println(redmineUrl + strconv.Itoa(issue))
		}
	}
}

func AppendIfMissing(slice []int, i int) []int {
    for _, ele := range slice {
        if ele == i {
            return slice
        }
    }
    return append(slice, i)
}

//@see http://www.dotnetperls.com/remove-duplicates-slice
func removeDuplicates(elements []int) []int {
    // Use map to record duplicates as we find them.
    encountered := map[int]bool{}
    result := []int{}

    for v := range elements {
	if encountered[elements[v]] == true {
	    // Do not add duplicate.
	} else {
	    // Record this element as an encountered element.
	    encountered[elements[v]] = true
	    // Append to result slice.
	    result = append(result, elements[v])
	}
    }
    // Return the new slice.
    return result
}
