package types

// Module names
const (
	// ModuleName is the name of the module
	ModuleName = "starname"
	// DomainStore key defines the store key used to store domains information
	DomainStoreKey = "starname"
	// RouterKey defines the path used to interact with the domain module
	RouterKey    = ModuleName
	QuerierAlias = "starname"
	// QuerierRoute defines the query path used to interact with the domain module
	QuerierRoute = ModuleName
	// DefaultParamSpace defines the key for the default param space
	DefaultParamSpace = ModuleName
)

// Events attribute keys
const (
	AttributeKeyDomainName  = "domain_name"
	AttributeKeyAccountName = "account_name"
	AttributeKeyDomainType  = "domain_type"
	AttributeKeyFeePaid     = "paid_fees"
	AttributeKeyFeePayer    = "fee_payer"
	AttributeKeyOwner       = "owner"

	AttributeKeyNewCertificate          = "new_certificate"
	AttributeKeyDeletedCertificate      = "deleted_certificate"
	AttributeKeyNewResources            = "new_resources"
	AttributeKeyNewMetadata             = "new_metadata"
	AttributeKeyTransferAccountNewOwner = "new_account_owner"
	AttributeKeyTransferAccountReset    = "transfer_account_reset"

	AttributeKeyTransferDomainNewOwner = "new_domain_owner"
	AttributeKeyTransferDomainFlag     = "transfer_domain_flag"
)
