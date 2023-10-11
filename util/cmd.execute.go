package util

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"golang.org/x/text/encoding/simplifiedchinese"
)

type Charset string

const (
	UTF8    = Charset("UTF-8")
	GB18030 = Charset("GB18030")
)

func ConvertByte2String(byte []byte, charset Charset) string {

	var str string
	switch charset {
	case GB18030:
		decodeBytes, _ := simplifiedchinese.GB18030.NewDecoder().Bytes(byte)
		str = string(decodeBytes)
	case UTF8:
		fallthrough
	default:
		str = string(byte)
	}

	return str
}

func ExecuteCommand(cmdStr string) (results string, err error) {
	// Create the command object
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.CommandContext(ctx, "cmd", "/C", cmdStr) // Replace with your command and arguments
	} else {
		cmd = exec.CommandContext(ctx, cmdStr) // Replace with your command and arguments
	}

	// Set up the output pipes
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	// Start the command
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	// Create scanners for stdout and stderr
	// Set output pipe encoding to UTF-8
	stdoutScanner := bufio.NewScanner(stdout)
	stderrScanner := bufio.NewScanner(stderr)

	outputChan := make(chan string)
	// Start goroutines to read stdout and stderr asynchronously
	go func() {
		defer close(outputChan)
		for stdoutScanner.Scan() {
			// line := stdoutScanner.Text()
			line := ConvertByte2String(stdoutScanner.Bytes(), GB18030)
			fmt.Println("stdout:", line)
			outputChan <- line
			// Process stdout line as needed
		}
	}()

	go func() {
		for stderrScanner.Scan() {
			// line := stderrScanner.Text()
			garbledStr := ConvertByte2String(stderrScanner.Bytes(), GB18030)
			fmt.Println("stderr:", garbledStr)

			// Process stderr line as needed
		}
	}()
	// Read output from the channel
	outbuilder := strings.Builder{}
	for line := range outputChan {
		fmt.Println("Output:", line)
		outbuilder.WriteString(line + "\n")
	}
	results = outbuilder.String()
	// Wait for the command to finish or a timeout
	timeout := time.After(30 * time.Second)
	select {
	case <-timeout:
		// Timeout occurred
		fmt.Println("Command timed out")
	case <-cmdDone(cmd):
		// Command completed successfully
		fmt.Println("Command completed")
	}

	fmt.Println("Done")
	return
}

// Helper function to wait for the command to finish
func cmdDone(cmd *exec.Cmd) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		cmd.Wait()
		close(done)
	}()
	return done
}
