<template>
  <div>
    <h1>Game</h1>
    <div>
      {{ positions }}
    </div>
    <div id="container"></div>
  </div>
</template>

<script>
import { mapActions, mapGetters } from 'vuex';
import * as Three from 'three';
export default {
  name: 'Game',
  static() {
    return {
      camera: null,
      scene: null,
      renderer: null,
      mesh: null,
      now: undefined,
      createdBox: false
    };
  },
  computed: {
    ...mapGetters({
      positions: 'positions'
    })
  },
  methods: {
    ...mapActions({
      changeButtonState: 'changeButtonState'
    }),
    init() {
      const container = document.getElementById('container');

      this.camera = new Three.PerspectiveCamera(70, container.clientWidth / container.clientHeight, 0.01, 10);
      this.camera.position.z = 1;
      this.scene = new Three.Scene();
      this.addBox(0.1, 0.1, 0.1, 0, 0.1, 0.1);
      this.renderer = new Three.WebGLRenderer({ antialias: true });
      this.renderer.setSize(container.clientWidth, container.clientHeight);
      container.appendChild(this.renderer.domElement);
    },
    animate(timeStamp) {
      if (this.start === undefined) {
        this.start = timeStamp;
      }
      const elapsed = timeStamp - this.start;
      requestAnimationFrame(this.animate);
      this.mesh.rotation.x += 0.01;
      this.mesh.rotation.y += 0.02;
      if (elapsed > 3000 && !this.createdBox) {
        this.addBox(0.2, 0.2, 0.2, -0.1, -0.1, -0.1);
        this.createdBox = true;
      }
      this.renderer.render(this.scene, this.camera);
    },
    addBox(x, y, z, posX, posY, posZ) {
      const geometry = new Three.BoxGeometry(x, y, z);
      const material = new Three.MeshNormalMaterial();

      this.mesh = new Three.Mesh(geometry, material);
      this.mesh.position.x = posX;
      this.mesh.position.y = posY;
      this.mesh.position.z = posZ;
      this.scene.add(this.mesh);
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
    this.animate();
  },
  unmounted() {
    document.removeEventListener('keydown', this.keyDown);
    document.removeEventListener('keyup', this.keyUp);
  }
};
</script>

<style scoped>
  #container {
    width: 100%;
    height: 75vh;
  }
</style>
