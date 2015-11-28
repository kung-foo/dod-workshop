package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	log "github.com/Sirupsen/logrus"
	"github.com/docopt/docopt-go"
	"github.com/mediocregopher/radix.v2/pool"
	"github.com/mediocregopher/radix.v2/redis"
)

type Cluster struct {
	sync.Mutex
	master *pool.Pool
	slaves []*pool.Pool
}

type Config struct {
	Master string
	Slaves []string
}

var usage = `
Reverse Redis Proxy
Usage:
    rrproxy [options] --master=<master>
    rrproxy [options] --config=<config>
    rrproxy -h | --help | --version

Options:
    -h, --help                    Show this screen.
    --version                     Show version.
    -p, --port=<port>             Listen on [default: 8888].
    -m, --master=<master>         Master redis instance.
    -c, --config=<config>         Configuration file.
`

const (
	// PoolSize is the number of redis connection to initially open
	PoolSize = 10

	// ForcedLatency is the amount of time we stall the redis server in order
	// to simulate a loaded box
	ForcedLatency = time.Duration(1 * time.Millisecond)
)

func main() {
	mainEx(os.Args[1:])
}

func mainEx(argv []string) {
	var err error
	args, err := docopt.Parse(usage, argv, true, "1.0", false)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	if args["--master"] != nil {
		config.Master = args["--master"].(string)
	}

	if args["--config"] != nil {
		if _, err := toml.DecodeFile(args["--config"].(string), &config); err != nil {
			log.Fatal(err)
		}

		if config.Master == "" {
			log.Fatalf("'Master' not specified in %v", args["--config"])
		}
	}

	log.Warn("This is a toy server for use the Oslo Day of Docker Workshop 2015. Do not use in production!")

	cluster, err := NewCluster(&config)
	if err != nil {
		log.Fatal(err)
	}

	// setup a basic ping response for testing
	http.HandleFunc("/db/_ping", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "OK")
	})

	// init the GET and POST handler
	http.Handle("/db/", http.StripPrefix("/db/", cluster))

	log.Infof("Starting server on *:%v", args["--port"])

	// start up this elegant bit of late-night, øl fueled engineering
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", args["--port"]), nil))
}

func NewCluster(config *Config) (*Cluster, error) {
	master, err := pool.New("tcp", fmt.Sprintf("%s:6379", config.Master), PoolSize)
	if err != nil {
		return nil, err
	}
	c := Cluster{
		master: master,
		slaves: make([]*pool.Pool, len(config.Slaves)),
	}

	for i, slave := range config.Slaves {
		conn, err := pool.New("tcp", fmt.Sprintf("%s:6379", slave), PoolSize)
		if err != nil {
			return nil, err
		}
		c.slaves[i] = conn
	}

	log.Infof("Connected 1 master and %d slaves", len(c.slaves))

	return &c, nil
}

func (c *Cluster) doGET(w http.ResponseWriter, r *http.Request) {
	var err error

	// path is already read because we are using http.StripPrefix
	path := r.URL.Path

	if path == "" {
		http.NotFound(w, r)
		return
	}

	slave := c.getReadSlave()

	// this makes it easy to debug with curl -v
	w.Header().Add("X-RRPROXY-SERVER", slave.Addr)

	// get a redis client from the pool
	client, err := slave.Get()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// make sure we return it
	defer slave.Put(client)

	var resp *redis.Resp

	// check to see if the key exists
	resp = client.Cmd("EXISTS", path)
	v, _ := resp.Int()

	if v == 0 {
		http.NotFound(w, r)
		return
	}

	// here we introduce some artifical latency. without it, and with redis
	// running on the same laptop as the server, other layers become the bottle
	// neck rather than redis, which kinda defeats the purpose of this lab
	if ForcedLatency > 0 {
		resp = client.Cmd("DEBUG", "sleep", ForcedLatency.Seconds())
		if resp.Err != nil {
			log.Fatal(resp.Err)
		}
	}

	// get the raw value as bytes and write it out
	resp = client.Cmd("GET", path)
	data, err := resp.Bytes()
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func (c *Cluster) doPOST(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	if path == "" {
		http.Error(w, "empty path not allowed", http.StatusBadRequest)
		return
	}

	// no streaming, just read everything into memory. yay for non-production
	// code!!!
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = c.master.Cmd("SET", path, data).Err
	if err != nil {
		log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	io.WriteString(w, fmt.Sprintf("%d", len(data)))
}

func (c *Cluster) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		c.doGET(w, r)
	case "POST":
		c.doPOST(w, r)
	default:
		http.Error(w, fmt.Sprintf("%s not supported", r.Method), http.StatusMethodNotAllowed)
	}
}

func (c *Cluster) getReadSlave() *pool.Pool {
	// special case for a single node system
	if len(c.slaves) == 0 {
		return c.master
	}

	// TODO: should the rand object be per goroutine? Otherwise we have a shared
	// lock in the global rand
	//r := rand.New(rand.NewSource(time.Now().UnixNano()))

	return c.slaves[rand.Int()%len(c.slaves)]
}
