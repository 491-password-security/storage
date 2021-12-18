var crypto = require('../crypto.js');
var ot = require('../ot.js');

const MOD = crypto.constants.MOD;
const GEN = crypto.constants.GEN;
const Number = crypto.Number;

const ONE = new Number('1');

function OT(choice, m_0, m_1) {
    let receiver = new ot.ObliviousTransferReceiver(choice, null, null);
    let sender = new ot.ObliviousTransferSender(m_0, m_1, null, null);

    let C = sender.C;

    receiver.generateKeys(C);

    let receiverKey = receiver.keys[receiver.choice];

    sender.generateKeys(receiverKey);

    var [e_0, e_1] = sender.encryptMessages();

    let result = receiver.readMessage([e_0, e_1]);

    return result; // returns Number
}

function F(k, bits) {
    let exp = new Number('1');
    for (var i = 0; i < 256; i++) {
        if (bits[i] == '1') {
            exp = exp.multiply(k[i]).mod(MOD);
        }
    }
    return GEN.modPow(exp, MOD);
}

function OPRF(k,bits) {
    let a = crypto.util.generatePRFKey(256);

    let client_prod = new Number('1');
    let server_prod = new Number('1');
    for (var i = 0; i < 256; i++) {
        let m_0 = a[i];
        let m_1 = a[i].multiply(k[i]).mod(MOD);
        
        server_prod = server_prod.multiply(m_0).mod(MOD);

        let client_reveal = OT(parseInt(bits[i]), m_0, m_1);

        client_prod = client_prod.multiply(client_reveal).mod(MOD);
    }

    let server_prod_inv = server_prod.modInverse(MOD);

    let exp = server_prod_inv.multiply(client_prod).mod(MOD);
    return GEN.modPow(exp, MOD);
}

let pwd = 'helloworld';
let x = new Number(crypto.util.hash(pwd), 16);
let k = crypto.util.generatePRFKey(256);
let bits = crypto.codec.hex2Bin(x.hex);

// console.log(F(k, bits).decimal);
// console.log(OPRF(k, bits).decimal);

for (var i = 0; i < 256; i++) {
    console.log("new Number(" + "\"" + k[i].hex + "\"),");
}