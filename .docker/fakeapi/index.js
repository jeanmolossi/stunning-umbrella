const express = require('express');

const app = express();
app.use(express.json())

const routes = express.Router()

routes.get('/ping', function (req, res) {
    return res.status(200).json({ message: "pong" })
})

async function delay(sleepTime) {
    await new Promise((resolve => {
        setTimeout(resolve, sleepTime)
    }))
}

routes.put('/api/recruiter/:id/access-level', async function (req, res) {
    const {id} = req.params

    console.log("Receiving request to ID ", id)

    const { nivelAcesso } = req.body

    const delayTime = Math.round(Math.random() * 500)

    if(!['admin','recruiter'].includes(nivelAcesso)){
        return res.status(402).json({ message: "Bad request" })
    }

    const intId = +id;

    const notFoundChance = Math.round(Math.random() * 10);

    if (notFoundChance <= 1) {
        await delay(delayTime)
        return res.status(404).json({
            message: "Not found"
        })
    }

    if (intId === 2) {
        await delay(5000)
        return res.status(504).json({
            message: "Bad Gateway"
        })
    } else if (intId === 4) {
        await delay(delayTime)
        return res.status(422).json({message: "Unprocessable Entity"})
    } else {
        await delay(delayTime)
        console.log("Responds to request ID ", id)
        return res.status(200).json({message: "Ok, foi"})
    }
})

app.use(routes);

app.listen(8080, () => {
    console.log('Running fake api on port 8080')
})