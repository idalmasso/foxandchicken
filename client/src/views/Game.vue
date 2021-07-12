<template>
  <h1>Game</h1>
  <div>
    {{ positions }}
    <button
      @mousedown="setAcceleration({ accelX: 0, accelY: 1 })"
      @mouseup="setAcceleration({ accelX: 0, accelY: 0 })"
    >
      up
    </button>
    <button
      @mousedown="setAcceleration({ accelX: 1, accelY: 0 })"
      @mouseup="setAcceleration({ accelX: 0, accelY: 0 })"
    >
      toleft
    </button>
    <button
      @mousedown="setAcceleration({ accelX: -1, accelY: 0 })"
      @mouseup="setAcceleration({ accelX: 0, accelY: 0 })"
    >
      right
    </button>
    <button
      @mousedown="setAcceleration({ accelX: 0, accelY: -1 })"
      @mouseup="setAcceleration({ accelX: 0, accelY: 0 })"
    >
      down
    </button>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';

export default {
  computed: {
    ...mapGetters({
      positions: 'positions'
    })
  },
  methods: {
    ...mapActions({
      setAcceleration: 'setAcceleration'
    })
  }
};
</script>

<style></style>
