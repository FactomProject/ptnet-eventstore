### TicTacToe contract

Contract Declaration
--------------------

```
{
    "inputs": [
        {
            "address": "|DEPOSITOR|",
            "amount": 1
        }
    ],
    "outputs": [
        {
            "address": "|DEPOSITOR|",
            "amount": 1
        },
        {
            "address": "|PLAYERX|",
            "amount": 1
        },
        {
            "address": "|PLAYERO|",
            "amount": 1
        }
    ],
    "blockheight": 60221409,
    "salt": "|RANDOM|",
    "contractid": "|ContractID|",
    "schema": "octoe-v1",
    "state": [ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1 ],
    "actions": {
        "ENDO": [ 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 0, 1 ],
        "ENDX": [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 1 ],
        "EXEC": [ 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, -1, -1, -1 ],
        "O00": [ -1, 0, 0, 0, 0, 0, 0, 0, 0, -1, 1, 0, 0, 0 ],
        "O01": [ 0, -1, 0, 0, 0, 0, 0, 0, 0, -1, 1, 0, 0, 0 ],
        "O02": [ 0, 0, -1, 0, 0, 0, 0, 0, 0, -1, 1, 0, 0, 0 ],
        "O10": [ 0, 0, 0, -1, 0, 0, 0, 0, 0, -1, 1, 0, 0, 0 ],
        "O11": [ 0, 0, 0, 0, -1, 0, 0, 0, 0, -1, 1, 0, 0, 0 ],
        "O12": [ 0, 0, 0, 0, 0, -1, 0, 0, 0, -1, 1, 0, 0, 0 ],
        "O20": [ 0, 0, 0, 0, 0, 0, -1, 0, 0, -1, 1, 0, 0, 0 ],
        "O21": [ 0, 0, 0, 0, 0, 0, 0, -1, 0, -1, 1, 0, 0, 0 ],
        "O22": [ 0, 0, 0, 0, 0, 0, 0, 0, -1, -1, 1, 0, 0, 0 ],
        "WINO": [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 1, 0, 0 ],
        "WINX": [ 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 1, 0 ],
        "X00": [ -1, 0, 0, 0, 0, 0, 0, 0, 0, 1, -1, 0, 0, 0 ],
        "X01": [ 0, -1, 0, 0, 0, 0, 0, 0, 0, 1, -1, 0, 0, 0 ],
        "X02": [ 0, 0, -1, 0, 0, 0, 0, 0, 0, 1, -1, 0, 0, 0 ],
        "X10": [ 0, 0, 0, -1, 0, 0, 0, 0, 0, 1, -1, 0, 0, 0 ],
        "X11": [ 0, 0, 0, 0, -1, 0, 0, 0, 0, 1, -1, 0, 0, 0 ],
        "X12": [ 0, 0, 0, 0, 0, -1, 0, 0, 0, 1, -1, 0, 0, 0 ],
        "X20": [ 0, 0, 0, 0, 0, 0, -1, 0, 0, 1, -1, 0, 0, 0 ],
        "X21": [ 0, 0, 0, 0, 0, 0, 0, -1, 0, 1, -1, 0, 0, 0 ],
        "X22": [ 0, 0, 0, 0, 0, 0, 0, 0, -1, 1, -1, 0, 0, 0 ]
    },
    "guards": [
        [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0 ],
        [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0, 0 ],
        [
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            0,
            -1,
            0,
            0,
            0,
            0
        ]
    ],
    "conditions": [
        [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1 ],
        [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0 ],
        [ 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, -1, 0, 0 ]
    ]
}
```
event0:
-------

First event triggers the EXEC 'BEGIN' transaction
and includes a json-serialized payload containging contract declaration

```
{
    "timestamp": 1541078199752125017,
    "schema": "octoe-v1",
    "action": "EXEC",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1 ],
    "output": [ 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0 ],
    "payload": "deadBeefewogICAgImlucHV0cyI6IFsKICAgICAgICB7CiAgICAgICAgICAgICJhZGRyZXNzIj="
}
```

event1:
-------

```
{
    "timestamp": 1541078199752724107,
    "schema": "octoe-v1",
    "action": "X11",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 1, 0, 0, 0 ],
    "output": [ 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0 ],
    "payload": "e30="
}
```

event2:
-------

```
{
    "timestamp": 1541078199752826526,
    "schema": "octoe-v1",
    "action": "O01",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0 ],
    "output": [ 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0 ],
    "payload": "e30="
}
```

event3:
-------

```
{
    "timestamp": 1541078199752920020,
    "schema": "octoe-v1",
    "action": "X00",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 1, 0, 1, 1, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0 ],
    "output": [ 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0 ],
    "payload": "e30="
}
```

event4:
-------

```
{
    "timestamp": 1541078199753010166,
    "schema": "octoe-v1",
    "action": "O02",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0 ],
    "output": [ 0, 0, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0 ],
    "payload": "e30="
}
```

event5:
-------

```
{
    "timestamp": 1541078199753095273,
    "schema": "octoe-v1",
    "action": "X22",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 0, 0, 0, 1, 0, 1, 1, 1, 1, 0, 1, 0, 0, 0 ],
    "output": [ 0, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0 ],
    "payload": "e30="
}
```

event6:
-------

```
{
    "timestamp": 1541078199753199866,
    "schema": "octoe-v1",
    "action": "WINX",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 0, 0, 0, 1, 0, 1, 1, 1, 0, 1, 0, 0, 0, 0 ],
    "output": [ 0, 0, 0, 1, 0, 1, 1, 1, 0, 0, 0, 0, 1, 0 ],
    "payload": "e30="
}
```


### Option Contract

Contract Declaration
--------------------

```
{
    "inputs": [ { "address": "|DEPOSITOR|", "amount": 1 } ],
    "outputs": [
        {
            "address": "|DEPOSITOR|",
            "amount": 1
        },
        {
            "address": "|PUBKEY1|",
            "amount": 1
        },
        {
            "address": "|PUBKEY2|",
            "amount": 1
        }
    ],
    "blockheight": 60221409,
    "salt": "|RANDOM|",
    "contractid": "|ContractID|",
    "schema": "option-v1",
    "state": [ 1, 1, 1, 1, 1 ],
    "actions": {
        "EXEC": [ 0, -1, -1, 0, -1 ],
        "FAIL": [ -1, 0, 0, -1, 1 ],
        "HALT": [ 0, 0, 0, -1, 0 ],
        "OPT_0": [ -1, 1, 0, 0, 0 ],
        "OPT_1": [ -1, 0, 1, 0, 0 ]
    },
    "guards": [
        [ 0, 0, 0, -1, 0 ],
        [ 0, 0, 0, -1, 0 ],
        [ 0, 0, 0, -1, 0 ]
    ],
    "conditions": [
        [ 0, 0, 0, 0, -1 ],
        [ 0, -1, 0, 0, 0 ],
        [ 0, 0, -1, 0, 0 ]
    ]
}
```

event0:
-------

```
{
    "timestamp": 1541078199753397543,
    "schema": "option-v1",
    "action": "EXEC",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 1, 1, 1, 1, 1 ],
    "output": [ 1, 0, 0, 1, 0 ],
    "payload": "deadBeefewogICAgImlucHV0cyI6IFsKICAgICAgICB7CiAgICAgICAgICAgICJhZGRyZXNzIj="
}
```

event1:
-------

```
{
    "timestamp": 1541078199753569309,
    "schema": "option-v1",
    "action": "OPT_0",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 1, 0, 0, 1, 0 ],
    "output": [ 0, 1, 0, 1, 0 ],
    "payload": "e30="
}
```

event1:
-------
This action is called HALT (and also happens to halt the contract)

```
{
    "timestamp": 1541078199753664258,
    "schema": "option-v1",
    "action": "HALT",
    "oid": "|ContractID|",
    "value": 1,
    "input": [ 0, 1, 0, 1, 0 ],
    "output": [ 0, 1, 0, 0, 0 ],
    "payload": "e30="
}
```
