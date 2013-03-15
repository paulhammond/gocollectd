# gocollectd

This is an implementation of the [Collectd Binary Protocol][protocol] in [Go][go].

## Installation

Install with `go get`:

    go get github.com/paulhammond/gocollectd

You probably want to import gocollectd into your code with a different name:

    import (
      collectd "github.com/paulhammond/gocollectd"
    )
    // now you can call collectd.Parse() etc

## Usage

In short, gocollectd will parse a binary stream of data and return values.
However there are some details.

Some collectd plugins send multiple values for the same timestamp at once. For
example, the [interface plugin][interface] sends both `tx` and `rx` in one
packet. Some applications will want to mirror the behavior of collectd's RRD
plugin and store both values in one file. To allow that, parsing data returns
Packets instead of Values.

    b := []byte{…}
    packets := collectd.Parse(b)

    packet = packets[0]
    fmt.Println(packet.Hostname)       // "laptop.lan"
    fmt.Println(packet.Time)           // A go time value
    fmt.Println(packet.Plugin)         // "Load"
    fmt.Println(packet.ValueCount)     // 3
    fmt.Println(packet.ValueNames())   // { "load1", "load5", "load15" }

Collectd values are sent as one of the RRD types: `Counter`, `Guage`,
`Derive` or `Absolute`. This, in turn, means that they are sent as `int`s,
`uint`s or `float`s. You have a few options on how to handle this:

    // just get the bytes
    fmt.Println(packet.ValueBytes())   // [][]byte{ … }

    // or use the Number interface type
    numbers := packet.ValueNumbers()
    fmt.Println(numbers) // { 1.13, 0.89, 0.60 }

    // if you don't care about exact precision because you're about to average
    // lots of numbers, you can convert everything to a float:
    f := numbers[0].Float64()  // f == float64(1.13)

    // or you can use go type assertions
    if i, ok := numbers[0].(int64); ok {
       // do something with i, knowing it's an int64
    }

The most common use case is to read collectd data directly from the network.
A simple server implementation is provided that sends received packets on a
channel:

    c := make(chan collectd.Packet)
    go collectd.Listen("127.0.0.1:25827", c)
    for {
      packet := <-c
      // do something with the packet
    }

An example is of using this server is provided in
[gocollectd-example](gocollectd-example/example.go)

## Known issues

Signed and encrypted collectd packets are not currently supported.

## References

  * [go][go]
  * [Collectd][collectd]
  * [Collectd Binary Protocol][collectd]

[go]: http://golang.org/
[collectd]: https://collectd.org/
[protocol]: https://collectd.org/wiki/index.php/Binary_protocol
[interface]: https://collectd.org/wiki/index.php/Plugin:Interface

## License

Copyright (c) 2013 Paul Hammond. gocollectd is available under the MIT
license, see [LICENSE.txt](LICENSE.txt) for details