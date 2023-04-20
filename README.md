# hakarungo
A go-based time tracker

It keeps records of the time when one of the files tracked by git is modified and sums up the duration (currently, the maximum is 5 minutes) as a work time.

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
$ hakarungo hakaru <path-to-directory1> <path-to-directory2> ...
```

### Example

```
$ hakarungo hakaru . ~/wd/direcotry2 ../directory3
2023/04/15 20:00:14 modified file:  app/models/hoge.rb
2023/04/15 20:00:14 work time:  0s
2023/04/15 20:00:15 modified file:  app/services/fuga_service.rb
2023/04/15 20:00:15 work time:  1.133917784s
2023/04/15 20:08:13 modified file:  app/services/fuga_service.rb
2023/04/15 20:08:13 work time:  5m1.133917784s
2023/04/15 20:08:14 modified file:  app/models/hoge.rb
2023/04/15 20:08:14 work time:  5m2.734506335s
2023/04/15 20:09:26 modified file:  app/models/piyo.rb
2023/04/15 20:09:26 work time:  6m14.632756347s
2023/04/15 20:09:27 modified file:  app/services/fuga_service.rb
2023/04/15 20:09:27 work time:  6m15.273970107s
^C2023/04/15 20:09:30 Signal accepted: interrupt # you can terminate tracker by Ctrl + C
2023/04/15 20:09:30 Directories is /path/to/directory1, /path/to/directory2, /path/to/directory3
2023/04/15 20:09:30 Working time is 6m15.273970107s
```
