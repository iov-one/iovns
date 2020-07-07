package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	keys2 "github.com/cosmos/cosmos-sdk/crypto/keys"
	"github.com/gorilla/mux"
	"github.com/iov-one/iovns/cmd/faucet/pkg"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
)

func main() {
	// setup configuration
	conf, err := pkg.NewConfiguration()
	if err != nil {
		log.Fatalf("configuration: %s", err)
	}
	// setup node
	node, err := rpchttp.New(conf.TendermintRPC, "/websocket")
	kb := keys2.NewInMemory()
	if err := kb.ImportPrivKey("faucet", conf.Armor, conf.Passphrase); err != nil {
		log.Fatalf("keybase: %v", err)
	}
	// setup tx manager
	txManager := pkg.NewTxManager(*conf, node).WithKeybase(kb)
	if err := txManager.Init(); err != nil {
		log.Fatalf("tx manager: %v", err)
	}

	// Wait for ListenAndServe goroutine to close.
	r := mux.NewRouter()
	faucet := pkg.NewFaucetHandler(txManager)
	r.Handle("/credit", faucet)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		return
	})
	server := &http.Server{Addr: conf.Port, Handler: r}

	go func() {
		log.Print("server started")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("http server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
