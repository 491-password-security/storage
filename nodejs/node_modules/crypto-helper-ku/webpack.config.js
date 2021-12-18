const path = require('path');

module.exports = {
    entry: './crypto.js',
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'crypto.js',
        library: {
            name: 'crypto-helper',
            type: 'umd'
        }
    }
}