<template>
  <h1>Game</h1>
  <div>
      {{ positions }}
  </div>
</template>

<script>
import { mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters({
      positions: 'positions'
    })
  }
};
</script>

<style>

</style>
