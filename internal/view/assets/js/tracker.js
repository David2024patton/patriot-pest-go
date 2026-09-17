/* ============================================================
   TRACKER.JS - first-party pageview + click beacons for the
   Patriot Pest Control marketing site. Vanilla, no dependencies.

   Fires:
     - page_view  POST /api/track/view   on every load
     - click      POST /api/track/event  delegated on a/button clicks

   Identity:
     - ppc_vid : persistent visitor id (localStorage, random hex)
     - ppc_sid : session id (sessionStorage), rotated after 30 min
                 of inactivity (last activity in localStorage ppc_last)

   Respects:
     - navigator.doNotTrack === "1"  -> sends nothing
     - ppc_admin=1 cookie (set at admin login) -> sends nothing,
       so the owners' own visits never pollute the dashboard
   ============================================================ */
(function () {
  "use strict";

  /* ---------- opt-outs ---------- */
  if (navigator.doNotTrack === "1") return;
  if (/(?:^|;\s*)ppc_admin=1(?:;|$)/.test(document.cookie)) return;

  /* ---------- storage helpers (private mode may throw) ---------- */
  function lsGet(k) { try { return localStorage.getItem(k); } catch (e) { return null; } }
  function lsSet(k, v) { try { localStorage.setItem(k, v); } catch (e) {} }
  function ssGet(k) { try { return sessionStorage.getItem(k); } catch (e) { return null; } }
  function ssSet(k, v) { try { sessionStorage.setItem(k, v); } catch (e) {} }

  function hex(n) {
    var b = new Uint8Array(n);
    if (window.crypto && crypto.getRandomValues) {
      crypto.getRandomValues(b);
    } else {
      for (var i = 0; i < n; i++) b[i] = Math.floor(Math.random() * 256);
    }
    var s = "";
    for (var j = 0; j < b.length; j++) s += ("0" + b[j].toString(16)).slice(-2);
    return s;
  }

  function validId(v) { return typeof v === "string" && /^[0-9a-f]{64}$/.test(v); }

  /* ---------- visitor id (persistent) ---------- */
  var vid = lsGet("ppc_vid");
  if (!validId(vid)) { vid = hex(32); lsSet("ppc_vid", vid); }

  /* ---------- session id (30 min inactivity rotates it) ---------- */
  var SESSION_TTL = 30 * 60 * 1000;
  var sid = ssGet("ppc_sid");
  var last = parseInt(lsGet("ppc_last") || "0", 10);
  var nowMs = Date.now();
  if (!validId(sid) || (nowMs - last) > SESSION_TTL) {
    sid = hex(32);
    ssSet("ppc_sid", sid);
  }
  function touch() { lsSet("ppc_last", String(Date.now())); }
  touch();

  /* ---------- small utils ---------- */
  function cap(s, n) {
    s = String(s == null ? "" : s);
    return s.length > n ? s.slice(0, n) : s;
  }
  function qs(name) {
    var m = new RegExp("[?&]" + name + "=([^&#]*)").exec(location.search);
    return m ? decodeURIComponent(m[1].replace(/\+/g, " ")) : "";
  }
  function pagePath() { return cap(location.pathname || "/", 512); }

  /* ---------- transport: sendBeacon, fetch+keepalive fallback ---------- */
  function send(url, data) {
    var body;
    try { body = JSON.stringify(data); } catch (e) { return; }
    try {
      if (navigator.sendBeacon) {
        navigator.sendBeacon(url, new Blob([body], { type: "application/json" }));
        return;
      }
    } catch (e) {}
    try {
      fetch(url, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: body,
        keepalive: true
      });
    } catch (e) {}
  }

  /* ---------- page view ---------- */
  send("/api/track/view", {
    vid: vid,
    sid: sid,
    path: pagePath(),
    ref: cap(document.referrer || "", 2048),
    utm_source: cap(qs("utm_source"), 64),
    utm_medium: cap(qs("utm_medium"), 64),
    utm_campaign: cap(qs("utm_campaign"), 64),
    lang: cap(navigator.language || "", 16)
  });

  /* ---------- clicks (delegated: catches dynamically added links too) ---------- */
  document.addEventListener("click", function (ev) {
    touch();
    var el = ev.target && ev.target.closest ? ev.target.closest("a,button") : null;
    if (!el) return;

    var rawHref = el.getAttribute("href") || "";
    var hrefPath = "";
    if (rawHref) {
      try {
        var u = new URL(rawHref, location.href);
        if (u.origin === location.origin) hrefPath = u.pathname || "/";
      } catch (e) {}
    }

    var label = el.getAttribute("aria-label") ||
      (el.textContent || "").replace(/\s+/g, " ").trim() ||
      el.id || el.value || hrefPath;
    label = cap(label, 120);

    send("/api/track/event", {
      vid: vid,
      sid: sid,
      path: pagePath(),
      kind: "click",
      el: cap(el.tagName || "", 16),
      label: label,
      href: cap(hrefPath, 512)
    });
  }, false);
})();
