package view

// pageFreeQuote — the dedicated free quote page (/free-quote), added
// 2026-11-12 at David's order after Skyler asked whether the site had
// one. Until now every "Get a Free Quote" button on the site landed on
// the contact form; this page is the real quote request: what pest,
// where the home is, and how to reach the customer, posted to the same
// lead capture as the contact form (CRM contact + staff alert).
// Data: .Success (string), .Errors ([]string), .OldName/.OldEmail/
// .OldPhone/.OldAddress/.OldCity/.OldZip (prior input), .PestOptions
// (template.HTML select options), .Csrf (template.HTML hidden input),
// .AnalyticsEvent ("generate_lead" fires the gtag conversion event).
const pageFreeQuote = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">FREE QUOTE // NO OBLIGATION</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Get a Free Quote</h1>
    <p class="lead">Tell us what is bugging you and where you are. We call you back with a straight price for your home. No obligation and no pressure.</p>
    <p style="color:var(--khaki);line-height:1.7">Patriot Pest Control serves Eastern Washington and North Idaho.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin-bottom:.8rem">How it works</h2>
    <ol style="color:var(--khaki);line-height:1.8;margin:.6rem 0 0 1.2rem;max-width:46rem">
      <li><b style="color:var(--cream)">Send the form below.</b> It takes about a minute.</li>
      <li><b style="color:var(--cream)">We call you back</b> within one business day, usually much faster, with the exact price for your home.</li>
      <li><b style="color:var(--cream)">If you like the price,</b> we book your first visit. If not, no hard feelings.</li>
    </ol>

    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin:2.2rem 0 .8rem">The plans at a glance</h2>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem">Bronze is $49 per month for exterior-only protection. Silver is $59 per month for interior and exterior. Gold is $79 per month for priority interior and exterior. Every plan includes free re-treatments between visits and a 90-day warranty. Your quote confirms the exact price for your home before you commit to anything.</p>
    <p style="margin-top:1rem"><a class="btn btn-ghost" href="/prices">See the Plans ▸</a></p>
  </div>
</section>

<section class="block">
  <div class="wrap">
    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin-bottom:.8rem">Request your free quote</h2>
    <div class="form-panel" style="max-width:34rem">
      {{if .Success}}<div class="notice success">{{.Success}}</div>
        {{if eq .AnalyticsEvent "generate_lead"}}
        <script>
          if (window.gtag) {
            gtag('event', 'generate_lead', {
              'event_category': 'conversion',
              'event_label': 'Free Quote Form'
            });
          }
          if (window.fbq) {
            fbq('track', 'Lead', { 'content_name': 'Free Quote Form' });
          }
        </script>
        {{end}}
      {{end}}
      {{if .Errors}}<div class="notice error">
        <strong>Please fix the following:</strong>
        <ul>{{range .Errors}}<li>{{.}}</li>{{end}}</ul>
      </div>{{end}}

      <form method="post" action="/free-quote" class="form-stack">
        {{raw .Csrf}}
        <div class="field">
          <label for="fq-name">Name</label>
          <input type="text" id="fq-name" name="name" required autocomplete="name" value="{{.OldName}}">
        </div>
        <div class="form-row">
          <div class="field">
            <label for="fq-phone">Phone</label>
            <input type="tel" id="fq-phone" name="phone" required autocomplete="tel" value="{{.OldPhone}}">
          </div>
          <div class="field">
            <label for="fq-email">Email (optional)</label>
            <input type="email" id="fq-email" name="email" autocomplete="email" value="{{.OldEmail}}">
          </div>
        </div>
        <div class="field">
          <label for="fq-address">Service address (optional)</label>
          <input type="text" id="fq-address" name="address" autocomplete="street-address" value="{{.OldAddress}}">
        </div>
        <div class="form-row">
          <div class="field">
            <label for="fq-city">City (optional)</label>
            <input type="text" id="fq-city" name="city" autocomplete="address-level2" value="{{.OldCity}}">
          </div>
          <div class="field">
            <label for="fq-zip">ZIP</label>
            <input type="text" id="fq-zip" name="zip" maxlength="10" autocomplete="postal-code" value="{{.OldZip}}">
          </div>
        </div>
        <div class="field">
          <label for="fq-pest">What pest is bugging you?</label>
          <select id="fq-pest" name="pest">{{raw .PestOptions}}</select>
        </div>
        <div class="field">
          <label for="fq-message">Anything we should know? (optional)</label>
          <textarea id="fq-message" name="message" placeholder="Where you are seeing them, how bad it is, and when you would like us to come out...">{{.OldMessage}}</textarea>
        </div>
        <div class="field">
          <label for="fq-sms-consent" style="display:flex;gap:.6rem;align-items:flex-start;font-weight:400;line-height:1.5">
            <input type="checkbox" id="fq-sms-consent" name="sms_consent" value="yes" style="margin-top:.2rem;width:1.15rem;height:1.15rem;flex:none">
            <span>Yes, Patriot Pest Control may text me about my quote, my appointments, and my service. Message frequency varies. Message and data rates may apply. Reply STOP to end, HELP for help. Consent is not a condition of purchase.</span>
          </label>
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary">Get My Free Quote ▸</button>
        </div>
      </form>
      <p style="color:var(--khaki);font-size:.85rem;line-height:1.6;margin-top:1.2rem">We call you back about your quote whether or not you check the texting box. If you check it, we may also text you. Message frequency varies. Message and data rates may apply. Reply STOP to end messages at any time, HELP for help. Consent is not a condition of purchase. We never share or sell your mobile number or your consent to anyone else. See our <a href="/privacy-policy" style="color:var(--orange)">Privacy Policy</a> and <a href="/terms-of-use" style="color:var(--orange)">Terms of Use</a>.</p>
      <p style="color:var(--khaki);line-height:1.7;margin-top:1rem">Rather talk now? Call <a href="tel:+15094715767" style="color:var(--orange)">(509) 471-5767</a>.</p>
    </div>
  </div>
</section>`

func init() {
	RegisterPageTemplate("free-quote", pageFreeQuote)
}
