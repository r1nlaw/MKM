<template>
  <div>
    <h2>Rocket Force</h2>
    <div v-if="force">
      <p>Thrust Y: {{ force.thrustY }}</p>
      <p>Gravity Y: {{ force.gravityY }}</p>
      <p>Result Force Y: {{ force.resultFY }}</p>
    </div>
    <div v-else>
      <p>Waiting for data...</p>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    force() {
      return this.$store.state.force; // Берем force из Vuex
    }
  },
  mounted() {
    this.initializeWebSocket(); // Подключаемся к WebSocket при монтировании компонента
  },
  beforeUnmount() {
    if (this.socket) {
      this.socket.close();  // Закрываем WebSocket при уничтожении компонента
    }
  },
  methods: {
    initializeWebSocket() {
      // Открытие соединения с WebSocket
      this.socket = new WebSocket('ws://localhost:8086/force');

      this.socket.onopen = () => {
        console.log('WebSocket подключен');
        // Отправляем запрос на начальные данные о силе
        this.socket.send(JSON.stringify({ action: 'getForce' }));
      };

      this.socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        // Обновляем данные о силе в Vuex
        if (data && data.force) {
          this.$store.commit('updateForce', data.force);
        }
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket ошибка:', error);
      };

      this.socket.onclose = () => {
        console.log('WebSocket закрыт');
      };
    }
  }
};
</script>

<style scoped>
h2 {
  color: #4CAF50;
}
p {
  font-size: 1.2em;
}
</style>
