<template>
  <div class="Login">
    <text-adder
      :showLabel="true"
      :textRequest="'Name'"
      :submitLabel="'Login'"
      @text-added="tryLogin"
    ></text-adder>
  </div>
</template>

<script>
import { mapActions } from 'vuex';
import TextAdder from './ui/TextAdder.vue';
export default {
  components: { TextAdder },
  name: 'Login',
  data() {
    return {
      name: ''
    };
  },
  methods: {
    ...mapActions(['login']),
    tryLogin(text) {
      console.log(text);
      this.login(text);
    }
  }
};
</script>

<style>

</style>
