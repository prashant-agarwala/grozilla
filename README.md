# Grozilla
The Grozilla is a simple implementation that allows downloading of video,audio,package or zip files parallely and
efficiently using light weight go routines.

### Usage:

```
grozilla [-m] [-n] [-r] [-t] [-N] [download link]
```

### Installation:
```
  $ go get github.com/gophergala2016/grozilla
```
This will download grozilla to $GOPATH/src/github.com/gophergala2016/grozilla. In this directory run go build to create the grozilla binary.

### Description:

The utility allows to parallely download any downloadble file from the download link. Following are the customization flags which a client can
give

``` -n routines ```
	Used to specify number of go routines (default 10).

``` -r ```
	Used to resume pending download (which was stopped due to sudden exit).

``` -t time ```
	Used to specify maximum time in seconds it will wait to establish a connection (default 900).

``` -m attempts ```
	Used to specify maximum attempt to establish a connection (default 1).

``` -N nolimit ```
	Used to override maximum connection limit of 20


### Example:

![Grozilla-image](https://github.com/gophergala2016/grozilla/blob/master/screenshot_grozilla.jpg "grozilla")

### Coming Soon

- Smooth UI
- Additional Flags
	- for output filename
	- additional log messages
	- header info
	- retry delay
	- and some more

### References

- https://github.com/cheggaaa/pb
