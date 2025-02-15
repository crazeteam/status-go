package activity

import (
	"context"
	"database/sql"
	"math/big"
	"testing"
	"time"

	"go.uber.org/mock/gomock"

	eth "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/event"

	"github.com/status-im/status-go/appdatabase"
	"github.com/status-im/status-go/multiaccounts/accounts"
	ethclient "github.com/status-im/status-go/rpc/chain/ethclient"
	mock_rpcclient "github.com/status-im/status-go/rpc/mock/client"
	"github.com/status-im/status-go/services/wallet/bigint"
	"github.com/status-im/status-go/services/wallet/common"
	"github.com/status-im/status-go/services/wallet/thirdparty"
	"github.com/status-im/status-go/services/wallet/token"
	mock_token "github.com/status-im/status-go/services/wallet/token/mock/token"
	"github.com/status-im/status-go/services/wallet/transfer"
	"github.com/status-im/status-go/services/wallet/walletevent"
	"github.com/status-im/status-go/t/helpers"
	"github.com/status-im/status-go/transactions"
	"github.com/status-im/status-go/walletdatabase"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

const shouldNotWaitTimeout = 19999 * time.Second

// mockCollectiblesManager implements the collectibles.ManagerInterface
type mockCollectiblesManager struct {
	mock.Mock
}

func (m *mockCollectiblesManager) FetchAssetsByCollectibleUniqueID(ctx context.Context, uniqueIDs []thirdparty.CollectibleUniqueID, asyncFetch bool) ([]thirdparty.FullCollectibleData, error) {
	args := m.Called(uniqueIDs)
	res := args.Get(0)
	if res == nil {
		return nil, args.Error(1)
	}
	return res.([]thirdparty.FullCollectibleData), args.Error(1)
}

func (m *mockCollectiblesManager) FetchCollectionSocialsAsync(contractID thirdparty.ContractID) error {
	args := m.Called(contractID)
	res := args.Get(0)
	if res == nil {
		return args.Error(1)
	}
	return nil
}

type testState struct {
	service          *Service
	eventFeed        *event.Feed
	tokenMock        *mock_token.MockManagerInterface
	collectiblesMock *mockCollectiblesManager
	close            func()
	pendingTracker   *transactions.PendingTxTracker
	chainClient      *transactions.MockChainClient
	rpcClient        *mock_rpcclient.MockClientInterface
}

func setupTestService(tb testing.TB) (state testState) {
	db, err := helpers.SetupTestMemorySQLDB(walletdatabase.DbInitializer{})
	require.NoError(tb, err)

	appDB, err := helpers.SetupTestMemorySQLDB(appdatabase.DbInitializer{})
	require.NoError(tb, err)
	accountsDB, err := accounts.NewDB(appDB)
	require.NoError(tb, err)

	state.eventFeed = new(event.Feed)
	mockCtrl := gomock.NewController(tb)
	state.tokenMock = mock_token.NewMockManagerInterface(mockCtrl)
	state.collectiblesMock = &mockCollectiblesManager{}

	state.chainClient = transactions.NewMockChainClient()
	state.rpcClient = mock_rpcclient.NewMockClientInterface(mockCtrl)
	state.rpcClient.EXPECT().AbstractEthClient(gomock.Any()).DoAndReturn(func(chainID common.ChainID) (ethclient.BatchCallClient, error) {
		return state.chainClient.AbstractEthClient(chainID)
	}).AnyTimes()

	// Ensure we process pending transactions as needed, only once
	pendingCheckInterval := time.Second
	state.pendingTracker = transactions.NewPendingTxTracker(db, state.rpcClient, nil, state.eventFeed, pendingCheckInterval)

	state.service = NewService(db, accountsDB, state.tokenMock, state.collectiblesMock, state.eventFeed, state.pendingTracker)
	state.service.debounceDuration = 0
	state.close = func() {
		require.NoError(tb, state.pendingTracker.Stop())
		require.NoError(tb, db.Close())
		defer mockCtrl.Finish()
	}

	return state
}

type arg struct {
	chainID         common.ChainID
	tokenAddressStr string
	tokenIDStr      string
	tokenID         *big.Int
	tokenAddress    *eth.Address
}

// insertStubTransfersWithCollectibles will insert nil if tokenIDStr is empty
func insertStubTransfersWithCollectibles(t *testing.T, db *sql.DB, args []arg) (fromAddresses, toAddresses []eth.Address) {
	trs, fromAddresses, toAddresses := transfer.GenerateTestTransfers(t, db, 0, len(args))
	for i := range args {
		trs[i].ChainID = args[i].chainID
		if args[i].tokenIDStr == "" {
			args[i].tokenID = nil
		} else {
			args[i].tokenID = new(big.Int)
			args[i].tokenID.SetString(args[i].tokenIDStr, 0)
		}
		args[i].tokenAddress = new(eth.Address)
		*args[i].tokenAddress = eth.HexToAddress(args[i].tokenAddressStr)
		transfer.InsertTestTransferWithOptions(t, db, trs[i].To, &trs[i], &transfer.TestTransferOptions{
			TokenAddress: *args[i].tokenAddress,
			TokenID:      args[i].tokenID,
		})
	}
	return fromAddresses, toAddresses
}

func TestService_UpdateCollectibleInfo(t *testing.T) {
	state := setupTestService(t)
	defer state.close()

	args := []arg{
		{5, "0xA2838FDA19EB6EED3F8B9EFF411D4CD7D2DE0313", "0x0D", nil, nil},
		{5, "0xA2838FDA19EB6EED3F8B9EFF411D4CD7D2DE0313", "0x762AD3E4934E687F8701F24C7274E5209213FD6208FF952ACEB325D028866949", nil, nil},
		{5, "0xA2838FDA19EB6EED3F8B9EFF411D4CD7D2DE0313", "0x762AD3E4934E687F8701F24C7274E5209213FD6208FF952ACEB325D028866949", nil, nil},
		{5, "0x3d6afaa395c31fcd391fe3d562e75fe9e8ec7e6a", "", nil, nil},
		{5, "0xA2838FDA19EB6EED3F8B9EFF411D4CD7D2DE0313", "0x0F", nil, nil},
	}
	fromAddresses, toAddresses := insertStubTransfersWithCollectibles(t, state.service.db, args)

	ch := make(chan walletevent.Event)
	sub := state.eventFeed.Subscribe(ch)

	// Expect one call for the fungible token
	state.tokenMock.EXPECT().LookupTokenIdentity(uint64(5), eth.HexToAddress("0x3d6afaa395c31fcd391fe3d562e75fe9e8ec7e6a"), false).Return(
		&token.Token{
			ChainID: 5,
			Address: eth.HexToAddress("0x3d6afaa395c31fcd391fe3d562e75fe9e8ec7e6a"),
			Symbol:  "STT",
		},
	).Times(1)
	state.collectiblesMock.On("FetchAssetsByCollectibleUniqueID", []thirdparty.CollectibleUniqueID{
		{
			ContractID: thirdparty.ContractID{
				ChainID: args[4].chainID,
				Address: *args[4].tokenAddress},
			TokenID: &bigint.BigInt{Int: args[4].tokenID},
		}, {
			ContractID: thirdparty.ContractID{
				ChainID: args[1].chainID,
				Address: *args[1].tokenAddress},
			TokenID: &bigint.BigInt{Int: args[1].tokenID},
		}, {
			ContractID: thirdparty.ContractID{
				ChainID: args[0].chainID,
				Address: *args[0].tokenAddress},
			TokenID: &bigint.BigInt{Int: args[0].tokenID},
		},
	}).Return([]thirdparty.FullCollectibleData{
		{
			CollectibleData: thirdparty.CollectibleData{
				ID: thirdparty.CollectibleUniqueID{
					ContractID: thirdparty.ContractID{
						ChainID: args[4].chainID,
						Address: *args[4].tokenAddress},
					TokenID: &bigint.BigInt{Int: args[4].tokenID},
				},
				Name:     "Test 4",
				ImageURL: "test://url/4"},
			CollectionData: nil,
		}, {
			CollectibleData: thirdparty.CollectibleData{
				ID: thirdparty.CollectibleUniqueID{
					ContractID: thirdparty.ContractID{
						ChainID: args[1].chainID,
						Address: *args[1].tokenAddress},
					TokenID: &bigint.BigInt{Int: args[1].tokenID},
				},
				Name:     "Test 1",
				ImageURL: "test://url/1"},
			CollectionData: nil,
		},
		{
			CollectibleData: thirdparty.CollectibleData{
				ID: thirdparty.CollectibleUniqueID{
					ContractID: thirdparty.ContractID{
						ChainID: args[0].chainID,
						Address: *args[0].tokenAddress},
					TokenID: &bigint.BigInt{Int: args[0].tokenID},
				},
				Name:     "Test 0",
				ImageURL: "test://url/0"},
			CollectionData: nil,
		},
	}, nil).Once()

	state.service.FilterActivityAsync(0, append(fromAddresses, toAddresses...), allNetworksFilter(), Filter{}, 0, 10)

	filterResponseCount := 0
	var updates []EntryData

	for i := 0; i < 2; i++ {
		select {
		case res := <-ch:
			switch res.Type {
			case EventActivityFilteringDone:
				payload, err := walletevent.GetPayload[FilterResponse](res)
				require.NoError(t, err)
				require.Equal(t, ErrorCodeSuccess, payload.ErrorCode)
				require.Equal(t, 5, len(payload.Activities))
				filterResponseCount++
			case EventActivityFilteringUpdate:
				err := walletevent.ExtractPayload(res, &updates)
				require.NoError(t, err)
			}
		case <-time.NewTimer(shouldNotWaitTimeout).C:
			require.Fail(t, "timeout while waiting for event")
		}
	}

	// FetchAssetsByCollectibleUniqueID will receive only unique ids, while number of entries can be bigger
	require.Equal(t, 1, filterResponseCount)
	require.Equal(t, 4, len(updates))
	require.Equal(t, "Test 4", *updates[0].NftName)
	require.Equal(t, "test://url/4", *updates[0].NftURL)
	require.Equal(t, "Test 1", *updates[1].NftName)
	require.Equal(t, "test://url/1", *updates[1].NftURL)
	require.Equal(t, "Test 1", *updates[2].NftName)
	require.Equal(t, "test://url/1", *updates[2].NftURL)
	require.Equal(t, "Test 0", *updates[3].NftName)
	require.Equal(t, "test://url/0", *updates[3].NftURL)

	sub.Unsubscribe()
}

func TestService_UpdateCollectibleInfo_Error(t *testing.T) {
	state := setupTestService(t)
	defer state.close()

	args := []arg{
		{5, "0xA2838FDA19EB6EED3F8B9EFF411D4CD7D2DE0313", "0x762AD3E4934E687F8701F24C7274E5209213FD6208FF952ACEB325D028866949", nil, nil},
		{5, "0xA2838FDA19EB6EED3F8B9EFF411D4CD7D2DE0313", "0x0D", nil, nil},
	}

	ch := make(chan walletevent.Event, 4)
	sub := state.eventFeed.Subscribe(ch)

	fromAddresses, toAddresses := insertStubTransfersWithCollectibles(t, state.service.db, args)

	state.collectiblesMock.On("FetchAssetsByCollectibleUniqueID", mock.Anything).Return(nil, thirdparty.ErrChainIDNotSupported).Once()

	state.service.FilterActivityAsync(0, append(fromAddresses, toAddresses...), allNetworksFilter(), Filter{}, 0, 5)

	filterResponseCount := 0
	updatesCount := 0

	for i := 0; i < 2; i++ {
		select {
		case res := <-ch:
			switch res.Type {
			case EventActivityFilteringDone:
				payload, err := walletevent.GetPayload[FilterResponse](res)
				require.NoError(t, err)
				require.Equal(t, ErrorCodeSuccess, payload.ErrorCode)
				require.Equal(t, 2, len(payload.Activities))
				filterResponseCount++
			case EventActivityFilteringUpdate:
				updatesCount++
			}
		case <-time.NewTimer(20 * time.Millisecond).C:
			// We wait to ensure the EventActivityFilteringUpdate is never sent
		}
	}

	require.Equal(t, 1, filterResponseCount)
	require.Equal(t, 0, updatesCount)

	sub.Unsubscribe()
}

func setupTransactions(t *testing.T, state testState, txCount int, testTxs []transactions.TestTxSummary) (allAddresses []eth.Address, pendings []transactions.PendingTransaction, ch chan walletevent.Event, cleanup func()) {
	ch = make(chan walletevent.Event, 4)
	sub := state.eventFeed.Subscribe(ch)

	pendings = transactions.MockTestTransactions(t, state.chainClient, testTxs)
	for _, p := range pendings {
		allAddresses = append(allAddresses, p.From, p.To)
	}

	txs, fromTrs, toTrs := transfer.GenerateTestTransfers(t, state.service.db, len(pendings), txCount)
	for i := range txs {
		transfer.InsertTestTransfer(t, state.service.db, txs[i].To, &txs[i])
	}

	allAddresses = append(append(allAddresses, fromTrs...), toTrs...)

	state.tokenMock.EXPECT().LookupTokenIdentity(gomock.Any(), gomock.Any(), gomock.Any()).Return(
		&token.Token{
			ChainID: 5,
			Address: eth.Address{},
			Symbol:  "ETH",
		},
	).AnyTimes()

	state.tokenMock.EXPECT().LookupToken(gomock.Any(), gomock.Any()).Return(
		&token.Token{
			ChainID: 5,
			Address: eth.Address{},
			Symbol:  "ETH",
		}, true,
	).AnyTimes()

	return allAddresses, pendings, ch, func() {
		sub.Unsubscribe()
	}
}

func getValidateSessionUpdateHasNewOnTopFn(t *testing.T) func(payload SessionUpdate) bool {
	return func(payload SessionUpdate) bool {
		require.NotNil(t, payload.HasNewOnTop)
		require.True(t, *payload.HasNewOnTop)
		return false
	}
}

// validateSessionUpdateEvent expects will give up early if checkPayloadFn return true and not wait for expectCount
func validateSessionUpdateEvent(t *testing.T, ch chan walletevent.Event, filterResponseCount *int, expectCount int, checkPayloadFn func(payload SessionUpdate) bool) (pendingTransactionUpdate, sessionUpdatesCount int) {
	for sessionUpdatesCount < expectCount {
		select {
		case res := <-ch:
			switch res.Type {
			case transactions.EventPendingTransactionUpdate:
				pendingTransactionUpdate++
			case EventActivitySessionUpdated:
				payload, err := walletevent.GetPayload[SessionUpdate](res)
				require.NoError(t, err)

				if checkPayloadFn != nil && checkPayloadFn(*payload) {
					return
				}

				sessionUpdatesCount++
			case EventActivityFilteringDone:
				(*filterResponseCount)++
			}
		case <-time.NewTimer(shouldNotWaitTimeout).C:
			require.Fail(t, "timeout while waiting for EventActivitySessionUpdated")
		}
	}
	return
}

type extraExpect struct {
	offset    *int
	errorCode *ErrorCode
}

func getOptionalExpectations(e *extraExpect) (expectOffset int, expectErrorCode ErrorCode) {
	expectOffset = 0
	expectErrorCode = ErrorCodeSuccess

	if e != nil {
		if e.offset != nil {
			expectOffset = *e.offset
		}
		if e.errorCode != nil {
			expectErrorCode = *e.errorCode
		}
	}
	return
}

func validateFilteringDone(t *testing.T, ch chan walletevent.Event, resCount int, checkPayloadFn func(payload FilterResponse), extra *extraExpect) (filterResponseCount int) {
	for filterResponseCount < 1 {
		select {
		case res := <-ch:
			switch res.Type {
			case EventActivityFilteringDone:
				payload, err := walletevent.GetPayload[FilterResponse](res)
				require.NoError(t, err)

				expectOffset, expectErrorCode := getOptionalExpectations(extra)

				require.Equal(t, expectErrorCode, payload.ErrorCode)
				require.Equal(t, resCount, len(payload.Activities))

				require.Equal(t, expectOffset, payload.Offset)
				filterResponseCount++

				if checkPayloadFn != nil {
					checkPayloadFn(*payload)
				}
			}
		case <-time.NewTimer(shouldNotWaitTimeout).C:
			require.Fail(t, "timeout while waiting for EventActivityFilteringDone")
		}
	}
	return
}
