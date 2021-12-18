var sjcl = require('sjcl');
var secrets = require('secrets.js-grempe');


function random(bits, returnBits=false) {
    var rand = sjcl.random.randomWords(bits/32);
    return (returnBits) ? rand : sjcl.codec.hex.fromBits(rand);
}

function hash(input) {
    var out = sjcl.hash.sha256.hash(input);
    return sjcl.codec.hex.fromBits(out);
}

function encrypt(key, plaintext) {
    key = sjcl.codec.hex.toBits(key);
    plaintext = sjcl.codec.hex.toBits(plaintext);

    var aes = new sjcl.cipher.aes(key);
    var iv = random(128, returnBits=true);
    var ciphertext = sjcl.mode.ccm.encrypt(aes, plaintext ,iv);

    return {
        iv: sjcl.codec.hex.fromBits(iv), 
        ciphertext: sjcl.codec.hex.fromBits(ciphertext)
    };
}

function decrypt(key, iv, ciphertext) {
    key = sjcl.codec.hex.toBits(key);
    iv = sjcl.codec.hex.toBits(iv);
    ciphertext = sjcl.codec.hex.toBits(ciphertext);

    var aes = new sjcl.cipher.aes(key);
    var plaintext = sjcl.mode.ccm.decrypt(aes, ciphertext, iv);

    return sjcl.codec.hex.fromBits(plaintext);
}

// secret will already be a hex string, being the output of a hash function
function share(secret, t, n) {
    return secrets.share(secret, n, t);
}

function combine(shares) {
    return secrets.combine(shares);
}

function newShare(id, shares) {
    return secrets.newShare(id, shares);
}

module.exports = {random, hash, encrypt, decrypt, share, combine, newShare};