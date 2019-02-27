const express = require('express');
const path = require('path');
const http = require('http');

const webpack = require("webpack");
const webpackConfig = require("./webpack.config.dev.js");
const compiler = webpack(webpackConfig);

const app = new express();

app.use(
    require("webpack-dev-middleware")(compiler, {
        noInfo: true,
        publicPath: "/"
    })
);

app.use(require("webpack-hot-middleware")(compiler));

app.use('/', express.static(path.join(__dirname, "resources")))

http.createServer(app).listen(3002);
console.log(`Simulator started at port 3002...`);
