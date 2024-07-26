<script lang="ts">
	import { onMount } from 'svelte';
	// import "./global.css";

	let balls = [];
	let width: number;
	let height: number;
	let simulationRunning = false;
	let interval: number;

	const fetchData = async () => {
		const res = await fetch('http://localhost:8000/simulation');
		if (res.ok) {
			const data = await res.json();
			balls = data.balls;
			width = data.width;
			height = data.height;
			draw();
		} else {
			console.error('Failed to fetch data:', res.status);
		}
	};

	const draw = () => {
		const canvas = <HTMLCanvasElement> document.getElementById('simulationCanvas');
		const ctx = canvas.getContext('2d');

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

	const startSimulation = () => {
		simulationRunning = true;
		fetchData();
		interval = setInterval(fetchData, 16);
	};

	const stopSimulation = () => {
		simulationRunning = false;
		clearInterval(interval);
	};

	onMount(() => {
		// Initialize with simulation started
    startSimulation();
	});
</script>

<main>
	<div class="description">
		<h1>Multi-body collision</h1>
		<p>
			<!-- This simulation demonstrates how multiple balls collide within a confined space. Click "Start" -->
			<!-- to begin the simulation and "Stop" to pause it. -->
		</p>
	</div>

	<div class="controls">
		<!-- <button on:click={startSimulation} disabled={simulationRunning}>Start</button> -->
		<!-- <button on:click={stopSimulation} disabled={!simulationRunning}>Stop</button> -->
	</div>
	<canvas id="simulationCanvas" {width} {height}></canvas>
</main>
