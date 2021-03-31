'use strict';
var createError = require("http-errors");
const express = require('express');
const cors = require("cors");
const port = process.env.SERVER_PORT;

const testAPIRouter = require("./routes/testAPI");

var app = express();

app.use(cors()); // hopefully we don't need

app.get('/', function(req, res) {
    res.send('Hello, World :)');
});

app.use("/testAPI", testAPIRouter);

app.listen(port, function() {
    console.log('listening on port ' + port);
});