<template>
  <div class="particles">
    <div
      class="particle"
      v-for="(p, i) in particles"
      :key="i"
      :style="{
        left: p.x + '%',
        animationDelay: p.delay + 's',
        animationDuration: p.duration + 's',
        width: p.size + 'px',
        height: p.size + 'px',
        opacity: p.opacity
      }"
    ></div>
  </div>
</template>

<script>
export default {
  name: 'Particles',
  data() {
    return {
      particles: Array.from({ length: 50 }, () => ({
        x: Math.random() * 100,
        delay: Math.random() * 1,
        duration: 10 + Math.random() * 20,
        size: 1 + Math.random() * 3,
        opacity: Math.random() * 0.5 + 0.3
      }))
    }
  }
}
</script>

<style scoped>
.particles {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  pointer-events: none;
  overflow: hidden;
  z-index: 0;
}

.particle {
  position: absolute;
  bottom: -10px;
  border-radius: 50%;
  background: var(--neon-cyan);
  box-shadow: 0 0 10px var(--neon-cyan);
  animation: float-up linear infinite;
}

@keyframes float-up {
  0% { transform: translateY(0) translateX(0); opacity: 1; }
  100% { transform: translateY(-100vh) translateX(50px); opacity: 0; }
}
</style>
