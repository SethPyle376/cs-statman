const { getMatch, getUsersMatchIDs, getUserMatchData, getRecentMatches, getPopularPlayers } = require('./../lib/StatmanClient')

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
        },
        getUserMatchData: async (parent, args, context, info) => {
            const matchData = await getUserMatchData(args.userID)
            return matchData.playerMatchData
        },
        getRecentMatches: async (parent, args, context, info) => {
            const matchData = await getRecentMatches()
            return matchData.matchData
        },
        getPopularPlayers: async (parent, args, context, info) => {
            const data = await getPopularPlayers()
            return data.popularPlayerData
        }
    }
}