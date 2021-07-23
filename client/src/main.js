import { createApp } from 'vue';
import App from './App.vue';
import router from './router';
import store from './store';

const app = createApp(App);
app.config.globalProperties.$showLog = false;
app.config.globalProperties.$serverPort = 3000;
store.$serverPort = 3000;
app.use(store).use(router).mount('#app');
