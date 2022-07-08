package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"
	
	gofathom "github.com/bbars/go-fathom"
	"github.com/bbars/limitedpool"
)

func corsMiddleware(w http.ResponseWriter, r *http.Request, allowOrigin string, next http.HandlerFunc) {
	w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
	if r.Method == http.MethodOptions {
		return
	}
	next(w, r)
}

func main() {
	var listen string
	var tbDir string
	var poolSize uint
	var allowOrigin string
	var maxTime time.Duration
	var help bool
	flag.StringVar(&listen, "listen", ":80", "HTTP listen [host]:port")
	flag.StringVar(&tbDir, "tbDir", "./tables", "Path to the directory containing Tablebase files")
	flag.UintVar(&poolSize, "poolSize", 4, "Pool size of concurrent Fathom instances")
	flag.StringVar(&allowOrigin, "allowOrigin", "*", "Value for HTTP header Access-Control-Allow-Origin")
	flag.DurationVar(&maxTime, "maxTime", 0, "Max time limit")
	flag.BoolVar(&help, "help", false, "Show usage info")
	
	flag.Parse()
	if help {
		flag.Usage()
		return
	}
	
	log.Println("listen", listen)
	log.Println("tbDir", tbDir)
	log.Println("poolSize", poolSize)
	log.Println("allowOrigin", allowOrigin)
	log.Println("maxTime", maxTime)
	
	fathomPool := limitedpool.New(int(poolSize), func() interface{} {
		res, err := gofathom.NewFathom(tbDir)
		if err != nil {
			panic(err)
		}
		return res
	})
	// test Fathom:
	fathomPool.Put(fathomPool.Get(context.TODO()))
	
	httpHandlers := HttpHandlers{
		ctx: context.Background(),
		fathomPool: fathomPool,
		maxTime: maxTime,
	}
	
	http.HandleFunc("/wdl", func (w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r, allowOrigin, httpHandlers.HandleWDL)
	})
	http.HandleFunc("/root", func (w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r, allowOrigin, httpHandlers.HandleRoot)
	})
	http.HandleFunc("/root-dtz", func (w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r, allowOrigin, httpHandlers.HandleRootDTZ)
	})
	http.HandleFunc("/root-wdl", func (w http.ResponseWriter, r *http.Request) {
		corsMiddleware(w, r, allowOrigin, httpHandlers.HandleRootWDL)
	})
	
	log.Fatal(http.ListenAndServe(listen, nil))
}
