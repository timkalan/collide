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
  let fps: number = 60;
  let gravity: number = 0.0;
  let size: number = 10;
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
      fps = data.fps;
      gravity = data.gravity;
      size = data.size;
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

  const changeGravity = async () => {
    try {
      const res = await fetch('http://localhost:8000/gravity', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json; charser=UTF-8'
        },
        body: JSON.stringify({
          gravity: (<HTMLInputElement>document.getElementById('gravity')).value
        })
      });
      if (res.ok) {
        console.log('Changed gravity');
        await fetchData();
      } else {
        console.error('Failed to change gravity', res.status);
      }
    } catch (err) {
      console.error('Failed to change gravity', err);
    }
  };

  const changeSize = async () => {
    try {
      const res = await fetch('http://localhost:8000/size', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json; charser=UTF-8'
        },
        body: JSON.stringify({
          size: (<HTMLInputElement>document.getElementById('size')).value
        })
      });
      if (res.ok) {
        console.log('Changed size');
        await fetchData();
      } else {
        console.error('Failed to change size', res.status);
      }
    } catch (err) {
      console.error('Failed to change size', err);
    }
  };

  const changeBallNumber = async () => {
    const quantity = parseInt((<HTMLInputElement>document.getElementById('quantity')).value);
    if (quantity < 1 || quantity > 500) {
      console.error('Invalid quantity:', quantity);
      return;
    }
    try {
      const res = await fetch('http://localhost:8000/changeballnumber', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json; charset=UTF-8'
        },
        body: JSON.stringify({ n: quantity })
      });
      if (res.ok) {
        console.log('Ball number changed');
        await fetchData();
      } else {
        console.error('Failed to change ball number:', res.status);
      }
    } catch (err) {
      console.error('Failed to change ball number:', err);
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
  <div class="buttons">
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
    <p>Gravity:</p>
    <input
      type="range"
      min="0.0"
      max="0.1"
      step="0.01"
      value={gravity}
      class="slider"
      id="gravity"
      on:change={changeGravity}
    />
    <p>Size:</p>
    <input
      type="range"
      min="1"
      max="20"
      value={size}
      class="slider"
      id="size"
      on:change={changeSize}
    />
  </div>
  <div class="canvas">
    <canvas id="simulationCanvas" {width} {height}></canvas>
  </div>
  <div class="balls">
    <label for="quantity">Quantity (between 1 and 500):</label>
    <input
      type="number"
      id="quantity"
      name="quantity"
      min="1"
      max="500"
      on:change={changeBallNumber}
    />
  </div>
  <div class="infoview">
    <p>FPS: {fps.toFixed(0)}</p>
  </div>
</main>
