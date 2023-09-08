This is a [QOTD server](http://tools.ietf.org/html/rfc865), it is pretty
full featured.

> You can get *anywhere* in ten minutes if you drive fast enough.

There's also a sister project, [QOTD
Client](https://github.com/nixterrimus/qotd-client).  While netcat works
good for testing the client is more full featured.

## Try it out

- Install the server: `go get -u github.com/nixterrimus/qotd-server`
- Start the server: `qotd-server https://raw.githubusercontent.com/nixterrimus/qotd-server/master/wisdom.txt`
- In another terminal make a request: `nc localhost 3333`

## Specifying a File

The QOTD server needs a path to a file or an HTTP URL that contains
quotes, this argument is mandatory and should specified as the first
argument.  Here are some examples:

- `qotd-server funny.txt`
- `qotd-server`
https://raw.githubusercontent.com/nixterrimus/qotd-server/master/wisdom.txt

## Quote File Format

This QOTD server has a mandatory files argument of a quote file, for
example to run with the `wisdom.txt` file included in the git repository
you would use: `qotd-server wisdom.txt`.

The quote file format is exactly the same as the [fortune file
format](http://en.wikipedia.org/wiki/Fortune_(Unix)#Fortune_files), that
is to say:

- It is a text file
- Each quote is separated by a `%` character of its own line

## The QOTD Protocol / RFC 865

RFC 865 defines the [Quote of the Day
Protocol](http://tools.ietf.org/html/rfc865). The specification is
really short.  Here's the gist of it:

- The Server listens on TCP (port 17, by convention)
- The Server may also listen on UDP (also port 17, by convention)
- On connection a quote is served
- The connection is closed immediately after a quote is serve
- Quotes should be less than 512 characters long

This server **can** be started RFC 865 compliant mode with the
`--strict` flag.

## Strict Mode

This QOTD server can be started in strict mode, it must be run as a super
user because it binds to port 17.  This ensures that the server is RFC
865 compliant.

Strict mode does the following things:

- Starts on port 17
- Listens on both TCP and UDP
- Serves quotes that are 512 characters or less

## TCP and UDP

By default the server listens on both TCP and UDP on the same port.  If
that doesn't fit your needs, you have a couple of flags:

- `no-tcp` - Starts the server and only listens on UDP
- `no-udp` - Starts the server and only listens on TCP

You can test the TCP interface with netcat: `nc localhost 3333` And you can test on UDP
as well: `echo -n " " | nc -4u -w1 localhost 3333`.

## IPv6

IPv6 is supported on operating systems that support it.  You'll be
pleased to know that it worked nicely.

## Interface

By default `0.0.0.0` but not forever.

## Service Discovery

The QOTD service is, by default advertised over
[MDNS](http://en.wikipedia.org/wiki/Multicast_DNS).  This makes it
possible for clients to find the server with zero config, pretty neat!

If that doesn't work for you, can specify `--no-mdns` to turn off
advertising the service.

Since some networks restrict publishing multicast DNS it's possible that
this feature just won't work on your network.  . Notably, multicast cannot be used 
in any sort of cloud, or shared infrastructure environment.

## Next Steps / Project Goals

- [X] Accept file as CLI argument
- [ ] Accept Quotes via Standard In
- [X] Accept a URL for a quote file
- [X] Accept port as CLI argument
- [X] Server should run in 865 compliant or non-compliant mode (512
character truncation)
- [X] Server should listen on UDP
- [ ] Server should be well tested
- [ ] Installable through homebrew
- [ ] Advertise QOTD over MDNS
- [ ] Never serve clients the same quote twice

## Install through homebrew

Notes:
  - https://github.com/Homebrew/homebrew-binary
  - https://github.com/Homebrew/homebrew/pull/20895
  - https://github.com/Homebrew/homebrew/pull/21810/files
  - https://github.com/Homebrew/homebrew/issues/23703

## Author & Thanks

This was written by [Nick Rowe](http://dcxn.com) in the summer of
2014.

Thanks to Rafe Colton for encouraging me to get started with go and
providing guidance on getting started.

Thanks to Parker Moore for hacking on UDP support.
addition
.
.
.
.
.
.