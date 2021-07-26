package main

import (
	"golang.org/x/net/context/ctxhttp"
	"golang.org/x/net/context" 
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/aws/aws-xray-sdk-go/xray"

	"github.com/fasthttp/router"
    "github.com/valyala/fasthttp"
)

func init() {
	xray.Configure(xray.Config{
		DaemonAddr:     "xray-service.default:2000",
		LogLevel:       "info",
	})
}


func main() {
    fh := xray.NewFastHTTPInstrumentor(nil)
    r := router.New()
    
    r.GET("/{trace?}", middleware("x-ray-sample-middle-k8s", middle, fh))

    fasthttp.ListenAndServe(":8080", r.Handler)
}

func middleware(name string, h fasthttp.RequestHandler, fh xray.FastHTTPHandler) fasthttp.RequestHandler {
    f := func(ctx *fasthttp.RequestCtx) {
        h(ctx)
    }
    return fh.Handler(xray.NewFixedSegmentNamer(name), f)
}

func middle(ctx *fasthttp.RequestCtx){
	trace := ctx.UserValue("trace")

	if(trace == "all"){
		url := `http://x-ray-sample-back-k8s.default.svc.cluster.local`
	    resp,err := traceUrl(ctx,url)
		if err != nil {
	        fmt.Println("请求失败:", err.Error())
	        return
	    }
	    fmt.Fprintf(ctx,"Middle Server->"+string(resp))

	}else{
		fmt.Fprintf(ctx,"Middle Server")
	}
}

var tr = &http.Transport{
	MaxIdleConns: 20,
	IdleConnTimeout: 30 * time.Second,
}

func traceUrl(ctx context.Context,url string) ([]byte, error) {
    resp, err := ctxhttp.Get(ctx, xray.Client(&http.Client{Transport: tr}), url)
    if err != nil {
      return nil, err
    }
    return ioutil.ReadAll(resp.Body)
}