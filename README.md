# hakarungo
A go-based time tracker

It keeps records of the time when one of the files in files tracked by git is modified and sums up the duration (currently, the maximum is 5 minutes) as a work time.

## Setup

Make sure `$GOPATH/bin` is added to `$PATH`.

```
$ git clone https://github.com/38tter/hakarungo
$ cd hakarungo
$ go install
$ go build
```

## How to use

```
$ hakarungo hakaru <path-to-directory>
```

### Example

```
$ hakarungo hakaru .
2023/04/15 19:37:44 modified file:  app/models/hoge.rb
2023/04/15 19:37:44 work time:  0s
2023/04/15 19:37:50 modified file:  app/services/fuga_service.rb
2023/04/15 19:37:50 work time:  5.979647051s
^C2023/04/15 19:38:03 Signal accepted: interrupt # you can terminate tracker with Ctrl+C
2023/04/15 19:38:03 Directories is .
2023/04/15 19:38:03 Working time is 5.979647051s
```
