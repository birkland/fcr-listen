# Simple Fedora Messaging Client

Connects via STOMP to a topic or queue, and displays message contents.  This is primarily intended as a diagnostic or educational tool

## Quick Start

Download the executable for your OS from the [releases page](https://github.com/birkland/fcr-listen/releases)

-_or_-

If you have `go` installed, you may fetch and build via

    go get -u github.com/birkland/fcr-listen

If you run with no arguments, it will attempt to connect to a STOMP messaging enndpoint on localhost (port 61613).  This will connect to a locally-running Fedora configured with default settings.

    fcr-listen
    
## Usage

Use the `-h` or `--help` flag to see the options and defaults:

<pre>
$ fcr-listen -h
Usage of /path/to/fcr-listen:
  -host string
        STOMP connection host or IP (default "localhost")
  -port int
        STOMP connection port (default 61613)
  -subscribe string
        Queue or topic to subscribe to (default "/topic/fedora")
</pre>
    
Override whatever defaults you wish

If no connection is possibe (i.e. Fedora or its messaging bus are down), it will try to connect every three seconds in an infinite loop, printing to STDERR each time:
<pre>
$ fcr-listen                                                                                                                                                                                    
2017/10/12 15:07:08 main.go:39: Could not connect! dial tcp [::1]:61613: connectex: No connection could be made because the target machine actively refused it.
2017/10/12 15:07:11 main.go:39: Could not connect! dial tcp [::1]:61613: connectex: No connection could be made because the target machine actively refused it.
</pre>

If successful, you'll see a logging message (sent to STDERR)
<pre>
# fcr-listen
2017/10/12 19:08:30 main.go:65: Subscribed to /topic/fedora
</pre>

For each message recieved, the following will be printed to STDOUT
* Message headers, with keys and values separated by a single space
  * When run on a tty supporting color, these will be printed in cyan.  Headers that begin with `org.fcrepo` will be in bold
* Message body
  * When run on a tty supporting color, this will be printed in orange

Fedora message bodies are JSON, printed on a single line.  The output can be piped to other tools for further enhancement, such as [jq](https://stedolan.github.io/jq/)

## Examples

To ignore headers and pretty print the JSON bodies of message:

    fcr-listen | grep '{' | jq .
