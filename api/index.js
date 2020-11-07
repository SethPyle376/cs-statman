const { loadSchemaSync, GraphQLFileLoader, addResolversToSchema } = require('graphql-tools')
const { join } = require('path')
const express = require('express')
const { graphqlHTTP } = require('express-graphql')
const { merge } = require('lodash')
const cors = require('cors')
require('dotenv').config()

const {resolvers: MatchResolvers} = require('./resolvers/MatchResolvers')

const schema = loadSchemaSync(join(__dirname, './graphql/*.graphql'), {
    loaders: [
        new GraphQLFileLoader()
    ]
})

const resolversWithSchema = addResolversToSchema({
    schema,
    resolvers: merge(MatchResolvers)
})

const app = express()

app.use(cors(),
    graphqlHTTP({
        schema: resolversWithSchema
    })
)

app.listen(process.env.PORT || 8082)