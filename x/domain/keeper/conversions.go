package keeper

// getDomainPrefixKey returns the domain prefix byte key
func getDomainPrefixKey(domainName string) []byte {
	return []byte(domainName)
}

// getAccountKey returns the account byte key by its name
func getAccountKey(accountName string) []byte {
	return []byte(accountName)
}

// accountKeyToString converts account key bytes to string
func accountKeyToString(accountKeyBytes []byte) string {
	return string(accountKeyBytes)
}
