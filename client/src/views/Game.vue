<template>
  <div>
    <h1>Game</h1>
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
      meshes: null,
      now: undefined,
      createdBox: false
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
      changeButtonState: 'changeButtonState'
    }),
    init() {
      const container = document.getElementById('container');

      this.camera = new Three.PerspectiveCamera(70, container.clientWidth / container.clientHeight, 0.01, 10);
      this.camera.position.z = 10;
      this.scene = new Three.Scene();
      this.meshes = [];
      for (const username in this.positions) {
        const position = this.positions[username].position;
        if (username === this.username) {
          this.meshes[username] = this.addBox(0.1, 0.1, 0.1, position.x, position.y, 0.1);
        } else {
          this.meshes[username] = this.addSphere(0.2, position.x, position.y, 0.1);
        }
      }
      this.renderer = new Three.WebGLRenderer({ antialias: true });
      this.renderer.setSize(container.clientWidth, container.clientHeight);
      container.appendChild(this.renderer.domElement);
    },
    animate(timeStamp) {
      if (this.start === undefined) {
        this.start = timeStamp;
      }
      // const elapsed = timeStamp - this.start;
      requestAnimationFrame(this.animate);
      for (const username in this.positions) {
        const position = this.positions[username].position;
        if (typeof this.meshes[username] === 'undefined') {
          if (username === this.username) {
            this.meshes[username] = this.addBox(0.1, 0.1, 0.1, position.x, position.y, 0.1);
          } else {
            this.meshes[username] = this.addSphere(0.2, position.x, position.y, 0.1);
          }
        } else {
          this.meshes[username].position.x = position.x;
          this.meshes[username].position.y = position.y;
          this.meshes[username].rotation.x += 0.01;
          this.meshes[username].rotation.y += 0.02;
        }
      }
      for (const username in this.meshes) {
        if (typeof this.positions[username] === 'undefined') {
          console.log('REMOVING ' + username);
          this.scene.remove(this.meshes[username]);
          this.meshes.splice(username, 1);
        }
      }
      this.renderer.render(this.scene, this.camera);
    },
    addBox(x, y, z, posX, posY, posZ) {
      console.log('adding a box' + posX);
      const geometry = new Three.BoxGeometry(x, y, z);
      const material = new Three.MeshNormalMaterial();
      const mesh = new Three.Mesh(geometry, material);
      mesh.position.x = posX;
      mesh.position.y = posY;
      mesh.position.z = posZ;
      this.scene.add(mesh);
      return mesh;
    },
    addSphere(radius, posX, posY, posZ) {
      console.log('adding a sphere' + posX);
      const geometry = new Three.SphereGeometry(radius);
      const material = new Three.MeshNormalMaterial();
      const mesh = new Three.Mesh(geometry, material);
      mesh.position.x = posX;
      mesh.position.y = posY;
      mesh.position.z = posZ;
      this.scene.add(mesh);
      return mesh;
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
