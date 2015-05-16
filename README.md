# Twitter Domain Counter #

Counts the number of occurances of domain names in Twitter sample
stream. Also links to the tweets that contain a specfic domain.

## URL Exapnding ##

The following shorteners are expanded,

* [bit.ly](https://bit.ly)
* [goo.gl](https://goo.gl)

> t.co links are automatically expanded by twitter

## Building ##

Install a recent version of the Golang toolchain.

First compile the React modules. You may need to install
[react-tools](https://www.npmjs.com/package/react-tools) via `npm` for
this to work

```sh
$ jsx react/ res/data/react/ --extension jsx --no-cache-dir
```

`go get` [this](https://github.com/vasuman/go-res-pack) package that
is used to pack the resource files. Then from the root of this repo,

```sh
$ go generate ./res/
```

Finally, build the server,

```sh
$ go build ./cmd/TwitDomCount
```


