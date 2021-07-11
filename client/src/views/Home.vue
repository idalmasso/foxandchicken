<template>
  <div class="home">
    <h1>Welcome!</h1>
    <Login />
  </div>
</template>

<script>
// @ is an alias to /src
import Login from '@/components/Login.vue';

export default {
  name: 'Home',
  components: {
    Login
  }
};
</script>
