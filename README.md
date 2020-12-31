# Tweetrod (still in heavy development)

> twitter's analytics API is expensive, this lib extracts impressions and engagements analytics from a profile or tweet

### building

make sure you have a recent version of go

```sh
$ go install https://github.com/navicstein/tweetrod
```

### usage

for viewing the available options

```sh
❯ tweetrod --help
Usage of tweetrod:
  -password string
    	your twitter password (default "xxxx")
  -use-proxy
    	whether to use reverse proxy at http://127.0.0.1:8080
  -username string
    	your twitter username (default "navicstein")

```

start scapping

```sh
tweetrod --use-proxy --username my_twitter_username --password 1292883
```

### enabling tracing

create a `.rod` file with these contents in `nano $(echo $PWD)/.rod ` or any other code editor you like

```txt
trace
show
```

- `show` enables the emulated browser
- `trace` enables logging

tweetrod will attempt to use a diffrent emulation device anytime the command is run
```
IPhoneX
IPad
GalaxyNoteII
IPadPro
```


## known bugs/issues 
- no error reporting on failed login creds (causes other things to fail)
- can't harvest engagements due to changing anchor tags 💡 (work around available)
- can't change IP address from proxy 💡 (work around available)
- and more

> my IP address was blocked because of agressive testing while developing this lib (where's my VPN?), can't ensure it consitency for now
