<template>
  <the-application-bar />
  <router-view class="content" />
</template>
<script>
import TheApplicationBar from './components/TheApplicationBar.vue';
export default {
  components: { TheApplicationBar },
  data() {
    return {
      rooms: {
        name: 'Rooms'
      }
    };
  }
};
</script>
<style>
body {
  width: 100%;
  height: 100vh;
  margin: 0;
  padding: 0;
  border: 0;
  font-size: 100%;
  font: inherit;
  vertical-align: baseline;
}
#app {
  position: relative;
  padding: 2px;
}
.content {
  width: 80%;
  margin-right: 15%;
  margin-left: 15%;
}
</style>
