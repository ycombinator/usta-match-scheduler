const express = require('express');
const proxy = require('express-http-proxy');

const app = express();

const apiHostPort = process.env.API_HOST_PORT || "localhost:8000"
app.use('/api', proxy(apiHostPort, {
    proxyReqPathResolver: (req) => {
        return req.url;
    },
}));

app.use(express.static("dist"))

app.listen(3000, function () {
  console.log('Application listening on port 3000!');
});