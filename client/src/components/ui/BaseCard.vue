<template>
  <div class="box-div">
    <header v-if="$slots.header">
      <slot name="header"> </slot>
    </header>
    <main>
      <slot></slot>
    </main>
    <footer>
      <slot name="footer"></slot>
    </footer>
  </div>
</template>

<script>
export default {};
</script>

<style scoped>
header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
div .box-div {
  margin: 2rem auto;
  max-width: 50rem;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.26);
  padding: 1rem;
  background-color: rgb(198, 224, 224);
}
</style>
