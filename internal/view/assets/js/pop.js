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
})();
