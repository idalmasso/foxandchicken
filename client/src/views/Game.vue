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
      createdBox: false,
      cameraStart: 0,
      lerpDuration: 0,
      isLerping: false,
      vectorEnd: null
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
      while (container.hasChildNodes()) {
        container.removeChild(container.lastChild);
      }
      this.cameraStart = 0;
      this.lerpDuration = 25;
      this.isLerping = false;
      this.vectorEnd = new Three.Vector3();
      this.camera = new Three.PerspectiveCamera(70, container.clientWidth / container.clientHeight, 0.01, 12);
      // this.camera.position.x = 50;
      // this.camera.position.y = 50;
      this.camera.position.z = 10;
      this.scene = new Three.Scene();
      this.addBackground(20, 20);
      this.meshes = [];
      for (const username in this.positions) {
        const position = this.positions[username].position;
        this.addObject(position.x, position.y, username);
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
          this.addObject(position.x, position.y, username);
        } else {
          this.meshes[username].position.x = position.x;
          this.meshes[username].position.y = position.y;
          this.meshes[username].rotation.x += 0.01;
          this.meshes[username].rotation.y += 0.02;
        }
        if (username === this.username) {
          if (!this.isLerping && (this.camera.position.x !== position.x || this.camera.position.y !== position.y)) {
            this.cameraStart = timeStamp;
            this.cameraEnd = this.cameraStart + 20;
            this.vectorEnd.x = position.x;
            this.vectorEnd.y = position.y;
            this.vectorEnd.z = this.camera.position.z;
            this.isLerping = true;
          }
          if (this.cameraStart !== 0) {
            this.camera.position.lerp(this.vectorEnd, (timeStamp - this.cameraStart) / this.lerpDuration);
            this.meshes[username].position.x = this.camera.position.x;
            this.meshes[username].position.y = this.camera.position.y;
            if (timeStamp > this.cameraStart + this.lerpDuration) {
              this.isLerping = false;
            }
          }
        }
      }
      for (const username in this.meshes) {
        if (typeof this.positions[username] === 'undefined') {
          this.$showLog && console.log('REMOVING ' + username);
          this.scene.remove(this.meshes[username]);
          this.meshes.splice(username, 1);
        }
      }
      this.renderer.render(this.scene, this.camera);
    },
    addObject(posX, posY, username) {
      if (username === this.username) {
        this.meshes[username] = this.addBox(0.15, 0.15, 0.15, posX, posY, 0.1);
      } else {
        this.meshes[username] = this.addSphere(0.15, posX, posY, 0.1);
      }
      var canvas = document.createElement('canvas');
      canvas.width = 256;
      canvas.height = 256;
      var ctx = canvas.getContext('2d');
      ctx.font = '44pt Arial';
      ctx.fillStyle = 'white';
      ctx.textAlign = 'center';
      ctx.fillText(username, 128, 44);
      var tex = new Three.Texture(canvas);
      tex.needsUpdate = true;
      var spriteMat = new Three.SpriteMaterial({ map: tex });
      var sprite = new Three.Sprite(spriteMat);
      this.meshes[username].add(sprite);
    },
    addBox(x, y, z, posX, posY, posZ) {
      this.$showLog && console.log('adding a box');
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
      this.$showLog && console.log('adding a sphere');
      const geometry = new Three.SphereGeometry(radius);
      const material = new Three.MeshNormalMaterial();
      const mesh = new Three.Mesh(geometry, material);
      mesh.position.x = posX;
      mesh.position.y = posY;
      mesh.position.z = posZ;
      this.scene.add(mesh);
      return mesh;
    },
    addBackground(sizeX, sizeY) {
      this.$showLog && console.log('adding background');
      const geometry = new Three.BoxGeometry(sizeX, sizeY, 0.1);
      const material = new Three.MeshBasicMaterial({ color: 0x344522, wireframe: false });
      const mesh = new Three.Mesh(geometry, material);
      mesh.position.x = sizeX / 2;
      mesh.position.y = sizeY / 2;
      mesh.position.z = -1;
      this.scene.add(mesh);
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
