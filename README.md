# idxdo

`idxdo` is a a CLI tool for interacting with
[`IDX`](https://developers.idx.xyz/learn/overview/), an identity protocol for
open applications.

## Overview

`IDX` uses a decentralized index to hold records associated with a `DID`. The
data are stored using [`Ceramic`](https://blog.ceramic.network/what-is-ceramic/)
to allow structured data publishing and sharing, as well as other important
properties such as access control and version control.

Interacting with the records in `Ceramic` is non-trivial --- there is no
official web portal that you can go to to get all your records and look up
others' records. What is more important is that, a decentralized system should
allow any participant to verify and validate the records themselves. Because
these records are owned by the partcipants of the system, anyone should be able
to archive any record that are interesting to them and be able to prove its
validity.

`idxdo` is a CLI tool that aims to reduce the complexity of interacting with IDX
records, and provides an independent way to verify and validate records
published. `idxdo` does not depend on any code or library from `IDX` or
`Ceramic` projects.

## Installation

Execuring

`go get github.com/phsiao/idxdo`

would install the command `idxdo` in your `$GOPATH/bin`.

## Usage

### Example: GitCoin Passport

#### Get your IDX StreamID

GitCoin Passport uses the PKH method for DID. So you first need to compute the
IDX index StreamID for your DID. For example,

```
$ idxdo idx id pkh --account 0x6C1e268Fd076B5EaD3774F26D65f21A21D369179
k2t6wyfsu4pg062qh6tvm5zkb3qe6e7i59s592zrl4knu0vb7ykz0s18g5i5pv
```

`k2t6wyfsu4pg062qh6tvm5zkb3qe6e7i59s592zrl4knu0vb7ykz0s18g5i5pv` is the StreamID
associated with account `0x6C1e268Fd076B5EaD3774F26D65f21A21D369179` using PKH
method.

#### Get your IDX Index

After you have your IDX index StreamID you can then get the document by
executing

```
idxdo idx state <your StreamID from previous step>
```

The `content` secion of the output json is the list of your identity documents.
The key `kjzl6cwe1jw148h1e14jb5fkf55xmqhmyorp29r9cq356c7ou74ulowf8czjlzs`
represents GitCoin passport documents that is associated with this identity.

#### Show your identity document

You can run

```
idxdo idx record <your StreamID from previous step>
```

to go through the `content` in the index and interpret what documents they are.
