var crypto = require('../crypto.js');
var ot = require('../ot.js');

const MOD = crypto.constants.MOD;
const GEN = crypto.constants.GEN;
const Number = crypto.Number;

// two functions for sending and receiving objects over a network
// arguments can change (to use sockets etc.)
function sendCallback(address, value) {
    // send value to address
    return 1;
}

function receiveCallback() {
    return 1;
}

let pwd = 'hello';

// input to PRF
let x = new Number(crypto.util.hash(pwd), 16);
let bits = crypto.codec.hex2Bin(x.hex);

let client_prod = new Number('1')
for (let i = 0; i < 256; i++) {
    let choice = parseInt(bits[i]);

    let receiver = new ot.ObliviousTransferReceiver(choice,
        sendCallback,
        receiveCallback
    );
    let reveal = receiver.start();
    client_prod = client_prod.multiply(reveal).mod(MOD);
}

// receive inverse of server_prod from server
let server_prod_inv = receiveCallback();

// compute g^(server_prod_inv * client_prod)
let exp = server_prod_inv.multiply(client_prod).mod(MOD);
let result = GEN.modPow(exp, MOD);