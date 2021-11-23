# TODO Service

A simple service to add to-do task with simple options like: title, comments, labels, due date, and done.

_Written for technical assigment purposes._

## Run

for running the service you can run:

```
make run
```

## Deployment
To deploy the service on k8s, you should make container image and push it to a container registry.

## Environment variables

|   Environment Name    |       Description         | Required  |  Default      |
|-----------------------|---------------------------|---------- |---------------|
| HOST                  | Application host          | false     | 127.0.0.1     |
| PORT                  | Application port          | false     | 6666          |
| ENV                   | Environment name          | false     | development   |
| JAEGER_ENDPOINT       | Tracing endpoint          | false     | empty         |

## Test

For test purposes you may use grpcui, it can be installed by running:

```
go install github.com/fullstorydev/grpcui/cmd/grpcui@latest
```

This installs the command into the bin sub-folder of wherever your $GOPATH environment variable points. If this
directory is already in your $PATH, then you should be good to go.

make sure to run the service on `development` environment for testing. then run:
```
grpcui -plaintext 127.0.0.1:6666
```
it will show you a URL to put into a browser in order to access the web UI.

## Want to memory profile the code? Please be aware!
If you are going to profile the code please be aware that you can not get a memory profile in production environment.
That's because the memory profiling functionality is disabled in this environment by simply setting the variable of
[runtime.MemProfileRate](https://golang.org/pkg/runtime/#pkg-variables) to zero which makes it disable.

The reason behind disabling it, was getting rid of 179999 automatic `int` allocation in runtime. The memory profiler
samples heap allocations. It will show function calls allocations. Recording all allocation and unwinding the stack
trace would be expensive, therefore a sampling technique is used. The sampling process relies on a pseudo random number
generator based on [exponential distribution](https://en.wikipedia.org/wiki/Exponential_distribution) to sample only a
fraction of allocations. The generated numbers define the distance between samples in terms of allocated memory size.
This means that only allocations that cross the next random sampling point will be sampled. The sampling rate defines
the mean of exponential distribution . The default value is 512 KB which is specified by the `runtime.MemProfileRate`.
So the memory allocation profiler is always active. Setting the `runtime.MemProfileRate` to 0, will turn off the memory
sampling entirely and Also saving ~1.4 MB in 64bit machines.