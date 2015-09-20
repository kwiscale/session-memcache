/*
This package provides memcached support for session on Kwiscale Framework.

To use it, you need to import the package in your application to let runtime call the "init()" function.

	import _ "gopkg.in/kwiscale/session-memcached"

Kwiscale will now be able to use session-memcached it you set your App configuration to "memcached":

	app := kwiscale.NewApp(&kwiscale.Config{
		SessionEgnine : "memcached",
		SessionOptions: SessionOptions{
			"servers" : "192.168.1.5:11211"
		}
	})

The "servers" options could be a coma separated list in string or a slices of strings:

	app := kwiscale.NewApp(&kwiscale.Config{
		SessionEgnine : "memcached",
		SessionOptions: SessionOptions{
			"servers" : []string{
				"192.168.1.5:11211",
				"192.168.1.6:11211",
				"192.168.1.7:11211",
				"192.168.1.8:11211",
			}
		}
	})

*/
package kwiscalesessionmemcached

import (
	"errors"
	"strings"

	"github.com/bradfitz/gomemcache/memcache"
	gsm "github.com/bradleypeabody/gorilla-sessions-memcache"
	"gopkg.in/kwiscale/framework.v0" // github.com/kwiscale/framework
)

func init() {
	kwiscale.RegisterSessionEngine("memcached", &MemcachedSessionEngine{})
}

// Implement ISessionStore
type MemcachedSessionEngine struct {
	memcache *memcache.Client
	store    *gsm.MemcacheStore
	secret   []byte
	name     string
	prefix   string
}

// Register session-memcached in kwiscale.
func (mc *MemcachedSessionEngine) Init() {
	mc.store = gsm.NewMemcacheStore(mc.memcache, mc.prefix, mc.secret)
}

// Name set the name of session (ie. the session cookie name too).
func (mc *MemcachedSessionEngine) Name(name string) {
	mc.name = name
}

// SetOptions is used here to set servers addresses.
// Example:
//
//		app := kwiscale.NewApp(&kwiscale.Config{
//			SessionEngine: "memcached",
//			SessionengineOptions: SessionEngineOptions{
//				"servers" : "192.168.1.1:11211,192.168.1.2:11211", // required
//				"prefix"  : "prefix_string", // optionnal
//			}
//		})
//
// The "servers" list could be a []string or a coma separated list.
func (mc *MemcachedSessionEngine) SetOptions(options kwiscale.SessionEngineOptions) {
	if servers, ok := options["servers"]; ok {
		switch servers := servers.(type) {
		case string:
			mc.setServers(strings.Split(servers, ",")...)
		case []string:
			mc.setServers(servers...)
		}
	}

	if prefix, ok := options["prefix"]; ok {
		mc.prefix, _ = prefix.(string)
	}
}

// Set secret []byte to encode session.
func (mc *MemcachedSessionEngine) SetSecret(s []byte) {
	mc.secret = s
}

// Get a value from the session.
func (mc *MemcachedSessionEngine) Get(handler kwiscale.IBaseHandler, key interface{}) (interface{}, error) {
	session, err := mc.store.Get(handler.GetRequest(), mc.name)
	if err != nil {
		return nil, err
	}

	if val, ok := session.Values[key]; ok {
		return val, nil
	}

	return nil, errors.New("empty session")
}

// Set value to the session.
func (mc *MemcachedSessionEngine) Set(handler kwiscale.IBaseHandler, key, value interface{}) {
	session, _ := mc.store.Get(handler.GetRequest(), mc.name)
	session.Values[key] = value
	session.Save(handler.GetRequest(), handler.GetResponse())
}

// Clean removes the entire values in session.
func (mc *MemcachedSessionEngine) Clean(handler kwiscale.IBaseHandler) {
	session, _ := mc.store.Get(handler.GetRequest(), mc.name)
	session.Values = make(map[interface{}]interface{})
	session.Save(handler.GetRequest(), handler.GetResponse())
}

// SetServers is called by "Init" method to create a memcache client used by the store.
func (mc *MemcachedSessionEngine) setServers(servers ...string) {
	mc.memcache = memcache.New(servers...)
}
