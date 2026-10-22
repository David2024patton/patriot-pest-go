package view

// pageContact — templates/pages/contact.php. Contact + free-quote form.
// Data: .Success (string), .Errors ([]string), .OldName/.OldEmail/.OldPhone/
// .OldMessage/.OldZip/.OldPest (prior input for re-population), .Csrf
// (template.HTML hidden input), .PestOptions (template.HTML picker options),
// .AnalyticsEvent ("generate_lead" fires the gtag conversion event).
const pageContact = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">COMMS // CONTACT</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Get your <em>free quote.</em></h1>
    <p class="lead">Call the line or send the form. We respond within one business day (usually much faster). Same-day service available across all three states.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="split" style="align-items:start">
      <div>
        <div class="panel">
          <h2 style="font-family:var(--display);color:var(--cream);font-size:1.17rem">Direct Lines</h2>
          <dl class="kv" style="margin-top:.8rem">
            <dt>WA · ID · OR</dt><dd><a href="tel:+15094715767" style="color:var(--orange)">(509) 471-5767</a></dd>
            <dt>Email</dt><dd><a href="mailto:info@patriotpest.pro" style="color:var(--orange)">info@patriotpest.pro</a></dd>
            <dt>HQ</dt><dd>Deer Park, WA 99006</dd>
            <dt>Hours</dt><dd>Mon–Sat 7a–7p · Sun closed · 24/7 online scheduling</dd>
          </dl>
        </div>
        <div class="panel">
          <h2 style="font-family:var(--display);color:var(--cream);font-size:1.17rem">Why Patriot?</h2>
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
            if (window.fbq) {
              fbq('track', 'Lead', { 'content_name': 'Contact Form' });
            }
          </script>
          {{end}}
        {{end}}
        {{if .Errors}}<div class="notice error">
          <strong>Please fix the following:</strong>
          <ul>{{range .Errors}}<li>{{.}}</li>{{end}}</ul>
        </div>{{end}}

        <form method="post" action="/contact" class="form-stack">
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
              <input type="tel" id="phone" name="phone" required autocomplete="tel" value="{{.OldPhone}}">
            </div>
          </div>
          <div class="form-row">
            <div class="field">
              <label for="pest">What pest is bugging you?</label>
              <select id="pest" name="pest">{{raw .PestOptions}}</select>
            </div>
            <div class="field">
              <label for="zip">ZIP</label>
              <input type="text" id="zip" name="zip" maxlength="10" autocomplete="postal-code" value="{{.OldZip}}">
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
