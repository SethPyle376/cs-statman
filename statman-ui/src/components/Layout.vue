<template>
  <div>
    <v-app-bar app clipped-left>
      <v-app-bar-nav-icon @click.stop="drawer = !drawer" color="primary"/>
      <v-toolbar-title class="primary--text">Statman</v-toolbar-title>
    </v-app-bar>
    <v-navigation-drawer v-model="drawer" app clipped>
    </v-navigation-drawer>
    <v-main>
      <router-view/>
    </v-main>
  </div>
</template>

<script>
export default {
  data() {
    return {
      drawer: false
    }
  }
}
</script>