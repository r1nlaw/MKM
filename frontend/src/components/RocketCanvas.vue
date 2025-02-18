<template>
    <div>
      <canvas ref="canvas" width="800" height="600"></canvas>
    </div>
  </template>
  
  <script>
  export default {
    data() {
      return {
        rocket: {
          x: 400,  // Начальная позиция по оси X
          y: 550,  // Начальная позиция по оси Y
          width: 20, // Ширина ракеты
          height: 60, // Высота ракеты
        },
        surfaceY: 580, // Высота поверхности
      };
    },
    mounted() {
      this.drawRocketAndSurface();
    },
    methods: {
        drawRocketAndSurface() {
            const canvas = this.$refs.canvas;
            const ctx = canvas.getContext('2d');

            // Рисуем поверхность
            ctx.fillStyle = '#006400';
            ctx.fillRect(0, this.surfaceY, canvas.width, canvas.height - this.surfaceY);

            // Рисуем ракету
            ctx.fillStyle = '#FF0000';
            ctx.fillRect(this.rocket.x - this.rocket.width / 2, this.rocket.y - this.rocket.height, this.rocket.width, this.rocket.height);

            // Анимация ракеты
            if (this.rocket.y > this.surfaceY - this.rocket.height) {
            this.rocket.y -= 1; // Сдвигаем ракету вверх
            }

            // Перерисовываем каждый кадр
            requestAnimationFrame(this.drawRocketAndSurface);
        }
    }

  };
  </script>
  
  <style scoped>
  canvas {
    border: 1px solid #000;
  }
  </style>
  