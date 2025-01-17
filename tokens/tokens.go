package tokens

import (
    "context"
    "fmt"
    "math/big"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/ethclient"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/rocketpool-go/rocketpool"
)


// Token balances
type Balances struct {
    ETH *big.Int    `json:"eth"`
    NETH *big.Int   `json:"neth"`
    RETH *big.Int   `json:"reth"`
}


// Get token balances of an address
func GetBalances(rp *rocketpool.RocketPool, address common.Address, opts *bind.CallOpts) (Balances, error) {

    // Get call options block number
    var blockNumber *big.Int
    if opts != nil { blockNumber = opts.BlockNumber }

    // Data
    var wg errgroup.Group
    var ethBalance *big.Int
    var nethBalance *big.Int
    var rethBalance *big.Int

    // Load data
    wg.Go(func() error {
        var err error
        ethBalance, err = rp.Client.BalanceAt(context.Background(), address, blockNumber)
        return err
    })
    wg.Go(func() error {
        var err error
        nethBalance, err = GetNETHBalance(rp, address, opts)
        return err
    })
    wg.Go(func() error {
        var err error
        rethBalance, err = GetRETHBalance(rp, address, opts)
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return Balances{}, err
    }

    // Return
    return Balances{
        ETH: ethBalance,
        NETH: nethBalance,
        RETH: rethBalance,
    }, nil

}


// Get a token contract's ETH balance
func contractETHBalance(rp *rocketpool.RocketPool, tokenContract *rocketpool.Contract, opts *bind.CallOpts) (*big.Int, error) {
    var blockNumber *big.Int
    if opts != nil { blockNumber = opts.BlockNumber }
    return rp.Client.BalanceAt(context.Background(), *(tokenContract.Address), blockNumber)
}


// Get a token's total supply
func totalSupply(tokenContract *rocketpool.Contract, tokenName string, opts *bind.CallOpts) (*big.Int, error) {
    totalSupply := new(*big.Int)
    if err := tokenContract.Call(opts, totalSupply, "totalSupply"); err != nil {
        return nil, fmt.Errorf("Could not get %s total supply: %w", tokenName, err)
    }
    return *totalSupply, nil
}


// Get a token balance
func balanceOf(tokenContract *rocketpool.Contract, tokenName string, address common.Address, opts *bind.CallOpts) (*big.Int, error) {
    balance := new(*big.Int)
    if err := tokenContract.Call(opts, balance, "balanceOf", address); err != nil {
        return nil, fmt.Errorf("Could not get %s balance of %s: %w", tokenName, address.Hex(), err)
    }
    return *balance, nil
}


// Transfer tokens to an address
func transfer(client *ethclient.Client, tokenContract *rocketpool.Contract, tokenName string, to common.Address, amount *big.Int, opts *bind.TransactOpts) (*types.Receipt, error) {
    txReceipt, err := tokenContract.Transact(opts, "transfer", to, amount)
    if err != nil {
        return nil, fmt.Errorf("Could not transfer %s to %s: %w", tokenName, to.Hex(), err)
    }
    return txReceipt, nil
}

