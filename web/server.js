const express = require('express');
const proxy = require('express-http-proxy');

const webpack = require('webpack');
const webpackDevMiddleware = require('webpack-dev-middleware');

const config = require('./webpack.config.js');
const compiler = webpack(config);

const app = express();

app.use('/api', proxy("localhost:8000", {
    proxyReqPathResolver: (req) => {
        return req.url;
    },
}));

app.use(express.static("dist"))
app.use(webpackDevMiddleware(compiler, {
  publicPath: config.devServer.static.directory
}));

app.listen(3000, function () {
  console.log('Application listening on port 3000!');
});