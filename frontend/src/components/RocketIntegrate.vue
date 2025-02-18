<template>
  <div>
    <h2>Rocket Integration</h2>
    <div v-if="rocketState">
      <p>Rocket X: {{ rocketState.x }}</p>
      <p>Rocket Y: {{ rocketState.y }}</p>
      <p>Velocity X: {{ rocketState.vx }}</p>
      <p>Velocity Y: {{ rocketState.vy }}</p>
      <p>Acceleration X: {{ rocketState.ax }}</p>
      <p>Acceleration Y: {{ rocketState.ay }}</p>
      <p>Fuel: {{ rocketState.fuel }}</p>
      <p>Mass: {{ rocketState.mass }}</p>
      <p>Thrust: {{ rocketState.thrust }}</p>
    </div>
    <div v-else>
      <p>Loading rocket state...</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      rocketState: null,  // Инициализируем rocketState как null
      socket: null,       // Инициализация WebSocket
    };
  },
  computed: {
    rocketData() {
      return this.rocketState ? this.rocketState : null;  // Проверяем, что rocketState не null
    }
  },
  mounted() {
    this.initializeWebSocket();  // Подключаемся к WebSocket при монтировании компонента
  },
  beforeUnmount() {
    if (this.socket) {
      this.socket.close();  // Закрываем WebSocket при уничтожении компонента
    }
  },
  methods: {
    initializeWebSocket() {
      // Открытие соединения с WebSocket
      this.socket = new WebSocket('ws://localhost:8086/rocket');

      this.socket.onopen = () => {
        console.log('WebSocket подключен');
        // Отправляем запрос на начальное состояние ракеты
        this.socket.send(JSON.stringify({ action: 'getState' }));
      };

      this.socket.onmessage = (event) => {
        const data = JSON.parse(event.data);
        // Обновляем состояние ракеты, полученное через WebSocket
        if (data && data.rocketState) {
          this.rocketState = data.rocketState;
        }
      };

      this.socket.onerror = (error) => {
        console.error('WebSocket ошибка:', error);
      };

      this.socket.onclose = () => {
        console.log('WebSocket закрыт');
      };
    },
  },
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
