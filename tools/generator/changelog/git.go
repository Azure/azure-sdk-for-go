package changelog

import (
	"log"
	"os/exec"
)

func gitAddAll() error {
	log.Printf("Executing `git add *`")
	c := exec.Command("git", "add", "*")
	return c.Run()
}

func gitStash() error {
	log.Printf("Executing `git stash`")
	c := exec.Command("git", "stash")
	return c.Run()
}

func gitStashPop() error {
	log.Printf("Executing `git stash pop`")
	c := exec.Command("git", "stash", "pop")
	return c.Run()
}

func gitResetHead() error {
	log.Printf("Executing `git reset HEAD`")
	c := exec.Command("git", "reset", "HEAD")
	return c.Run()
}
