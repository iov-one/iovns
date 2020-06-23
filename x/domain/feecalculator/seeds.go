package feecalculator

const (
	DefaultSeedSuffix      = "default"
	DomainOpenSeedSuffix   = "open"
	DomainClosedSeedSuffix = "closed"
)

func buildSeedID(args ...string) string {
	var str string
	for _, a := range args {
		str = a + "_"
	}
	return str
}
