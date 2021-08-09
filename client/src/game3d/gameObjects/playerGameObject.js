import HealthBar from './healthBar/healthBar';
import * as Three from 'three';

export default class PlayerGameObject {
  gameObject = null;
  mesh = null;
  textSprite = null;
  healthBarCanvas = null;
  healthBar = null;
  username = '';
  characterType = '';
  hitPoints = 0;
  constructor(username, posX, posY, characterType, hitPoints) {
    this.username = username;
    this.characterType = characterType;
    this.hitPoints = hitPoints;
    this.gameObject = new Three.Object3D();
    this.gameObject.position.x = posX;
    this.gameObject.position.y = posY;
    this.gameObject.position.z = 0.1;
    var canvas = document.createElement('canvas');
    canvas.width = 256;
    canvas.height = 256;
    var ctx = canvas.getContext('2d');
    ctx.font = '44pt Arial';
    ctx.fillStyle = this.getTextColorStyle();
    ctx.textAlign = 'center';
    ctx.fillText(username, 128, 44);
    var tex = new Three.Texture(canvas);
    tex.needsUpdate = true;
    this.textSprite = new Three.Sprite(new Three.SpriteMaterial({ map: tex }));
    this.textSprite.position.set(0, 1, 1);
    this.gameObject.add(this.textSprite);
    canvas = document.createElement('canvas');
    this.healthBarCanvas = canvas;
    canvas.width = 256;
    canvas.height = 256;
    ctx = canvas.getContext('2d');
    this.healthBar = this.buildHealthBar(hitPoints, ctx);
    this.createHealthBarSpriteAndDisplay();
  }

  update(position, timeStamp) {
    this.gameObject.position.x = position.x;
    this.gameObject.position.y = position.y;
    this.mesh.rotation.x += 0.01;
    this.mesh.rotation.y += 0.02;
    this.updateHealthBar(position.hitpoints);
  }

  getTextColorStyle() {
    return 'white';
  }

  createHealthBarSpriteAndDisplay() {
    var ctx = this.healthBarCanvas.getContext('2d');
    ctx.clearRect(0, 0, 256, 256);
    this.healthBar.show();
    var tex = new Three.Texture(this.healthBarCanvas);
    tex.needsUpdate = true;
    var spriteMat = new Three.SpriteMaterial({ map: tex });
    this.healthBarSprite = new Three.Sprite(spriteMat);
    this.healthBarSprite.position.set(0.1, -1.5, 1);
    this.gameObject.add(this.healthBarSprite);
  }

  updateHealthBar(health) {
    if (this.healthBar.getHealth() !== health) {
      this.gameObject.remove(this.healthBarSprite);
      this.healthBar.updateHealth(health);
      this.createHealthBarSpriteAndDisplay();
    }
  }

  setPosition(x, y) {
    this.gameObject.position.x = x;
    this.gameObject.position.y = y;
  }

  buildHealthBar(health, ctx) {
    const hbwidth = 200;
    const hbHeight = 30;
    const x = 0;
    const y = -2;
    return new HealthBar(x, y, hbwidth, hbHeight, health, 'green', 'red', ctx);
  }
}
