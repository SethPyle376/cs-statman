<template>
  <div style="max-height: 100%" class="overflow-y-auto">
    <v-toolbar dark class="primary--text" color="#3c3836">Recent Matches</v-toolbar>
    <v-list three-line>
      <template v-for="item in this.matchData">
        <v-list-item :key="item.matchID">
          <v-list-item-content>
            <v-list-item-title class="primary--text">{{ item.map }}</v-list-item-title>
            <v-list-item-subtitle class="primary--text">{{ item.date }}</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-list>
  </div>
</template>

<script>
import gql from 'graphql-tag'

const RECENT_MATCHES_QUERY = gql`
  query recentMatches {
    getRecentMatches {
      matchID,
      map,
      date
    }
  }
`

export default {
  data() {
    return {
      matchData : []
    }
  },
  methods: {
    getMatches: function() {
      this.$apollo.query({
        query: RECENT_MATCHES_QUERY
      }).then(data => {
        console.log(data)
        this.matchData = data.data.getRecentMatches
      })
    }
  },
  created: function() {
    this.getMatches()
  }
}
</script>