<template>
  <div>
    <h1>Посадка ракеты</h1>
    <img :src="rocketImageUrl" alt="Rocket Image" />
    <p>Тяга: {{ thrust }} N</p>
    <p>Нажмите стрелки вверх/вниз для изменения тяги.</p>
  </div>
</template>

<script>
export default {
data() {
  return {
    rocketImageUrl: "", // Ссылка на изображение ракеты
    thrust: 1000,  // Начальная тяга 
    intervalId: null, // Идентификатор интервала
  };
},
mounted() {
  this.fetchRocketImage(); // При монтировании компонента запрашиваем изображение
  this.startRocketUpdates(); // Начинаем обновлять ракету каждую миллисекунду
  window.addEventListener('keydown', this.handleKeyPress); // Обрабатываем нажатия клавиш
},
beforeUnmount() {
  clearInterval(this.intervalId); // Останавливаем интервал перед размонтированием компонента
  window.removeEventListener('keydown', this.handleKeyPress); // Убираем обработчик событий
},
methods: {
  // Запрос изображения ракеты с сервера
  fetchRocketImage() {
    fetch("http://localhost:8086/physics/rocket-image")
      .then((response) => response.blob()) 
      .then((blob) => {
        this.rocketImageUrl = URL.createObjectURL(blob); // Преобразуем ответ в изображение
      })
      .catch((error) => console.error("Error fetching rocket image:", error));
  },
  
  // Начало обновления ракеты 
  startRocketUpdates() {
    this.intervalId = setInterval(() => {
      this.fetchRocketImage();
    }, 16); // Каждые 16 мс 
  },

  // Обработка нажатия клавиш
  handleKeyPress(event) {
    if (event.key === "ArrowUp") {
      this.increaseThrust(); // Увеличить тягу
    } else if (event.key === "ArrowDown") {
      this.decreaseThrust(); // Уменьшить тягу
    }
  },
  
  // Увеличение тяги на 5% с округлением до целого
  increaseThrust() {
    this.thrust = Math.round(this.thrust * 1.05);
    this.updateRocketThrust();
  },

  // Уменьшение тяги на 5% с округлением до целого
  decreaseThrust() {
    if (this.thrust > 0) {
      this.thrust = Math.round(this.thrust * 0.95);
    }
    this.updateRocketThrust();
  },

  // Обновление тяги ракеты на сервере
  updateRocketThrust() {
    fetch("http://localhost:8086/physics/update-thrust", {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({ thrust: this.thrust }),
    })
    .then(response => response.json())
    .then(data => console.log("Thrust updated:", data))
    .catch(error => console.error("Error updating thrust:", error));
  },
},
};
</script>
