const express = require('express')
const app = express()
const port = 8000

app.get('/thisroute', (req, res) => {
    console.log('okayy')
    res.json({"arst": " arst"})
})

app.get('/okayy', (req, res) => {
    console.log('okayarsty')
    res.json({"ADMIN_SECRET": " youtube.com"})
})
app.get('/', (req, res) => {
    console.log('heae')
    res.send('Hello World!')
})

app.listen(port, () => {
    console.log(`Example app listening at http://localhost:${port}`)
})