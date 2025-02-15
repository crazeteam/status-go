package requests

import (
	"errors"

	"github.com/status-im/status-go/eth-node/crypto"
	"github.com/status-im/status-go/eth-node/types"
)

var ErrRequestToJoinCommunityInvalidCommunityID = errors.New("request-to-join-community: invalid community id")
var ErrRequestToJoinCommunityNoAddressesToReveal = errors.New("request-to-join-community: no addresses to reveal")
var ErrRequestToJoinCommunityMissingPassword = errors.New("request-to-join-community: password is necessary when sending a list of addresses")
var ErrRequestToJoinNoAirdropAddress = errors.New("request-to-join-community: airdropAddress is necessary when sending a list of addresses")
var ErrRequestToJoinNoAirdropAddressAmongAddressesToReveal = errors.New("request-to-join-community: airdropAddress must be in the set of addresses to reveal")
var ErrRequestToJoinCommunityInvalidSignature = errors.New("request-to-join-community: invalid signature")

type RequestToJoinCommunity struct {
	CommunityID          types.HexBytes   `json:"communityId"`
	ENSName              string           `json:"ensName"`
	AddressesToReveal    []string         `json:"addressesToReveal"`
	Signatures           []types.HexBytes `json:"signatures"` // the order of signatures should match the order of addresses
	AirdropAddress       string           `json:"airdropAddress"`
	ShareFutureAddresses bool             `json:"shareFutureAddresses"`
}

func (j *RequestToJoinCommunity) Validate() error {
	if len(j.CommunityID) == 0 {
		return ErrRequestToJoinCommunityInvalidCommunityID
	}

	if len(j.AddressesToReveal) == 0 {
		return ErrRequestToJoinCommunityNoAddressesToReveal
	}

	if j.AirdropAddress == "" {
		return ErrRequestToJoinNoAirdropAddress
	}

	found := false
	for _, address := range j.AddressesToReveal {
		if address == j.AirdropAddress {
			found = true
			break
		}
	}

	if !found {
		return ErrRequestToJoinNoAirdropAddressAmongAddressesToReveal
	}

	for _, signature := range j.Signatures {
		if len(signature) != crypto.SignatureLength {
			return ErrRequestToJoinCommunityInvalidSignature
		}
	}

	return nil
}
