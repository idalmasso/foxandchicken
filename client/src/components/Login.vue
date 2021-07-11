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
label {
  padding-right: 1rem;
  display: block;
}
button {
  margin-top: 1rem;
  width: 10%;
  border-top-right-radius: 8px;
  border-bottom-right-radius: 8px;
  background-color: darksalmon;
  padding: 8px;
}
input {
  box-sizing: border-box;
  border-top-left-radius: 8px;
  border-bottom-left-radius: 8px;
  width: 80%;
  padding: 8px 20px;
}
input:focus {
  outline: 0;
  background-color: wheat;
}
</style>
