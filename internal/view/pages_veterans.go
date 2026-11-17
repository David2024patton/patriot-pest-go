package view

// pageVeterans — the public Veterans Day offer page. Skyler's offer
// (2026-11-05): veterans get their first month free on any annual plan,
// book by Veterans Day, Nov 11, 2026. This page is the confirmation link
// the veterans-deals lists (Military.com, Veteran.com, WDVA) require.
const pageVeterans = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">VETERANS DAY // THANK YOU FOR YOUR SERVICE</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Veterans: your first month is <em>free.</em></h1>
    <p class="lead">Patriot Pest Control is veteran-owned. This Veterans Day, any U.S. military veteran who signs up for an annual pest control plan gets the first month free. Book by Wednesday, November 11, 2026.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap">
    <div class="steps" style="display:grid;grid-template-columns:repeat(auto-fit,minmax(220px,1fr));gap:1.2rem">
      <div class="card step"><span class="plan-tier">THE OFFER</span><h3 style="font-family:var(--display);color:var(--cream)">First Month Free</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Sign up for any 12-month plan (Bronze $49/mo, Silver $59/mo, or Gold $79/mo) and your first month is on us. One offer per household, for new annual plan signups, and it cannot be combined with other offers.</p></div>
      <div class="card step"><span class="plan-tier">WHO QUALIFIES</span><h3 style="font-family:var(--display);color:var(--cream)">U.S. Military Veterans</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Show any one of these when you sign up: a VA health care ID card, a DD Form 214, a military or retiree ID card, or the veteran designation on your state driver's license.</p></div>
      <div class="card step"><span class="plan-tier">HOW TO REDEEM</span><h3 style="font-family:var(--display);color:var(--cream)">Use Code VETSDAY</h3><p style="color:var(--khaki);line-height:1.7;margin-top:.5rem">Call or text (509) 471-5767 and mention code VETSDAY, or sign up online and put VETSDAY in the notes. Book by November 11, 2026.</p></div>
    </div>
    <div class="promise" style="margin-top:1.6rem"><b>Veteran-owned:</b> Patriot Pest Control is owned and run by Skyler Rose, a U.S. Military Veteran, based in Deer Park, Washington. Serving Spokane, Spokane Valley, Cheney, Liberty Lake, Airway Heights, Medical Lake, Deer Park, and Mead in Washington, and Coeur d'Alene, Post Falls, Hayden, and Rathdrum in Idaho.</div>
  </div>
</section>

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">Claim your <em>free month.</em></h2>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="tel:+15094715767"><svg class="ic" viewBox="0 0 24 24" aria-hidden="true"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6A19.79 19.79 0 0 1 2.12 4.18 2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72c.13.96.36 1.9.7 2.81a2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45c.91.34 1.85.57 2.81.7A2 2 0 0 1 22 16.92z"/></svg> Call or Text (509) 471-5767</a>
      <a class="btn btn-ghost" href="/signup">Sign Up ▸</a>
    </div>
  </div>
</section>`

func init() {
	RegisterPageTemplate("veterans", pageVeterans)
}
