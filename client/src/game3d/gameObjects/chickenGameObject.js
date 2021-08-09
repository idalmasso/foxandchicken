import PlayerGameObject from './playerGameObject';
import * as Three from 'three';

export default class ChickenGameObject extends PlayerGameObject {
  constructor(username, posX, posY, characterType, hitPoints) {
    super(username, posX, posY, characterType, hitPoints);
    this.mesh = this.addSphere(1, 0, 0, 0);
    this.gameObject.add(this.mesh);
  }

  addSphere(radius, posX, posY, posZ) {
    this.$showLog && console.log('adding a sphere');
    const geometry = new Three.SphereGeometry(radius, 48);
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
