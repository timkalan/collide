<script lang="ts">
  import { onMount } from 'svelte';
  // import "./global.css";

  interface Ball {
    r: number;
    x: number;
    y: number;
    vx: number;
    vy: number;
    color: string;
  }

  let balls: Ball[] = [];
  let width = 800;
  let height = 600;
  let paused: boolean;
  let war: boolean;
  let interval: number;

  const fetchData = async () => {
    const res = await fetch('http://localhost:8000/simulation');
    if (res.ok) {
      const data = await res.json();
      balls = data.balls;
      width = data.width;
      height = data.height;
      paused = data.paused;
      war = data.war;
      draw();
    } else {
      console.error('Failed to fetch data:', res.status);
    }
  };

  const draw = () => {
    const canvas = <HTMLCanvasElement>document.getElementById('simulationCanvas');
    const ctx = canvas.getContext('2d')!;

    // Clear the canvas
    ctx.clearRect(0, 0, width, height);

    // Draw each ball
    balls.forEach((ball) => {
      ctx.beginPath();
      ctx.arc(ball.x, ball.y, ball.r, 0, Math.PI * 2);
      ctx.fillStyle = ball.color;
      ctx.fill();
      ctx.closePath();
    });
  };

  const startSimulation = async () => {
    await fetchData();
    interval = setInterval(fetchData, 16);
    try {
      const res = await fetch('http://localhost:8000/resume', {
        method: 'POST'
      });
      if (res.ok) {
        console.log('Server resumed');
        await fetchData();
      } else {
        console.error('Failed to resume server:', res.status);
      }
    } catch (err) {
      console.error('Failed to resume server:', err);
    }
  };

  const stopSimulation = async () => {
    clearInterval(interval);
    try {
      const res = await fetch('http://localhost:8000/pause', {
        method: 'POST'
      });
      if (res.ok) {
        console.log('Server paused');
        await fetchData();
      } else {
        console.error('Failed to pause server:', res.status);
      }
    } catch (err) {
      console.error('Failed to pause server:', err);
    }
  };

  const startWar = async () => {
    try {
      const res = await fetch('http://localhost:8000/startwar', {
        method: 'POST'
      });
      if (res.ok) {
        console.log('War started');
        await fetchData();
      } else {
        console.error('Failed to start war', res.status);
      }
    } catch (err) {
      console.error('Failed to start war', err);
    }
  };

  const stopWar = async () => {
    try {
      const res = await fetch('http://localhost:8000/stopwar', {
        method: 'POST'
      });
      if (res.ok) {
        console.log('War ended');
        await fetchData();
      } else {
        console.error('Failed to end war', res.status);
      }
    } catch (err) {
      console.error('Failed to end war', err);
    }
  };
  onMount(async () => {
    await fetchData();
    interval = setInterval(fetchData, 16);
  });
</script>

<main>
  <div class="title">
    <h1>Multi-body collision</h1>
  </div>
  <div class="description">
    <p>
      <!-- This simulation demonstrates how multiple balls collide within a confined space. Click "Start" -->
      <!-- to begin the simulation and "Stop" to pause it. -->
    </p>
  </div>

  <div class="controls">
    {#if !paused}
      <button on:click={stopSimulation}>Stop</button>
    {:else}
      <button on:click={startSimulation}>Start</button>
    {/if}
    {#if !war}
      <button on:click={startWar}>Fight!</button>
    {:else}
      <button on:click={stopWar}>Reset</button>
    {/if}
  </div>
  <canvas id="simulationCanvas" {width} {height}></canvas>
</main>
