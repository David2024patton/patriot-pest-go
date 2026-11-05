// Pop UI behavior for pop.patriotpest.pro: tile rail arrows and the popup
// system (quick view layer + quote form layer). Rules from the pop_ui skill:
// popups are centered, opening never focuses a text field, Escape peels one
// layer at a time, the quote form does not close on a backdrop tap.
(function () {
  "use strict";

  // ---- Rails ------------------------------------------------------------
  // Pairs are keyed by element, not by id, so a page shown inside the
  // viewer can carry rails whose ids match the base page's without
  // colliding. Arrows fade at the ends and vanish when the row fits.
  var railPairs = [];
  function railState(pair) {
    if (!pair || !pair.rail) return;
    var max = pair.rail.scrollWidth - pair.rail.clientWidth;
    if (max <= 4) { pair.box.classList.add("gone"); return; }
    pair.box.classList.remove("gone");
    if (pair.prev) pair.prev.classList.toggle("off", pair.rail.scrollLeft <= 4);
    if (pair.next) pair.next.classList.toggle("off", pair.rail.scrollLeft >= max - 4);
  }
  function registerRails(scope) {
    scope.querySelectorAll("[data-rail-prev],[data-rail-next]").forEach(function (btn) {
      var id = btn.getAttribute("data-rail-prev") || btn.getAttribute("data-rail-next");
      var rail = scope.querySelector('[id="' + id + '"]');
      if (!rail) return;
      var pair = null;
      for (var i = 0; i < railPairs.length; i++) if (railPairs[i].rail === rail) pair = railPairs[i];
      if (!pair) {
        pair = { rail: rail, prev: null, next: null, box: btn.parentElement };
        railPairs.push(pair);
        rail.addEventListener("scroll", function () { railState(pair); }, { passive: true });
      }
      if (btn.hasAttribute("data-rail-prev")) pair.prev = btn; else pair.next = btn;
      btn._railPair = pair;
      railState(pair);
    });
  }

  // ---- Top pills: scroll, never navigate (streaming pattern) -------------
  // On the home page a pill smooth-scrolls to its section. From any other
  // page the pill links home with the hash, and the load handler below
  // finishes the scroll. There is no Home pill: the brand is the way home.
  function scrollToSection(id) {
    var el = document.getElementById(id);
    if (!el) return false;
    el.scrollIntoView({ behavior: "smooth", block: "start" });
    return true;
  }
  if (window.location.pathname === "/" && window.location.hash.length > 1) {
    var hashTarget = document.getElementById(window.location.hash.slice(1));
    if (hashTarget) {
      // Images and fonts shift the layout after first paint, so one early
      // scroll lands short. Re-scroll on load and once more shortly after,
      // unless the visitor starts scrolling on their own.
      var hashDone = false;
      var cancelHash = function () { hashDone = true; };
      window.addEventListener("wheel", cancelHash, { passive: true });
      window.addEventListener("touchstart", cancelHash, { passive: true });
      var applyHash = function () {
        if (hashDone) return;
        hashTarget.scrollIntoView();
      };
      setTimeout(applyHash, 80);
      window.addEventListener("load", function () {
        applyHash();
        setTimeout(applyHash, 500);
      });
    }
  }

  // ---- Popup stack --------------------------------------------------------
  var stack = [];

  function openPopup(el) {
    if (!el || stack.indexOf(el) !== -1) return;
    el.hidden = false;
    stack.push(el);
    document.body.classList.add("pop-locked");
    var x = el.querySelector(".pop-x");
    if (x) x.focus({ preventScroll: true });
  }

  function closeTop() {
    var el = stack.pop();
    if (!el) return;
    el.hidden = true;
    if (!stack.length) document.body.classList.remove("pop-locked");
  }

  function closePopup(el) {
    var i = stack.indexOf(el);
    if (i === -1) return;
    stack.splice(i, 1);
    el.hidden = true;
    if (!stack.length) document.body.classList.remove("pop-locked");
  }

  // Quick view: fill the shared popup from a hidden source block.
  // More Info opens layer 2 (the bigger window) filled from the source's
  // .qv-more block, like alphaflux.net. If a source has no .qv-more, the
  // button falls back to the source's data-more page.
  var qv = document.getElementById("pop-qv");
  var moreLayer = document.getElementById("pop-more");
  var currentSrc = null;
  function fillQuickView(btn) {
    var src = document.getElementById(btn.getAttribute("data-qv"));
    if (!src || !qv) return;
    currentSrc = src;
    var body = qv.querySelector(".pop-qv-body");
    body.innerHTML = src.innerHTML;
    var embedded = body.querySelector(".qv-more");
    if (embedded) embedded.remove();
    qv.querySelector(".pop-popup").classList.toggle("has-art", !!src.querySelector(".pop-qv-art"));
    var more = qv.querySelector(".pop-qv-more");
    var hasMore = !!src.querySelector(".qv-more") || !!src.getAttribute("data-more");
    more.style.display = hasMore ? "" : "none";
    openPopup(qv);
  }

  var moreBtn = qv ? qv.querySelector(".pop-qv-more") : null;
  if (moreBtn) {
    moreBtn.addEventListener("click", function () {
      if (!currentSrc) return;
      var block = currentSrc.querySelector(".qv-more");
      if (block && moreLayer) {
        var head = "";
        var lab = currentSrc.querySelector(".pop-lab");
        var h3 = currentSrc.querySelector("h3");
        if (lab) head += '<div class="pop-lab">' + lab.innerHTML + "</div>";
        if (h3) head += "<h3>" + h3.innerHTML + "</h3>";
        moreLayer.querySelector(".pop-more-body").innerHTML = head + block.innerHTML;
        openPopup(moreLayer);
      } else if (currentSrc.getAttribute("data-more")) {
        openPage(currentSrc.getAttribute("data-more"));
      }
    });
  }

  // Backdrop tap closes ordinary popups only. The quote form keeps what was
  // typed, so it closes by X or Escape only.
  document.querySelectorAll(".pop-scrim").forEach(function (scrim) {
    scrim.addEventListener("click", function (e) {
      if (e.target === scrim && scrim.id !== "pop-quote") closePopup(scrim);
    });
  });

  // Escape peels one layer at a time, top first.
  document.addEventListener("keydown", function (e) {
    if (e.key === "Escape") closeTop();
  });

  // A page can ask for a popup to arrive open (a form that came back with
  // errors stays in its popup, per the failure-stays-in-popup law).
  document.querySelectorAll(".pop-scrim[data-autopen]").forEach(function (el) {
    openPopup(el);
  });

  // ---- Page viewer: in-site pages open as a popup over the base page ---
  // The base page never changes. The viewer leaves a blurred sliver of it
  // on every side and closes from one big X (David's rule, 2026-11-05).
  var pageEl = document.getElementById("pop-page");
  var pageBody = pageEl ? pageEl.querySelector(".pop-page-body") : null;
  function openPage(url) {
    if (!pageEl || !pageBody) { window.location.href = url; return; }
    closeEm();
    pageBody.innerHTML = '<div class="pop-page-loading">Loading...</div>';
    openPopup(pageEl);
    fetch(url, { credentials: "same-origin" }).then(function (r) { return r.text(); }).then(function (html) {
      var doc = new DOMParser().parseFromString(html, "text/html");
      var main = doc.querySelector("main");
      if (!main) { pageBody.innerHTML = '<div class="pop-page-loading">That page could not load here. <a href="' + url + '">Open it directly</a></div>'; return; }
      pageBody.innerHTML = main.innerHTML;
      pageBody.scrollTop = 0;
      registerRails(pageBody);
    }).catch(function () {
      pageBody.innerHTML = '<div class="pop-page-loading">That page could not load here. <a href="' + url + '">Open it directly</a></div>';
    });
  }

  // Forms inside the viewer post in place; the answer renders in the viewer.
  document.addEventListener("submit", function (e) {
    var form = e.target;
    if (!pageBody || !pageBody.contains(form)) return;
    if ((form.getAttribute("method") || "get").toLowerCase() !== "post") return;
    e.preventDefault();
    var body = new URLSearchParams(new FormData(form));
    fetch(form.getAttribute("action") || window.location.pathname, {
      method: "POST", credentials: "same-origin",
      headers: { "Content-Type": "application/x-www-form-urlencoded" },
      body: body.toString()
    }).then(function (r) { return r.text(); }).then(function (html) {
      var doc = new DOMParser().parseFromString(html, "text/html");
      var main = doc.querySelector("main");
      if (main) { pageBody.innerHTML = main.innerHTML; pageBody.scrollTop = 0; registerRails(pageBody); }
    }).catch(function () {});
  });

  // ---- Delegated clicks: one listener covers the shell, the base page,
  // and anything the viewer injects later. -------------------------------
  document.addEventListener("click", function (e) {
    var t = e.target;
    if (!t.closest) return;
    var railBtn = t.closest("[data-rail-prev],[data-rail-next]");
    if (railBtn && railBtn._railPair) {
      var pair = railBtn._railPair;
      var dir = railBtn.hasAttribute("data-rail-prev") ? -1 : 1;
      pair.rail.scrollBy({ left: dir * Math.min(pair.rail.clientWidth * 0.85, 560), behavior: "smooth" });
      return;
    }
    var closer = t.closest("[data-pop-close]");
    if (closer) { closePopup(closer.closest(".pop-scrim")); return; }
    var qvBtn = t.closest("[data-qv]");
    if (qvBtn) { fillQuickView(qvBtn); return; }
    var opener = t.closest("[data-pop-open]");
    if (opener) {
      var name = opener.getAttribute("data-pop-open");
      if (name === "booking") { openBooking(opener.getAttribute("data-offer")); return; }
      openPopup(document.getElementById("pop-" + name));
      return;
    }
    var scroller = t.closest("[data-scroll]");
    if (scroller) {
      var sid = scroller.getAttribute("data-scroll");
      if (pageEl && !pageEl.hidden) closePopup(pageEl);
      if (window.location.pathname === "/" && scrollToSection(sid)) {
        e.preventDefault();
        history.replaceState(null, "", "#" + sid);
      }
      return;
    }
    var a = t.closest("a[href]");
    if (!a) return;
    var href = a.getAttribute("href");
    if (!href || href.charAt(0) !== "/" || href.charAt(1) === "/") return;
    if (a.hasAttribute("download")) return;
    var target = a.getAttribute("target");
    if (target && target !== "_self") return;
    if (href === "/") {
      if (pageEl && !pageEl.hidden) { e.preventDefault(); closePopup(pageEl); window.scrollTo({ top: 0, behavior: "smooth" }); }
      return;
    }
    if (href.indexOf("/#") === 0) {
      if (window.location.pathname === "/") {
        var hid = href.slice(2);
        if (pageEl && !pageEl.hidden) closePopup(pageEl);
        if (scrollToSection(hid)) { e.preventDefault(); history.replaceState(null, "", "#" + hid); }
      }
      return;
    }
    e.preventDefault();
    openPage(href);
  });

  // ---- Booking popup: plan picker with yearly lock-in pricing ---------
  var BK_PLANS = { Bronze: 49, Silver: 59, Gold: 79 };
  var bkEl = document.getElementById("pop-booking");
  var bkOffer = null;
  var bkPlan = "Gold";
  var bkBilling = "yearly";
  function bkMoney(n) { return "$" + n.toLocaleString("en-US"); }
  function bkRender() {
    if (!bkEl) return;
    var m = BK_PLANS[bkPlan];
    var yearFull = m * 12, yearPrice = yearFull;
    if (bkOffer === "veterans") yearPrice = m * 11;
    if (bkOffer === "save25") yearPrice = yearFull - 25;
    var yBig = document.getElementById("bk-yearly-big");
    var yNote = document.getElementById("bk-yearly-note");
    var mBig = document.getElementById("bk-monthly-big");
    var mNote = document.getElementById("bk-monthly-note");
    if (yBig) yBig.innerHTML = (yearPrice !== yearFull ? '<span class="bk-was">' + bkMoney(yearFull) + "</span>" : "") + bkMoney(yearPrice);
    if (yNote) yNote.textContent = bkOffer === "veterans" ? "One payment for 12 months. First month free for veterans." : bkOffer === "save25" ? "One payment for 12 months. $25 off applied." : "One payment for 12 months.";
    if (mBig) mBig.innerHTML = bkOffer === "veterans" ? '<span class="bk-was">' + bkMoney(m) + "</span>$0" : bkOffer === "save25" ? '<span class="bk-was">' + bkMoney(m) + "</span>" + bkMoney(m - 25) : bkMoney(m);
    if (mNote) mNote.textContent = bkOffer === "veterans" ? "First month free, then " + bkMoney(m) + "/mo on the 12-month plan." : bkOffer === "save25" ? "First month with $25 off, then " + bkMoney(m) + "/mo." : bkMoney(m) + "/mo on the 12-month plan.";
    bkEl.querySelectorAll(".bk-plan").forEach(function (row) {
      var radio = row.querySelector('input[type="radio"]');
      row.classList.toggle("on", !!radio && radio.value === bkPlan);
    });
    bkEl.querySelectorAll(".bk-bill").forEach(function (b) { b.classList.toggle("on", b.getAttribute("data-billing") === bkBilling); });
    var pi = document.getElementById("bk-plan-input"); if (pi) pi.value = bkPlan;
    var bi = document.getElementById("bk-billing-input"); if (bi) bi.value = bkBilling;
    var fi = document.getElementById("bk-promo-input"); if (fi) fi.value = bkOffer === "veterans" ? "VETSDAY" : bkOffer === "save25" ? "SAVE25" : "";
    var line = document.getElementById("bk-offer-line");
    if (line) {
      if (bkOffer === "veterans") { line.textContent = "Veterans Day offer applied: first month free on your annual plan. Code VETSDAY."; line.hidden = false; }
      else if (bkOffer === "save25") { line.textContent = "$25 off your first service applied. Code SAVE25."; line.hidden = false; }
      else line.hidden = true;
    }
    var lab = document.getElementById("bk-lab");
    if (lab) lab.textContent = bkOffer === "veterans" ? "Veterans Day // Book Your Plan" : "Book Your Plan";
  }
  function openBooking(offer) {
    if (!bkEl) { openPopup(document.getElementById("pop-quote")); return; }
    bkOffer = offer || null;
    var done = document.getElementById("bk-done");
    var main = bkEl.querySelector(".bk-main");
    if (done) done.hidden = true;
    if (main) main.hidden = false;
    var errs = document.getElementById("bk-errors");
    if (errs) errs.hidden = true;
    bkRender();
    openPopup(bkEl);
  }
  if (bkEl) {
    bkEl.querySelectorAll('input[name="bk-plan-radio"]').forEach(function (r) {
      r.addEventListener("change", function () { bkPlan = r.value; bkRender(); });
    });
    bkEl.querySelectorAll(".bk-bill").forEach(function (b) {
      b.addEventListener("click", function () { bkBilling = b.getAttribute("data-billing"); bkRender(); });
    });
    var bkPest = document.getElementById("bk-pest");
    if (bkPest && !bkPest.value) bkPest.value = "Other";
    var bkForm = document.getElementById("bk-form");
    if (bkForm) bkForm.addEventListener("submit", function (e) {
      e.preventDefault();
      var errs = document.getElementById("bk-errors");
      var btn = document.getElementById("bk-submit");
      if (btn) { btn.disabled = true; btn.textContent = "Booking..."; }
      var body = new URLSearchParams(new FormData(bkForm));
      fetch("/signup", { method: "POST", credentials: "same-origin", headers: { "Content-Type": "application/x-www-form-urlencoded" }, body: body.toString() })
        .then(function (r) { return r.text(); })
        .then(function (html) {
          if (btn) { btn.disabled = false; btn.textContent = "Book My Plan"; }
          if (html.indexOf("notice success") !== -1) {
            var main = bkEl.querySelector(".bk-main");
            var done = document.getElementById("bk-done");
            var dt = document.getElementById("bk-done-text");
            if (dt) dt.textContent = "Your " + bkPlan + " plan (" + (bkBilling === "yearly" ? "pay for the year" : "pay monthly") + ") is booked in" + (bkOffer === "veterans" ? " with the veterans first month free" : bkOffer === "save25" ? " with $25 off your first service" : "") + ". The Patriot team will text you to confirm your first visit and set up payment. Nothing was charged today.";
            if (main) main.hidden = true;
            if (done) done.hidden = false;
          } else {
            var msgs = [];
            var doc = new DOMParser().parseFromString(html, "text/html");
            doc.querySelectorAll(".notice.error li").forEach(function (li) { msgs.push(li.textContent); });
            if (errs) {
              errs.innerHTML = msgs.length ? msgs.map(function (x) { var d = document.createElement("div"); d.textContent = x; return d.outerHTML; }).join("") : "That did not go through. Check your details and try again, or call (509) 471-5767.";
              errs.hidden = false;
            }
          }
        })
        .catch(function () {
          if (btn) { btn.disabled = false; btn.textContent = "Book My Plan"; }
          if (errs) { errs.textContent = "That did not go through. Please call or text (509) 471-5767 and we will book you by phone."; errs.hidden = false; }
        });
    });
    bkRender();
  }

  registerRails(document);
  window.addEventListener("resize", function () { railPairs.forEach(railState); });
  window.addEventListener("load", function () { setTimeout(function () { railPairs.forEach(railState); }, 400); });

  // ---- Promos: holiday window first, standing offer otherwise ----------
  var PROMOS = [
    { id: "veterans", code: "VETSDAY", lab: "Veterans Day", qr: "/assets/img/qr-veterans.png",
      from: new Date(2026, 10, 5, 0, 0, 0), to: new Date(2026, 10, 11, 23, 59, 59),
      title: "Veterans: your first month is free",
      text: "For Veterans Day, veterans get the first month free on any annual plan. Book by November 11 and mention your service. Thank you for yours.",
      line: "Veterans Day offer applied: first month free on an annual plan for veterans." },
    { id: "save25", code: "SAVE25", lab: "Special Offer", qr: "/assets/img/qr-save25.png",
      from: null, to: null,
      title: "$25 off your first service",
      text: "New around here? Take $25 off your first service with Patriot Pest Control. Book now and the code rides along with your request.",
      line: "$25 off your first service applied with code SAVE25." }
  ];
  function promoInWindow(p) {
    if (!p.from) return true;
    var now = new Date();
    return now >= p.from && now <= p.to;
  }
  function activePromo() {
    for (var i = 0; i < PROMOS.length; i++) if (promoInWindow(PROMOS[i])) return PROMOS[i];
    return PROMOS[PROMOS.length - 1];
  }
  function promoByCode(code) {
    for (var i = 0; i < PROMOS.length; i++) if (PROMOS[i].code === code) return PROMOS[i];
    return null;
  }
  var currentPromo = null;
  var quoteEl = document.getElementById("pop-quote");
  function applyPromoToQuote(p) {
    currentPromo = p;
    var line = document.getElementById("pop-promo-line");
    if (line && p) { line.textContent = p.line; line.hidden = false; }
  }
  if (quoteEl) {
    var qform = quoteEl.querySelector("form");
    if (qform) qform.addEventListener("submit", function () {
      if (!currentPromo) return;
      var msg = qform.querySelector('[name="message"]');
      if (msg && msg.value.indexOf(currentPromo.code) === -1) {
        msg.value = (msg.value + "\nPromo code: " + currentPromo.code).trim();
      }
    });
  }

  // Arriving from a QR scan or a promo link: open the booking popup with
  // the offer applied (David's rule: the QR goes to booking, not a form).
  var params = new URLSearchParams(window.location.search);
  var arrivedWithPromo = false;
  if (params.get("book") === "1") {
    var p = promoByCode((params.get("promo") || "").toUpperCase()) || activePromo();
    openBooking(p && p.id === "veterans" ? "veterans" : p && p.id === "save25" ? "save25" : null);
    arrivedWithPromo = true;
  }

  // First-visit promo popup, on a timer, once per visitor per promo.
  var promoEl = document.getElementById("pop-promo");
  if (promoEl && !arrivedWithPromo) {
    var ap = activePromo();
    var seenKey = "pop-promo-seen-" + ap.id;
    var seen = false;
    try { seen = !!window.localStorage.getItem(seenKey); } catch (e) { seen = false; }
    if (!seen) {
      document.getElementById("promo-lab").textContent = ap.lab;
      document.getElementById("promo-title").textContent = ap.title;
      document.getElementById("promo-text").textContent = ap.text;
      document.getElementById("promo-code").textContent = ap.code;
      document.getElementById("promo-qr").setAttribute("src", ap.qr);
      setTimeout(function () {
        if (!stack.length) openPopup(promoEl);
        try { window.localStorage.setItem(seenKey, "1"); } catch (e) {}
      }, 18000);
    }
    var promoBook = document.getElementById("promo-book");
    if (promoBook) promoBook.addEventListener("click", function () {
      closePopup(promoEl);
      openBooking(ap.id === "veterans" ? "veterans" : ap.id === "save25" ? "save25" : null);
      try { window.localStorage.setItem(seenKey, "1"); } catch (e) {}
    });
    promoEl.querySelectorAll("[data-pop-close]").forEach(function (btn) {
      btn.addEventListener("click", function () {
        try { window.localStorage.setItem(seenKey, "1"); } catch (e) {}
      });
    });
  }

  // ---- Settings: gear button, text and tile size, on this device --------
  var settingsEl = document.getElementById("pop-settings");
  var gearBtn = document.getElementById("pop-gear");
  if (gearBtn && settingsEl) gearBtn.addEventListener("click", function () { openPopup(settingsEl); });
  var savedSettings = { fs: "", ts: "" };
  try { savedSettings = JSON.parse(window.localStorage.getItem("pop-settings") || "{}"); } catch (e) { savedSettings = {}; }
  if (typeof savedSettings.fs !== "string") savedSettings.fs = "";
  if (typeof savedSettings.ts !== "string") savedSettings.ts = "";
  function applySettings() {
    document.documentElement.classList.toggle("pop-fs-lg", savedSettings.fs === "lg");
    document.documentElement.classList.toggle("pop-fs-xl", savedSettings.fs === "xl");
    document.body.classList.toggle("pop-tiles-lg", savedSettings.ts === "lg");
    document.querySelectorAll("#set-text button").forEach(function (b) { b.classList.toggle("on", b.getAttribute("data-fs") === savedSettings.fs); });
    document.querySelectorAll("#set-tiles button").forEach(function (b) { b.classList.toggle("on", b.getAttribute("data-ts") === savedSettings.ts); });
    railPairs.forEach(railState);
  }
  document.querySelectorAll("#set-text button").forEach(function (b) {
    b.addEventListener("click", function () {
      savedSettings.fs = b.getAttribute("data-fs");
      try { window.localStorage.setItem("pop-settings", JSON.stringify(savedSettings)); } catch (e) {}
      applySettings();
    });
  });
  document.querySelectorAll("#set-tiles button").forEach(function (b) {
    b.addEventListener("click", function () {
      savedSettings.ts = b.getAttribute("data-ts");
      try { window.localStorage.setItem("pop-settings", JSON.stringify(savedSettings)); } catch (e) {}
      applySettings();
    });
  });
  var topBtn = document.getElementById("set-top");
  if (topBtn) topBtn.addEventListener("click", function () {
    if (settingsEl) closePopup(settingsEl);
    window.scrollTo({ top: 0, behavior: "smooth" });
  });
  applySettings();

  // ---- Emily: the site assistant -----------------------------------------
  var emPanel = document.getElementById("em-panel");
  var emBackdrop = document.getElementById("em-backdrop");
  var emBody = document.getElementById("em-body");
  var emChips = document.getElementById("em-chips");
  var emForm = document.getElementById("em-form");
  var emText = document.getElementById("em-text");
  var emGreeted = false;
  var EM_PESTS = [
    ["bed bugs", "bed-bugs"], ["yellow jackets", "yellow-jackets"], ["box elder bugs", "box-elder-bugs"],
    ["fruit flies", "fruit-flies"], ["pantry pests", "pantry-pests"], ["pack rats", "pack-rats"],
    ["fleas", "fleas-ticks"], ["ticks", "fleas-ticks"], ["ants", "ants"], ["spiders", "spiders"],
    ["rodents", "rodents"], ["mice", "rodents"], ["rats", "rodents"], ["termites", "termites"],
    ["cockroaches", "cockroaches"], ["roaches", "cockroaches"], ["wasps", "wasps"], ["hornets", "hornets"],
    ["bees", "bees"], ["mosquitoes", "mosquitoes"], ["silverfish", "silverfish"], ["crickets", "crickets"],
    ["stink bugs", "stink-bugs"], ["scorpions", "scorpions"], ["bats", "bats"], ["gophers", "gophers"], ["moles", "moles"]
  ];
  var EM_TOWNS = ["spokane valley", "spokane", "cheney", "liberty lake", "airway heights", "medical lake", "deer park", "mead", "coeur d'alene", "post falls", "hayden", "rathdrum", "hermiston", "milton-freewater"];
  function emSay(text, actions) {
    var d = document.createElement("div");
    d.className = "em-msg";
    d.textContent = text;
    if (actions && actions.length) {
      var row = document.createElement("div");
      row.className = "em-actions";
      actions.forEach(function (a) {
        var el;
        if (a.href) { el = document.createElement("a"); el.setAttribute("href", a.href); }
        else { el = document.createElement("button"); el.setAttribute("type", "button"); }
        el.className = "pop-btn" + (a.primary ? " primary" : "");
        el.textContent = a.label;
        if (a.open === "booking") el.addEventListener("click", function () {
          closeEm();
          openBooking(null);
        });
        row.appendChild(el);
      });
      d.appendChild(row);
    }
    emBody.appendChild(d);
    emBody.scrollTop = emBody.scrollHeight;
  }
  function emAnswer(raw) {
    var q = " " + raw.toLowerCase() + " ";
    var bookActions = [
      { label: "Call or Text (509) 471-5767", href: "tel:+15094715767", primary: true },
      { label: "Book Now", open: "booking" }
    ];
    if (/veteran|military|discount|promo|offer|deal|coupon/.test(q)) {
      var ap = activePromo();
      if (ap.id === "veterans") {
        emSay("For Veterans Day, veterans get the first month free on any annual plan. Book by November 11. Want me to get you started?", bookActions);
      } else {
        emSay("Right now new customers get $25 off the first service with code SAVE25. Want me to get you booked with it on?", bookActions);
      }
      return;
    }
    if (/book|appoint|schedul|come out|sign up|get started/.test(q)) {
      emSay("Easy. Emily, our AI receptionist, books real appointments right on a call against the live calendar. Or send the quick form and the team calls you back. Which way do you want to go?", bookActions);
      return;
    }
    if (/price|cost|plan|how much|quote/.test(q)) {
      emSay("Three plans, no hidden fees. Bronze is $49 a month for exterior protection, Silver is $59 for interior plus exterior, and Gold is $79 for priority interior and exterior. Every plan rides on a 90-day guarantee.", [
        { label: "See the Plans", href: "/prices", primary: true },
        { label: "Book Now", open: "booking" }
      ]);
      return;
    }
    if (/guarantee|warranty/.test(q)) {
      emSay("Every mission is backed by our 90-day guarantee. If the pests come back inside that window, so do we.", bookActions);
      return;
    }
    if (/human|person|real person|someone|talk to/.test(q)) {
      emSay("Of course. Skyler and Megan's line is (509) 471-5767. Call or text any time and a person picks it up from there.", [
        { label: "Call or Text (509) 471-5767", href: "tel:+15094715767", primary: true }
      ]);
      return;
    }
    for (var i = 0; i < EM_PESTS.length; i++) {
      if (q.indexOf(EM_PESTS[i][0]) !== -1) {
        var name = EM_PESTS[i][0];
        emSay("Yes, we handle " + name + ". Tap through for the full briefing on what you are up against and how we treat it, or book and we will come look.", [
          { label: "See " + name, href: "/pest/" + EM_PESTS[i][1], primary: true },
          { label: "Book Now", open: "booking" }
        ]);
        return;
      }
    }
    for (var t = 0; t < EM_TOWNS.length; t++) {
      if (q.indexOf(EM_TOWNS[t]) !== -1) {
        emSay("Yes, we run routes there. Same-day service is available when it can't wait, and every job carries the 90-day guarantee.", bookActions);
        return;
      }
    }
    if (/area|where|serve|location|city|cities|town/.test(q)) {
      emSay("We cover Eastern Washington and North Idaho: Spokane, Spokane Valley, Cheney, Liberty Lake, Airway Heights, Medical Lake, Deer Park, Mead, Coeur d'Alene, Post Falls, Hayden, and Rathdrum.", bookActions);
      return;
    }
    if (/safe|kid|pet|child|dog|cat/.test(q)) {
      emSay("Good question for the team. Treatments are applied by our own techs, and Skyler's crew will walk you through exactly what is used and any prep for kids and pets before work starts. A quick call settles it.", bookActions);
      return;
    }
    if (/hour|open|when/.test(q)) {
      emSay("Fastest is a call or text at (509) 471-5767. Same-day service is available when it can't wait, and the form gets an answer within one business day, usually much faster.", bookActions);
      return;
    }
    if (/hi|hello|hey/.test(q)) {
      emSay("Hey there. I can walk you through prices, pests, areas, or get you booked. What are we hunting today?");
      return;
    }
    emSay("I can help with prices and plans, what pests we treat, the areas we serve, deals, and booking an appointment. Try one of these, or ask me straight.", null);
  }
  function emChip(label, question) {
    var b = document.createElement("button");
    b.type = "button";
    b.className = "em-chip";
    b.textContent = label;
    b.addEventListener("click", function () {
      emEcho(question);
      emAnswer(question);
    });
    emChips.appendChild(b);
  }
  function emEcho(text) {
    var d = document.createElement("div");
    d.className = "em-msg me";
    d.textContent = text;
    emBody.appendChild(d);
    emBody.scrollTop = emBody.scrollHeight;
  }
  function openEm() {
    if (!emPanel) return;
    if (emBackdrop) emBackdrop.hidden = false;
    emPanel.hidden = false;
    if (!emGreeted) {
      emGreeted = true;
      emSay("Hi, I'm Emily, the Patriot Pest Control assistant. I know our prices, our pests, and our towns, and I can get you booked. What's going on at your place?");
      emChip("Prices & plans", "What do the plans cost?");
      emChip("Book an appointment", "I want to book an appointment");
      emChip("What pests?", "What pests do you treat?");
      emChip("Areas served", "What areas do you serve?");
      emChip("Any deals?", "Do you have any deals?");
    }
  }
  function closeEm() { if (emPanel) emPanel.hidden = true; if (emBackdrop) emBackdrop.hidden = true; }
  if (emBackdrop) emBackdrop.addEventListener("click", closeEm);
  var askBtn = document.getElementById("pop-ask");
  if (askBtn) askBtn.addEventListener("click", openEm);
  var emClose = document.getElementById("em-close");
  if (emClose) emClose.addEventListener("click", closeEm);
  if (emForm) emForm.addEventListener("submit", function (e) {
    e.preventDefault();
    var v = emText.value.trim();
    if (!v) return;
    emText.value = "";
    emEcho(v);
    emAnswer(v);
  });
})();
