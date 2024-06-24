const canvas = document.getElementById("simulationCanvas");
const ctx = canvas.getContext("2d");

function drawBall(ball) {
  ctx.beginPath();
  ctx.arc(ball.x, ball.y, ball.r, 0, Math.PI * 2);
  ctx.fillStyle = ball.color;
  ctx.fill();
  ctx.closePath();
}

function drawBalls(balls) {
  ctx.clearRect(0, 0, canvas.width, canvas.height);
  balls.forEach((ball) => drawBall(ball));
}

async function fetchSimulationData() {
  const response = await fetch("/simulation");
  return response.json();
}

async function updateSimulation() {
  const simData = await fetchSimulationData();
  drawBalls(simData.balls);
  requestAnimationFrame(updateSimulation);
}

updateSimulation();
