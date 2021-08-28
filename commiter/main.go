package main

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/chyroc/goexec"
)

func main() {
	fmt.Println("action.sh::commiter")

	addFiles := flag.String("add", "", "add files, split by `,`")
	msg := flag.String("msg", "", "msg")
	branch := flag.String("branch", "", "branch")
	force := flag.Bool("force", false, "force")
	flag.Parse()

	fmt.Println("+ setup action user")
	setupActionUser()

	fmt.Println("+ git add files")
	gitAddFiles(strings.Split(*addFiles, ","))

	fmt.Println("+ get changed files")
	changedFiles := gitGetChangedFiles()
	fmt.Println("+ files=", changedFiles)

	if len(changedFiles) > 0 {
		fmt.Printf("+changed files %d\n", len(changedFiles))

		fmt.Println("+ git get branched")
		currentBranch, branches := gitGetBranches()

		createBranch := !branches[*branch]
		fmt.Printf("+ current_branch: %s, new_branch: %s, checkout branch: %v\n", currentBranch, *branch, createBranch)
		if *branch != "" && *branch != currentBranch {
			gitNewBranch(*branch, createBranch)
		}

		fmt.Println("+ commit")
		gitCommit(*msg)

		fmt.Println("+ push")
		gitPush(*branch, *force)
	} else {
		fmt.Println("+ no changed files, skip commit and push")
	}
}

// by: https://github.community/t/github-actions-bot-email-address/17204/5
func setupActionUser() {
	assert(goexec.New("git", "config", "--global", "user.name", "github-actions[bot]").RunInStream())
	assert(goexec.New("git", "config", "--global", "user.email", "41898282+github-actions[bot]@users.noreply.github.com").RunInStream())
}

func gitGetBranches() (current string, res map[string]bool) {
	res = map[string]bool{}
	out, _, err := goexec.New("git", "branch").RunInTee()
	assert(err)
	for _, v := range strings.Split(out, "\n") {
		v = strings.TrimSpace(v)
		if v == "" {
			continue
		}
		if strings.HasPrefix(v, "* ") {
			current = v[2:]
		} else {
			res[v] = true
		}
	}
	return current, res
}

func gitGetChangedFiles() (res []string) {
	// 	git diff HEAD --name-only --diff-filter=AMCR
	out, _, err := goexec.New("git", "diff", "HEAD", "--name-only", "--diff-filter=AMCR").RunInTee()
	assert(err)
	for _, v := range strings.Split(out, "\n") {
		v = strings.TrimSpace(v)
		res = append(res, v)
	}
	return res
}

func gitAddFiles(files []string) {
	for _, file := range files {
		file = strings.TrimSpace(file)
		_ = goexec.New("git", "add", file).RunInStream()
	}
}

func gitNewBranch(branch string, create bool) {
	args := []string{"git", "checkout"}
	if create {
		args = append(args, "-b")
	}
	args = append(args, branch)
	assert(goexec.New(args...).RunInStream())
}

func gitCommit(msg string) {
	assert(goexec.New("git", "commit", "-a", "-m", msg).RunInStream())
}

func gitPush(branch string, force bool) {
	args := []string{"git", "push", "--set-upstream", "origin"}
	if branch != "" {
		args = append(args, branch)
	}
	if force {
		args = append(args, "-f")
	}
	assert(goexec.New(args...).RunInStream())
}

func assert(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
