<template>
  <div>
    <div class="header-div">
      <h1>Game</h1>
      <button @click="leaveRoom">Leave room</button>
    </div>
    <div id="container"></div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import FoxChickenScene from '../game3d/foxchickensScene';
export default {
  name: 'Game',
  static() {
    return {
      foxChickenScene: null
    };
  },
  computed: {
    ...mapGetters({
      positions: 'positions',
      username: 'username'
    })
  },
  methods: {
    ...mapActions({
      changeButtonState: 'changeButtonState',
      leaveRoom: 'leaveRoom'
    }),
    init() {
      const container = document.getElementById('container');
      this.foxChickenScene = new FoxChickenScene(
        container,
        this.positions,
        this.username
      );
    },
    animate(timeStamp) {
      // const elapsed = timeStamp - this.start;
      requestAnimationFrame(this.animate);
      this.foxChickenScene.update(timeStamp, this.positions);
    },

    keyboardHandler(event, pressed) {
      const arrows = code => {
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
          case 'KeyE':
          case 'Space':
            this.changeButtonState({ button: 'action', isPressed: pressed });
            break;
          default:
            this.$showLog && console.log(code);
        }
      };
      switch (event.code) {
        case 'Space':
          break;
        default:
          arrows(event.code);
      }
    },
    keyDown(event) {
      this.keyboardHandler(event, true);
    },
    keyUp(event) {
      this.keyboardHandler(event, false);
    }
  },
  mounted() {
    document.addEventListener('keydown', this.keyDown);
    document.addEventListener('keyup', this.keyUp);
    this.init();
    requestAnimationFrame(this.animate);
  },
  unmounted() {
    document.removeEventListener('keydown', this.keyDown);
    document.removeEventListener('keyup', this.keyUp);
    this.animating = false;
  }
};
</script>

<style scoped>
#container {
  width: 100%;
  height: 75vh;
}
.header-div {
  display: flex;
  justify-content: space-between;
}
.header-div > button {
  margin-left: auto;
  border-radius: 15px;
  height: 60px;
  align-self: center;
  background-color: cyan;
}
</style>
