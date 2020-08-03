package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/iovns/x/configuration"
	"time"
)

type Configuration struct {
	fees *configuration.Fees
	conf *configuration.Config
}

func NewConfiguration(fees *configuration.Fees, conf *configuration.Config) Configuration {
	return Configuration{
		fees: fees,
		conf: conf,
	}
}
func (c Configuration) GetFees(_ sdk.Context) *configuration.Fees {
	return c.fees
}

func (c Configuration) GetConfiguration(_ sdk.Context) configuration.Config {
	return *c.conf
}

func (c Configuration) IsOwner(_ sdk.Context, addr sdk.AccAddress) bool {
	return c.conf.Configurer.Equals(addr)
}

func (c Configuration) GetValidDomainNameRegexp(_ sdk.Context) string {
	return c.conf.ValidDomainName
}

func (c Configuration) GetDomainRenewDuration(_ sdk.Context) time.Duration {
	return c.conf.DomainRenewalPeriod
}

func (c Configuration) GetDomainGracePeriod(c_ sdk.Context) time.Duration {
	return c.conf.DomainGracePeriod
}
