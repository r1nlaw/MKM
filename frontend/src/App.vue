<template>
  <div id="app">
    <h1>Rocket Simulation</h1>
    <RocketCanvas />
    <RocketForce />
    <RocketTrajectory />
    <RocketIntegrate />
  </div>
</template>

<script>
import RocketForce from './components/RocketForce.vue';
import RocketTrajectory from './components/RocketTrajectory.vue';
import RocketIntegrate from './components/RocketIntegrate.vue';
import RocketCanvas from './components/RocketCanvas.vue';


export default {
  components: {
    RocketForce,
    RocketTrajectory,
    RocketIntegrate,
    RocketCanvas
  },
};
</script>

<style>
#app {
  font-family: Avenir, Helvetica, Arial, sans-serif;
  text-align: center;
  color: #2c3e50;
}
</style>
