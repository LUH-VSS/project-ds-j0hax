# All-In-One Demonstration (Track 1)

This is a very simple demonstration of mapreduce-ish processing using *n* input files and *m* output files. These are read by corresponding *n* `ReaderWorker` goroutines and *m* `WriterWorker` goroutines.

## Usage

```console
$ ./project-ds-j0hax file_1.txt file_2.txt file_n.txt
```

Use *n* files as input. The *m* WriterWorkers are currently hard-coded.

## To-Do
- [ ] Prevent deadlocks when all readers have finished
- [ ] Sort Outputs
- [ ] Less jank: actually conform to a MapReduce pattern