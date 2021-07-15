<template>
  <div>
    <h1>Game</h1>
    <div>
      {{ positions }}
    </div>
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
      changeButtonState: 'changeButtonState'
    }),
    keyboardHandler(event, pressed) {
      const arrows = (code) => {
        switch (code) {
          case 'ArrowUp':
          case 'KeyW':
            this.changeButtonState({ button: 'up', isPressed: pressed });
            break;
          case 'ArrowDown':
          case 'KeyS':
            this.changeButtonState({ button: 'down', isPressed: pressed });
            break;
          case 'ArrowLeft':
          case 'KeyA':
            this.changeButtonState({ button: 'left', isPressed: pressed });
            break;
          case 'ArrowRight':
          case 'KeyD':
            this.changeButtonState({ button: 'right', isPressed: pressed });
            break;
          default:
            console.log(code);
        }
      };
      switch (event.code) {
        case 'Space':
          break;
        default:
          arrows(event.code);
      }
    }
  },
  mounted() {
    document.addEventListener('keydown', (event) => this.keyboardHandler(event, true));
    document.addEventListener('keyup', (event) => this.keyboardHandler(event, false));
  },
  unmounted() {
    document.removeEventListener('keydown', (event) => this.keyboardHandler(event, true));
    document.addEventListener('keyup', (event) => this.keyboardHandler(event, false));
  }
};
</script>

<style></style>
