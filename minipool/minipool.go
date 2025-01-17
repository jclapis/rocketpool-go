package minipool

import (
    "fmt"
    "math/big"
    "sync"

    "github.com/ethereum/go-ethereum/accounts/abi/bind"
    "github.com/ethereum/go-ethereum/common"
    "golang.org/x/sync/errgroup"

    "github.com/rocket-pool/rocketpool-go/rocketpool"
    rptypes "github.com/rocket-pool/rocketpool-go/types"
)


// Settings
const (
    MinipoolAddressBatchSize = 50
    MinipoolDetailsBatchSize = 20
)


// Minipool details
type MinipoolDetails struct {
    Address common.Address              `json:"address"`
    Exists bool                         `json:"exists"`
    Pubkey rptypes.ValidatorPubkey      `json:"pubkey"`
    WithdrawalTotalBalance *big.Int     `json:"withdrawalTotalBalance"`
    WithdrawalNodeBalance *big.Int      `json:"withdrawalNodeBalance"`
    Withdrawable bool                   `json:"withdrawable"`
    WithdrawalProcessed bool            `json:"withdrawalProcessed"`
}


// Get all minipool details
func GetMinipools(rp *rocketpool.RocketPool, opts *bind.CallOpts) ([]MinipoolDetails, error) {
    minipoolAddresses, err := GetMinipoolAddresses(rp, opts)
    if err != nil {
        return []MinipoolDetails{}, err
    }
    return loadMinipoolDetails(rp, minipoolAddresses, opts);
}


// Get unprocessed minipool details
func GetUnprocessedMinipools(rp *rocketpool.RocketPool, opts *bind.CallOpts) ([]MinipoolDetails, error) {
    minipoolAddresses, err := GetUnprocessedMinipoolAddresses(rp, opts)
    if err != nil {
        return []MinipoolDetails{}, err
    }
    return loadMinipoolDetails(rp, minipoolAddresses, opts);
}


// Get a node's minipool details
func GetNodeMinipools(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) ([]MinipoolDetails, error) {
    minipoolAddresses, err := GetNodeMinipoolAddresses(rp, nodeAddress, opts)
    if err != nil {
        return []MinipoolDetails{}, err
    }
    return loadMinipoolDetails(rp, minipoolAddresses, opts);
}


// Load minipool details
func loadMinipoolDetails(rp *rocketpool.RocketPool, minipoolAddresses []common.Address, opts *bind.CallOpts) ([]MinipoolDetails, error) {

    // Load minipool details in batches
    details := make([]MinipoolDetails, len(minipoolAddresses))
    for bsi := 0; bsi < len(minipoolAddresses); bsi += MinipoolDetailsBatchSize {

        // Get batch start & end index
        msi := bsi
        mei := bsi + MinipoolDetailsBatchSize
        if mei > len(minipoolAddresses) { mei = len(minipoolAddresses) }

        // Load details
        var wg errgroup.Group
        for mi := msi; mi < mei; mi++ {
            mi := mi
            wg.Go(func() error {
                minipoolAddress := minipoolAddresses[mi]
                minipoolDetails, err := GetMinipoolDetails(rp, minipoolAddress, opts)
                if err == nil { details[mi] = minipoolDetails }
                return err
            })
        }
        if err := wg.Wait(); err != nil {
            return []MinipoolDetails{}, err
        }

    }

    // Return
    return details, nil

}


// Get all minipool addresses
func GetMinipoolAddresses(rp *rocketpool.RocketPool, opts *bind.CallOpts) ([]common.Address, error) {

    // Get minipool count
    minipoolCount, err := GetMinipoolCount(rp, opts)
    if err != nil {
        return []common.Address{}, err
    }

    // Load minipool addresses in batches
    addresses := make([]common.Address, minipoolCount)
    for bsi := uint64(0); bsi < minipoolCount; bsi += MinipoolAddressBatchSize {

        // Get batch start & end index
        msi := bsi
        mei := bsi + MinipoolAddressBatchSize
        if mei > minipoolCount { mei = minipoolCount }

        // Load addresses
        var wg errgroup.Group
        for mi := msi; mi < mei; mi++ {
            mi := mi
            wg.Go(func() error {
                address, err := GetMinipoolAt(rp, mi, opts)
                if err == nil { addresses[mi] = address }
                return err
            })
        }
        if err := wg.Wait(); err != nil {
            return []common.Address{}, err
        }

    }

    // Return
    return addresses, nil

}


// Get unprocessed minipool addresses
func GetUnprocessedMinipoolAddresses(rp *rocketpool.RocketPool, opts *bind.CallOpts) ([]common.Address, error) {

    // Get minipool count
    minipoolCount, err := GetUnprocessedMinipoolCount(rp, opts)
    if err != nil {
        return []common.Address{}, err
    }

    // Load minipool addresses in batches
    addresses := make([]common.Address, minipoolCount)
    for bsi := uint64(0); bsi < minipoolCount; bsi += MinipoolAddressBatchSize {

        // Get batch start & end index
        msi := bsi
        mei := bsi + MinipoolAddressBatchSize
        if mei > minipoolCount { mei = minipoolCount }

        // Load addresses
        var wg errgroup.Group
        for mi := msi; mi < mei; mi++ {
            mi := mi
            wg.Go(func() error {
                address, err := GetUnprocessedMinipoolAt(rp, mi, opts)
                if err == nil { addresses[mi] = address }
                return err
            })
        }
        if err := wg.Wait(); err != nil {
            return []common.Address{}, err
        }

    }

    // Return
    return addresses, nil

}


// Get a node's minipool addresses
func GetNodeMinipoolAddresses(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) ([]common.Address, error) {

    // Get minipool count
    minipoolCount, err := GetNodeMinipoolCount(rp, nodeAddress, opts)
    if err != nil {
        return []common.Address{}, err
    }

    // Load minipool addresses in batches
    addresses := make([]common.Address, minipoolCount)
    for bsi := uint64(0); bsi < minipoolCount; bsi += MinipoolAddressBatchSize {

        // Get batch start & end index
        msi := bsi
        mei := bsi + MinipoolAddressBatchSize
        if mei > minipoolCount { mei = minipoolCount }

        // Load addresses
        var wg errgroup.Group
        for mi := msi; mi < mei; mi++ {
            mi := mi
            wg.Go(func() error {
                address, err := GetNodeMinipoolAt(rp, nodeAddress, mi, opts)
                if err == nil { addresses[mi] = address }
                return err
            })
        }
        if err := wg.Wait(); err != nil {
            return []common.Address{}, err
        }

    }

    // Return
    return addresses, nil

}


// Get a node's validating minipool pubkeys
func GetNodeValidatingMinipoolPubkeys(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) ([]rptypes.ValidatorPubkey, error) {

    // Get minipool count
    minipoolCount, err := GetNodeValidatingMinipoolCount(rp, nodeAddress, opts)
    if err != nil {
        return []rptypes.ValidatorPubkey{}, err
    }

    // Load pubkeys in batches
    pubkeys := make([]rptypes.ValidatorPubkey, minipoolCount)
    for bsi := uint64(0); bsi < minipoolCount; bsi += MinipoolAddressBatchSize {

        // Get batch start & end index
        msi := bsi
        mei := bsi + MinipoolAddressBatchSize
        if mei > minipoolCount { mei = minipoolCount }

        // Load pubkeys
        var wg errgroup.Group
        for mi := msi; mi < mei; mi++ {
            mi := mi
            wg.Go(func() error {
                minipoolAddress, err := GetNodeValidatingMinipoolAt(rp, nodeAddress, mi, opts)
                if err != nil {
                    return err
                }
                pubkey, err := GetMinipoolPubkey(rp, minipoolAddress, opts)
                if err != nil {
                    return err
                }
                pubkeys[mi] = pubkey
                return nil
            })
        }
        if err := wg.Wait(); err != nil {
            return []rptypes.ValidatorPubkey{}, err
        }

    }

    // Return
    return pubkeys, nil

}


// Get a minipool's details
func GetMinipoolDetails(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.CallOpts) (MinipoolDetails, error) {

    // Data
    var wg errgroup.Group
    var exists bool
    var pubkey rptypes.ValidatorPubkey
    var withdrawalTotalBalance *big.Int
    var withdrawalNodeBalance *big.Int
    var withdrawable bool
    var withdrawalProcessed bool

    // Load data
    wg.Go(func() error {
        var err error
        exists, err = GetMinipoolExists(rp, minipoolAddress, opts)
        return err
    })
    wg.Go(func() error {
        var err error
        pubkey, err = GetMinipoolPubkey(rp, minipoolAddress, opts)
        return err
    })
    wg.Go(func() error {
        var err error
        withdrawalTotalBalance, err = GetMinipoolWithdrawalTotalBalance(rp, minipoolAddress, opts)
        return err
    })
    wg.Go(func() error {
        var err error
        withdrawalNodeBalance, err = GetMinipoolWithdrawalNodeBalance(rp, minipoolAddress, opts)
        return err
    })
    wg.Go(func() error {
        var err error
        withdrawable, err = GetMinipoolWithdrawable(rp, minipoolAddress, opts)
        return err
    })
    wg.Go(func() error {
        var err error
        withdrawalProcessed, err = GetMinipoolWithdrawalProcessed(rp, minipoolAddress, opts)
        return err
    })

    // Wait for data
    if err := wg.Wait(); err != nil {
        return MinipoolDetails{}, err
    }

    // Return
    return MinipoolDetails{
        Address: minipoolAddress,
        Exists: exists,
        Pubkey: pubkey,
        WithdrawalTotalBalance: withdrawalTotalBalance,
        WithdrawalNodeBalance: withdrawalNodeBalance,
        Withdrawable: withdrawable,
        WithdrawalProcessed: withdrawalProcessed,
    }, nil

}


// Get the minipool count
func GetMinipoolCount(rp *rocketpool.RocketPool, opts *bind.CallOpts) (uint64, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return 0, err
    }
    minipoolCount := new(*big.Int)
    if err := rocketMinipoolManager.Call(opts, minipoolCount, "getMinipoolCount"); err != nil {
        return 0, fmt.Errorf("Could not get minipool count: %w", err)
    }
    return (*minipoolCount).Uint64(), nil
}


// Get a minipool address by index
func GetMinipoolAt(rp *rocketpool.RocketPool, index uint64, opts *bind.CallOpts) (common.Address, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAddress := new(common.Address)
    if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getMinipoolAt", big.NewInt(int64(index))); err != nil {
        return common.Address{}, fmt.Errorf("Could not get minipool %d address: %w", index, err)
    }
    return *minipoolAddress, nil
}


// Get the unprocessed minipool count
func GetUnprocessedMinipoolCount(rp *rocketpool.RocketPool, opts *bind.CallOpts) (uint64, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return 0, err
    }
    minipoolCount := new(*big.Int)
    if err := rocketMinipoolManager.Call(opts, minipoolCount, "getUnprocessedMinipoolCount"); err != nil {
        return 0, fmt.Errorf("Could not get minipool count: %w", err)
    }
    return (*minipoolCount).Uint64(), nil
}


// Get an unprocessed minipool address by index
func GetUnprocessedMinipoolAt(rp *rocketpool.RocketPool, index uint64, opts *bind.CallOpts) (common.Address, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAddress := new(common.Address)
    if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getUnprocessedMinipoolAt", big.NewInt(int64(index))); err != nil {
        return common.Address{}, fmt.Errorf("Could not get minipool %d address: %w", index, err)
    }
    return *minipoolAddress, nil
}


// Get a node's minipool count
func GetNodeMinipoolCount(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) (uint64, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return 0, err
    }
    minipoolCount := new(*big.Int)
    if err := rocketMinipoolManager.Call(opts, minipoolCount, "getNodeMinipoolCount", nodeAddress); err != nil {
        return 0, fmt.Errorf("Could not get node %s minipool count: %w", nodeAddress.Hex(), err)
    }
    return (*minipoolCount).Uint64(), nil
}


// Get a node's minipool address by index
func GetNodeMinipoolAt(rp *rocketpool.RocketPool, nodeAddress common.Address, index uint64, opts *bind.CallOpts) (common.Address, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAddress := new(common.Address)
    if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getNodeMinipoolAt", nodeAddress, big.NewInt(int64(index))); err != nil {
        return common.Address{}, fmt.Errorf("Could not get node %s minipool %d address: %w", nodeAddress.Hex(), index, err)
    }
    return *minipoolAddress, nil
}


// Get a node's validating minipool count
func GetNodeValidatingMinipoolCount(rp *rocketpool.RocketPool, nodeAddress common.Address, opts *bind.CallOpts) (uint64, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return 0, err
    }
    minipoolCount := new(*big.Int)
    if err := rocketMinipoolManager.Call(opts, minipoolCount, "getNodeValidatingMinipoolCount", nodeAddress); err != nil {
        return 0, fmt.Errorf("Could not get node %s validating minipool count: %w", nodeAddress.Hex(), err)
    }
    return (*minipoolCount).Uint64(), nil
}


// Get a node's validating minipool address by index
func GetNodeValidatingMinipoolAt(rp *rocketpool.RocketPool, nodeAddress common.Address, index uint64, opts *bind.CallOpts) (common.Address, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAddress := new(common.Address)
    if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getNodeValidatingMinipoolAt", nodeAddress, big.NewInt(int64(index))); err != nil {
        return common.Address{}, fmt.Errorf("Could not get node %s validating minipool %d address: %w", nodeAddress.Hex(), index, err)
    }
    return *minipoolAddress, nil
}


// Get a minipool address by validator pubkey
func GetMinipoolByPubkey(rp *rocketpool.RocketPool, pubkey rptypes.ValidatorPubkey, opts *bind.CallOpts) (common.Address, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return common.Address{}, err
    }
    minipoolAddress := new(common.Address)
    if err := rocketMinipoolManager.Call(opts, minipoolAddress, "getMinipoolByPubkey", pubkey[:]); err != nil {
        return common.Address{}, fmt.Errorf("Could not get validator %s minipool address: %w", pubkey.Hex(), err)
    }
    return *minipoolAddress, nil
}


// Check whether a minipool exists
func GetMinipoolExists(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.CallOpts) (bool, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return false, err
    }
    exists := new(bool)
    if err := rocketMinipoolManager.Call(opts, exists, "getMinipoolExists", minipoolAddress); err != nil {
        return false, fmt.Errorf("Could not get minipool %s exists status: %w", minipoolAddress.Hex(), err)
    }
    return *exists, nil
}


// Get a minipool's validator pubkey
func GetMinipoolPubkey(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.CallOpts) (rptypes.ValidatorPubkey, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return rptypes.ValidatorPubkey{}, err
    }
    pubkey := new(rptypes.ValidatorPubkey)
    if err := rocketMinipoolManager.Call(opts, pubkey, "getMinipoolPubkey", minipoolAddress); err != nil {
        return rptypes.ValidatorPubkey{}, fmt.Errorf("Could not get minipool %s pubkey: %w", minipoolAddress.Hex(), err)
    }
    return *pubkey, nil
}


// Get a minipool's total balance at withdrawal
func GetMinipoolWithdrawalTotalBalance(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.CallOpts) (*big.Int, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return nil, err
    }
    balance := new(*big.Int)
    if err := rocketMinipoolManager.Call(opts, balance, "getMinipoolWithdrawalTotalBalance", minipoolAddress); err != nil {
        return nil, fmt.Errorf("Could not get minipool %s withdrawal total balance: %w", minipoolAddress.Hex(), err)
    }
    return *balance, nil
}


// Get a minipool's node balance at withdrawal
func GetMinipoolWithdrawalNodeBalance(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.CallOpts) (*big.Int, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return nil, err
    }
    balance := new(*big.Int)
    if err := rocketMinipoolManager.Call(opts, balance, "getMinipoolWithdrawalNodeBalance", minipoolAddress); err != nil {
        return nil, fmt.Errorf("Could not get minipool %s withdrawal node balance: %w", minipoolAddress.Hex(), err)
    }
    return *balance, nil
}


// Check whether a minipool is withdrawable
func GetMinipoolWithdrawable(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.CallOpts) (bool, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return false, err
    }
    withdrawable := new(bool)
    if err := rocketMinipoolManager.Call(opts, withdrawable, "getMinipoolWithdrawable", minipoolAddress); err != nil {
        return false, fmt.Errorf("Could not get minipool %s withdrawable status: %w", minipoolAddress.Hex(), err)
    }
    return *withdrawable, nil
}


// Check whether a minipool's validator withdrawal has been processed
func GetMinipoolWithdrawalProcessed(rp *rocketpool.RocketPool, minipoolAddress common.Address, opts *bind.CallOpts) (bool, error) {
    rocketMinipoolManager, err := getRocketMinipoolManager(rp)
    if err != nil {
        return false, err
    }
    processed := new(bool)
    if err := rocketMinipoolManager.Call(opts, processed, "getMinipoolWithdrawalProcessed", minipoolAddress); err != nil {
        return false, fmt.Errorf("Could not get minipool %s withdrawal processed status: %w", minipoolAddress.Hex(), err)
    }
    return *processed, nil
}


// Get contracts
var rocketMinipoolManagerLock sync.Mutex
func getRocketMinipoolManager(rp *rocketpool.RocketPool) (*rocketpool.Contract, error) {
    rocketMinipoolManagerLock.Lock()
    defer rocketMinipoolManagerLock.Unlock()
    return rp.GetContract("rocketMinipoolManager")
}

