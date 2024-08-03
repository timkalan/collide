<script lang="ts">
  import { onMount, onDestroy } from 'svelte';

  interface Point {
    x: number;
    y: number;
  }

  interface Square {
    topLeft: Point;
    bottomRight: Point;
    velocity: number;
    weight: string;
  }

  let smallSquare: Square;
  let bigSquare: Square;
  let width = 800;
  let height = 600;
  let numCollisions = 0;
  let multiplier = 1;
  let paused: boolean;
  let interval: number;

  const fetchData = async () => {
    const res = await fetch('http://localhost:8000/pillision');
    if (res.ok) {
      const data = await res.json();
      smallSquare = data.smallSquare;
      bigSquare = data.bigSquare;
      width = data.width;
      height = data.height;
      numCollisions = data.numCollisions;
      multiplier = data.multiplier;
      paused = data.paused;
      draw();
    } else {
      console.error('Failed to fetch data:', res.status);
    }
  };

  const draw = () => {
    const canvas = <HTMLCanvasElement>document.getElementById('pillisionCanvas');
    const ctx = canvas.getContext('2d')!;

    // Clear the canvas
    ctx.clearRect(0, 0, width, height);

    ctx.beginPath();
    ctx.rect(
      bigSquare.topLeft.x,
      bigSquare.topLeft.y,
      bigSquare.bottomRight.x - bigSquare.topLeft.x,
      bigSquare.bottomRight.y - bigSquare.topLeft.y
    );
    ctx.fillStyle = '#b16286';
    ctx.fill();
    ctx.closePath();

    ctx.beginPath();
    ctx.rect(
      smallSquare.topLeft.x,
      smallSquare.topLeft.y,
      smallSquare.bottomRight.x - smallSquare.topLeft.x,
      smallSquare.bottomRight.y - smallSquare.topLeft.y
    );
    ctx.fillStyle = '#fabd2f';
    ctx.fill();
    ctx.closePath();
  };

  const startSimulation = async () => {
    await fetchData();
    interval = setInterval(fetchData, 16);
    try {
      const res = await fetch('http://localhost:8000/resumepi', {
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
      const res = await fetch('http://localhost:8000/pausepi', {
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

  const changeWeight = async () => {
    try {
      const res = await fetch('http://localhost:8000/weight', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json; charser=UTF-8'
        },
        body: JSON.stringify({
          weight: (<HTMLInputElement>document.getElementById('weight')).value
        })
      });
      if (res.ok) {
        console.log('Changed weight');
        await fetchData();
      } else {
        console.error('Failed to change weight', res.status);
      }
    } catch (err) {
      console.error('Failed to change weight', err);
    }
  };

  onMount(async () => {
    await fetchData();
  });

  onDestroy(async () => {
    clearInterval(interval);
    await stopSimulation();
  });
</script>

<main>
  <!-- <div class="title"> -->
  <!--   <h1>Multi-body collision</h1> -->
  <!-- </div> -->
  <div class="controls">
    {#if !paused}
      <button on:click={stopSimulation}>Reset</button>
    {:else}
      <button on:click={startSimulation}>Start</button>
    {/if}
    <ul>
      <li>
        <p>Weight:</p>
      </li>
      <li>
        <input
          type="range"
          min="1"
          max="5"
          step="1"
          value={multiplier}
          class="slider"
          id="weight"
          on:change={changeWeight}
        />
      </li>
    </ul>
    <p>Collisions: {numCollisions}</p>
  </div>
  <div class="canvas">
    <canvas id="pillisionCanvas" {width} {height}></canvas>
  </div>
</main>
