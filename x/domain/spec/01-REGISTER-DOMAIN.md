# Register domain spec

### Closed domain
**prereqs:**

- `domain.Name` must not exist already
- `time(CurrentBlock)` < `domain.gracePeriodUntil`
- anyone is allowed to `domain.Name`
- `domain.Name` must match the regexp stored in `configuration.ValidDomain`

use case (1 signature, potentially 2):

- user register-domain
    - message might set the `domain.Admin` to another address different from the user address, default to user
    - message might set a "broker" name stored in the domain name, default to "".
    - message might set "fee payer" to another address, default to user
    - user signs here
- if fee payer is different from user, the fee payer have to sign the transaction as well
    - fee payer signs here

**fees:**

- from product fee for register-domain

**after:**

- the `Domain` is put inside the KVStore
- `domain.ValidUntil` becomes `time(CurrentBlock)` + `configuration.DomainRenewalPeriod`
- `domain.gracePeriodUntil` becomes `domain.ValidUntil` + `configuration.DomainGracePeriod`
- An `Account` associated with the domain is registered with an empty `account.Name` equal to `""`
    - `account.Owner` is `domain.Admin`
    - `account.ValidUntil` is max unixtime

### Open domain

use case (1 signature, potentially 2):

- user register-domain
    - message must set `domain.`Type **== "open"**

**after:**

- for the newly created empty account, `account.ValidUntil` is `domain.ValidUntil`

## Controller
- Validates the domain name against `conf.ValidDomainName`
- Validates that domain does not exist

