<template>
  <div>
    <h1>Посадка ракеты</h1>
    <div class="rocket-container">
      <div class="charts-left">
        <apexchart type="line" :options="altitudeChartOptions" :series="altitudeSeries" class="chart"></apexchart>
        <apexchart type="line" :options="velocityChartOptions" :series="velocitySeries" class="chart"></apexchart>
      </div>

      <img :src="rocketImageUrl" alt="Ракета" class="rocket-image" />

      <div class="charts-right">
        <apexchart type="line" :options="accelerationChartOptions" :series="accelerationSeries" class="chart"></apexchart>
        <apexchart type="line" :options="fuelChartOptions" :series="fuelSeries" class="chart"></apexchart>
      </div>
    </div>

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
import VueApexCharts from "vue3-apexcharts";

class RingBuffer {
  constructor(size) {
    this.size = size;
    this.buffer = new Array(size);
    this.index = 0;
    this.filled = false;
  }

  push(item) {
    this.buffer[this.index] = item;
    this.index = (this.index + 1) % this.size;
    if (this.index === 0) this.filled = true;
  }

  toArray() {
    return this.filled
      ? this.buffer.slice(this.index).concat(this.buffer.slice(0, this.index))
      : this.buffer.slice(0, this.index);
  }
}

export default {
  components: { apexchart: VueApexCharts },
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
      altitudeData: new RingBuffer(10),
      velocityData: new RingBuffer(100),
      accelerationData: new RingBuffer(100),
      fuelData: new RingBuffer(10),
      time: 0,
    };
  },
  computed: {
    formattedThrust() { return this.thrust.toFixed(2); },
    formattedAcceleration() { return this.acceleration.toFixed(2); },
    formattedY() { return this.y.toFixed(2); },
    formattedVelocity() { return this.velocity.toFixed(2); },
    formattedFuel() { return this.fuel.toFixed(2); },
    formattedEnergy() { return this.energy.toFixed(2); },
    formattedMass() { return this.mass.toFixed(2); },
    formattedDrag() { return this.drag.toFixed(2); },
    formattedLosses() { return this.losses.toFixed(2); },

    altitudeChartOptions() { return this.getChartOptions("Высота (м)"); },
    velocityChartOptions() { return this.getChartOptions("Скорость (м/с)"); },
    accelerationChartOptions() { return this.getChartOptions("Ускорение (м/с²)"); },
    fuelChartOptions() { return this.getChartOptions("Расход топлива (кг)"); },

    altitudeSeries() { return [{ name: "Высота", data: this.altitudeData.toArray() }]; },
    velocitySeries() { return [{ name: "Скорость", data: this.velocityData.toArray() }]; },
    accelerationSeries() { return [{ name: "Ускорение", data: this.accelerationData.toArray() }]; },
    fuelSeries() { return [{ name: "Расход топлива", data: this.fuelData.toArray() }]; }
  },
  mounted() {
    this.startRocketImageUpdates();
    this.startRocketDataUpdates();
    window.addEventListener('keydown', this.handleKeyPress);
  },
  beforeUnmount() {
    clearInterval(this.rocketImageIntervalId);
    clearInterval(this.rocketDataIntervalId);
    window.removeEventListener('keydown', this.handleKeyPress);
  },
  methods: {
    fetchRocketImage() {
      fetch("http://localhost:8086/physics/rocket-image")
      .then(response => response.blob())
      .then(blob => { this.rocketImageUrl = URL.createObjectURL(blob); })
      .catch(error => console.error("Error fetching rocket image:", error));
    },
    startRocketImageUpdates() {
      this.rocketImageIntervalId = setInterval(this.fetchRocketImage, 40);
    },
    startRocketDataUpdates() {
      this.rocketDataIntervalId = setInterval(this.updateRocketData, 40);
    },
    handleKeyPress(event) {
      if (event.key === "ArrowUp") this.increaseThrust();
      else if (event.key === "ArrowDown") this.decreaseThrust();
    },
    increaseThrust() { this.thrust = Math.round(this.thrust * 1.05); this.updateRocketThrust(); },
    decreaseThrust() { if (this.thrust > 0) this.thrust = Math.round(this.thrust * 0.95); this.updateRocketThrust(); },
    updateRocketThrust() {
      fetch("http://localhost:8086/physics/update-thrust", {
        method: 'POST', headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ thrust: this.thrust }),
      })
      .catch(error => console.error("Error updating thrust:", error));
    },
    updateRocketData() {
      fetch("http://localhost:8086/physics/update-data")
        .then(response => response.json())
        .then(data => {
          this.acceleration = data.acceleration;
          this.y = data.y;
          this.velocity = data.velocity_y;
          this.fuel = data.fuel_mass;
          this.mass = data.mass;
          this.drag = data.drag;
          this.energy = data.energy;
          this.losses = data.losses;

          this.altitudeData.push({ x: this.time, y: data.y });
          this.velocityData.push({ x: this.time, y: data.velocity_y });
          this.accelerationData.push({ x: this.time, y: data.acceleration });
          this.fuelData.push({ x: this.time, y: data.fuel_mass });

          this.time += 0.04;
        })
        .catch(error => console.error("Error updating rocket data:", error));
    },
    getChartOptions(title) {
      return {
        chart: { animations: { enabled: true, easing: "linear" }, background: "#1e1e1e" },
        title: { text: title, align: "center", style: { color: "#ffffff", fontSize: "16px" } },
        xaxis: { 
          title: { text: "Время (с)" },
          labels: { 
            style: { colors: "#ffffff" }, 
            formatter: (value) => Number(value).toFixed(2) 
          },
        },
        yaxis: { 
          labels: { 
            style: { colors: "#ffffff" },
            formatter: (value) => Number(value).toFixed(2) 
          } 
        },
        grid: { show: true, borderColor: "#444" },
        stroke: { curve: "smooth" },
        theme: { mode: "dark" },
        colors: ["#ff5733"]
      };
    }
  },
};
</script>


<style>
.rocket-container {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 40px;
}

.charts-left, .charts-right {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.chart {
  width: 800px;  
  height: 250px; 
}

.rocket-image {
  max-width: 1450px;
  max-height: 1420px;
  height: auto;
}

.data-row {
  display: flex;
  flex-wrap: wrap;
  gap: 40px;
  margin-top: 20px;
}
</style>