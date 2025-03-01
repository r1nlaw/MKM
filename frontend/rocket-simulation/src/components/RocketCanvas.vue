<template>
  <div>
    <h1>Посадка ракеты</h1>
    <img :src="rocketImageUrl" alt="" />
    <div class="data-row">
      <p>Тяга: {{ formattedThrust }}</p>
      <p>Высота: {{ formattedY }}</p>
      <p>Скорость: {{ formattedVelocity }}</p>
      <p>Ускорение: {{ formattedAcceleration }}</p>
      <p>Кол-во топлива: {{ formattedFuel }}</p>
      <p>Сохранение энергии: {{ formattedEnergy }}</p>
      <p>Масса: {{ formattedMass }}</p>
      <p>Сопротивление: {{ formattedDrag }}</p>
      <p>Потеря из-за внешних сил от Циолковского: {{ formattedLosses }}</p>
    </div>
  </div>
</template>

<script>
export default {
  data() {
    return {
      rocketImageUrl: "", 
      thrust: 75000, 
      acceleration: 0,
      y: 0,
      velocity: 0,
      fuel: 0,
      energy: 0,
      mass: 0,
      drag: 0,
      losses: 0,
      intervalId: null, 
    };
  },
  computed: {
      formattedThrust() {
        return this.thrust.toFixed(2);
      },
      formattedAcceleration() {
        return this.acceleration.toFixed(2);
      },
      formattedY() {
        return this.y.toFixed(2);
      },
      formattedVelocity() {
        return this.velocity.toFixed(2);
      },
      formattedFuel() {
        return this.fuel.toFixed(2);
      },
      formattedEnergy() {
        return this.energy.toFixed(2);
      },
      formattedMass() {
        return this.mass.toFixed(2);
      },
      formattedDrag() {
        return this.drag.toFixed(2);
      },
      formattedLosses() {
        return this.losses.toFixed(2);
      },
    },
  mounted() {
    this.fetchRocketImage(); 
    this.startRocketUpdates(); 
    this.startRocketDataUpdates();
    window.addEventListener('keydown', this.handleKeyPress); // Обрабатываем нажатия клавиш
  },
  beforeUnmount() {
    clearInterval(this.intervalId); // Останавливаем интервал перед размонтированием компонента
    clearInterval(this.rocketDataIntervalId);
    window.removeEventListener('keydown', this.handleKeyPress); // Убираем обработчик событий
  },
  methods: {
    // Запрос изображения ракеты с сервера
    fetchRocketImage() {
      fetch("http://localhost:8086/physics/rocket-image")
      .then(response => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.blob();
      })
      .then(blob => {
        this.rocketImageUrl = URL.createObjectURL(blob); // Преобразуем ответ в изображение
      })
      .catch(error => {
        console.error("Error fetching rocket image:", error);
      });
    },
    
    // Начало обновления ракеты 
    startRocketUpdates() {
      this.intervalId = setInterval(() => {
        this.fetchRocketImage();
      }, 40); 
    },
    startRocketDataUpdates() {
      this.rocketDataIntervalId = setInterval(() => {
        this.updateRocketdataThrust(); 
      }, 40); 
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

    // Уменьшение тяги на 5% с округлением до целого, но не меньше 1
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
        body: JSON.stringify({ 
          thrust: this.thrust,
        }),
      })
      .then(response => response.json())
      .then(data => console.log("Thrust updated:", data))
      .catch(error => console.error("Error updating thrust:", error));
    },
    // получаем данные от сервера
    updateRocketdataThrust() {
      fetch("http://localhost:8086/physics/update-data", {
        method: 'GET', 
      })
      .then(response => {
        if (!response.ok) {
          throw new Error("Network response was not ok");
        }
        return response.json();
      })
      .then(data => {
        this.acceleration = data.acceleration;
        this.y = data.y;
        this.velocity = data.velocity_y;
        this.fuel = data.fuel_mass;
        this.mass = data.mass;
        this.drag = data.drag;
        this.energy = data.energy;
        this.losses = data.losses;
      })
      .catch(error => {
        console.error("Error updating rocket data:", error);
      });
    },
  },
}

</script>

<style>
.data-row {
  display: flex;
  flex-wrap: wrap;
  gap: 40px;
  margin-left: 0px;
  margin-top: 50px;
}

.data-row p {
  margin: 0; 
}
</style>
