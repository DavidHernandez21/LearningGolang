package main

import (
	"flag"
	"fmt"
	stdlog "log"
	"net/http"
	"os"
	"os/signal"
	myloggy "src/github.com/DavidHernandez21/src/github.com/DavidHernandez21/justForfunc/context/mylog"
	"syscall"
	"time"
	// "github.com/campoy/justforfunc/09-context/log.
)

func main() {
	port := flag.Int("p", 5450, "port to listen to")
	flag.Parse()
	http.HandleFunc("/", myloggy.Decorate(handler))
	stdlog.Printf("Listening on port %v", *port)
	stopSignalHandler()
	stdlog.Fatal(http.ListenAndServe(fmt.Sprintf("127.0.0.1:%v", *port), nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	myloggy.Println(ctx, "handler started")
	defer myloggy.Println(ctx, "handler ended")

	select {
	case <-time.After(1 * time.Second):
		fmt.Fprintln(w, "hello")
	case <-ctx.Done():
		err := ctx.Err()
		myloggy.Println(ctx, err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		os.Exit(2)

	}
}

func stopSignalHandler() {

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		sig := <-signalChan
		stdlog.Printf("- Ctrl+C pressed, exiting\n Signal recieved: %v\n", sig)
		os.Exit(0)
	}()

}
