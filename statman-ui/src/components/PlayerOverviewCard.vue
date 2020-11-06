<template>
  <v-card width=350px height=500px class="mx-4 mb-16" color="#3c3836">
    <v-card-text class="primary--text">Player Overview</v-card-text>
    <v-list color="#3c3836">
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="title primary--text">
            {{ this.name }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
    <v-divider/>
    <v-list nav dense color="#3c3836">
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="primary--text">
            Appears in: {{ this.matchCount }} matches
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="primary--text">
            Kills: {{ this.killCount }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="primary--text">
            Deaths: {{ this.deathCount }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="primary--text">
            KDR: {{ this.kdr }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
      <v-list-item>
        <v-list-item-content>
          <v-list-item-title class="primary--text">
            ADR: {{ this.adr }}
          </v-list-item-title>
        </v-list-item-content>
      </v-list-item>
    </v-list>
  </v-card>
</template>

<script>
export default {
  props: {
    matchData: {
      type: Array
    }
  },
  computed: {
    name: function() {
      return this.matchData[0].playerData.name
    },
    matchCount: function() {
      return this.matchData.length
    },
    killCount: function() {
      var accum = 0
      this.matchData.forEach(match => {
        accum += match.playerData.kills
      })
      return accum
    },
    deathCount: function() {
      var accum = 0
      this.matchData.forEach(match => {
        accum += match.playerData.deaths
      })
      return accum
    },
    kdr: function() {
      return parseFloat(this.killCount / this.deathCount).toFixed(2)
    },
    adr: function() {
      var adrAccum = 0
      var mapCount = 0
      this.matchData.forEach(match => {
        adrAccum += match.playerData.adr
        mapCount++
      })
      return parseFloat(adrAccum / mapCount).toFixed(2)
    }
  }
}
</script>