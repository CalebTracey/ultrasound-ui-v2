/* eslint-disable @typescript-eslint/no-var-requires */
const express = require('express')
const cors = require('cors')
const path = require('path')

const app = express()
const PORT = process.env.PORT || 80

app.use(
    cors({
        origin: [
            'http://localhost:80',
            'https://ultrasound-api.herokuapp.com',
            'https://ultrasound-ui.herokuapp.com',
        ],
        methods: ['GET', 'POST', 'DELETE', 'UPDATE', 'PUT', 'PATCH'],
    })
)

const server = app.listen(PORT, () => {
    const { port } = server.address()
    console.log(`Server listening on port ${port}`)
})

app.use(express.static(path.join(__dirname, 'build')))

// This route serves the React app
app.get('/*', (req, res) => {
    console.log(__dirname)
    res.sendFile(path.resolve(__dirname, 'build', 'index.html'))
})
