package main

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/coreos/go-semver/semver"
)

var (
	Major = "major"
	Minor = "minor"
	Patch = "patch"

	// Answers.
	Yes  = "y"
	No   = "n"
	Quit = "q"
)

func gitDescribeTags() (string, error) {
	cmd := exec.Command("git", "describe", "--tags", "--abbrev=0")

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", err
	}

	return strings.TrimSpace(out.String()), nil
}

func main() {
	tag, err := gitDescribeTags()
	if err != nil {
		log.Fatal(err)
	}

	var prefix string
	if strings.HasPrefix(tag, "v") {
		prefix = "v"
		tag = tag[1:]
	}

	var bumped bool
	sv := semver.Must(semver.NewVersion(tag))

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Enter a version [major|minor|patch]: ")
	for scanner.Scan() {
		txt := scanner.Text()

		if bumped {
			switch txt {
			case Yes:
				fmt.Println(sv)
				return
			case Quit:
				return
			case No:
			default:
			}
			fmt.Print("Enter a version [major|minor|patch]: ")
			sv = semver.Must(semver.NewVersion(tag))
			bumped = false
		} else {
			bumped = true
			switch txt {
			case Major:
				sv.BumpMajor()
			case Minor:
				sv.BumpMinor()
			case Patch:
				sv.BumpPatch()
			case Quit:
				return
			case No:
				return
			default:
				bumped = false
			}
			if bumped {
				fmt.Printf("Bump version: %s%s -> %s%s? [ynq]\n", prefix, tag, prefix, sv.String())
			}
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println(err)
	}

}
