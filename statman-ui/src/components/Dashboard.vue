<template>
  <div>
    <v-card class="mx-auto my-12" max-width="400" color="#3c3836">
      <v-card-text class="primary--text">Maps Played</v-card-text>
      <MapDistributionChart v-if="loaded" :matchData="this.loadedData" />
    </v-card>
    <v-card class="mx-auto my-12" max-width="400" color="#3c3836">
      <v-card-text class="primary--text">Total Kills/Deaths by Map</v-card-text>
      <MapStats v-if="loaded" :matchData="this.loadedData" />
    </v-card>
  </div>
</template>

<script>
import gql from 'graphql-tag'
import MapDistributionChart from './charts/MapDistributionChart'
import MapStats from './charts/MapStats'

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
  components: {
    MapDistributionChart,
    MapStats
  }, 
  data() {
    return {
      loadedData: {},
      loaded: false
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
        this.loadedData = data.getUserMatchData
        this.loaded=true
        return data.getUserMatchData
      }
    }
  }
}
</script>