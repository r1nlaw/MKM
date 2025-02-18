<template>
    <div class="vector-container">
      <h2>Vector Data</h2>
      <p><strong>X:</strong> {{ vector.x }}</p>
      <p><strong>Y:</strong> {{ vector.y }}</p>
      <button @click="fetchVectorData">Обновить вектор</button>
    </div>
  </template>
  
  <script>
  export default {
    computed: {
      vector() {
        return this.$store.state.vector; // Берем данные о векторе из Vuex
      }
    },
    mounted() {
      this.initializeWebSocket(); // Подключаемся к WebSocket при монтировании компонента
    },
    beforeUnmount() {
      if (this.socket) {
        this.socket.close(); // Закрываем WebSocket при уничтожении компонента
      }
    },
    methods: {
      initializeWebSocket() {
        // Открытие соединения с WebSocket
        this.socket = new WebSocket('ws://localhost:8086/vector');
  
        this.socket.onopen = () => {
          console.log('WebSocket подключен');
          // Отправляем запрос на начальные данные о векторе
          this.socket.send(JSON.stringify({ action: 'getVector' }));
        };
  
        this.socket.onmessage = (event) => {
          const data = JSON.parse(event.data);
          // Обновляем данные о векторе в Vuex
          if (data && data.vector) {
            this.$store.commit('updateVector', data.vector);
          }
        };
  
        this.socket.onerror = (error) => {
          console.error('WebSocket ошибка:', error);
        };
  
        this.socket.onclose = () => {
          console.log('WebSocket закрыт');
        };
      },
      fetchVectorData() {
        // Отправляем запрос на обновление данных о векторе
        if (this.socket && this.socket.readyState === WebSocket.OPEN) {
          this.socket.send(JSON.stringify({ action: 'getVector' }));
        }
      }
    }
  };
  </script>
  
  <style scoped>
  .vector-container {
    padding: 10px;
    border: 1px solid #ccc;
    border-radius: 5px;
    width: 200px;
    margin: 10px;
  }
  button {
    margin-top: 10px;
    padding: 5px 10px;
    background-color: #007bff;
    color: white;
    border: none;
    cursor: pointer;
  }
  button:hover {
    background-color: #0056b3;
  }
  </style>
  