package deposit

import (
    "testing"

    "github.com/rocket-pool/rocketpool-go/deposit"
    "github.com/rocket-pool/rocketpool-go/node"
    "github.com/rocket-pool/rocketpool-go/settings"
    "github.com/rocket-pool/rocketpool-go/utils/eth"

    "github.com/rocket-pool/rocketpool-go/tests/testutils/evm"
)


func TestDeposit(t *testing.T) {

    // State snapshotting
    if err := evm.TakeSnapshot(); err != nil { t.Fatal(err) }
    t.Cleanup(func() { if err := evm.RevertSnapshot(); err != nil { t.Fatal(err) } })

    // Make deposit
    opts := userAccount.GetTransactor()
    opts.Value = eth.EthToWei(10)
    if _, err := deposit.Deposit(rp, opts); err != nil {
        t.Fatal(err)
    }

    // Get & check deposit pool balance
    if balance, err := deposit.GetBalance(rp, nil); err != nil {
        t.Error(err)
    } else if balance.Cmp(opts.Value) != 0 {
        t.Error("Incorrect deposit pool balance")
    }

    // Get & check deposit pool excess balance
    if excessBalance, err := deposit.GetExcessBalance(rp, nil); err != nil {
        t.Error(err)
    } else if excessBalance.Cmp(opts.Value) != 0 {
        t.Error("Incorrect deposit pool excess balance")
    }

}


func TestAssignDeposits(t *testing.T) {

    // State snapshotting
    if err := evm.TakeSnapshot(); err != nil { t.Fatal(err) }
    t.Cleanup(func() { if err := evm.RevertSnapshot(); err != nil { t.Fatal(err) } })

    // Disable deposit assignments
    if _, err := settings.SetAssignDepositsEnabled(rp, false, ownerAccount.GetTransactor()); err != nil { t.Fatal(err) }

    // Make user deposit
    userDepositOpts := userAccount.GetTransactor()
    userDepositOpts.Value = eth.EthToWei(32)
    if _, err := deposit.Deposit(rp, userDepositOpts); err != nil { t.Fatal(err) }

    // Register node
    if _, err := node.RegisterNode(rp, "Australia/Brisbane", nodeAccount.GetTransactor()); err != nil { t.Fatal(err) }

    // Make node deposit
    nodeDepositOpts := nodeAccount.GetTransactor()
    nodeDepositOpts.Value = eth.EthToWei(16)
    if _, err := node.Deposit(rp, 0, nodeDepositOpts); err != nil { t.Fatal(err) }

    // Re-enable deposit assignments
    if _, err := settings.SetAssignDepositsEnabled(rp, true, ownerAccount.GetTransactor()); err != nil { t.Fatal(err) }

    // Get initial deposit pool balance
    balance1, err := deposit.GetBalance(rp, nil)
    if err != nil {
        t.Fatal(err)
    }

    // Assign deposits
    if _, err := deposit.AssignDeposits(rp, userAccount.GetTransactor()); err != nil {
        t.Fatal(err)
    }

    // Get & check updated deposit pool balance
    balance2, err := deposit.GetBalance(rp, nil)
    if err != nil {
        t.Fatal(err)
    } else if balance2.Cmp(balance1) != -1 {
        t.Error("Deposit pool balance did not decrease after assigning deposits")
    }

}

