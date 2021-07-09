<template>
  <the-application-bar/>
  <router-view/>
</template>
<script>
import TheApplicationBar from './components/TheApplicationBar.vue'
export default {
  components: { TheApplicationBar },
  data () {
    return {
      rooms: {
        name: 'Rooms'
      }
    }
  }
}
</script>
<style>
#app {
  margin: 0;
  padding: 0;
  border: 0;
  font-size: 100%;
  font: inherit;
  vertical-align: baseline;
}

#nav {
  padding: 30px;
}

#nav a {
  font-weight: bold;
  color: #2c3e50;
}

#nav a.router-link-exact-active {
  color: #42b983;
}
</style>
