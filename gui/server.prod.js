const express = require('express');
const path = require('path');
const http = require('http');

const app = new express();

app.use('/', express.static(path.join(__dirname, "build")))

app.get('/*', (req, res) => {   
    res.sendFile(path.join(__dirname, "dist/index.html"));
})

const server = http.createServer(app).listen(3002);
console.log(`Simulator started at port 3002...`);
