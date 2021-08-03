class HealthBar {
  constructor(x, y, w, h, maxHealth, color, emptyColor, context) {
    this.x = x;
    this.y = y;
    this.w = w;
    this.h = h;
    this.maxHealth = maxHealth;
    this.maxWidth = w;
    this.health = maxHealth;
    this.color = color;
    this.emptyColor = emptyColor;
    this.context = context;
  }

  show() {
    this.context.lineWidth = 4;
    this.context.strokeStyle = '#333';
    this.context.fillStyle = this.emptyColor;
    this.context.fillRect(this.x, this.y, this.maxWidth, this.h);
    this.context.fillStyle = this.color;
    this.context.fillRect(this.x, this.y, this.w, this.h);
    this.context.strokeRect(this.x, this.y, this.maxWidth, this.h);
  }

  updateHealth(val) {
    this.health = val;
    this.w = (this.health / this.maxHealth) * this.maxWidth;
  }

  getHealth() {
    return this.health;
  }
}

export default HealthBar;
