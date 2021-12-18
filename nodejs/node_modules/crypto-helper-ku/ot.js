var crypto = require('./crypto.js');

// OT should work by initiaiting the sender and receiver and then just calling start for each

let Number = crypto.Number;

const MOD = crypto.constants.MOD;
const GEN = crypto.constants.GEN;

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