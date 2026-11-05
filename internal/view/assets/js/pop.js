// Pop UI behavior for pop.patriotpest.pro: tile rail arrows and the popup
// system (quick view layer + quote form layer). Rules from the pop_ui skill:
// popups are centered, opening never focuses a text field, Escape peels one
// layer at a time, the quote form does not close on a backdrop tap.
(function () {
  "use strict";

  // ---- Rails ------------------------------------------------------------
  document.querySelectorAll("[data-rail-prev],[data-rail-next]").forEach(function (btn) {
    btn.addEventListener("click", function () {
      var id = btn.getAttribute("data-rail-prev") || btn.getAttribute("data-rail-next");
      var rail = document.getElementById(id);
      if (!rail) return;
      var dir = btn.hasAttribute("data-rail-prev") ? -1 : 1;
      rail.scrollBy({ left: dir * Math.min(rail.clientWidth * 0.85, 560), behavior: "smooth" });
    });
  });

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
  document.querySelectorAll("[data-scroll]").forEach(function (a) {
    a.addEventListener("click", function (e) {
      var id = a.getAttribute("data-scroll");
      if (window.location.pathname === "/" && scrollToSection(id)) {
        e.preventDefault();
        history.replaceState(null, "", "#" + id);
      }
    });
  });
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
  document.querySelectorAll("[data-qv]").forEach(function (btn) {
    btn.addEventListener("click", function () {
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
    });
  });

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
        window.location.href = currentSrc.getAttribute("data-more");
      }
    });
  }

  // Named popups (quote form).
  document.querySelectorAll("[data-pop-open]").forEach(function (btn) {
    btn.addEventListener("click", function () {
      openPopup(document.getElementById("pop-" + btn.getAttribute("data-pop-open")));
    });
  });

  // Close buttons peel their own popup.
  document.querySelectorAll("[data-pop-close]").forEach(function (btn) {
    btn.addEventListener("click", function () {
      closePopup(btn.closest(".pop-scrim"));
    });
  });

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

  // ---- Rail arrows: fade at the ends, gone when the row fits ------------
  var railPairs = {};
  document.querySelectorAll("[data-rail-prev],[data-rail-next]").forEach(function (btn) {
    var id = btn.getAttribute("data-rail-prev") || btn.getAttribute("data-rail-next");
    if (!railPairs[id]) railPairs[id] = { prev: null, next: null, box: btn.parentElement };
    if (btn.hasAttribute("data-rail-prev")) railPairs[id].prev = btn; else railPairs[id].next = btn;
  });
  function railState(id) {
    var rail = document.getElementById(id);
    var pair = railPairs[id];
    if (!rail || !pair) return;
    var max = rail.scrollWidth - rail.clientWidth;
    if (max <= 4) { pair.box.classList.add("gone"); return; }
    pair.box.classList.remove("gone");
    if (pair.prev) pair.prev.classList.toggle("off", rail.scrollLeft <= 4);
    if (pair.next) pair.next.classList.toggle("off", rail.scrollLeft >= max - 4);
  }
  Object.keys(railPairs).forEach(function (id) {
    var rail = document.getElementById(id);
    if (rail) rail.addEventListener("scroll", function () { railState(id); }, { passive: true });
    railState(id);
  });
  window.addEventListener("resize", function () { Object.keys(railPairs).forEach(railState); });
  window.addEventListener("load", function () { setTimeout(function () { Object.keys(railPairs).forEach(railState); }, 400); });

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

  // Arriving from a QR scan or a promo link: open the quote with the code on.
  var params = new URLSearchParams(window.location.search);
  var arrivedWithPromo = false;
  if (params.get("book") === "1" && quoteEl) {
    var p = promoByCode((params.get("promo") || "").toUpperCase()) || activePromo();
    applyPromoToQuote(p);
    openPopup(quoteEl);
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
      applyPromoToQuote(ap);
      closePopup(promoEl);
      if (quoteEl) openPopup(quoteEl);
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
    Object.keys(railPairs).forEach(railState);
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
        if (a.open === "quote") el.addEventListener("click", function () {
          closeEm();
          if (quoteEl) openPopup(quoteEl);
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
      { label: "Book Now", open: "quote" }
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
        { label: "Book Now", open: "quote" }
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
          { label: "Book Now", open: "quote" }
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
  function closeEm() { if (emPanel) emPanel.hidden = true; }
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
