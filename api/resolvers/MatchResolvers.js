const { getMatch } = require('./../lib/StatmanClient')

module.exports.resolvers = {
    Query: {
        getMatch: async (parent, args, context, info) => {
            const match = await getMatch(args.matchID)
            console.log(match)
            return match.matchInfo
        }
    }
}