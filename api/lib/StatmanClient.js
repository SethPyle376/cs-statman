const fs = require('fs')
const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader')

const PROTO_PATH = __dirname + '/../protos/matchinfo.proto'

var packageDefinition = protoLoader.loadSync(PROTO_PATH, 
    {
        keepCase: true,
        longs: String,
        enums: String,
        defaults: true,
        oneofs: true
   }
)

var protoDescriptor = grpc.loadPackageDefinition(packageDefinition)
// The protoDescriptor object has the full package hierarchy
var statman = grpc.loadPackageDefinition(packageDefinition).statman

const grpc_host = process.env.STATMAN_HOST || 'localhost'
const grpc_port = process.env.STATMAN_PORT || '4000'

var ssl_creds = {}

if (process.env.SECURE_GRPC === 'true') {
    ssl_creds = grpc.credentials.createSsl(
        fs.readFileSync(process.env.TLS_CA_CERT),
        fs.readFileSync(process.env.TLS_CLIENT_KEY),
        fs.readFileSync(process.env.TLS_CLIENT_CERT)
    )
} else {
    ssl_creds = grpc.credentials.createInsecure()
}

var client = new statman.Statman(`${grpc_host}:${grpc_port}`, ssl_creds);

const getMatch = (matchID) => {
    return new Promise((resolve, reject) => {
        client.GetMatch({matchID}, (err, reply) => {
            if (err) {
                reject(err)
                console.error("ERROR: " + err)
            } else {
                resolve(reply)
            }
        })
    })
}

const getUsersMatchIDs = (playerID) => {
    return new Promise((resolve, reject) => {
        client.GetPlayerMatchIDs({ playerID }, (err, reply) => {
            if (err) {
                reject(err)
            } else {
                resolve(reply)
            }
        })
    })
}

const getUserMatchData = (playerID) => {
    return new Promise((resolve, reject) => {
        client.GetPlayerMatchData({ playerID }, (err, reply) => {
            if (err) {
                reject(err)
            } else {
                resolve(reply)
            }
        })
    })
}

const getRecentMatches = () => {
    return new Promise((resolve, reject) => {
        client.GetRecentMatches({}, (err, reply) => {
            if (err) {
                reject(err)
            } else {
                resolve(reply)
            }
        })
    })
}

const getPopularPlayers = () => {
    return new Promise((resolve, reject) => {
        client.getPopularPlayers({}, (err, reply) => {
            if (err) {
                console.log(err)
                reject (err)
            } else {
                resolve(reply)
            }
        })
    })
}

module.exports = {
    getMatch,
    getUsersMatchIDs,
    getUserMatchData,
    getRecentMatches,
    getPopularPlayers
}