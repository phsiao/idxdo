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
records, and provides an independent way to verify and validate published
records. `idxdo` does not depend on any code or library from `IDX` or `Ceramic`
projects.

## Installation

Executing

`$ go get github.com/phsiao/idxdo`

would install the command `idxdo` in your `$GOPATH/bin`.

## Usage

### Example: Gitcoin Passport backup

---

Gitcoin Passport issues stamps for identities that they can verify about you,
such as your Facebook user id, Gmail email address, and Twitter handle. Gitcoin
then calculates trust score from the stamps --- it would filter out stamps that
share the same external identities with other Gitcoin Passports to mitigate
obvious Syble attacks.

The Passport document is stored in a Ceramic stream contains links to stamps
that are in Verifiable Credential format. Your Passport document is linked from
an IDX index that is created from your Ethereum address.

So knowing an Ethereum address would allow anyone to retrieve its Gitcoin
Passport document if it has one. The document can then be backed up. You can
also use it as a proof that you have identities that are verified by Gitcoin,
and if the verifiers trust Gitcoin then they might be able to trust you to have
those identities without asking you to prove it to them again. Note that the
stamp does not reveal your identity in other community, just that you have an
identity that Gitcoin was able to verify.

Executing the command below would try to download and dump the Gitcoin Passport
document to stdout. You can back up your Passport and inspect what stamps
Gitcoin have issued to you.

```
$ idxdo gp dump --account <your etherum address starts with 0x...>
```

### Example: Gitocin Passport data walk through

---

You can use `idxdo` to inspect IDX documents that are interesting to you,
whether they belong to you or not, and who issued them. Below is an example of
performing exploration using the functionalities provided by `idxdo`.

#### Get your IDX StreamID

Gitcoin Passport uses the PKH method for DID. So you first need to compute the
IDX index StreamID for your DID.

```
$ idxdo idx id pkh --account 0x6C1e268Fd076B5EaD3774F26D65f21A21D369179
k2t6wyfsu4pg062qh6tvm5zkb3qe6e7i59s592zrl4knu0vb7ykz0s18g5i5pv
```

`k2t6wyfsu4pg062qh6tvm5zkb3qe6e7i59s592zrl4knu0vb7ykz0s18g5i5pv` is the StreamID
associated with account `0x6C1e268Fd076B5EaD3774F26D65f21A21D369179` using PKH
method. To get the IDX StreamID of the account that you are interested in,
replace the `--account` argument with the account that you are interested in.

Different IDX streams can use differnet DID methods, for example, `3id` and
`key` are two other DID methods that are supported and used by other dapps as
identity in IDX. `idxdo` currently only support PKH method but can be extended
to support other methods.

#### Get your IDX Index

After you have your IDX index StreamID you can then get the document by
executing

```
$ idxdo idx state <your IDX StreamID from previous step>
```

The `content` secion of the json output is the list of your identity documents
by their key/value pairs. The key
`kjzl6cwe1jw148h1e14jb5fkf55xmqhmyorp29r9cq356c7ou74ulowf8czjlzs` is used to
indicate that the linked document is a Gitcoin Passport associated with the
identity.

The record value starts with `ceramic://` and for querying StreamID using
`idxdo` you need to remove the scheme and prefix.

#### Show your identity document

You can run

```
$ idxdo idx record <your IDX StreamID from previous step>
```

to go through the `content` of the index and let `idxdo` interpret what they are
for you. Currently `idxdo` only interprets Gitcoin Passport.

### Inspect CID, StreamID, and Stream State/Content

`idxdo` also has command to help you understand CID, StreamID, and inspect
Stream state and its content. Use `-h` or `--help` to show usage at different
category of the CLI.

## Alternatives

[cerscan](https://cerscan.com/) by the Orbis team is a web application allowing
you to query streams that they have indexed. It is also a very useful tool to
learn about IDX, Ceramic, and Gitcoin Passport.
