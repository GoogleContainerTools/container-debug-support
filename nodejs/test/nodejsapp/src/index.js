const express = require('express')
const { echo } = require('./utils');
const os = require('os');

const app = express()
const port = 3000

app.get('/', (req, res) => res.send(echo(`Hello from ${os.platform()}/${os.arch()}!\n`)))

app.listen(port, () => console.log(`Example app listening on port ${port}!`))
