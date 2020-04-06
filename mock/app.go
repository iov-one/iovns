package mock

import (
	"github.com/iov-one/iovnsd/app"
	"github.com/tendermint/tendermint/libs/log"
)
import kv "github.com/tendermint/tm-db"

// App is an application mock, it generates a mock application, exposing
// the keepers in order to make them accessible from different test cases.
// It also exposes applications public methods and fields.
type App struct {
	*app.NewApp
}

// NewApp is the constructor for App
func NewApp() App {
	db := kv.NewMemDB()
	newApp := app.NewInitApp(log.NewNopLogger(), db, nil, true, 0, nil)
	return App{newApp}
}
