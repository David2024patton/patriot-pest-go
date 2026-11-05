/* Bug Blaster — the Patriot minigame, illustrated edition.
 * Bugs crawl in from the edges toward the snack pile in the middle.
 * Tap anywhere and little Skyler runs over, drops a BAIT BOMB, and hustles
 * back. The bait draws every bug within reach while the fuse burns (2.5s);
 * when it blows, everything inside the blast radius is eliminated. If 10
 * bugs reach the snack pile first, the round ends.
 * Vanilla canvas, no dependencies, no external assets. Only runs after the
 * visitor presses start, and pauses when hidden or off screen. */
(function () {
  'use strict';

  var BB_DEBUG = false;
  function dbg(event, data) {
    if (BB_DEBUG && window.console) console.log('[bugblaster]', event, data || '');
  }

  var canvas = document.getElementById('bugblaster');
  if (!canvas) { dbg('init-skip', 'no canvas'); return; }
  var overlay = document.getElementById('bugblaster-start');
  var hud = document.getElementById('bugblaster-hud');
  var ctx = canvas.getContext('2d');

  /* ---------------- tunables ---------------- */
  var FUSE_S = 2.5;        // bait bomb fuse, seconds
  var BLAST_R = 100;       // blast radius, px
  var LURE_R = 170;        // bait attraction radius, px
  var MAX_BUGS = 30;       // hard cap so phones never melt
  var MAX_PARTS = 130;     // particle cap
  var MAX_BOMBS = 4;       // live bombs at once
  var MAX_QUEUED = 3;      // taps waiting for Skyler
  var SNACK_LIMIT = 10;    // bugs at the pile that end the round
  var SKYLER_SPEED = 300;  // px/s
  var PILE_R = 24;         // snack pile reach radius

  /* ---------------- canvas sizing ---------------- */
  var W = 0, H = 0, DPR = 1;
  function resize() {
    DPR = Math.min(window.devicePixelRatio || 1, 2);
    var r = canvas.getBoundingClientRect();
    W = Math.max(280, Math.floor(r.width));
    H = Math.max(300, Math.floor(r.height));
    canvas.width = W * DPR;
    canvas.height = H * DPR;
    ctx.setTransform(DPR, 0, 0, DPR, 0, 0);
    skyler.home.x = W / 2; skyler.home.y = H / 2 + 46;
    if (!skyler.moving) { skyler.x = skyler.home.x; skyler.y = skyler.home.y; }
    dbg('resize', { w: W, h: H, dpr: DPR });
  }

  /* ---------------- state ---------------- */
  var bugs = [], bombs = [], blasts = [], parts = [], queue = [];
  var running = false, over = false, visible = true;
  var last = 0, spawnT = 0, spawnMs = 900;
  var kills = 0, snacks = 0;

  var skyler = {
    x: 0, y: 0, home: { x: 0, y: 0 },
    tx: 0, ty: 0, moving: false, returning: false,
    face: 1, walkPh: 0, dropX: 0, dropY: 0
  };

  function rand(a, b) { return a + Math.random() * (b - a); }
  function clamp(v, a, b) { return v < a ? a : (v > b ? b : v); }

  function updateHud() {
    if (hud) hud.textContent = 'BUGS ELIMINATED: ' + kills + '  ·  SNACKS LOST: ' + snacks + '/' + SNACK_LIMIT;
  }

  /* ---------------- bugs ---------------- */
  var BUG_TYPES = [
    { kind: 'ant',    w: 0.45, size: [12, 14], spd: [46, 70] },
    { kind: 'beetle', w: 0.35, size: [13, 16], spd: [36, 54] },
    { kind: 'spider', w: 0.20, size: [15, 18], spd: [28, 42] }
  ];
  function pickType() {
    var r = Math.random(), acc = 0;
    for (var i = 0; i < BUG_TYPES.length; i++) {
      acc += BUG_TYPES[i].w;
      if (r <= acc) return BUG_TYPES[i];
    }
    return BUG_TYPES[0];
  }

  function spawnBug() {
    if (bugs.length >= MAX_BUGS) return;
    var t = pickType();
    var edge = Math.floor(Math.random() * 4), x, y;
    if (edge === 0) { x = rand(0, W); y = -22; }
    else if (edge === 1) { x = W + 22; y = rand(0, H); }
    else if (edge === 2) { x = rand(0, W); y = H + 22; }
    else { x = -22; y = rand(0, H); }
    var ramp = Math.min(1.8, 1 + kills / 260);
    bugs.push({
      kind: t.kind, x: x, y: y,
      ang: Math.atan2(H / 2 - y, W / 2 - x),
      spd: rand(t.spd[0], t.spd[1]) * ramp,
      size: rand(t.size[0], t.size[1]),
      wob: rand(0, 6.28), legPh: rand(0, 6.28),
      jx: rand(-26, 26), jy: rand(-26, 26),
      dying: false, deathT: 0, alpha: 1
    });
  }

  function nearestLure(g) {
    var best = null, bd = LURE_R;
    for (var i = 0; i < bombs.length; i++) {
      var d = Math.hypot(bombs[i].x - g.x, bombs[i].y - g.y);
      if (d < bd) { bd = d; best = bombs[i]; }
    }
    return best;
  }

  function turnToward(g, want, dt) {
    var diff = want - g.ang;
    while (diff > Math.PI) diff -= Math.PI * 2;
    while (diff < -Math.PI) diff += Math.PI * 2;
    var maxTurn = 3.4 * dt;
    g.ang += clamp(diff, -maxTurn, maxTurn);
  }

  function stepBugs(dt, t) {
    for (var i = bugs.length - 1; i >= 0; i--) {
      var g = bugs[i];
      if (g.dying) {
        g.deathT += dt;
        g.alpha = Math.max(0, 1 - g.deathT / 0.7);
        if (g.alpha <= 0) bugs.splice(i, 1);
        continue;
      }
      g.wob += dt * 3;
      g.legPh += dt * g.spd * 0.55;

      var lure = nearestLure(g);
      var tx, ty, eating = false;
      if (lure) {
        tx = lure.x; ty = lure.y;
        eating = Math.hypot(tx - g.x, ty - g.y) < 18;
      } else {
        tx = W / 2 + g.jx; ty = H / 2 + g.jy;
      }
      turnToward(g, Math.atan2(ty - g.y, tx - g.x), dt);
      var sp = g.spd * (eating ? 0.15 : 1);
      g.x += (Math.cos(g.ang) * sp + Math.cos(g.ang + 1.57) * Math.sin(g.wob) * 12) * dt;
      g.y += (Math.sin(g.ang) * sp + Math.sin(g.ang + 1.57) * Math.sin(g.wob) * 12) * dt;

      // reached the snack pile (and no bait is pulling it away)
      if (!lure && Math.hypot(g.x - W / 2, g.y - H / 2) <= PILE_R) {
        puff(g.x, g.y, '#d9b06c', 6);
        bugs.splice(i, 1);
        if (!over) {
          snacks++;
          updateHud();
          dbg('snack-lost', { snacks: snacks });
          if (snacks >= SNACK_LIMIT) roundOver();
        }
        continue;
      }
      if (g.x < -60 || g.x > W + 60 || g.y < -60 || g.y > H + 60) bugs.splice(i, 1);
    }
  }

  /* ---------------- bombs ---------------- */
  function placeBomb(x, y) {
    if (bombs.length >= MAX_BOMBS) return;
    bombs.push({ x: x, y: y, fuse: FUSE_S, maxFuse: FUSE_S });
    dbg('bomb-placed', { x: Math.round(x), y: Math.round(y) });
  }

  function detonate(b) {
    blasts.push({ x: b.x, y: b.y, r: 10, maxR: BLAST_R, life: 0.45, maxLife: 0.45 });
    burst(b.x, b.y);
    var hit = 0;
    for (var j = 0; j < bugs.length; j++) {
      var g = bugs[j];
      if (!g.dying && Math.hypot(g.x - b.x, g.y - b.y) <= BLAST_R) {
        g.dying = true; g.deathT = 0; kills++; hit++;
      }
    }
    if (hit) {
      updateHud();
      if (kills % 12 === 0) spawnMs = Math.max(340, spawnMs - 60); // waves get hotter
    }
    dbg('detonate', { hits: hit, kills: kills });
  }

  function burst(x, y) {
    for (var i = 0; i < 16 && parts.length < MAX_PARTS; i++) {
      var a = rand(0, 6.28), sp = rand(60, 250);
      parts.push({ x: x, y: y, vx: Math.cos(a) * sp, vy: Math.sin(a) * sp, life: rand(0.3, 0.7), maxLife: 0.7, col: '#ffb35c', r: rand(2, 4) });
    }
    for (var k = 0; k < 5 && parts.length < MAX_PARTS; k++) {
      parts.push({ x: x + rand(-14, 14), y: y + rand(-14, 14), vx: rand(-16, 16), vy: rand(-30, -8), life: rand(0.5, 0.9), maxLife: 0.9, col: '#e8e2d4', r: rand(6, 12), puff: true });
    }
  }

  function puff(x, y, col, n) {
    for (var i = 0; i < n && parts.length < MAX_PARTS; i++) {
      var a = rand(0, 6.28), sp = rand(20, 70);
      parts.push({ x: x, y: y, vx: Math.cos(a) * sp, vy: Math.sin(a) * sp, life: rand(0.25, 0.5), maxLife: 0.5, col: col, r: rand(2, 3.5) });
    }
  }

  /* ---------------- Skyler ---------------- */
  function stepSkyler(dt) {
    if (!skyler.moving) {
      if (queue.length) {
        var q = queue.shift();
        skyler.dropX = q.x; skyler.dropY = q.y;
        skyler.tx = q.x; skyler.ty = q.y;
        skyler.moving = true; skyler.returning = false;
      }
      return;
    }
    var dx = skyler.tx - skyler.x, dy = skyler.ty - skyler.y;
    var d = Math.hypot(dx, dy);
    skyler.walkPh += dt * 11;
    if (dx !== 0) skyler.face = dx > 0 ? 1 : -1;
    if (d <= Math.max(6, SKYLER_SPEED * dt)) {
      skyler.x = skyler.tx; skyler.y = skyler.ty;
      if (!skyler.returning) {
        placeBomb(skyler.dropX, skyler.dropY);
        skyler.tx = skyler.home.x; skyler.ty = skyler.home.y;
        skyler.returning = true;
      } else {
        skyler.moving = false; skyler.returning = false;
      }
    } else {
      skyler.x += dx / d * SKYLER_SPEED * dt;
      skyler.y += dy / d * SKYLER_SPEED * dt;
    }
  }

  /* ---------------- shared drawing helpers ---------------- */
  function gloss(x, y, rx, ry, rot) {
    ctx.save();
    ctx.globalAlpha = 0.20;
    ctx.fillStyle = '#ffffff';
    ctx.beginPath(); ctx.ellipse(x, y, rx, ry, rot || 0, 0, 6.29); ctx.fill();
    ctx.restore();
  }
  function cuteEyes(x1, y1, x2, y2, r, dying) {
    if (dying) {
      ctx.strokeStyle = '#ffe9c9'; ctx.lineWidth = 1.6;
      var pts = [[x1, y1], [x2, y2]];
      for (var i = 0; i < 2; i++) {
        var ex = pts[i][0], ey = pts[i][1], er = r * 0.8;
        ctx.beginPath();
        ctx.moveTo(ex - er, ey - er); ctx.lineTo(ex + er, ey + er);
        ctx.moveTo(ex + er, ey - er); ctx.lineTo(ex - er, ey + er);
        ctx.stroke();
      }
      return;
    }
    ctx.fillStyle = '#ffffff';
    ctx.beginPath(); ctx.arc(x1, y1, r, 0, 6.29); ctx.fill();
    ctx.beginPath(); ctx.arc(x2, y2, r, 0, 6.29); ctx.fill();
    ctx.fillStyle = '#241f14';
    ctx.beginPath(); ctx.arc(x1 + r * 0.25, y1, r * 0.52, 0, 6.29); ctx.fill();
    ctx.beginPath(); ctx.arc(x2 + r * 0.25, y2, r * 0.52, 0, 6.29); ctx.fill();
  }
  function star(cx, cy, r, col) {
    ctx.fillStyle = col;
    ctx.beginPath();
    for (var i = 0; i < 10; i++) {
      var rad = i % 2 === 0 ? r : r * 0.45;
      var a = -Math.PI / 2 + i * Math.PI / 5;
      var px = cx + Math.cos(a) * rad, py = cy + Math.sin(a) * rad;
      if (i === 0) ctx.moveTo(px, py); else ctx.lineTo(px, py);
    }
    ctx.closePath(); ctx.fill();
  }

  /* ---------------- bug drawing (illustrated, canvas-drawn) ---------------- */
  function drawAnt(g) {
    var s = g.size, sw = Math.sin(g.legPh) * s * 0.3;
    ctx.strokeStyle = '#17140f'; ctx.lineWidth = Math.max(1.3, s * 0.1);
    ctx.lineCap = 'round';
    for (var l = -1; l <= 1; l++) {
      var lx = l * s * 0.28, ph = l === 0 ? sw : -sw;
      ctx.beginPath(); ctx.moveTo(lx, -s * 0.2); ctx.lineTo(lx + ph, -s * 0.95); ctx.lineTo(lx + ph * 1.5, -s * 1.25); ctx.stroke();
      ctx.beginPath(); ctx.moveTo(lx, s * 0.2); ctx.lineTo(lx - ph, s * 0.95); ctx.lineTo(lx - ph * 1.5, s * 1.25); ctx.stroke();
    }
    // gaster (big back segment)
    ctx.fillStyle = '#2b2620';
    ctx.beginPath(); ctx.ellipse(-s * 0.95, 0, s * 0.85, s * 0.62, 0, 0, 6.29); ctx.fill();
    gloss(-s * 1.15, -s * 0.25, s * 0.4, s * 0.18, -0.5);
    // thorax
    ctx.fillStyle = '#38312a';
    ctx.beginPath(); ctx.arc(0, 0, s * 0.42, 0, 6.29); ctx.fill();
    // head
    ctx.fillStyle = '#2b2620';
    ctx.beginPath(); ctx.arc(s * 0.78, 0, s * 0.52, 0, 6.29); ctx.fill();
    gloss(s * 0.62, -s * 0.22, s * 0.22, s * 0.1, -0.4);
    // elbowed antennae
    ctx.strokeStyle = '#17140f'; ctx.lineWidth = Math.max(1.1, s * 0.07);
    ctx.beginPath(); ctx.moveTo(s * 1.1, -s * 0.18); ctx.lineTo(s * 1.5, -s * 0.55); ctx.lineTo(s * 1.95, -s * 0.4); ctx.stroke();
    ctx.beginPath(); ctx.moveTo(s * 1.1, s * 0.18); ctx.lineTo(s * 1.5, s * 0.55); ctx.lineTo(s * 1.95, s * 0.4); ctx.stroke();
    cuteEyes(s * 0.92, -s * 0.2, s * 0.92, s * 0.2, s * 0.17, g.dying);
  }

  function drawSpider(g) {
    var s = g.size;
    ctx.strokeStyle = '#241f2c'; ctx.lineWidth = Math.max(1.4, s * 0.11);
    ctx.lineCap = 'round';
    for (var side = -1; side <= 1; side += 2) {
      for (var l = 0; l < 4; l++) {
        var bx = (l - 1.5) * s * 0.3, sw = Math.sin(g.legPh + l * 1.1) * s * 0.22;
        var kx = bx + sw, ky = side * s * 1.05;
        ctx.beginPath(); ctx.moveTo(bx, side * s * 0.3);
        ctx.lineTo(kx, ky); ctx.lineTo(kx + sw * 0.6, ky + side * s * 0.5); ctx.stroke();
      }
    }
    // abdomen
    ctx.fillStyle = '#463d55';
    ctx.beginPath(); ctx.arc(-s * 0.5, 0, s * 0.82, 0, 6.29); ctx.fill();
    gloss(-s * 0.75, -s * 0.3, s * 0.36, s * 0.16, -0.5);
    // cephalothorax
    ctx.fillStyle = '#332c40';
    ctx.beginPath(); ctx.arc(s * 0.45, 0, s * 0.52, 0, 6.29); ctx.fill();
    gloss(s * 0.3, -s * 0.2, s * 0.2, s * 0.09, -0.4);
    cuteEyes(s * 0.72, -s * 0.19, s * 0.72, s * 0.19, s * 0.16, g.dying);
  }

  function drawBeetle(g) {
    var s = g.size, sw = Math.sin(g.legPh) * s * 0.28;
    ctx.strokeStyle = '#3a2410'; ctx.lineWidth = Math.max(1.3, s * 0.1);
    ctx.lineCap = 'round';
    for (var l = -1; l <= 1; l++) {
      var lx = l * s * 0.42, ph = l === 0 ? sw : -sw;
      ctx.beginPath(); ctx.moveTo(lx, -s * 0.3); ctx.lineTo(lx + ph, -s * 0.95); ctx.stroke();
      ctx.beginPath(); ctx.moveTo(lx, s * 0.3); ctx.lineTo(lx - ph, s * 0.95); ctx.stroke();
    }
    // shell
    ctx.fillStyle = '#7a4a26';
    ctx.beginPath(); ctx.ellipse(-s * 0.1, 0, s * 1.05, s * 0.72, 0, 0, 6.29); ctx.fill();
    ctx.strokeStyle = '#54300f'; ctx.lineWidth = 1.4;
    ctx.beginPath(); ctx.ellipse(-s * 0.1, 0, s * 1.05, s * 0.72, 0, 0, 6.29); ctx.stroke();
    // wing split + shading bands
    ctx.strokeStyle = 'rgba(46,26,8,0.55)'; ctx.lineWidth = 1.3;
    ctx.beginPath(); ctx.moveTo(-s * 1.12, 0); ctx.lineTo(s * 0.78, 0); ctx.stroke();
    ctx.beginPath(); ctx.moveTo(-s * 0.55, -s * 0.62); ctx.quadraticCurveTo(-s * 0.4, 0, -s * 0.55, s * 0.62); ctx.stroke();
    gloss(-s * 0.35, -s * 0.3, s * 0.45, s * 0.16, -0.35);
    // head
    ctx.fillStyle = '#4a2c14';
    ctx.beginPath(); ctx.arc(s * 0.92, 0, s * 0.4, 0, 6.29); ctx.fill();
    // antennae
    ctx.strokeStyle = '#3a2410'; ctx.lineWidth = Math.max(1.1, s * 0.07);
    ctx.beginPath(); ctx.moveTo(s * 1.2, -s * 0.14); ctx.quadraticCurveTo(s * 1.7, -s * 0.5, s * 2.05, -s * 0.3); ctx.stroke();
    ctx.beginPath(); ctx.moveTo(s * 1.2, s * 0.14); ctx.quadraticCurveTo(s * 1.7, s * 0.5, s * 2.05, s * 0.3); ctx.stroke();
    cuteEyes(s * 1.0, -s * 0.17, s * 1.0, s * 0.17, s * 0.14, g.dying);
  }

  function drawBug(g) {
    ctx.save();
    ctx.globalAlpha = g.alpha;
    ctx.translate(g.x, g.y);
    ctx.rotate(g.ang + (g.dying ? Math.PI : 0));
    if (g.kind === 'ant') drawAnt(g);
    else if (g.kind === 'spider') drawSpider(g);
    else drawBeetle(g);
    ctx.restore();
  }

  /* ---------------- Skyler drawing ---------------- */
  function drawSkyler(t) {
    var bob = skyler.moving ? Math.abs(Math.sin(skyler.walkPh)) * 2.2 : Math.sin(t * 2) * 0.7;
    var legSw = skyler.moving ? Math.sin(skyler.walkPh) * 6 : 0;
    ctx.save();
    ctx.translate(skyler.x, skyler.y - bob);
    // ground shadow
    ctx.fillStyle = 'rgba(20,16,8,0.25)';
    ctx.beginPath(); ctx.ellipse(0, bob + 1, 15, 4.5, 0, 0, 6.29); ctx.fill();
    ctx.scale(skyler.face, 1);
    // legs + boots
    ctx.strokeStyle = '#31405a'; ctx.lineWidth = 5.5; ctx.lineCap = 'round';
    ctx.beginPath(); ctx.moveTo(-3, -16); ctx.lineTo(-3 + legSw * 0.7, -1); ctx.stroke();
    ctx.beginPath(); ctx.moveTo(3, -16); ctx.lineTo(3 - legSw * 0.7, -1); ctx.stroke();
    ctx.fillStyle = '#241f14';
    ctx.fillRect(-6 + legSw * 0.7, -3.5, 7, 3.5);
    ctx.fillRect(0 - legSw * 0.7, -3.5, 7, 3.5);
    // sprayer tank on his back
    ctx.fillStyle = '#e8e2d4';
    ctx.beginPath();
    if (ctx.roundRect) ctx.roundRect(-13, -36, 7, 17, 3); else ctx.rect(-13, -36, 7, 17);
    ctx.fill();
    ctx.fillStyle = '#f4772e'; ctx.fillRect(-13, -36, 7, 4);
    // torso: green work shirt
    ctx.fillStyle = '#3f6b35';
    ctx.beginPath();
    if (ctx.roundRect) ctx.roundRect(-10, -38, 20, 23, 6); else ctx.rect(-10, -38, 20, 23);
    ctx.fill();
    star(4.5, -30, 3.2, '#f4772e'); // chest star
    // arm holding the spray wand, angled down-forward
    ctx.strokeStyle = '#3f6b35'; ctx.lineWidth = 5;
    ctx.beginPath(); ctx.moveTo(8, -34); ctx.lineTo(15, -22); ctx.stroke();
    ctx.strokeStyle = '#f2c49b'; ctx.lineWidth = 4;
    ctx.beginPath(); ctx.moveTo(15, -22); ctx.lineTo(17, -19); ctx.stroke();
    ctx.strokeStyle = '#8a8f96'; ctx.lineWidth = 2.4;
    ctx.beginPath(); ctx.moveTo(17, -19); ctx.lineTo(24, -10); ctx.stroke();
    ctx.fillStyle = '#f4772e';
    ctx.beginPath(); ctx.arc(24, -10, 2.2, 0, 6.29); ctx.fill();
    // other arm swings as he walks
    ctx.strokeStyle = '#3f6b35'; ctx.lineWidth = 5;
    ctx.beginPath(); ctx.moveTo(-8, -34); ctx.lineTo(-12, -22 + legSw * 0.5); ctx.stroke();
    // head
    ctx.fillStyle = '#f2c49b';
    ctx.beginPath(); ctx.arc(0, -46, 8.5, 0, 6.29); ctx.fill();
    // friendly face
    ctx.fillStyle = '#241f14';
    ctx.beginPath(); ctx.arc(-2.8, -47, 1.15, 0, 6.29); ctx.fill();
    ctx.beginPath(); ctx.arc(2.8, -47, 1.15, 0, 6.29); ctx.fill();
    ctx.strokeStyle = '#241f14'; ctx.lineWidth = 1.2;
    ctx.beginPath(); ctx.arc(0.5, -43.5, 3.4, 0.25, Math.PI - 0.6); ctx.stroke();
    // cap with brim + orange star
    ctx.fillStyle = '#3f6b35';
    ctx.beginPath(); ctx.arc(0, -49.5, 8.5, Math.PI, 0); ctx.closePath(); ctx.fill();
    ctx.fillRect(2, -51.5, 11, 3);
    star(0, -53.5, 3.4, '#f4772e');
    ctx.restore();
  }

  /* ---------------- scenery + effects drawing ---------------- */
  function drawPile(t) {
    ctx.save();
    ctx.translate(W / 2, H / 2);
    // house zone ring
    ctx.strokeStyle = 'rgba(244,119,46,0.4)'; ctx.lineWidth = 1.6;
    ctx.setLineDash([6, 6]);
    ctx.beginPath(); ctx.arc(0, 0, 34, 0, 6.29); ctx.stroke();
    ctx.setLineDash([]);
    // crumbs + a cookie + a cheese wedge: the snack pile
    ctx.fillStyle = '#d9b06c';
    ctx.beginPath(); ctx.ellipse(0, 5, 22, 9, 0, 0, 6.29); ctx.fill();
    ctx.fillStyle = '#c98f4e';
    ctx.beginPath(); ctx.arc(-8, 0, 8, 0, 6.29); ctx.fill();
    ctx.fillStyle = '#5a3a1a';
    [[-11, -2], [-5, 2], [-8, -5]].forEach(function (c) {
      ctx.beginPath(); ctx.arc(c[0], c[1], 1.5, 0, 6.29); ctx.fill();
    });
    ctx.fillStyle = '#f4c542';
    ctx.beginPath(); ctx.moveTo(6, -8); ctx.lineTo(17, 1); ctx.lineTo(5, 3); ctx.closePath(); ctx.fill();
    ctx.fillStyle = 'rgba(160,110,20,0.7)';
    ctx.beginPath(); ctx.arc(9, -2, 1.8, 0, 6.29); ctx.fill();
    ctx.beginPath(); ctx.arc(12, 0.5, 1.3, 0, 6.29); ctx.fill();
    ctx.fillStyle = '#e5c07a';
    [[-18, 8], [15, 9], [2, 11], [-3, -9]].forEach(function (c) {
      ctx.beginPath(); ctx.arc(c[0], c[1], 2, 0, 6.29); ctx.fill();
    });
    ctx.restore();
  }

  function drawBomb(b, t) {
    var frac = Math.max(0, b.fuse / b.maxFuse);
    ctx.save();
    ctx.translate(b.x, b.y);
    // blast radius preview
    ctx.strokeStyle = 'rgba(244,119,46,0.30)'; ctx.lineWidth = 1.4;
    ctx.setLineDash([5, 7]);
    ctx.beginPath(); ctx.arc(0, 0, BLAST_R, 0, 6.29); ctx.stroke();
    ctx.setLineDash([]);
    ctx.fillStyle = 'rgba(244,119,46,0.05)';
    ctx.beginPath(); ctx.arc(0, 0, BLAST_R, 0, 6.29); ctx.fill();
    var pulse = frac < 0.25 ? 1 + Math.sin(t * 22) * 0.07 : 1;
    ctx.scale(pulse, pulse);
    // bomb body
    ctx.fillStyle = '#2b2b30';
    ctx.beginPath(); ctx.arc(0, 0, 11, 0, 6.29); ctx.fill();
    gloss(-4, -4.5, 4.2, 2, -0.5);
    // the bait: a cheese wedge lure parked on top
    ctx.fillStyle = '#f4c542';
    ctx.beginPath(); ctx.moveTo(-6, -13); ctx.lineTo(5, -18); ctx.lineTo(6, -10); ctx.closePath(); ctx.fill();
    ctx.fillStyle = 'rgba(160,110,20,0.8)';
    ctx.beginPath(); ctx.arc(0, -13.5, 1.7, 0, 6.29); ctx.fill();
    ctx.beginPath(); ctx.arc(3.4, -11.5, 1.2, 0, 6.29); ctx.fill();
    // fuse spark
    var sx = 5 + Math.sin(t * 31) * 1.6, sy = -19 + Math.cos(t * 27) * 1.4;
    star(sx, sy, 3.4 + Math.sin(t * 40) * 0.8, '#ffb35c');
    // fuse timer ring: shrinks as the fuse burns, green -> amber -> red
    ctx.strokeStyle = frac > 0.5 ? '#58b368' : (frac > 0.25 ? '#f4a72e' : '#e5484d');
    ctx.lineWidth = 3.4; ctx.lineCap = 'round';
    ctx.beginPath(); ctx.arc(0, 0, 17.5, -Math.PI / 2, -Math.PI / 2 + frac * Math.PI * 2); ctx.stroke();
    ctx.restore();
  }

  function draw(t) {
    ctx.clearRect(0, 0, W, H);
    // faint grid, tactical vibe
    ctx.strokeStyle = 'rgba(138,122,90,0.14)'; ctx.lineWidth = 1;
    for (var gx = 0; gx <= W; gx += 44) { ctx.beginPath(); ctx.moveTo(gx, 0); ctx.lineTo(gx, H); ctx.stroke(); }
    for (var gy = 0; gy <= H; gy += 44) { ctx.beginPath(); ctx.moveTo(0, gy); ctx.lineTo(W, gy); ctx.stroke(); }

    drawPile(t);
    var i;
    for (i = 0; i < bombs.length; i++) drawBomb(bombs[i], t);
    for (i = 0; i < bugs.length; i++) drawBug(bugs[i]);
    drawSkyler(t);

    for (i = 0; i < blasts.length; i++) {
      var e = blasts[i], a = Math.max(0, e.life / e.maxLife);
      ctx.save();
      ctx.globalAlpha = a * 0.5;
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
      ctx.globalAlpha = Math.max(0, p.life / p.maxLife) * (p.puff ? 0.55 : 1);
      ctx.fillStyle = p.col;
      ctx.beginPath(); ctx.arc(p.x, p.y, p.r, 0, 6.29); ctx.fill();
    }
    ctx.restore();
  }

  /* ---------------- simulation step ---------------- */
  function step(dt, t) {
    if (!over) {
      spawnT += dt * 1000;
      if (spawnT >= spawnMs) { spawnT = 0; spawnBug(); }
    }
    stepBugs(dt, t);
    stepSkyler(dt);
    var i;
    for (i = bombs.length - 1; i >= 0; i--) {
      bombs[i].fuse -= dt;
      if (bombs[i].fuse <= 0) { detonate(bombs[i]); bombs.splice(i, 1); }
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
      if (p.puff) p.r += dt * 14;
      if (p.life <= 0) parts.splice(i, 1);
    }
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

  /* ---------------- round control ---------------- */
  function startRound() {
    resize(); // popup is open by now, so the arena has its real size
    bugs = []; bombs = []; blasts = []; parts = []; queue = [];
    kills = 0; snacks = 0; spawnMs = 900; spawnT = 0;
    over = false;
    skyler.moving = false; skyler.returning = false;
    skyler.x = skyler.home.x; skyler.y = skyler.home.y;
    updateHud();
    if (overlay) overlay.style.display = 'none';
    running = true;
    last = performance.now();
    requestAnimationFrame(loop);
    dbg('round-start', { w: W, h: H });
  }

  function roundOver() {
    if (over) return;
    over = true;
    queue = [];
    dbg('round-over', { kills: kills, snacks: snacks });
    setTimeout(function () {
      if (!overlay) return;
      var p = overlay.querySelector('p');
      if (p) p.innerHTML = 'The bugs got the snacks!<br>Final score: BUGS ELIMINATED: ' + kills;
      var btn = overlay.querySelector('button');
      if (btn) btn.textContent = 'PLAY AGAIN';
      overlay.style.display = '';
    }, 650);
  }

  /* ---------------- wiring ---------------- */
  canvas.addEventListener('pointerdown', function (ev) {
    if (!running || over) return;
    ev.preventDefault();
    var r = canvas.getBoundingClientRect();
    var x = ev.clientX - r.left, y = ev.clientY - r.top;
    if (queue.length < MAX_QUEUED) {
      queue.push({ x: x, y: y });
      dbg('tap-queued', { x: Math.round(x), y: Math.round(y), queued: queue.length });
    }
  });

  new IntersectionObserver(function (entries) {
    visible = entries[0].isIntersecting;
  }, { threshold: 0.1 }).observe(canvas);

  window.addEventListener('resize', resize);
  resize();

  if (overlay) {
    var btn = overlay.querySelector('button');
    if (btn) btn.addEventListener('click', startRound);
  }
  dbg('init', 'ready');
})();
