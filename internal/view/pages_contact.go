package view

// pageContact — templates/pages/contact.php. Contact + free-quote form.
// Data: .Success (string), .Errors ([]string), .OldName/.OldEmail/.OldPhone/
// .OldMessage (prior input for re-population), .Csrf (template.HTML hidden
// input), .AnalyticsEvent ("generate_lead" fires the gtag conversion event).
const pageContact = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">COMMS // CONTACT</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Get your <em>free quote.</em></h1>
    <p class="lead">Call the line or send the form. We respond within one business day (usually much faster). Same-day service available across all four states.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="split" style="align-items:start">
      <div>
        <div class="panel">
          <h3 style="font-family:var(--display);color:var(--cream)">Direct Lines</h3>
          <dl class="kv" style="margin-top:.8rem">
            <dt>WA · ID · OR</dt><dd><a href="tel:+15094715767" style="color:var(--orange)">(509) 471-5767</a></dd>
            <dt>Arizona</dt><dd><a href="tel:+16027558414" style="color:var(--orange)">(602) 755-8414</a></dd>
            <dt>Email</dt><dd><a href="mailto:info@patriotpest.pro" style="color:var(--orange)">info@patriotpest.pro</a></dd>
            <dt>HQ</dt><dd>Deer Park, WA 99006</dd>
            <dt>Hours</dt><dd>Mon–Fri 9a–5p · Sat–Sun 10a–4p · 24/7 line</dd>
          </dl>
        </div>
        <div class="panel">
          <h3 style="font-family:var(--display);color:var(--cream)">Why Patriot?</h3>
          <ul style="margin:.6rem 0 0 1.1rem;color:var(--khaki);line-height:1.8">
            <li>Same-day service available</li>
            <li>Free quotes, no hidden fees</li>
            <li>90-day warranty, 100% guaranteed</li>
            <li>Eco-friendly, family &amp; pet safe</li>
          </ul>
        </div>
      </div>

      <div class="form-panel">
        {{if .Success}}<div class="notice success">{{.Success}}</div>
          {{if eq .AnalyticsEvent "generate_lead"}}
          <script>
            if (window.gtag) {
              gtag('event', 'generate_lead', {
                'event_category': 'conversion',
                'event_label': 'Contact Form'
              });
            }
          </script>
          {{end}}
        {{end}}
        {{if .Errors}}<div class="notice error">
          <strong>Please fix the following:</strong>
          <ul>{{range .Errors}}<li>{{.}}</li>{{end}}</ul>
        </div>{{end}}

        <form method="post" action="/contact" class="form-stack" novalidate>
          {{raw .Csrf}}
          <div class="field">
            <label for="name">Name</label>
            <input type="text" id="name" name="name" required autocomplete="name" value="{{.OldName}}">
          </div>
          <div class="form-row">
            <div class="field">
              <label for="email">Email</label>
              <input type="email" id="email" name="email" required autocomplete="email" value="{{.OldEmail}}">
            </div>
            <div class="field">
              <label for="phone">Phone</label>
              <input type="tel" id="phone" name="phone" autocomplete="tel" value="{{.OldPhone}}">
            </div>
          </div>
          <div class="field">
            <label for="message">How can we help?</label>
            <textarea id="message" name="message" required placeholder="Tell us about the pest, your location, and preferred timing…">{{.OldMessage}}</textarea>
          </div>` + smsConsent + `
          <div class="form-actions">
            <button type="submit" class="btn btn-primary">Send Message ▸</button>
          </div>
        </form>
      </div>
    </div>
  </div>
</section>`
