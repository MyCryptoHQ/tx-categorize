# tx-categorize

### Description
tx-categorize is a tool to parse & auto-categorize Ethereum transactions based on some schema-defined traits.

##### TLDR of implementation:
You feed in a tx history data structure that roughly conforms to Etherscan's transaction history endpoint's data structure and an array of transaction type schema objects found here to categorize the transaction by applying each schema object to the transaction to determine which schemas fit, then selects the schema object with the highest priority.

##### As an example
[This transaction](https://etherscan.io/tx/0x87597adb363a41cc300b4aecdfa2e8a7d2b279771476155857edb2d94bbe81ac) will categorize the transaction types to be a `UNISWAP_V2_EXCHANGE` transaction, a `1INCH_V3_EXCHANGE` transaction and an `ERC_20_TRANSFER` transaction, but will select the `1INCH_V3_EXCHANGE` tx type because of it's higher priority.

A schema applied to a transaction can be validated based off of simple things like the `from` address of a tx is a specific address (as is the case with mining pool payout transactions or some exchange withdraws), or containing a specific `log topic` or `methodID`

### CLI
The cli tool is used to improve the processes of adding / updating / monitoring schemas and their test configs.


### Todos:
- [x] Add commands:
- - [x] Generate schema template
- - [x] Generate tests from schemaid and tx hash
- - [x] List all unique schema types
- - [x] Run a test attempt of categorizing a tx
- - [x] Validate a schema
- - [ ] Create schemas based on subgraphs (_optional_)
- [x] Create build ci to auto-update the s3 storage of schemas
- [ ] Write documentation on tx-categorize process and schema outline


Schema applied to a transaction can be validated based off simple things like the from address of a tx is a specific address (as is the case with mining pool payout transactions or some exchange withdraws), or containing a specific log topic or methodID