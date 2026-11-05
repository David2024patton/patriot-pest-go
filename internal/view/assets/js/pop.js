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
  var qv = document.getElementById("pop-qv");
  document.querySelectorAll("[data-qv]").forEach(function (btn) {
    btn.addEventListener("click", function () {
      var src = document.getElementById(btn.getAttribute("data-qv"));
      if (!src || !qv) return;
      qv.querySelector(".pop-qv-body").innerHTML = src.innerHTML;
      var more = qv.querySelector(".pop-qv-more");
      var href = src.getAttribute("data-more");
      if (href) {
        more.setAttribute("href", href);
        more.style.display = "";
      } else {
        more.style.display = "none";
      }
      openPopup(qv);
    });
  });

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
})();
