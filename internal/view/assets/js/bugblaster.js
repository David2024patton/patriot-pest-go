/* Bug Blaster — the classic Patriot minigame.
 * Bugs crawl in from the edges. Click/tap anywhere to drop a pest bomb;
 * it detonates and every bug inside the blast radius dies and fades away.
 * Vanilla canvas, no dependencies. Only runs after the visitor presses start,
 * and pauses when the tab is hidden or the game scrolls off screen. */
(function () {
  'use strict';
  var canvas = document.getElementById('bugblaster');
  if (!canvas) return;
  var overlay = document.getElementById('bugblaster-start');
  var hud = document.getElementById('bugblaster-hud');
  var ctx = canvas.getContext('2d');

  var W = 0, H = 0, DPR = 1;
  function resize() {
    DPR = Math.min(window.devicePixelRatio || 1, 2);
    var r = canvas.getBoundingClientRect();
    W = Math.max(280, Math.floor(r.width));
    H = Math.max(300, Math.floor(r.height));
    canvas.width = W * DPR;
    canvas.height = H * DPR;
    ctx.setTransform(DPR, 0, 0, DPR, 0, 0);
  }

  var bugs = [], bombs = [], blasts = [], parts = [];
  var running = false, visible = true, last = 0, spawnT = 0, kills = 0, spawnMs = 950;

  function rand(a, b) { return a + Math.random() * (b - a); }

  function spawnBug() {
    if (bugs.length >= 36) return;
    var edge = Math.floor(Math.random() * 4), x, y;
    if (edge === 0) { x = rand(0, W); y = -20; }
    else if (edge === 1) { x = W + 20; y = rand(0, H); }
    else if (edge === 2) { x = rand(0, W); y = H + 20; }
    else { x = -20; y = rand(0, H); }
    // crawl toward a jittered point near the middle
    var tx = W / 2 + rand(-W / 3, W / 3), ty = H / 2 + rand(-H / 3, H / 3);
    var dx = tx - x, dy = ty - y, d = Math.hypot(dx, dy) || 1;
    var sp = rand(28, 62) * (1 + kills / 220); // slow ramp with score
    bugs.push({
      x: x, y: y, vx: dx / d * sp, vy: dy / d * sp,
      size: rand(11, 19), wob: rand(0, 6.28),
      dying: false, deathT: 0, alpha: 1
    });
  }

  function dropBomb(x, y) {
    bombs.push({ x: x, y: y, fuse: 0.38, maxFuse: 0.38 });
  }

  function detonate(b) {
    var R = 95;
    blasts.push({ x: b.x, y: b.y, r: 8, maxR: R, life: 0.45, maxLife: 0.45 });
    for (var i = 0; i < 14; i++) {
      var a = rand(0, 6.28), sp = rand(60, 240);
      parts.push({ x: b.x, y: b.y, vx: Math.cos(a) * sp, vy: Math.sin(a) * sp, life: rand(0.3, 0.7), maxLife: 0.7 });
    }
    for (var j = 0; j < bugs.length; j++) {
      var g = bugs[j];
      if (!g.dying && Math.hypot(g.x - b.x, g.y - b.y) <= R) {
        g.dying = true; g.deathT = 0; kills++;
      }
    }
    if (hud) hud.textContent = 'BUGS ELIMINATED: ' + kills;
    if (kills > 0 && kills % 15 === 0) spawnMs = Math.max(380, spawnMs - 90); // waves get hotter
  }

  function drawBug(g, t) {
    ctx.save();
    ctx.globalAlpha = g.alpha;
    ctx.translate(g.x, g.y);
    ctx.rotate(Math.atan2(g.vy, g.vx) + (g.dying ? Math.PI : 0));
    var s = g.size, legSwing = Math.sin(t * 14 + g.wob) * s * 0.35;
    // legs
    ctx.strokeStyle = '#8a7a5a'; ctx.lineWidth = Math.max(1.2, s * 0.09);
    for (var l = -1; l <= 1; l++) {
      var lx = l * s * 0.45;
      ctx.beginPath(); ctx.moveTo(lx, -s * 0.28);
      ctx.lineTo(lx + legSwing * (l === 0 ? 1 : -1), -s * 0.85); ctx.stroke();
      ctx.beginPath(); ctx.moveTo(lx, s * 0.28);
      ctx.lineTo(lx - legSwing * (l === 0 ? 1 : -1), s * 0.85); ctx.stroke();
    }
    // body
    ctx.fillStyle = g.dying ? '#4a3b28' : '#241f14';
    ctx.beginPath(); ctx.ellipse(0, 0, s, s * 0.62, 0, 0, 6.29); ctx.fill();
    ctx.strokeStyle = g.dying ? '#6b5a3e' : '#f4772e'; ctx.lineWidth = 1.5;
    ctx.beginPath(); ctx.ellipse(0, 0, s, s * 0.62, 0, 0, 6.29); ctx.stroke();
    // head
    ctx.fillStyle = '#f4772e';
    ctx.beginPath(); ctx.arc(s * 0.95, 0, s * 0.32, 0, 6.29); ctx.fill();
    // antennae
    ctx.strokeStyle = '#8a7a5a'; ctx.lineWidth = 1.2;
    ctx.beginPath(); ctx.moveTo(s * 1.1, -s * 0.12); ctx.lineTo(s * 1.55, -s * 0.5); ctx.stroke();
    ctx.beginPath(); ctx.moveTo(s * 1.1, s * 0.12); ctx.lineTo(s * 1.55, s * 0.5); ctx.stroke();
    if (g.dying) {
      // X eyes on the dearly departed
      ctx.strokeStyle = '#f4772e'; ctx.lineWidth = 1.6;
      var ex = s * 0.95, er = s * 0.14;
      ctx.beginPath();
      ctx.moveTo(ex - er, -er); ctx.lineTo(ex + er, er);
      ctx.moveTo(ex + er, -er); ctx.lineTo(ex - er, er);
      ctx.stroke();
    }
    ctx.restore();
  }

  function step(dt, t) {
    spawnT += dt * 1000;
    if (spawnT >= spawnMs) { spawnT = 0; spawnBug(); }

    var i, g;
    for (i = bugs.length - 1; i >= 0; i--) {
      g = bugs[i];
      if (g.dying) {
        g.deathT += dt;
        g.alpha = Math.max(0, 1 - g.deathT / 0.7);
        if (g.alpha <= 0) bugs.splice(i, 1);
        continue;
      }
      g.wob += dt * 3;
      g.x += (g.vx + Math.sin(g.wob) * 14) * dt;
      g.y += (g.vy + Math.cos(g.wob * 0.8) * 14) * dt;
      // wrap stragglers back in so the field never empties
      if (g.x < -40 || g.x > W + 40 || g.y < -40 || g.y > H + 40) {
        bugs.splice(i, 1);
      }
    }
    for (i = bombs.length - 1; i >= 0; i--) {
      var b = bombs[i];
      b.fuse -= dt;
      if (b.fuse <= 0) { detonate(b); bombs.splice(i, 1); }
    }
    for (i = blasts.length - 1; i >= 0; i--) {
      var e = blasts[i];
      e.life -= dt;
      e.r += (e.maxR - e.r) * Math.min(1, dt * 14);
      if (e.life <= 0) blasts.splice(i, 1);
    }
    for (i = parts.length - 1; i >= 0; i--) {
      var p = parts[i];
      p.life -= dt; p.x += p.vx * dt; p.y += p.vy * dt;
      p.vx *= 0.96; p.vy *= 0.96;
      if (p.life <= 0) parts.splice(i, 1);
    }
  }

  function draw(t) {
    ctx.clearRect(0, 0, W, H);
    // faint grid, tactical vibe
    ctx.strokeStyle = 'rgba(138,122,90,0.14)'; ctx.lineWidth = 1;
    for (var gx = 0; gx <= W; gx += 44) { ctx.beginPath(); ctx.moveTo(gx, 0); ctx.lineTo(gx, H); ctx.stroke(); }
    for (var gy = 0; gy <= H; gy += 44) { ctx.beginPath(); ctx.moveTo(0, gy); ctx.lineTo(W, gy); ctx.stroke(); }

    var i;
    for (i = 0; i < bugs.length; i++) drawBug(bugs[i], t);
    for (i = 0; i < bombs.length; i++) {
      var b = bombs[i], pulse = 1 - b.fuse / b.maxFuse;
      ctx.save();
      ctx.strokeStyle = '#f4772e'; ctx.lineWidth = 2;
      ctx.beginPath(); ctx.arc(b.x, b.y, 10 + pulse * 8, 0, 6.29); ctx.stroke();
      ctx.fillStyle = '#f4772e';
      ctx.beginPath(); ctx.arc(b.x, b.y, 5, 0, 6.29); ctx.fill();
      // crosshair
      ctx.strokeStyle = 'rgba(244,119,46,0.7)'; ctx.lineWidth = 1.5;
      ctx.beginPath();
      ctx.moveTo(b.x - 16, b.y); ctx.lineTo(b.x - 8, b.y);
      ctx.moveTo(b.x + 8, b.y); ctx.lineTo(b.x + 16, b.y);
      ctx.moveTo(b.x, b.y - 16); ctx.lineTo(b.x, b.y - 8);
      ctx.moveTo(b.x, b.y + 8); ctx.lineTo(b.x, b.y + 16);
      ctx.stroke();
      ctx.restore();
    }
    for (i = 0; i < blasts.length; i++) {
      var e = blasts[i], a = Math.max(0, e.life / e.maxLife);
      ctx.save();
      ctx.globalAlpha = a * 0.55;
      ctx.fillStyle = '#f4772e';
      ctx.beginPath(); ctx.arc(e.x, e.y, e.r, 0, 6.29); ctx.fill();
      ctx.globalAlpha = a;
      ctx.strokeStyle = '#ffe9c9'; ctx.lineWidth = 3;
      ctx.beginPath(); ctx.arc(e.x, e.y, e.r, 0, 6.29); ctx.stroke();
      ctx.restore();
    }
    ctx.save();
    for (i = 0; i < parts.length; i++) {
      var p = parts[i];
      ctx.globalAlpha = Math.max(0, p.life / p.maxLife);
      ctx.fillStyle = '#ffb35c';
      ctx.fillRect(p.x - 2, p.y - 2, 4, 4);
    }
    ctx.restore();
  }

  function loop(ts) {
    if (!running) return;
    requestAnimationFrame(loop);
    if (!visible || document.hidden) { last = ts; return; }
    var t = ts / 1000, dt = Math.min(0.05, (ts - last) / 1000 || 0.016);
    last = ts;
    step(dt, t);
    draw(t);
  }

  canvas.addEventListener('pointerdown', function (ev) {
    if (!running) return;
    ev.preventDefault();
    var r = canvas.getBoundingClientRect();
    dropBomb(ev.clientX - r.left, ev.clientY - r.top);
  });

  new IntersectionObserver(function (entries) {
    visible = entries[0].isIntersecting;
  }, { threshold: 0.1 }).observe(canvas);

  window.addEventListener('resize', resize);
  resize();

  if (overlay) {
    var btn = overlay.querySelector('button');
    if (btn) btn.addEventListener('click', function () {
      overlay.style.display = 'none';
      running = true;
      last = performance.now();
      requestAnimationFrame(loop);
    });
  }
})();
