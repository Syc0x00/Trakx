# trakx [![Go Report Card](https://godoc.org/github.com/crimist/trakx?status.svg)](https://godoc.org/github.com/crimist/trakx) [![Go Report Card](https://goreportcard.com/badge/github.com/crimist/trakx)](https://goreportcard.com/report/github.com/crimist/trakx)

Fast bittorrent tracker

## Performance

Here's an HTTP tracker running on heroku free tier with the `fast` tag disabled.

![performance](img/performance.png)
![performance](img/stats.png)
![flame](img/flame.png)

As demonstrated by the flame graph almost all of the CPU usage is from handling the TCP connections. Here the databases save function made only 0.3% of the flame graph cpu time.

Memory wise the go GC only runs every 2 min ([the maximum GC period](https://github.com/golang/go/blob/895b7c85addfffe19b66d8ca71c31799d6e55990/src/runtime/proc.go#L4481-L4486)) at this level of traffic. The `inuse_space` delta between right after vs right before GC is 7.5%, basically this collection frequency would be sustained at `GOGC=8`.

## Install

go 1.13+ recommended for `sync.Pool` and `sync.RMutex` optimizations.

### Quick Install

```sh
git clone github.com/crimist/trakx
cd trakx

# install to go bin
go install .
trakx start # starting...

# or you can build it
go build .
./trakx start # starting...
```

### Modifying Config & HTML

The config can be updated at `~/.config/trakx/trakx.yaml`. Note you'll have to run trakx at least once to generate this file.

If you want to change the index or dmca HTML pages you can change the files in the `install/` folder and than rebuild / reinstall.

### Updating

You can simply pull and build / install to update:

```sh
git pull
go install .
trakx restart
```

### Netdata graph install

**Warning:** `install.sh` will overwrite `go_expvar.conf`. If you have other expvar programs in netdata you can manually merge the two files.

* Run `/etc/netdata/edit-config python.d.conf` and change the `go_expvar` setting to to `yes`
* Customize the url in `netdata/expvar.conf` if needed
* Install netdata plugins with `cd netdata; ./install.sh`

### Build Tags

You can build with different tags by using `go build/install -tags <tag> .`

* `fast` tag will build without IP, seeds, and leeches metrics which will reduce cpu and memory usage
* `heroku` tag will build the service for app engines, this means that when executed the binary will immediately run the tracker rather than the controller

## Notes

* If you're going to be serving a lot of clients on a non managed service take a look at [sysctl tuning](https://wiki.mikejung.biz/Sysctl_tweaks). This is especially important if you're running a TCP tracker
* Technically there's no guarantee that database saves work between go versions. By default I use `unsafe` to read raw memory so if they change struct padding or the internal slice structure it could break your save between versions (though this should never happen). You can change the encoding method to `encodeBinary()` to avoid this issue but it takes 3x more memory and is 7x slower. 
