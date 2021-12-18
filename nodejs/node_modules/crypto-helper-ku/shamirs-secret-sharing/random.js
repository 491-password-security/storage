var sjcl = require('../sjcl')

function sjcl_random(bits, returnBits=false) {
  sjcl.random.addEntropy(Math.random(), bits, 'Math.random()')
  var rand = sjcl.random.randomWords(bits/32);
  return (returnBits) ? rand : sjcl.codec.hex.fromBits(rand);
}

function randomBytes(size) {
  return Buffer.from(sjcl_random(size*32, 'hex'))
}

function random(size) {
  const r = randomBytes(32 + size)
  return r.slice(32)
}

module.exports = {
  random
}

