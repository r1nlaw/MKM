<template>
  <div id="app">
    <RocketCanvas />
  </div>
</template>

<script>
import RocketCanvas from './components/RocketCanvas.vue';

export default {
  name: 'App',
  components: {
    RocketCanvas,
  },
};
</script>

<style>
#app {
  text-align: center;
}
</style>
