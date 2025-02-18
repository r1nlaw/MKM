import { createStore } from 'vuex';  // Используем createStore для Vue 3

export default createStore({
  state: {
    rocketState: {
      x: 0,
      y: 0,
      vx: 0,
      vy: 0,
      ax: 0,
      ay: 0,
      fuel: 100,
      mass: 500,
      thrust: 100,
    },
    force: {
      ThrustY: 0,
      GravityY: 0,
      ResultFY: 0,
    },
    trajectory: [],
    vector: { x: 0, y: 0 }, // Добавили vector в state
  },
  mutations: {
    updateRocketState(state, newState) {
      state.rocketState = newState;
    },
    updateForce(state, newForce) {
      state.force = newForce;
    },
    updateTrajectory(state, newTrajectory) {
      state.trajectory = newTrajectory;
    },
    updateVector(state, newVector) {
      state.vector = newVector;
    },
  },
  actions: {
    async getForce({ commit }, rocketState) {
      const response = await fetch('http://localhost:8086/physics/force', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rocketState),
      });
      const data = await response.json();
      commit('updateForce', data);
    },
    async getTrajectory({ commit }, rocketState) {
      const response = await fetch('http://localhost:8086/physics/trajectory', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rocketState),
      });
      const data = await response.json();
      commit('updateTrajectory', data);
    },
    async integrate({ commit }, rocketState) {
      const response = await fetch('http://localhost:8086/physics/integrate', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rocketState),
      });
      const data = await response.json();
      commit('updateRocketState', data);
    },
    async getVector({ commit }, rocketState) {
      const response = await fetch('http://localhost:8086/physics/vector', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rocketState),
      });

      if (!response.ok) {
        throw new Error(`Ошибка HTTP: ${response.status}`);
      }

      const data = await response.json();
      commit('updateVector', data);
    },
  },
});
