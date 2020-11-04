<template>
  <div>
    <ul>
      <li v-for="match in this.matchData" :key="match.matchData.matchID">
        {{ JSON.stringify(match) }}
      </li>
    </ul>
  </div>
</template>

<script>
import gql from 'graphql-tag'

const MATCH_DATA_QUERY = gql`
  query getUserMatchData($input: String!) {
    getUserMatchData(userID: $input) {
      matchData {
          matchID
          map
          roundCount
      }
      playerData {
          steamID
          kills
          deaths
          adr
      }
    }
  }
`
export default {
  data() {
    return {
      matchData: Object
    }
  },
  apollo: {
    matchData: {
      query: MATCH_DATA_QUERY,
      variables() {
        return {
          input: "76561198273487074"
        }
      },
      update(data) {
        return data.getUserMatchData
      }
    }
  }
}
</script>