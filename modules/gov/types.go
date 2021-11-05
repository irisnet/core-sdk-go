package gov

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gogo/protobuf/proto"
	yaml "gopkg.in/yaml.v2"

	codectypes "github.com/irisnet/core-sdk-go/codec/types"
	"github.com/irisnet/core-sdk-go/types"
	"github.com/irisnet/core-sdk-go/types/errors"
)

const (
	ModuleName             = "gov"
	AttributeKeyProposalID = "proposal_id"
)

var (
	_ types.Msg = &MsgSubmitProposal{}
	_ types.Msg = &MsgDeposit{}
	_ types.Msg = &MsgVote{}
	_ types.Msg = &MsgVoteWeighted{}
)

// NewMsgSubmitProposal creates a new MsgSubmitProposal.
//nolint:interfacer
func NewMsgSubmitProposal(content Content, initialDeposit types.Coins, proposer types.AccAddress) (*MsgSubmitProposal, error) {
	m := &MsgSubmitProposal{
		InitialDeposit: initialDeposit,
		Proposer:       proposer.String(),
	}
	err := m.SetContent(content)
	if err != nil {
		return nil, err
	}
	return m, nil
}

func (m *MsgSubmitProposal) GetInitialDeposit() types.Coins { return m.InitialDeposit }

func (m *MsgSubmitProposal) GetProposer() types.AccAddress {
	proposer, _ := types.AccAddressFromBech32(m.Proposer)
	return proposer
}

func (m *MsgSubmitProposal) SetContent(content Content) error {
	msg, ok := content.(proto.Message)
	if !ok {
		return fmt.Errorf("can't proto marshal %T", msg)
	}
	any, err := codectypes.NewAnyWithValue(msg)
	if err != nil {
		return err
	}
	m.Content = any
	return nil
}

func (m *MsgSubmitProposal) GetContent() Content {
	content, ok := m.Content.GetCachedValue().(Content)
	if !ok {
		return nil
	}
	return content
}

func (m MsgSubmitProposal) Route() string { return ModuleName }

// Type implements Msg
func (m MsgSubmitProposal) Type() string { return "submit_proposal" }

// ValidateBasic implements Msg
func (m MsgSubmitProposal) ValidateBasic() error {
	if m.Proposer == "" {
		return errors.Wrapf(ErrInvalidProposer, "missing Proposer")
	}
	if !m.InitialDeposit.IsValid() {
		return errors.Wrapf(ErrInvalidCoin, "invalidCoins coins, %s", m.InitialDeposit.String())
	}
	if m.InitialDeposit.IsAnyNegative() {
		return errors.Wrapf(ErrInvalidCoin, "invalidCoins coins, %s", m.InitialDeposit.String())
	}

	content := m.GetContent()
	if content == nil {
		return errors.Wrapf(ErrMissingContent, "missing content")
	}

	if err := content.ValidateBasic(); err != nil {
		return err
	}

	return nil
}

// GetSigners implements Msg
func (m MsgSubmitProposal) GetSigners() []types.AccAddress {
	proposer, _ := types.AccAddressFromBech32(m.Proposer)
	return []types.AccAddress{proposer}
}

// String implements the Stringer interface
func (m MsgSubmitProposal) String() string {
	out, _ := yaml.Marshal(m)
	return string(out)
}

// UnpackInterfaces implements UnpackInterfacesMessage.UnpackInterfaces
func (m MsgSubmitProposal) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	var content Content
	return unpacker.UnpackAny(m.Content, &content)
}

// ValidateBasic implements Msg
func (msg MsgDeposit) ValidateBasic() error {
	if msg.Depositor == "" {
		return errors.Wrapf(ErrMissingProposer, "missing Proposer")
	}
	if !msg.Amount.IsValid() {
		return errors.Wrapf(ErrInvalidAmount, "invalidCoins coins, %s", msg.Amount.String())
	}
	if msg.Amount.IsAnyNegative() {
		return errors.Wrapf(ErrInvalidCoin, "invalidCoins coins, %s", msg.Amount.String())
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgDeposit) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSigners implements Msg
func (msg MsgDeposit) GetSigners() []types.AccAddress {
	depositor, _ := types.AccAddressFromBech32(msg.Depositor)
	return []types.AccAddress{depositor}
}

// ValidateBasic implements Msg
func (msg MsgVote) ValidateBasic() error {
	if msg.Voter == "" {
		return errors.Wrapf(ErrMissingProposer, "missing Proposer")
	}

	if !ValidVoteOption(msg.Option) {
		return errors.Wrapf(ErrInvalidVoteOption, "invalid vote option %s", msg.Option.String())
	}

	return nil
}

func ValidVoteOption(option VoteOption) bool {
	if option == OptionYes ||
		option == OptionAbstain ||
		option == OptionNo ||
		option == OptionNoWithVeto {
		return true
	}
	return false
}

// String implements the Stringer interface
func (msg MsgVote) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSigners implements Msg
func (msg MsgVote) GetSigners() []types.AccAddress {
	voter, _ := types.AccAddressFromBech32(msg.Voter)
	return []types.AccAddress{voter}
}

// ValidateBasic implements Msg
func (msg MsgVoteWeighted) ValidateBasic() error {
	if msg.Voter == "" {
		return errors.Wrap(errors.ErrInvalidAddress, msg.Voter)
	}

	if len(msg.Options) == 0 {
		return errors.Wrap(errors.ErrInvalidRequest, WeightedVoteOptions(msg.Options).String())
	}

	totalWeight := types.NewDec(0)
	usedOptions := make(map[VoteOption]bool)
	for _, option := range msg.Options {
		if !ValidWeightedVoteOption(option) {
			return errors.Wrap(ErrInvalidWeightedVoteOption, option.String())
		}
		totalWeight = totalWeight.Add(option.Weight)
		if usedOptions[option.Option] {
			return errors.Wrap(ErrDuplicatedVoteOption, "Duplicated vote option")
		}
		usedOptions[option.Option] = true
	}

	if totalWeight.GT(types.NewDec(1)) {
		return errors.Wrap(ErrInvalidTotalWeight, "Total weight overflow 1.00")
	}

	if totalWeight.LT(types.NewDec(1)) {
		return errors.Wrap(ErrInvalidTotalWeight, "Total weight lower than 1.00")
	}

	return nil
}

// String implements the Stringer interface
func (msg MsgVoteWeighted) String() string {
	out, _ := yaml.Marshal(msg)
	return string(out)
}

// GetSigners implements Msg
func (msg MsgVoteWeighted) GetSigners() []types.AccAddress {
	voter, _ := types.AccAddressFromBech32(msg.Voter)
	return []types.AccAddress{voter}
}

func (v WeightedVoteOption) String() string {
	out, _ := json.Marshal(v)
	return string(out)
}

// WeightedVoteOptions describes array of WeightedVoteOptions
type WeightedVoteOptions []WeightedVoteOption

func (v WeightedVoteOptions) String() (out string) {
	for _, opt := range v {
		out += opt.String() + "\n"
	}

	return strings.TrimSpace(out)
}

// ValidWeightedVoteOption returns true if the sub vote is valid and false otherwise.
func ValidWeightedVoteOption(option WeightedVoteOption) bool {
	if !option.Weight.IsPositive() || option.Weight.GT(types.NewDec(1)) {
		return false
	}
	return ValidVoteOption(option.Option)
}

func (q Proposal) Convert() interface{} {
	return QueryProposalResp{
		ProposalId: q.ProposalId,
		Status:     ProposalStatus_name[int32(q.Status)],
		FinalTallyResult: QueryTallyResultResp{
			Yes:        q.FinalTallyResult.Yes,
			Abstain:    q.FinalTallyResult.Abstain,
			No:         q.FinalTallyResult.No,
			NoWithVeto: q.FinalTallyResult.NoWithVeto,
		},
		SubmitTime:      q.SubmitTime,
		DepositEndTime:  q.DepositEndTime,
		TotalDeposit:    q.TotalDeposit,
		VotingStartTime: q.VotingStartTime,
		VotingEndTime:   q.VotingEndTime,
	}
}

type Proposals []Proposal

func (qs Proposals) Convert() interface{} {
	var res []QueryProposalResp
	for _, q := range qs {
		res = append(res, q.Convert().(QueryProposalResp))
	}
	return res
}

func (v Vote) Convert() interface{} {
	return QueryVoteResp{
		ProposalId: v.ProposalId,
		Voter:      v.Voter,
		Option:     int32(v.Option),
	}
}

type Votes []Vote

func (vs Votes) Convert() interface{} {
	var res []QueryVoteResp
	for _, v := range vs {
		res = append(res, v.Convert().(QueryVoteResp))
	}
	return res
}

func (q QueryParamsResponse) Convert() interface{} {
	return QueryParamsResp{
		VotingParams: votingParams{
			VotingPeriod: q.VotingParams.VotingPeriod,
		},
		DepositParams: depositParams{
			MinDeposit:       q.DepositParams.MinDeposit,
			MaxDepositPeriod: q.DepositParams.MaxDepositPeriod,
		},
		TallyParams: tallyParams{
			Quorum:        q.TallyParams.Quorum,
			Threshold:     q.TallyParams.Threshold,
			VetoThreshold: q.TallyParams.VetoThreshold,
		},
	}
}

func (d Deposit) Convert() interface{} {
	return QueryDepositResp(d)
}

type Deposits []Deposit

func (ds Deposits) Convert() interface{} {
	var res []QueryDepositResp
	for _, d := range ds {
		res = append(res, d.Convert().(QueryDepositResp))
	}
	return res
}

func (t TallyResult) Convert() interface{} {
	return QueryTallyResultResp(t)
}
