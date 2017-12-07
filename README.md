Command unleash runs multiple copies of child command until they all finish.

	Usage: unleash [flags] -- child-program [child args]
	  -n int
		number of child processes to start (defaults to number of CPUs) (default 4)
	  -r int
		max number of times to restart child if it fails (until any child exits with 0)
