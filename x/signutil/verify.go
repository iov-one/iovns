package signutil

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

func Verify(tx auth.StdTx, chainID string, accountNumber, sequence uint64) error {
	signatures := tx.GetSignatures()
	signers := tx.GetPubKeys()
	if len(signatures) == 0 {
		return fmt.Errorf("at least one signature must be present")
	}
	if len(signers) != len(signatures) {
		return fmt.Errorf("invalid number of signers (%d) and signatures (%d)", len(signers), len(signatures))
	}
	for i, sig := range signatures {
		signer := signers[i]
		message := auth.StdSignBytes(chainID, accountNumber, sequence, tx.Fee, tx.Msgs, tx.Memo)
		if !signer.VerifyBytes(message, sig) {
			return fmt.Errorf("invalid signature from address found at index %d, from address: %s", i, tx.GetSigners()[i])
		}
	}
	return nil
}
