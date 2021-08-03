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
import HealthBar from '../healthBar/healthBar';
export default {
  name: 'Game',
  static() {
    return {
      camera: null,
      scene: null,
      renderer: null,
      playersGameObjects: null,
      now: undefined,
      createdBox: false,
      cameraStart: 0,
      lerpDuration: 0,
      isLerping: false,
      vectorEnd: null,
      animating: false,
      healthBars: null
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
      this.cameraStart = 0;
      this.lerpDuration = 50;
      this.isLerping = false;
      this.animating = true;
      this.vectorEnd = new Three.Vector3();
      this.camera = new Three.PerspectiveCamera(
        70,
        container.clientWidth / container.clientHeight,
        0.01,
        12
      );
      this.camera.position.z = 10;
      this.scene = new Three.Scene();
      this.addBackground(20, 20);
      this.playersGameObjects = [];
      this.healthBars = [];
      for (const username in this.positions) {
        const position = this.positions[username].position;
        this.addObject(position.x, position.y, username, this.positions[username].hitpoints);
      }
      this.renderer = new Three.WebGLRenderer({ antialias: true });
      this.renderer.setSize(container.clientWidth, container.clientHeight);
      container.appendChild(this.renderer.domElement);
    },
    animate(timeStamp) {
      if (this.start === undefined) {
        this.start = timeStamp;
      }
      if (!this.animating) {
        return;
      }
      // const elapsed = timeStamp - this.start;
      requestAnimationFrame(this.animate);
      for (const username in this.positions) {
        const position = this.positions[username].position;
        if (typeof this.playersGameObjects[username] === 'undefined') {
          this.addObject(position.x, position.y, username, this.positions[username].hitpoints);
        } else {
          this.playersGameObjects[username].parentObject.position.x = position.x;
          this.playersGameObjects[username].parentObject.position.y = position.y;
          this.playersGameObjects[username].meshChild.rotation.x += 0.01;
          this.playersGameObjects[username].meshChild.rotation.y += 0.02;
          this.playersGameObjects[username].healthBar.updateHealth(this.positions[username].hitpoints);
          this.updateHealthBar(username);
        }

        if (username === this.username) {
          if (
            this.vectorEnd.x !== position.x ||
            this.vectorEnd.y !== position.y
          ) {
            this.isLerping = false;
            this.$showLog && console.log('Stopped for positions changed');
          }
          if (
            !this.isLerping &&
            (this.camera.position.x !== position.x ||
              this.camera.position.y !== position.y)
          ) {
            this.cameraStart = timeStamp;
            this.vectorEnd.x = position.x;
            this.vectorEnd.y = position.y;
            this.vectorEnd.z = this.camera.position.z;
            this.isLerping = true;
          }
          if (this.isLerping) {
            this.camera.position.lerp(
              this.vectorEnd,
              (timeStamp - this.cameraStart) / this.lerpDuration
            );
            this.playersGameObjects[username].parentObject.position.x = this.camera.position.x;
            this.playersGameObjects[username].parentObject.position.y = this.camera.position.y;
            if (timeStamp > this.cameraStart + this.lerpDuration) {
              this.isLerping = false;
              this.$showLog && console.log('Stopped for end of lerp');
            }
          }
        }
      }
      for (const username in this.playersGameObjects) {
        if (typeof this.positions[username] === 'undefined') {
          this.$showLog && console.log('REMOVING ' + username);
          this.scene.remove(this.playersGameObjects[username].parentObject);
          this.playersGameObjects.splice(username, 1);
        }
      }
      this.renderer.render(this.scene, this.camera);
    },
    addObject(posX, posY, username, hitpoints) {
      this.playersGameObjects[username] = {};
      if (username === this.username) {
        this.playersGameObjects[username].meshChild = this.addBox(1, 1, 1, 0, 0, 0);
      } else {
        this.playersGameObjects[username].meshChild = this.addSphere(1, 0, 0, 0);
      }
      this.playersGameObjects[username].parentObject = new Three.Object3D();
      this.playersGameObjects[username].parentObject.position.x = posX;
      this.playersGameObjects[username].parentObject.position.y = posY;
      this.playersGameObjects[username].parentObject.position.z = 0.1;
      this.playersGameObjects[username].parentObject.add(this.playersGameObjects[username].meshChild);
      this.scene.add(this.playersGameObjects[username].parentObject);
      var canvas = document.createElement('canvas');
      canvas.width = 256;
      canvas.height = 256;
      var ctx = canvas.getContext('2d');
      ctx.font = '44pt Arial';
      ctx.fillStyle = 'white';
      if (this.positions[username].charactertype === 'fox') {
        ctx.fillStyle = 'red';
      }
      ctx.textAlign = 'center';
      ctx.fillText(username, 128, 44);
      var tex = new Three.Texture(canvas);
      tex.needsUpdate = true;
      var spriteMat = new Three.SpriteMaterial({ map: tex });
      this.playersGameObjects[username].textSprite = new Three.Sprite(spriteMat);
      this.playersGameObjects[username].textSprite.position.set(0, 1, 1);
      this.playersGameObjects[username].parentObject.add(this.playersGameObjects[username].textSprite);
      canvas = document.createElement('canvas');
      this.playersGameObjects[username].healthBarCanvas = canvas;
      canvas.width = 256;
      canvas.height = 256;
      ctx = canvas.getContext('2d');
      this.playersGameObjects[username].healthBar = this.buildHealthBar(hitpoints, ctx);
      this.createHealthBarSpriteAndDisplay(username);
    },
    buildHealthBar(health, ctx) {
      const hbwidth = 200;
      const hbHeight = 30;
      const x = 0;
      const y = -2;
      return new HealthBar(x, y, hbwidth, hbHeight, health, 'green', 'red', ctx);
    },
    updateHealthBar(username, health) {
      if (this.playersGameObjects[username].healthBar.getHealth() !== health) {
        this.playersGameObjects[username].parentObject.remove(this.playersGameObjects[username].healthBarSprite);
        this.createHealthBarSpriteAndDisplay(username);
      }
    },
    createHealthBarSpriteAndDisplay(username) {
      var ctx = this.playersGameObjects[username].healthBarCanvas.getContext('2d');
      ctx.clearRect(0, 0, 256, 256);
      this.playersGameObjects[username].healthBar.show();
      var tex = new Three.Texture(this.playersGameObjects[username].healthBarCanvas);
      tex.needsUpdate = true;
      var spriteMat = new Three.SpriteMaterial({ map: tex });
      this.playersGameObjects[username].healthBarSprite = new Three.Sprite(spriteMat);
      this.playersGameObjects[username].healthBarSprite.position.set(0.1, -1.5, 1);
      this.playersGameObjects[username].parentObject.add(this.playersGameObjects[username].healthBarSprite);
    },
    addBox(x, y, z, posX, posY, posZ) {
      this.$showLog && console.log('adding a box');
      const geometry = new Three.BoxGeometry(x, y, z);
      const material = new Three.MeshNormalMaterial();
      const mesh = new Three.Mesh(geometry, material);
      mesh.position.x = posX;
      mesh.position.y = posY;
      mesh.position.z = posZ;
      return mesh;
    },
    addSphere(radius, posX, posY, posZ) {
      this.$showLog && console.log('adding a sphere');
      const geometry = new Three.SphereGeometry(radius, 48);
      const material = new Three.MeshNormalMaterial();
      const mesh = new Three.Mesh(geometry, material);
      mesh.position.x = posX;
      mesh.position.y = posY;
      mesh.position.z = posZ;
      return mesh;
    },
    addBackground(sizeX, sizeY) {
      this.$showLog && console.log('adding background');
      const geometry = new Three.BoxGeometry(sizeX, sizeY, 0);
      const material = new Three.MeshBasicMaterial({
        color: 0x344522,
        wireframe: false
      });
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
