package view

// pageSmsOptIn — the public SMS opt-in proof page (/sms-opt-in). Twilio's
// A2P 10DLC registry rejected Patriot's campaign (error 30909) because the
// reviewers could not verify how customers consent to texts: the opt-in is
// verbal (by phone and at service-plan start) and no proof was public.
// This page is that proof: the program, the three opt-in paths, the exact
// verbal script, and a working web opt-in form with the consent checkbox.
// Data: .Success (string), .Errors ([]string), .OldName/.OldPhone (prior
// input for re-population), .Csrf (template.HTML hidden input),
// .AnalyticsEvent ("generate_lead" fires the gtag conversion event).
const pageSmsOptIn = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">TEXTING // OPT-IN</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Text Message Opt-In</h1>
    <p class="lead">Patriot Pest Control texts customers who ask us to. This page explains the program and exactly how consent is given.</p>
    <p style="color:var(--khaki);line-height:1.7">Patriot Pest Control serves Eastern Washington and North Idaho.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin-bottom:.8rem">What we send</h2>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem">We send appointment confirmations and reminders, service updates about your visits, and occasional requests for a review. Review request texts are sent only to customers who opted in through the web form or in writing on the signed form.</p>
    <p style="color:var(--cream);line-height:1.7;margin-top:.6rem"><b>Nothing else.</b></p>

    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin:2.2rem 0 .8rem">How customers opt in</h2>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem">There are three ways to join the texting list:</p>
    <ol style="color:var(--khaki);line-height:1.8;margin:.6rem 0 0 1.2rem;max-width:46rem">
      <li><b style="color:var(--cream)">By phone:</b> call <a href="tel:+15094715767" style="color:var(--orange)">(509) 471-5767</a> to book service or start a service plan, and we ask for your OK out loud before we ever text you.</li>
      <li><b style="color:var(--cream)">When starting a service plan</b> with a team member, in person or by phone.</li>
      <li><b style="color:var(--cream)">On this page,</b> using the form below.</li>
    </ol>
  </div>
</section>

<section class="block">
  <div class="wrap">
    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin-bottom:.8rem">The exact words we use on the phone</h2>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem">Before a number is added to the texting list, the customer hears these words and answers out loud:</p>
    <div class="panel" style="border-left:4px solid var(--orange);max-width:46rem;margin-top:1rem">
      <p style="color:var(--cream);font-size:1.05rem;line-height:1.7;margin:0">Patriot Pest Control can text you at this number about your appointments and your service. Message frequency varies, and message and data rates may apply. Reply STOP to end them at any time, or HELP for help. Our Privacy Policy is at patriotpest.pro/privacy-policy and our Terms are at patriotpest.pro/terms-of-use. Consent is not a condition of purchase, and your mobile information is never shared or sold. Is it OK to text you?</p>
    </div>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem;margin-top:1rem">We only add a number to the texting list when the customer says yes.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin-bottom:.8rem">The paper service plan form</h2>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem">When a customer starts a service plan in person, consent can be given on our paper service plan form instead. The form carries a field for the customer mobile number, a description of the texting service, the message frequency, a note that message and data rates may apply, HELP and STOP instructions, links to our Privacy Policy and Terms, a consent checkbox, a signature line, and a date field. Signed forms are kept on file.</p>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem;margin-top:1rem">Below is a copy of the consent wording and fields exactly as they appear on that form, so anyone can see the form itself:</p>
    <div class="panel" style="max-width:46rem;margin-top:1rem">
      <p style="color:var(--cream);font-weight:700;margin:0 0 .8rem">Patriot Pest Control service plan form: text message consent</p>
      <p style="color:var(--khaki);line-height:1.7;margin:0 0 .6rem">Name: ______________________________</p>
      <p style="color:var(--khaki);line-height:1.7;margin:0 0 .6rem">Mobile number: ______________________________</p>
      <p style="color:var(--khaki);line-height:1.7;margin:0 0 .8rem">Patriot Pest Control sends text messages about your appointments and your service, including appointment confirmations and reminders and service updates about your visits. Message frequency varies. Message and data rates may apply. Reply HELP for help, STOP to end. Privacy Policy: patriotpest.pro/privacy-policy. Terms: patriotpest.pro/terms-of-use.</p>
      <p style="color:var(--cream);line-height:1.7;margin:0 0 .8rem"><span style="display:inline-block;width:1.05rem;height:1.05rem;border:2px solid var(--cream);border-radius:.2rem;vertical-align:-.15rem;margin-right:.5rem"></span>Yes, Patriot Pest Control may text me at the mobile number above. Consent is not a condition of purchase, and my mobile information is never shared or sold.</p>
      <p style="color:var(--khaki);line-height:1.7;margin:0 0 .6rem">Signature: ______________________________</p>
      <p style="color:var(--khaki);line-height:1.7;margin:0">Date: ______________________________</p>
    </div>

    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin:2.2rem 0 .8rem">The confirmation text</h2>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem">Every new opt-in receives this confirmation text:</p>
    <div class="panel" style="border-left:4px solid var(--orange);max-width:46rem;margin-top:1rem">
      <p style="color:var(--cream);font-size:1.05rem;line-height:1.7;margin:0">Patriot Pest Control: you are signed up for service texts. Msg frequency varies, msg and data rates may apply. Reply HELP for help, STOP to end.</p>
    </div>

    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin:2.2rem 0 .8rem">How consent is recorded</h2>
    <p style="color:var(--khaki);line-height:1.7;max-width:46rem">Every opt-in is logged with the date, the time, the method (web form, phone, or signed form), and the exact wording the customer agreed to. Phone calls may be recorded.</p>
  </div>
</section>

<section class="block">
  <div class="wrap">
    <h2 style="font-family:var(--display);color:var(--cream);font-size:1.6rem;margin-bottom:.8rem">Opt in on this page</h2>
    <div class="form-panel" style="max-width:34rem">
      {{if .Success}}<div class="notice success">{{.Success}}</div>
        {{if eq .AnalyticsEvent "generate_lead"}}
        <script>
          if (window.gtag) {
            gtag('event', 'generate_lead', {
              'event_category': 'conversion',
              'event_label': 'SMS Opt-In Form'
            });
          }
          if (window.fbq) {
            fbq('track', 'Lead', { 'content_name': 'SMS Opt-In Form' });
          }
        </script>
        {{end}}
      {{end}}
      {{if .Errors}}<div class="notice error">
        <strong>Please fix the following:</strong>
        <ul>{{range .Errors}}<li>{{.}}</li>{{end}}</ul>
      </div>{{end}}

      <form method="post" action="/sms-opt-in" class="form-stack">
        {{raw .Csrf}}
        <div class="field">
          <label for="name">Name</label>
          <input type="text" id="name" name="name" required autocomplete="name" value="{{.OldName}}">
        </div>
        <div class="field">
          <label for="phone">Mobile number</label>
          <input type="tel" id="phone" name="phone" required autocomplete="tel" value="{{.OldPhone}}">
        </div>
        <div class="field">
          <label for="sms_consent" style="display:flex;gap:.6rem;align-items:flex-start;font-weight:400;line-height:1.5">
            <input type="checkbox" id="sms_consent" name="sms_consent" value="yes" style="margin-top:.2rem;width:1.15rem;height:1.15rem;flex:none">
            <span>Yes, Patriot Pest Control may text me about my appointments and service. Message frequency varies. Message and data rates may apply. Reply STOP to end, HELP for help. Consent is not a condition of purchase.</span>
          </label>
        </div>
        <div class="form-actions">
          <button type="submit" class="btn btn-primary">Sign me up</button>
        </div>
      </form>
      <p style="color:var(--khaki);font-size:.85rem;line-height:1.6;margin-top:1.2rem">Message frequency varies. Message and data rates may apply. Reply STOP to end messages at any time, HELP for help. Consent is not a condition of purchase. We never share or sell your mobile number or your consent to anyone else. See our <a href="/privacy-policy" style="color:var(--orange)">Privacy Policy</a> and <a href="/terms-of-use" style="color:var(--orange)">Terms of Use</a>.</p>
    </div>
  </div>
</section>`

func init() {
	RegisterPageTemplate("sms-opt-in", pageSmsOptIn)
}
