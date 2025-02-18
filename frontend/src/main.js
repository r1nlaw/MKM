import { createApp } from 'vue';
import App from './App.vue';
import store from './store';  // Подключаем store

const app = createApp(App);
app.use(store);  // Подключаем store к приложению
app.mount('#app');
