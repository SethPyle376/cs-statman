<template>
  <div>
    <v-row>
      <v-card height=500px class="mx-4 mb-16" color="#3c3836">
        <v-card-text class="primary--text">Maps Played</v-card-text>
        <MapDistributionChart v-if="loaded" :matchData="this.loadedData" />
      </v-card>
      <v-card class="mx-4 mb-16" color="#3c3836">
        <v-card-text class="primary--text">Total Kills/Deaths by Map</v-card-text>
        <MapStats v-if="loaded" :matchData="this.loadedData" />
      </v-card>
      <v-card class="mx-4 mb-16" color="#3c3836">
        <v-card-text class="primary--text">ADR by Map</v-card-text>
        <AdrChart v-if="loaded" :matchData="this.loadedData" />
      </v-card>
    </v-row>
  </div>
</template>

<script>
import gql from 'graphql-tag'
import MapDistributionChart from './charts/MapDistributionChart'
import MapStats from './charts/MapStats'
import AdrChart from './charts/AdrByMapChart'

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
    MapStats,
    AdrChart
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
          input: this.$route.params.userID
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