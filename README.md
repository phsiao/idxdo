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
`Ceramic`.
