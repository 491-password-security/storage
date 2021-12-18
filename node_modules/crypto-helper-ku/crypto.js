global.Buffer = global.Buffer || require('buffer').Buffer

var sjcl = require('./sjcl');
var secrets = require('./shamirs-secret-sharing');
const { BigInteger } = require('jsbn');

class Number {
    constructor(str="1", radix=10) {
        this.bigInt = new BigInteger(str, radix);
    }

    get hex() {
        return this.bigInt.toString(16);
    }
    set hex(val) {
        this.bigInt = new BigInteger(val, 16);
    }

    get bytes() {
        return bigInt2Bytes(this.bigInt);
    }
    set bytes(val) {
        this.bigInt = new BigInteger(bytes2Hex(val), 16);
    }

    get decimal() {
        return this.bigInt.toString(10);
    }

    divide(other) {
        return new Number(this.bigInt.divide(other.bigInt).toString());
    }

    modPow(other, mod) {
        return new Number(this.bigInt.modPow(other.bigInt, mod.bigInt).toString());
    }

    mod(other) {
        return new Number(this.bigInt.mod(other.bigInt).toString());
    }

    modInverse(mod) {
        return new Number(this.bigInt.modInverse(mod.bigInt).toString());
    }

    multiply(other) {
        return new Number(this.bigInt.multiply(other.bigInt).toString());
    }

    compareTo(other){
        return this.bigInt.compareTo(other.bigInt);
    }

    subtract(other) {
        return new Number(this.bigInt.subtract(other.bigInt).toString());
    }

    toString() {
        return this.decimal
    }
}

const BIG_TWO = new Number('2');
const BIG_ONE = new Number('1');

// const MOD = new Number('104334873255401717971305551311108568981602782554133676271604158174023613565338436519535738349159664075981513545995816898351274759273689547803611869080590323788134546218679576525351375421659491479861062524332418185137628175629882792848502958254366030986728999054034830850220407425928535174607722203029578103539');
// const GEN = new Number('15309569078288033140294527228325069587420150399530450735556668091277116408023136181284430449588830517258893721878398739530623279778683647761572205172467420662396761999763043433000129229039419004108765113420973429371572791200022523422170732284615282345655002021445578558188416639692531759416866286539604862128');

const MOD = new Number('ffffffffffffffffc90fdaa22168c234c4c6628b80dc1cd129024e088a67cc74020bbea63b139b22514a08798e3404ddef9519b3cd3a431b302b0a6df25f14374fe1356d6d51c245e485b576625e7ec6f44c42e9a63a3620ffffffffffffffff', 16)
const GEN = new Number('2', 10); 


function random(bits, returnBits=false) {
    var rand = sjcl.random.randomWords(bits/32);
    return (returnBits) ? rand : sjcl.codec.hex.fromBits(rand);
}

function hash(input, returnBits=false) {
    var out = sjcl.hash.sha256.hash(input);
    return (returnBits) ? out : sjcl.codec.hex.fromBits(out);
}

function extendedHash(input, count) {
    let last_output = input.hex;
    let result = [];
    for (var i = 0; i < count; i++) {
        last_output = hash(last_output);
        result.push(last_output);
    }
    return new Number(result.join(''), 16);
}

function generatePRFKey(count) {
    let key = [];
    for (var i = 0; i < count; i++) {
        let pow = getBoundedBigInt(MOD);
        key.push(GEN.modPow(pow, MOD));
    }        
    return key;
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

function share(secret, t, n) {
    let hex_shares = [];
    let shares = secrets.split(Buffer.from(secret), { shares: n, threshold: t });
    for (let i = 0; i < shares.length; i++) {
        hex_shares.push(shares[i].toString('hex'));
    }
    return hex_shares;
}

function combine(shares, encoding='hex') {
    return secrets.combine(shares).toString(encoding);
}

function getBoundedBigInt(max) {
    let bits = max.bigInt.bitLength();
    let number = new Number(null, null);
    do {
        number.bigInt = new BigInteger(random(bits));
    } while (number.bigInt.compareTo(max) >= 0);
    return number;
}

async function getElGamalKeys(bits) {
    var eg = await elgamal.default.generateAsync(bits);
    return {
        p: eg.p,
        g: eg.g,
        x: eg.x,
        g_x: eg.y,
    };
}

function xor(u, v) {
    let length = Math.min(u.bytes.length, v.bytes.length);
    let resultNum = new Number()
    var result = [];
    for (var i = 0; i < length; i++) {
        result.push(u.bytes[i] ^ v.bytes[i]);   
    }
    resultNum.bytes = result;
    return resultNum;
}   

function hex2Bin(hex){
    var out = "";
    for(var c of hex) {
        switch(c) {
            case '0': out += "0000"; break;
            case '1': out += "0001"; break;
            case '2': out += "0010"; break;
            case '3': out += "0011"; break;
            case '4': out += "0100"; break;
            case '5': out += "0101"; break;
            case '6': out += "0110"; break;
            case '7': out += "0111"; break;
            case '8': out += "1000"; break;
            case '9': out += "1001"; break;
            case 'a': out += "1010"; break;
            case 'b': out += "1011"; break;
            case 'c': out += "1100"; break;
            case 'd': out += "1101"; break;
            case 'e': out += "1110"; break;
            case 'f': out += "1111"; break;
            default: return "";
        }
    }
    return out;
}

function hex2Bytes(hex) {
    if (hex.length % 2 != 0) {
        hex = '0' + hex;
    }
    return sjcl.codec.bytes.fromBits(sjcl.codec.hex.toBits(hex));
}

function bytes2Hex(byteArray) {
    return sjcl.codec.hex.fromBits(sjcl.codec.bytes.toBits(byteArray));
  }

function bytes2BigInt(bytes) {
    return new BigInteger(bytes2Hex(bytes), 16);
}

function bigInt2Bytes(bigInt) {
    return hex2Bytes(bigInt.toString(16));
}

module.exports.constants = {MOD, GEN};
module.exports.ss = {share, combine};
module.exports.util = {random, hash, extendedHash, getBoundedBigInt, getElGamalKeys, xor, generatePRFKey};
module.exports.aes = {encrypt, decrypt};
module.exports.codec = {hex2Bytes, hex2Bin, bytes2Hex, bytes2BigInt, bigInt2Bytes}
module.exports.Number = Number;

// OT

module.exports.ObliviousTransferReceiver = class ObliviousTransferReceiver {
    constructor(choice, sendCallback, receiveCallback) {
        if (choice != 0 && choice != 1) {
            throw new Error('Choice neither 1 nor 0. Enter a single integer (0 or 1) as the choice.');
        }
        this.sendCallback = sendCallback;
        this.receiveCallback = receiveCallback;
        this.choice = choice;
        this.keys = [];
        
        let temp = crypto.util.getBoundedBigInt(MOD);        
        this.k = GEN.modPow(temp, MOD);
    }

    start(address) {
        // get the random constant C from the sender
        let C = new Number(this.receiveCallback(), 16);

        // generate two keys and send the valid key to the sender
        this.generateKeys(C);
        this.sendCallback(address, this.keys[this.choice].hex);

        // receive the two encryptions from the sender
        let choices = this.receiveCallback();

        // decrypt the chosen message
        return this.readMessage(choices);
    }

    generateKeys(C) {
        // generate two random keys (as elements from multiplicative Z_p) also using C
        let choiceKey = GEN.modPow(this.k, MOD);
        let negChoiceKey = choiceKey.modInverse(MOD).multiply(C).mod(MOD);
        this.keys = [choiceKey, negChoiceKey];
    }

    readMessage(choices) {
        // choose one of the messages
        let pair = choices[this.choice];
        let hint = pair[0];
        let ciphertext = pair[1];

        // g^(r_sigma)^k = PK_sigma^(r_sigma)
        let key = hint.modPow(this.k, MOD);
        let xorKey = crypto.util.extendedHash(key, 4);

        let result = ciphertext.multiply(GEN.modPow(key, MOD).modInverse(MOD)).mod(MOD);

        // decrypt the ciphertext
        // return crypto.util.xor(xorKey, ciphertext);
        return result;
    }
}

module.exports.ObliviousTransferSender = class ObliviousTransferSender {
    constructor(m_0, m_1, sendCallback, receiveCallback) {
        this.m_0 = m_0;
        this.m_1 = m_1;
        this.sendCallback = sendCallback;
        this.receiveCallback = receiveCallback;

        // initiate random constants
        this.log_C = crypto.util.getBoundedBigInt(MOD);
        this.C = GEN.modPow(this.log_C, MOD);
        let temp_0 = crypto.util.getBoundedBigInt(MOD);
        let temp_1 = crypto.util.getBoundedBigInt(MOD);
        this.r_0 = GEN.modPow(temp_0, MOD);
        this.r_1 = GEN.modPow(temp_1, MOD);
    }

    start(address) {
        // send the constant value C to the receiver
        this.sendCallback(address, this.C.hex);

        // receive one key from receiver
        let receiverKey = new Number(this.receiveCallback(address), 16);

        // generate two keys based on the received key and the hidden random values
        this.generateKeys(receiverKey);

        // send the encrypted messages to the receiver
        let messages = this.encryptMessages();
        this.sendCallback(address, messages);
    }

    generateKeys(receiverKey) {
        // generate keys for each message based on receiver's key and the hidden random values
        this.key_0 = receiverKey;
        this.key_1 = this.key_0.modInverse(MOD).multiply(this.C).mod(MOD);

        let temp_0 = this.key_0.modPow(this.r_0, MOD);
        let temp_1 = this.key_1.modPow(this.r_1, MOD);

        this.key_0 = temp_0;
        this.key_1 = temp_1;
        this.keys = [this.key_0, this.key_1];
    }

    encryptMessages() {
        let ct_0 = GEN.modPow(this.key_0, MOD).multiply(this.m_0).mod(MOD);
        let ct_1 = GEN.modPow(this.key_1, MOD).multiply(this.m_1).mod(MOD);

        let e_0 = [GEN.modPow(this.r_0, MOD), ct_0];
        let e_1 = [GEN.modPow(this.r_1, MOD), ct_1];

        return [e_0, e_1];
    }
}