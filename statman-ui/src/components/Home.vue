<template>
  <div>
    <v-row>
      <v-card width=350px height=500px class="mx-4 mb-16">
        <PopularPlayers/>
      </v-card>
      <v-card width=350px height=500px class="mx-4 mb-16">
        <RecentMatches/>
      </v-card>
    </v-row>
  </div>
</template>

<script>
import RecentMatches from './RecentMatches'
import PopularPlayers from './PopularPlayers'
export default {
  components: {
    RecentMatches,
    PopularPlayers
  }
}
</script>