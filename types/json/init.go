package json

import (
	cconsensus "github.com/tendermint/tendermint/consensus"
	cryptoed25519 "github.com/tendermint/tendermint/crypto/ed25519"
	cryptosecp256k1 "github.com/tendermint/tendermint/crypto/secp256k1"
	cryptosr25519 "github.com/tendermint/tendermint/crypto/sr25519"
	pc "github.com/tendermint/tendermint/proto/tendermint/crypto"
	ctypes "github.com/tendermint/tendermint/types"
)

func init() {

	RegisterType(&cconsensus.NewRoundStepMessage{}, "tendermint/NewRoundStepMessage")
	RegisterType(&cconsensus.NewValidBlockMessage{}, "tendermint/NewValidBlockMessage")
	RegisterType(&cconsensus.ProposalMessage{}, "tendermint/Proposal")
	RegisterType(&cconsensus.ProposalPOLMessage{}, "tendermint/ProposalPOL")
	RegisterType(&cconsensus.BlockPartMessage{}, "tendermint/BlockPart")
	RegisterType(&cconsensus.VoteMessage{}, "tendermint/Vote")
	RegisterType(&cconsensus.HasVoteMessage{}, "tendermint/HasVote")
	RegisterType(&cconsensus.VoteSetMaj23Message{}, "tendermint/VoteSetMaj23")
	RegisterType(&cconsensus.VoteSetBitsMessage{}, "tendermint/VoteSetBits")

	// ed25519
	RegisterType(cryptoed25519.PubKey{}, cryptoed25519.PubKeyName)
	RegisterType(cryptoed25519.PrivKey{}, cryptoed25519.PrivKeyName)

	// secp256k1
	RegisterType(cryptosecp256k1.PubKey{}, cryptosecp256k1.PubKeyName)
	RegisterType(cryptosecp256k1.PrivKey{}, cryptosecp256k1.PrivKeyName)

	//sr25519
	RegisterType(cryptosr25519.PubKey{}, cryptosr25519.PubKeyName)
	RegisterType(cryptosr25519.PrivKey{}, cryptosr25519.PrivKeyName)

	//
	RegisterType((*pc.PublicKey)(nil), "tendermint.crypto.PublicKey")
	RegisterType((*pc.PublicKey_Ed25519)(nil), "tendermint.crypto.PublicKey_Ed25519")
	RegisterType((*pc.PublicKey_Secp256K1)(nil), "tendermint.crypto.PublicKey_Secp256K1")

	// event

	RegisterType(ctypes.EventDataNewBlock{}, "tendermint/event/NewBlock")
	RegisterType(ctypes.EventDataNewBlockHeader{}, "tendermint/event/NewBlockHeader")
	RegisterType(ctypes.EventDataNewEvidence{}, "tendermint/event/NewEvidence")
	RegisterType(ctypes.EventDataTx{}, "tendermint/event/Tx")
	RegisterType(ctypes.EventDataRoundState{}, "tendermint/event/RoundState")
	RegisterType(ctypes.EventDataNewRound{}, "tendermint/event/NewRound")
	RegisterType(ctypes.EventDataCompleteProposal{}, "tendermint/event/CompleteProposal")
	RegisterType(ctypes.EventDataVote{}, "tendermint/event/Vote")
	RegisterType(ctypes.EventDataValidatorSetUpdates{}, "tendermint/event/ValidatorSetUpdates")
	RegisterType(ctypes.EventDataString(""), "tendermint/event/ProposalString")

	RegisterType(&ctypes.DuplicateVoteEvidence{}, "tendermint/DuplicateVoteEvidence")
	RegisterType(&ctypes.LightClientAttackEvidence{}, "tendermint/LightClientAttackEvidence")
}
