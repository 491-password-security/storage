# crypto-helper

- [crypto-helper](#crypto-helper)
  - [Install](#install)
  - [Remarks](#remarks)
  - [Usage](#usage)
    - [Encryption](#encryption)
    - [Secret Sharing](#secret-sharing)
    - [Utilities](#utilities)

## Install

Clone the repository and run:
```
npm i <path>/crypto --save
```

It can then be imported as:
```js
var crypto = require('crypto-helper')
```

## Remarks

Since we never work with UTF-8 strings, except for hashing, all string inputs and outputs mentioned below are hexadecimal strings.

## Usage

### Encryption

`encrypt(key: string, plaintext: string) -> Object`
* Encrypts the plaintext using the given key, both as hexadecimal strings using AES. Key length should be 128, 192, or 256 bits. 
* Returns an object with two fields: `iv` and `ciphertext`, both as hexadecimal strings. `iv` is the *initialization vector* required for decryption.

`decrypt(key: string, iv: string, ciphertext: string) -> string`
* Returns the plaintext output of decryption, as a hexadecimal string, with the given key and IV. 

### Secret Sharing

These methods simply wrap the [secrets.js](https://github.com/grempe/secrets.js) library. See that for more details.

`share(secret: string, t: int, n: int) -> [string]`
* Generates `n` shares such that `t` of them are enough to reconstruct the `secret`. Outputs an array of hexadecimal strings. 

`combine(shares: [string]) -> string`
* Takes an array of shares in hexadecimal form and returns their reconstruction of the secret. Does not explicitly warn if given fewer shares than the threshold.

`newShare(id: int, shares: [string]) -> string`
* Adds a new share with the given id. The number of shares given must match the threshold, otherwise the new share will not be valid.

### Utilities

`random(bits: int) -> string`
* Returns random hexadecimal string of `bits` bits.
* It is secure, i.e. it can be used to generate keys etc.

`hash(input) -> string` 
* Returns hexadecimal string obtained by hashing (SHA-256) the input. The input format does not matter.


