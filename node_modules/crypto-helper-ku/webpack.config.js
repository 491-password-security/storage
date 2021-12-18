const path = require('path');

module.exports = {
    entry: './crypto.js',
    resolve: {
        fallback: { 
            "crypto": require.resolve("crypto-browserify"),
            "buffer": require.resolve("buffer/"),
            "stream": require.resolve("stream-browserify")
         }
    },
    output: {
        path: path.resolve(__dirname, 'dist'),
        filename: 'crypto.js',
        library: {
            name: 'crypto-helper',
            type: 'umd'
        }
    }
}