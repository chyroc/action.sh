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

	flag.Parse()

	setupActionUser()

	gitAddFiles(strings.Split(*addFiles, ","))
	gitCommit(*msg)

}

// by: https://github.community/t/github-actions-bot-email-address/17204/5
func setupActionUser() {
	// goexec.New("git","config","--global","user.name","").RunInStream()
	assert(goexec.New("git", "config", "--global", "user.email", "41898282+github-actions[bot]@users.noreply.github.com").RunInStream())
}

func gitAddFiles(files []string) {
	for _, file := range files {
		file = strings.TrimSpace(file)
		_ = goexec.New("git", "add", file).RunInStream()
	}
}

func gitCommit(msg string) {
	_ = goexec.New("git", "commit", "-a", "-m", msg).RunInStream()
}

func assert(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// git commit -am "commit-by-action: $(date)" || (echo "no commit" && exit 0)
// git push
