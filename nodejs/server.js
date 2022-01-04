const app = require('express')();
const http = require('http').Server(app);
const io = require('socket.io')(http);
const crypto = require('crypto-helper-ku');
const { encrypt } = require('crypto-helper-ku/sjcl');
const port = process.env.PORT || 3000;
const {MongoClient} = require('mongodb');


const Number = crypto.Number;

const MOD = crypto.constants.MOD;
const GEN = crypto.constants.GEN;

const OPRF_KEY = new Number(process.env.ENC_KEY, 16)

app.get('/', (req, res) => {
  res.sendFile(__dirname + '/index.html');
});

io.on('connection', (socket) => {
    socket.on("beginOPRF", (a) => {
        a = new Number(a, 16);
        let b = a.modPow(OPRF_KEY, MOD);
        socket.emit("respondOPRF", b.hex);
    })
});

http.listen(port, () => {
  console.log(`Socket.IO server running at http://localhost:${port}/`);
});

// 'use strict';


// // const uri = "mongodb://mongo:mongo@mongo:27017/?authSource=admin";
const uri = process.env.MONGO_URI;
// console.log(uri);

const client = new MongoClient(uri);
var users;
var col;

(async () => {
    try {
        await client.connect();
        console.log("connected");
        const db = client.db('s3');
        col = db.collection('s4');
        users = await col.find({}).toArray();
    } catch (err) {
        console.error(err);
    }
})()


// App
app.get('/save-password-share/:key/:value/:iv', (req, res) => {
    
    col.insertOne({
        key: req.params.key,
        value: req.params.value,
        iv: req.params.iv
    });

    res.send(req.params.key + ':' + req.params.value + ':' + req.params.iv);
});

app.get('/get-password-share/:key', (req, res) => {

    col.findOne({
        key: req.params.key
    }, (err, doc) => {
        if (err) {
            console.log(err);
            var encrypted = crypto.encrypt(crypto.hash("baran"), req.params.key)
            res.send(
                encrypted.ciphertext + 
                ':' + 
                encrypted.iv);
        } else {
            res.send(doc.value + ':' + doc.iv);
        }
    });
});
