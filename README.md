# Grozilla
The Grozilla is a simple implementation that allows downloading of video,audio,package or zip files parallely and
efficiently using light weight go routines.

### Usage:

```
grozilla [-m] [-n] [-r] [-t]  [download link]
```

### Installation:

- Install go package :
    ```
    $ go get github.com/cheggaaa/pb
    ```

- Build the project using :
  ```
    $ go build
  ```

### Description:

The utility allows to parallely download any downloadble file from the download link. Following are the customization flags which a client can
give

``` -m attempts ```
	Used to specify maximum attempt to establish a connection (default 1).

``` -n routines ```
	Used to specify number of go routines (default 10).

``` -r ```
	Used to resume pending download (which was stopped due to sudden exit).

``` -t time ```
	Used to specify maximum time in seconds it will wait to establish a connection (default 900).

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
