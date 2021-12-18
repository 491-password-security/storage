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

// key to be used in PRF
let key = crypto.util.generatePRFKey(256);

// random values to facilitate OPRF
let a = crypto.util.generatePRFKey(256);

let server_prod = new Number('1')
for (let i = 0; i < 256; i++) {
    let m_0 = a[i];
    let m_1 = a[i].multiply(k[i]).mod(MOD);
    server_prod = server_prod.multiply(m_0).mod(MOD);

    let sender = new ot.ObliviousTransferSender(m_0, m_1, 
        sendCallback, 
        receiveCallback
    );
    sender.start();
}

// send inverse of server_prod to client
let server_prod_inv = server_prod.modInverse(MOD);
sendCallback(address, server_prod_inv);