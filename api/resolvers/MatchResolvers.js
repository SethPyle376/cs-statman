const { getMatch, getUsersMatchIDs } = require('./../lib/StatmanClient')

module.exports.resolvers = {
    Query: {
        getMatch: async (parent, args, context, info) => {
            const match = await getMatch(args.matchID)
            if (match) {
                return match.matchInfo
            } else {
                return {}
            }
        },
        getUsersMatches: async (parent, args, context, info) => {
            const matches = await getUsersMatchIDs(args.userID)
            return {
                matchIDs: matches.matchIDs
            }
        }
    }
}