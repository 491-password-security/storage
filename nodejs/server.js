
'use strict';

const express = require('express');
const crypto = require('crypto-helper-ku');
const {MongoClient} = require('mongodb');

//const uri = "mongodb://mongo:mongo@mongo:27017/?authSource=admin";
const uri = process.env.MONGO_URI;
console.log(uri);

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
        console.log(users);
    } catch (err) {
        console.error(err);
    }
})()

// Constants
const PORT = 8080;
const HOST = '0.0.0.0';


// App
const app = express();
app.get('/save-password-share/:key/:value', (req, res) => {
    
    col.insertOne({
        key: req.params.key,
        value: req.params.value
    });

    res.send(req.params.key + ':' + req.params.value);
});

app.get('/get-password-share/:key', (req, res) => {

    col.findOne({
        key: req.params.key
    }, (err, doc) => {
        if (err) {
            console.log(err);
        } else {
            res.send(doc.value);
        }
    });
});

app.listen(PORT, HOST);
console.log(`Running on http://${HOST}:${PORT}`);