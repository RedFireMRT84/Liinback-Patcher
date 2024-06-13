package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/inancgumus/screen"
)

var header = "Liinback Patcher\nProgrammed by: RedFire\n\n"

func main() {
	clear()
	fmt.Printf(header)
	fmt.Printf("A. Start\nB. Exit the patcher\n\nChoose one:")
	var input string
	fmt.Scanln(&input)
	selectOption(input, begin, main)
}

func begin() {
	clear()
	fmt.Printf("%s", header)
	fmt.Printf("Well, hey there. Welcome to the patcher. The patcher will go ahead and get the requirements together.")
	fmt.Printf("A. Begin the patching!\nB. Exit the patcher\n\nChoose one:")
	var input string
	fmt.Scanln(&input)
	selectOption(input, downloadApp, begin)
}

func downloadApp() {
	clear()
	fmt.Printf(header)
	fmt.Println("Downloading the YouTube channel...")

	wadURL := "http://liinback2.atspace.tv/base.wad"
	err := downloadFile("base.wad", wadURL)
	if err != nil {
		handleError(err)
		return
	}
	patchApp()
}

func finish() {
	clear()
	fmt.Printf(header)
	fmt.Printf("Patching has been completed! The WAD file is ready for use.")
	fmt.Printf("\nPress any key to exit the patcher.")
	var input string
	fmt.Scanln(&input)
}

func handleError(err error) {
	clear()
	fmt.Printf(header)
	fmt.Println("An error has occurred. Either contact RedFire on Discord: redfire_mrt84 or email him:\n\nredfirestudiosmrtv@gmail.com for support, along with this error.")
	fmt.Printf("\nError: %s\n", err)
	fmt.Println("\nPress any key to exit this program.")
	var input string
	fmt.Scanln(&input)
	os.Exit(1)
}

func clear() {
	screen.Clear()
	screen.MoveTopLeft()
}

func selectOption(input string, gotoFunc func(), currentFunc func()) {
	switch input {
	case "A":
		gotoFunc()
	case "B":
		os.Exit(0)
	default:
		currentFunc()
	}
}

func downloadFile(filename, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func patchApp() {
	var cmd *exec.Cmd

	// Determine the current operating system
	switch strings.ToLower(runtime.GOOS) {
	case "windows":
		cmd = exec.Command("cmd", "/c", "freethewads.exe", "base.wad")
	case "linux":
		cmd = exec.Command("regionfree", "base.wad")
	default:
		handleError(fmt.Errorf("unsupported operating system: %s", runtime.GOOS))
		return
	}

	err := cmd.Run()
	if err != nil {
		handleError(err)
		return
	}

	// Rename the file after patching
	err = os.Rename("base.wad", "YouTube (Liinback).wad")
	if err != nil {
		handleError(err)
		return
	}

	finish()
}
