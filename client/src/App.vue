<template>
  <v-app>
    <v-app-bar app>
      Vanilla
    </v-app-bar>

    <v-content app>
      <v-container>
        <router-view></router-view>
      </v-container>
    </v-content>

    <v-footer>
      <v-bottom-navigation app>
        <v-btn to="/">
          <span>Home</span>
          <v-icon>fa-home</v-icon>
        </v-btn>
        <v-btn>
          <span>Heroes</span>
          <v-icon>fa-helmet-battle</v-icon>
        </v-btn>
        <v-btn>
          <span>Explore</span>
          <v-icon>fa-map-signs</v-icon>
        </v-btn>
        <v-btn>
          <span>Friends</span>
          <v-icon>fa-users-crown</v-icon>
        </v-btn>
        <v-btn>
          <span>Heroes</span>
          <v-icon>mdi-heart</v-icon>
        </v-btn>
      </v-bottom-navigation>
    </v-footer>
  </v-app>
</template>

<script>
// import HelloWorld from './components/HelloWorld'

export default {
  name: 'App',

  components: {
    // HelloWorld
  },

  watch: {
    '$router' (to, from) {
      console.log('router changed to', to)
    }
  },

  data () {
    return {
      activeBtn: 'Heroes'
    }
  }
}
</script>

<style scoped>
  v-app{
    width: 100%;
    height: 100%;
  }
</style>
