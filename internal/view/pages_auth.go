package view

// Page body templates for the passwordless login flow — Go ports of
// templates/auth/login.php and templates/auth/verify.php. Rendered through the
// shared marketing layout (dark shell), same as every public page.
// Data keys: Csrf (template.HTML hidden input), FlashError string, SentTo string.

func init() {
	RegisterPageTemplate("auth-login", pageAuthLogin)
	RegisterPageTemplate("auth-verify", pageAuthVerify)
}

const pageAuthLogin = `
<div class="authx">

  <!-- ===== left / brand panel ===== -->
  <div class="authx-brand">
    <div class="authx-radar" aria-hidden="true"></div>

    <div class="authx-topline">
      <span><span class="dot"></span>PATRIOT PEST CONTROL</span>
      <span id="authx-clock" aria-hidden="true">--:--</span>
    </div>

    <div>
      <h1 class="authx-headline">Welcome <em>back.</em></h1>
      <p class="authx-lede">Your appointments, your technician, your plan. Everything about your pest-free home, right here.</p>
      <p class="authx-new"><b>New around here?</b> Welcome aboard. Enter your details and we'll get you set up in seconds.</p>

      <div class="authx-photo" aria-hidden="true">
        <img id="authx-img" src="{{ asset "img/pests/ants.jpg" }}" alt="">
        <div class="scan"></div>
        <span class="tag" id="authx-tag">ON WATCH // ANTS</span>
      </div>

      <div class="authx-trust">
        <span>Licensed &amp; Insured</span><span>90-Day Warranty</span><span>Family &amp; Pet Safe</span>
      </div>
    </div>

    <div class="authx-call">
      <span class="lbl">Prefer to talk? Your local line</span><br>
      <a href="{{ .PhoneHref }}">☎ {{ .PhoneDisplay }} <small>{{ .PhoneLabel }}</small></a>
    </div>
  </div>

  <!-- ===== right / form panel ===== -->
  <div class="authx-form">
    <div class="authx-card2">
      <span class="step">ACCOUNT ACCESS</span>
      <h1>Sign in to your account</h1>
      <p class="sub">Enter the email, phone number, or account number on your account and we'll send you a secure code. No password to remember.</p>

      {{ if .FlashError }}<div class="notice error">{{ .FlashError }}</div>{{ end }}

      <form method="post" action="/login" novalidate>
        {{ .Csrf }}
        <div class="authx-field">
          <label for="identifier">Email, phone, or account number</label>
          <input type="text" id="identifier" name="identifier" autocomplete="username" required autofocus
                 placeholder="you@example.com  ·  (509) 555-0101  ·  1001">
          <div class="hint">Use whichever you have on file. We'll recognize you and email a one-time code.</div>
        </div>
        <button type="submit" class="authx-btn">Send My Secure Code ▸</button>
      </form>

      <p class="authx-foot">Having trouble? <a href="/help">Get help signing in</a></p>
    </div>
  </div>
</div>

<script>
(function () {
  // Live clock in the brand panel.
  var clock = document.getElementById('authx-clock');
  function tick() {
    var d = new Date();
    clock.textContent = d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }
  if (clock) { tick(); setInterval(tick, 30000); }

  // Rotate the "on watch" threat photo with a soft crossfade.
  var shots = [
    ['ants.jpg', 'ANTS'], ['spiders.jpg', 'SPIDERS'],
    ['rodents.jpg', 'RODENTS'], ['scorpions.jpg', 'SCORPIONS']
  ];
  var img = document.getElementById('authx-img');
  var tag = document.getElementById('authx-tag');
  if (img && tag) {
    var i = 0;
    setInterval(function () {
      img.style.opacity = '0';
      setTimeout(function () {
        i = (i + 1) % shots.length;
        img.src = '/assets/img/pests/' + shots[i][0];
        tag.textContent = 'ON WATCH // ' + shots[i][1];
        img.onload = function () { img.style.opacity = '1'; };
      }, 350);
    }, 4200);
  }
})();
</script>
`

const pageAuthVerify = `
<div class="authx">

  <!-- ===== left / brand panel ===== -->
  <div class="authx-brand">
    <div class="authx-radar" aria-hidden="true"></div>

    <div class="authx-topline">
      <span><span class="dot"></span>PATRIOT PEST CONTROL</span>
      <span id="authx-clock" aria-hidden="true">--:--</span>
    </div>

    <div>
      <h1 class="authx-headline">Check your <em>inbox.</em></h1>
      <p class="authx-lede">We've sent a secure 6-digit code to the contact on file. Enter it here and you're in.</p>
      <p class="authx-new">Didn't get it? Give it a minute, check spam, or <a href="/login" style="color:var(--orange)">request a new code</a>.</p>

      <div class="authx-photo" aria-hidden="true">
        <img src="{{ asset "img/pests/spiders.jpg" }}" alt="">
        <div class="scan"></div>
        <span class="tag">SECURE CHANNEL // OPEN</span>
      </div>

      <div class="authx-trust">
        <span>Licensed &amp; Insured</span><span>90-Day Warranty</span><span>Family &amp; Pet Safe</span>
      </div>
    </div>

    <div class="authx-call">
      <span class="lbl">Prefer to talk? Your local line</span><br>
      <a href="{{ .PhoneHref }}">☎ {{ .PhoneDisplay }} <small>{{ .PhoneLabel }}</small></a>
    </div>
  </div>

  <!-- ===== right / code panel ===== -->
  <div class="authx-form">
    <div class="authx-card2">
      <span class="step">ALMOST THERE</span>
      <h1>Enter your code</h1>
      <p class="sub">
        {{ if .SentTo }}We sent a 6-digit code to <span class="sent">{{ .SentTo }}</span>. It expires in a few minutes and works once.{{ else }}Type the 6-digit code we sent you. It expires in a few minutes and works once.{{ end }}
      </p>

      {{ if .FlashError }}<div class="notice error">{{ .FlashError }}</div>{{ end }}

      <form method="post" action="/login/verify" novalidate>
        {{ .Csrf }}
        <div class="authx-field">
          <label for="code">6-digit code</label>
          <input type="text" id="code" name="code" class="authx-code" inputmode="numeric" autocomplete="one-time-code"
                 pattern="[0-9]{6}" maxlength="6" required autofocus placeholder="••••••">
          <div class="hint">Can't find it? Check your spam folder.</div>
        </div>
        <button type="submit" class="authx-btn">Verify &amp; Sign In ▸</button>
      </form>

      <p class="authx-foot"><a href="/login">◂ Use a different email or phone</a></p>
    </div>
  </div>
</div>

<script>
(function () {
  var clock = document.getElementById('authx-clock');
  function tick() {
    var d = new Date();
    clock.textContent = d.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  }
  if (clock) { tick(); setInterval(tick, 30000); }
})();
</script>
`
