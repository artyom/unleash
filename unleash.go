// Command unleash runs multiple copies of child command until they all finish.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/artyom/autoflags"
	"golang.org/x/sync/errgroup"
)

func main() {
	args := struct {
		N int `flag:"n,number of child processes to start (defaults to number of CPUs)"`
		R int `flag:"r,max number of times to restart child if it fails (until any child exits with 0)"`
	}{N: runtime.NumCPU()}
	autoflags.Parse(&args)
	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	if err := run(args.N, args.R, flag.Args()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(n, restart int, args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("nothing to run")
	}
	if n < 1 {
		n = 1
	}
	if restart < 0 {
		restart = 0
	}
	g, ctx := errgroup.WithContext(context.Background())
	for ; n > 0; n-- {
		g.Go(func() error {
			var err error
			for i := 0; i < restart+1; i++ {
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Stdout = os.Stdout // output of multiple childs may be interleaved
				cmd.Stderr = os.Stderr
				if err = cmd.Run(); err == nil {
					return nil
				}
				select {
				case <-ctx.Done():
					return err
				case <-time.After(time.Second):
				}
			}
			return err
		})
	}
	return g.Wait()
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [flags] -- child-program [child args]\n", filepath.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}
