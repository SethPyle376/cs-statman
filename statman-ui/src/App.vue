<template>
  <v-app style="background: #282828">
    <router-view/>
  </v-app>
</template>

<script>

export default {
  name: 'App',

  components: {},

  data: () => ({
    //
  }),
};
</script>

<style scoped>
.v-application {
  background-color: #00a86b;
}
</style>
