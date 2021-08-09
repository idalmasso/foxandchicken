import PlayerGameObject from './playerGameObject';
import * as Three from 'three';

export default class FoxGameObject extends PlayerGameObject {
  constructor(username, posX, posY, characterType, hitPoints) {
    super(username, posX, posY, characterType, hitPoints);
    this.mesh = this.addBox(1, 1, 1, 0, 0, 0);
    this.gameObject.add(this.mesh);
  }

  addBox(x, y, z, posX, posY, posZ) {
    this.$showLog && console.log('adding a box');
    const geometry = new Three.BoxGeometry(x, y, z);
    const material = new Three.MeshNormalMaterial();
    const mesh = new Three.Mesh(geometry, material);
    mesh.position.x = posX;
    mesh.position.y = posY;
    mesh.position.z = posZ;
    return mesh;
  }

  getTextColorStyle() {
    return 'red';
  }
}
