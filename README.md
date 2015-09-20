# Kwiscale Memcache Session Engine

This package provides a Session Engine to use Memcache. 

# Installation

Use gopkg.in url to get a versionned package:

    go get -u gopkg.in/kwiscale/session-memcached.v0

Then you can import package in you application:

```go
package main

import (
    "gopkg.in/kwiscale/framework"
    _ "gopkg.in/kwiscae/session-memcached"
)

func main(){
    // configuration
    cfg := &kwiscale.Config{
        SessionEngine: "memcached",
        SessionEngineOptions: SessionEngineOptions{
            "prefix": "optional prefix",
            "servers" : []string{
                "192.168.1.5:11211", // at least one server
            }
        }
    }
}
```

Note that the port is mandatory (at this time).


