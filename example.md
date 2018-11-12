

```go
import (
        "encoding/json"
        "github.com/FactomProject/ptnet-eventstore/contract"
        "github.com/FactomProject/ptnet-eventstore/ptnet"
)
```


```go
	c := contract.TicTacToeContract()
```


```go
[]interface{}{c.ContractID, c.Schema, c.BlockHeight}
```




    [|ContractID| octoe-v1 60221409]




```go
contract.Exists(c.Schema, c.ContractID)
```




    false




```go
event, _ := contract.Create(c, contract.DEPOSITOR_SECRET)
[]interface{}{event.Action, event.InputState, "=>", event.OutputState}
```




    [EXEC [1 1 1 1 1 1 1 1 1 1 1 1 1 1] => [1 1 1 1 1 1 1 1 1 0 1 0 0 0]]




```go
contract.Exists(c.Schema, c.ContractID)
```




    true




```go
contract.IsHalted(c)
```




    false




```go
contract.CanRedeem(c, contract.DEPOSITOR)
```




    false




```go
func commit(action string, key string, payload []byte) (*ptnet.Event, error) {
    return contract.Commit(contract.Command{
		ChainID:    contract.CHAIN_ID,      // chain to write to
		ContractID: contract.CONTRACT_ID,   // uniqueid for this contract execution
		Schema:     ptnet.OctoeV1,          // state machine version
		Action:     action,                 // state machine action
		Amount:     1,                      // triggers input action 'n' times
		Payload:    payload,                // arbitrary data optionally included
		Privkey:    contract.Identity[key], // secret identity used to sign event
		Pubkey:     key,                    // public identity
    })
}
```


```go
payload, _ := json.Marshal([]string{"hello", "world"})
event1, _ := commit("X11", contract.PLAYERX, payload)

[]interface{}{event1.Action, event1.InputState, "=>", event1.OutputState}
```




    [X11 [1 1 1 1 1 1 1 1 1 0 1 0 0 0] => [1 1 1 1 0 1 1 1 1 1 0 0 0 0]]




```go
// state machine makes players take turns
event2, err := commit("X00", contract.PLAYERX, nil)
[]interface{}{event2.Action, event2.InputState, "=>", err}
```




    [X00 [1 1 1 1 0 1 1 1 1 1 0 0 0 0] => invalid output: -1 offset: 10]




```go
// contract makes players sign each event - here we try to sign with an incorrect key
event3, err := commit("O01", contract.PLAYERX, nil)

[]interface{}{event3.Action, event3.InputState, "=>", err}
```




    [O01 [1 1 1 1 0 1 1 1 1 1 0 0 0 0] => invalid output: -1 offset: 10]




```go
event4, err := commit("O01", contract.PLAYERO, nil)

[]interface{}{event4.Action, event4.InputState, "=>", event4.OutputState}
```




    [O01 [1 1 1 1 0 1 1 1 1 1 0 0 0 0] => [1 0 1 1 0 1 1 1 1 0 1 0 0 0]]




```go
// state machine ensures move is used only once
event5, err := commit("X11", contract.PLAYERX, nil)

[]interface{}{event5.Action, event5.InputState, "=>", err}
```




    [X11 [1 0 1 1 0 1 1 1 1 0 1 0 0 0] => invalid output: -1 offset: 4]




```go
event6, err := commit("X00", contract.PLAYERX, nil)

[]interface{}{event6.Action, event6.InputState, "=>", event6.OutputState}
```




    [X00 [1 0 1 1 0 1 1 1 1 0 1 0 0 0] => [0 0 1 1 0 1 1 1 1 1 0 0 0 0]]




```go
event7, err := commit("O02", contract.PLAYERO, nil)

[]interface{}{event7.Action, event7.InputState, "=>", event7.OutputState}
```




    [O02 [0 0 1 1 0 1 1 1 1 1 0 0 0 0] => [0 0 0 1 0 1 1 1 1 0 1 0 0 0]]




```go
event8, err := commit("X22", contract.PLAYERX, nil)

[]interface{}{event8.Action, event8.InputState, "=>", event8.OutputState}
```




    [X22 [0 0 0 1 0 1 1 1 1 0 1 0 0 0] => [0 0 0 1 0 1 1 1 0 1 0 0 0 0]]




```go
// contract depositor confirms the winner
event9, err := commit("WINX", contract.DEPOSITOR, nil)

[]interface{}{event9.Action, event9.InputState, "=>", event9.OutputState}
```




    [WINX [0 0 0 1 0 1 1 1 0 1 0 0 0 0] => [0 0 0 1 0 1 1 1 0 0 0 0 1 0]]




```go
// contract is only Redeemable after Halting
contract.IsHalted(c)
```




    true




```go
// only the winner can redeem
contract.CanRedeem(c, contract.PLAYERX)
```




    true




```go
contract.CanRedeem(c, contract.PLAYERO)
```




    false


