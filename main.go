package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"runtime/pprof"
	"strings"
	"sync"
)

var (
	Verbose bool
	Workers chan bool
	mu      sync.Mutex
	stdinCh chan string
)

func main() {
	var n int
	var profile string
	flag.IntVar(&n, "n", runtime.NumCPU(), "workers num")
	flag.StringVar(&profile, "profile", "", "prof filename")
	flag.BoolVar(&Verbose, "v", false, "verbose mode")
	flag.Parse()

	if profile != "" {
		verboseLog("CPU profile enabled", profile)
		f, err := os.Create(profile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	stdinCh = make(chan string, 4096)
	command := flag.Args()
	if len(command) == 0 {
		fmt.Fprintln(os.Stderr, "sub command required")
		os.Exit(1)
	}

	var start sync.WaitGroup
	var done sync.WaitGroup

	Workers = make(chan bool, n)
	for i := 0; i < n; i++ {
		done.Add(1)
		start.Add(1)
		go worker(command, &start, &done)
	}
	start.Wait()

	done.Add(1)
	go reader(os.Stdin, &done)

	done.Wait()
	verboseLog("finished")
}

func reader(src io.ReadCloser, done *sync.WaitGroup) {
	defer done.Done()
	defer close(stdinCh)

	scanner := bufio.NewScanner(src)
	for scanner.Scan() {
		b := scanner.Text()
		if len(Workers) == 0 {
			verboseLog("all commands are unavaiable")
			return
		}
		stdinCh <- b
	}
	verboseLog("input finished")
}

func worker(command []string, start, done *sync.WaitGroup) {
	Workers <- true
	defer func() {
		<-Workers
	}()
	defer done.Done()

	verboseLog("invoking command", strings.Join(command, " "))
	var cmd *exec.Cmd
	if len(command) == 1 {
		cmd = exec.Command(command[0])
	} else {
		cmd = exec.Command(command[0], command[1:len(command)]...)
	}

	_stdin, err := cmd.StdinPipe()
	if err != nil {
		start.Done()
		errorLog(err)
		return
	}
	stdin := bufio.NewWriter(_stdin)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		errorLog(err)
		return
	}
	start.Done()
	for {
		input, more := <-stdinCh
		if !more {
			verboseLog("worker done")
			break
		}
		_, err := fmt.Fprintln(stdin, input)
		if err != nil {
			errorLog("failed to write to STDIN", err)
			break
		}
	}
	stdin.Flush()
	_stdin.Close()
	cmd.Wait()
}

func verboseLog(args ...interface{}) {
	if !Verbose {
		return
	}
	mu.Lock()
	defer mu.Unlock()
	log.Println(args...)
}

func errorLog(args ...interface{}) {
	mu.Lock()
	defer mu.Unlock()
	_args := []interface{}{"[error]"}
	_args = append(_args, args...)
	log.Println(_args...)
}
