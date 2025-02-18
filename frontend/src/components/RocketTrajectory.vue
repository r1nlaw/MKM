<template>
  <div>
    <h2>Rocket Trajectory</h2>
    <div v-if="trajectory">
      <p>Trajectory X: {{ trajectory.x }}</p>
      <p>Trajectory Y: {{ trajectory.y }}</p>
      <p>Velocity X: {{ trajectory.vx }}</p>
      <p>Velocity Y: {{ trajectory.vy }}</p>
    </div>
    <div v-else>
      <p>Loading trajectory data...</p>
    </div>
  </div>
</template>

<script>
export default {
  computed: {
    trajectory() {
      return this.$store.state.trajectory;  // Берем trajectory из Vuex
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
      this.socket = new WebSocket('ws://localhost:8086/trajectory');

      this.socket.onopen = () => {
        console.log('WebSocket подключен');
        // Отправляем запрос на начальные данные о траектории
        this.socket.send(JSON.stringify({ action: 'getTrajectory' }));
      };

      this.socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        // Обновляем данные о траектории в Vuex
        if (data && data.trajectory) {
          this.$store.commit('updateTrajectory', data.trajectory);
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
