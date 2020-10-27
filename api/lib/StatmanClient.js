const grpc = require('grpc');
const protoLoader = require('@grpc/proto-loader')

const PROTO_PATH = __dirname + './../../pkg/protos/matchinfo.proto'

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

var client = new statman.Statman('localhost:4000', grpc.credentials.createInsecure());

const getMatch = (matchID) => {
    return new Promise((resolve, reject) => {
        client.GetMatch({matchID}, (err, reply) => {
            if (err) {
                reject(err)
            } else {
                resolve(reply)
            }
        })
    })
}

module.exports = {
    getMatch
}