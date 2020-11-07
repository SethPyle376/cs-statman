<template>
  <div style="max-height: 100%" class="overflow-y-auto">
    <v-toolbar dark class="primary--text" color="#3c3836">Popular Players</v-toolbar>
    <v-list three-line>
      <template v-for="item in this.playerData">
        <v-list-item :key="item.userID" :to="'/app/' + item.userID + '/dashboard'">
          <v-list-item-content>
            <v-list-item-title class="primary--text">{{ item.name }}</v-list-item-title>
            <v-list-item-subtitle class="primary--text">Appears in {{ item.count }} matches</v-list-item-subtitle>
          </v-list-item-content>
        </v-list-item>
      </template>
    </v-list>
  </div>
</template>

<script>
import gql from 'graphql-tag'

const POPULAR_PLAYERS_QUERY = gql`
  query popularPlayers {
    getPopularPlayers {
      userID,
      count,
      name
    }
  }
`

export default {
  data() {
    return {
      playerData: []
    }
  },
  methods: {
    getPlayers: function() {
      try {
        this.$apollo.query({ query: POPULAR_PLAYERS_QUERY}).then(data => this.playerData = data.data.getPopularPlayers).catch(error => console.log(error))
      } catch (error) {
        console.log(error)
      }
    }
  },
  created: function () {
    this.getPlayers()
  }
}
</script>