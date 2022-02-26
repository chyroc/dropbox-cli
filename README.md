# dropbox-cli

## Installation

```shell
go install github.com/chyroc/dropbox-cli@latest
```

## Usage

### Download file

```shell
# download `/path` dir(recursive) or file to `local/`
# and `path/1.txt` to `local/1.txt`
dropbox-cli download path local
```

```shell
# download `/path` dir(recursive) or file to `local/path`
# and `path/1.txt` to `local/path/1.txt`
dropbox-cli download path local/
```

### Upload file

```shell
# upload `local` dir(recursive) or file to `path/`
# and `local/1.txt` to `path/1.txt`
dropbox-cli upload ./local path
```

```shell
# upload `local` dir(recursive) or `local` to `path/local/`
# and `local/1.txt` to `path/local/1.txt`
dropbox-cli upload ./local path/
```

## TODO

- [ ] download check hash