export default class WebSocketService {
    constructor(url) {
      this.socket = new WebSocket(url);
      this.socket.onmessage = this.onMessage;
    }
  
    onMessage(event) {
      // Обработка полученных данных от сервера
      console.log("Data from server:", event.data);
      // Можете обновить состояние Vuex или компоненты здесь
    }
  
    sendMessage(message) {
      // Отправка данных на сервер
      this.socket.send(message);
    }
  
    close() {
      this.socket.close();
    }
  }
  