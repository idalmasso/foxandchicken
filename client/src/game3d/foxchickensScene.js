import * as Three from 'three';
import FoxGameObject from './gameObjects/foxGameObject';
import ChickenGameObject from './gameObjects/chickenGameObject';

export default class FoxChickenScene {
  camera = null;
  username = '';
  scene = null;
  renderer = null;
  playersGameObjects = null;
  now = undefined;
  createdBox = false;
  cameraStart = 0;
  lerpDuration = 0;
  isLerping = false;
  vectorEnd = null;
  animating = false;
  container = null;
  constructor(container, positions, username) {
    this.init(container, positions, username);
  }

  update(timeStamp, positions) {
    if (this.start === undefined) {
      this.start = timeStamp;
    }
    if (!this.animating) {
      return;
    }
    for (const username in positions) {
      const position = positions[username].position;
      if (typeof this.playersGameObjects[username] === 'undefined') {
        this.addObject(
          position.x,
          position.y,
          username,
          positions[username].hitpoints,
          positions[username].charactertype
        );
      } else {
        this.playersGameObjects[username].update(position, timeStamp);
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
          this.playersGameObjects[username].setPosition(
            this.camera.position.x,
            this.camera.position.y
          );

          if (timeStamp > this.cameraStart + this.lerpDuration) {
            this.isLerping = false;
            this.$showLog && console.log('Stopped for end of lerp');
          }
        }
      }
    }
    for (const username in this.playersGameObjects) {
      if (typeof positions[username] === 'undefined') {
        this.$showLog && console.log('REMOVING ' + username);
        this.scene.remove(this.playersGameObjects[username].gameObject);
        this.playersGameObjects.splice(username, 1);
      }
    }
    this.renderer.render(this.scene, this.camera);
  }

  init(container, positions, username) {
    this.username = username;
    this.container = container;
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
    for (const username in positions) {
      const position = positions[username].position;
      this.addObject(
        position.x,
        position.y,
        username,
        positions[username].hitpoints,
        positions[username].charactertype
      );
    }
    this.renderer = new Three.WebGLRenderer({ antialias: true });
    this.renderer.setSize(container.clientWidth, container.clientHeight);
    container.appendChild(this.renderer.domElement);
  }

  addObject(posX, posY, username, hitPoints, characterType) {
    this.playersGameObjects[username] = {};
    switch (characterType) {
      case 'fox':
        this.playersGameObjects[username] = new FoxGameObject(
          username,
          posX,
          posY,
          characterType,
          hitPoints
        );
        break;
      case 'chicken':
        this.playersGameObjects[username] = new ChickenGameObject(
          username,
          posX,
          posY,
          characterType,
          hitPoints
        );
        break;
    }
    this.scene.add(this.playersGameObjects[username].gameObject);
  }

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
  }
}
