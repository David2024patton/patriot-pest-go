package view

// Carrier-required SMS consent language, in one place.
//
// The US carriers' 10DLC review compares the consent text where a number is
// collected against the privacy policy and terms pages, and against the brand
// name shown on the site. So the business name below and the copy on the legal
// pages have to stay in step: if either changes, change both.
//
// Used by the contact form and the signup form, the two public pages that
// collect a phone number. Kept as one constant so they cannot drift apart.
const smsConsent = `<p class="sms-consent" style="color:var(--khaki);font-size:.8rem;line-height:1.5;margin:.4rem 0 0">By providing your phone number you agree to receive text messages from Patriot Pest Control Co. about your appointments, estimates and service. Message frequency varies. Msg &amp; data rates may apply. Reply STOP to opt out, HELP for help. See our <a href="/privacy-policy">Privacy Policy</a> and <a href="/terms-of-use">Terms</a>.</p>`
