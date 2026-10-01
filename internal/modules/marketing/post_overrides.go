package marketing

import "github.com/David2024patton/patriot-pest-go/internal/data"

// postBodyOverrides replaces the stored body_html of select blog posts with
// longer, audit-compliant articles. Keys are post slugs. The content below is
// static and authored in source, so it does not need the bluemonday pass that
// database-loaded bodies get. It uses only plain markup (p, h3, ul, li) that
// the post template renders inside .prose.
var postBodyOverrides = map[string]string{
	"rodent-proof-your-home": `<p>Every October, as the first hard frost settles over the Inland Northwest, rodents start house hunting. The fields around Spokane and Coeur d'Alene go cold and quiet, natural food dries up, and your heated home with its pantry, pet food, and cozy wall voids looks like the best real estate in the county. Mice and rats do not wander in by accident. They follow warmth, food odor, and shelter, and once one finds a way in, the rest of the family follows the scent trail.</p>
<p>Here is the part most people miss: traps and poison only handle the rodents already inside. They do nothing about the gap that let them in. Exclusion, finding and sealing every entry point, is the only permanent fix. Do the exclusion work and the problem stops coming back.</p>
<h3>Know what is moving in</h3>
<p>Around here you will mostly meet two mice. House mice are the classic indoor invader: small, gray-brown, and happiest living within a short distance of their nest, which is why they rarely travel far from your kitchen or pantry. Deer mice are the outdoors species common around Deer Park, Mead, and the rural edges of Spokane County. They move into garages, sheds, and crawl spaces when the weather turns.</p>
<p>The deer mouse deserves extra respect. The Washington State Department of Health identifies deer mice as the main carriers of the virus that causes hantavirus in people, a rare but potentially fatal respiratory disease. It spreads through dust stirred up from dried droppings and urine, not through bites. So where deer mice have been nesting, never sweep or vacuum dry droppings. Wet the area down with disinfectant first, then wipe it up while wearing gloves and a dust mask.</p>
<p>Speed matters because mice breed at an absurd pace. A University of Missouri Extension specialist notes that house mice can produce 8 to 10 litters a year, with 5 to 6 young per litter. The pair you hear in the wall in October can be a full-blown infestation by January.</p>
<h3>Find every gap</h3>
<p>A mouse can squeeze through a hole the width of a pencil, about 1/4 inch across, according to the CDC. Rats need only about half an inch. That means the gap around a dryer vent, a worn piece of weatherstripping, or a crack where the foundation meets the siding is a wide-open front door.</p>
<p>Walk the house with a flashlight, inside and out. Check where pipes, cables, and gas lines enter the building. Look at the seam between the foundation and the siding, attic and crawl space vents, the fireplace, and the line where floors meet walls. Inside, check under and behind cabinets and appliances, in closets and their corners, and along the sill plate in the basement or crawl space. Give the garage door seal a hard look too, since it is one of the most common rodent entries and the corners wear out first.</p>
<h3>Seal it with materials they cannot chew through</h3>
<p>For small holes, pack in steel wool or copper mesh and seal over it with caulk or mortar. For larger openings, use hardware cloth, metal sheeting, or cement. These are the materials the Washington State Department of Health recommends, and the reason is simple: mice and rats gnaw straight through spray foam, plastic, and wood filler to reopen a hole. Seal it right the first time or you will be sealing it again.</p>
<p>Fit brush-style door sweeps on exterior doors, make sure windows and screens close tight, and cap the chimney. Pay special attention to the garage. It is attached to the house, full of clutter to hide in, and its big door seal is usually the weakest link.</p>
<h3>Take away the food, water, and cover</h3>
<p>A sealed house with an open pantry still invites trouble. Store people food, pet food, bird seed, and garden seed indoors in hard containers with tight-fitting lids. Keep garbage in thick plastic or metal cans with tight lids and clean up spilled waste instead of letting it sit.</p>
<p>Outside, trim shrubs and tree limbs so nothing touches the house, clear junk and debris from around outbuildings, and pick up fallen fruit. Store firewood, lumber, and hay at least 12 inches off the ground and away from the foundation. A woodpile leaning against the house is a rodent hotel with room service.</p>
<h3>Trap what is already inside</h3>
<p>Seal first, then deal with the current residents, or new ones will just replace the ones you catch. Snap traps set perpendicular to the wall, trigger end facing the baseboard, work well in the spots where you have seen droppings. Use several, and bait with a dab of peanut butter, which beats cheese. Check them daily.</p>
<p>If you keep catching mice week after week, you missed an entry point. Go back to the gap hunt. And if crawling under the house with a caulk gun in November does not appeal to you, a professional exclusion visit covers the whole checklist in one trip.</p>`,
	"spiders-fall-guide": `<p>If your home suddenly has more spiders in September and October, you are not imagining it, and it is not an invasion. It is dating season. In the fall, male spiders leave their webs and hiding spots and wander in search of mates, and some of that wandering leads them through your doors, windows, and foundation cracks. Washington State University entomologists note that of the spiders people bring in for identification each autumn, nine times out of ten it is a male, because males are the ones out roaming.</p>
<p>Cold weather adds a second push. As nights drop toward freezing around Spokane and Coeur d'Alene, outdoor spiders look for somewhere that will not freeze, and they are doing exactly what we do: getting ready for winter. A warm garage, basement, or crawl space fits the bill. Hauling in the first loads of firewood often carries hidden spiders inside with it.</p>
<h3>Who is actually showing up</h3>
<p>The hobo spider is the region's most famous fall visitor. It is a funnel-web spider, introduced from Europe, that builds flat sheet-like webs with a funnel retreat in the corner, often in basements, window wells, and ground-floor storage areas. Hobo spiders are poor climbers, so you will find them at floor level, not on the ceiling. Males wander through homes and garages in late summer and early fall looking for females.</p>
<p>Despite the scary name and an older nickname, the aggressive house spider, hobos are shy and try to escape when disturbed. Decades-old research once theorized their venom was dangerous, but WSU scientists have never been able to replicate that finding. There is no conclusive evidence the hobo spider is harmful to humans, though like any spider it deserves respect and proper identification if a bite is ever suspected.</p>
<p>You will also see wolf spiders, the big fast hunters that chase down insects instead of building webs, moving inside as the weather cools. Yellow sac spiders leave gardens for indoor shelter in fall as their outdoor food supply disappears. They hunt small insects at night and tuck into silky sacs in room corners by day.</p>
<h3>The one spider to actually respect</h3>
<p>The female black widow is the only medically dangerous spider native to Washington, and it is most common east of the Cascades, which includes our area. Look for the glossy black body with the red hourglass mark underneath. The good news is that black widows are extremely timid and retreat when disturbed. Bites are rare and almost always happen when someone presses one against skin by accident, like reaching bare-handed into a woodpile. Wear leather gloves when moving firewood or clearing clutter, and that risk drops to near zero.</p>
<p>One more myth to retire: the brown recluse does not occur naturally in the Pacific Northwest. Most reported recluse bites here turn out to be something else entirely, often a skin infection. Spiders also cannot be reliably identified by color alone, so if a bite is ever a concern, capture the spider in a jar for proper identification rather than guessing.</p>
<h3>Keep them outside where they belong</h3>
<ul>
<li>Install brush-style door sweeps on exterior doors and keep window screens tight and hole-free.</li>
<li>Seal cracks and gaps around doors, windows, and the foundation with caulk.</li>
<li>Knock down webs and remove egg sacs from eaves, corners, and window frames before they hatch.</li>
<li>Cut back outdoor lighting or switch to yellow bug lights, since lights draw the insects spiders feed on.</li>
<li>Declutter basements, garages, and storage areas so spiders have fewer dark hiding spots.</li>
<li>Shake out and inspect firewood outside before carrying it in, and wear gloves while handling it.</li>
</ul>
<p>Found one inside? A cup and a stiff card will relocate it outdoors without drama. Remember that spiders are free pest control. Every web in your garden all summer was catching flies, moths, and mosquitoes. A few sensible barriers keep the fall wanderers outside while letting them keep working for you out there.</p>`,
	"why-ants-invade-in-spring": `<p>Every April, kitchens around Spokane and Coeur d'Alene get the same unwelcome visitor: a thin line of ants marching across the counter like they pay rent. It starts with one or two scouts. As soil temperatures climb, colonies wake up hungry and send workers out looking for food and water. A single scout that finds a reliable source lays down a pheromone trail on the way home, and within days the whole foraging force is following that chemical highway straight into your kitchen.</p>
<p>Understanding that trail is the key to beating them. The ants you see are a tiny fraction of the colony. Kill every visible ant with spray and the colony simply sends more, because the queen underground never stops laying eggs. To actually end an ant problem, you have to either eliminate the colony or cut off what is drawing it inside.</p>
<h3>The two ants you will meet in the Inland Northwest</h3>
<p>Odorous house ants are the small brown ants behind most kitchen invasions. They love sweets, nest just about anywhere, and are famously hard to eliminate because their colonies hold multiple queens and constantly split off new nests, a habit called budding. In cooler climates they often overwinter inside wall voids, which is why they can seem to appear from nowhere in spring. Crush one and it gives off a rotten coconut smell, which is how they got the name.</p>
<p>Carpenter ants are the big ones, black or reddish and up to half an inch long, and they are a different problem entirely. They tunnel into wood to build nests, starting with damp or decayed wood around leaks, windows, and decks but expanding into sound framing. Look for small piles of coarse sawdust below windowsills and baseboards, and listen for faint rustling inside walls at night. A carpenter ant problem is a structural problem, and it usually calls for professional help.</p>
<h3>Why the spray can is not the answer</h3>
<p>Contact sprays kill the ants you can see and leave the colony untouched. Worse, spraying repellent insecticides around odorous house ants can scatter the colony and trigger budding, turning one nest into several. University of Nebraska entomologists put it plainly: sprays might kill the ones you are seeing, but you are not getting rid of the problem because you are not getting rid of the nest.</p>
<p>Baits work the opposite way. Foraging ants carry the slow-acting bait back and feed it to the colony, including the queen. When the queen dies, the colony collapses. Enclosed bait stations are also the safer option around kids and pets, since the pesticide stays inside the station instead of drifting wherever a spray lands.</p>
<h3>The plan that actually works</h3>
<ul>
<li>Wipe down the trails. Ants navigate by pheromone scent, so clean counters, baseboards, and floors along their route with a household cleaner to erase the trail.</li>
<li>Vacuum up the ants you see instead of squishing them, which avoids spreading scent and alarm signals.</li>
<li>Store food, especially sweets, cereal, and pet food, in airtight containers. Do not leave pet bowls out all day.</li>
<li>Fix the moisture. Repair leaks, keep gutters clear, and ventilate crawl spaces. Damp wood is what invites carpenter ants.</li>
<li>Seal entry points with caulk around the foundation, windows, and where pipes enter the house. Trim vegetation touching the siding.</li>
<li>Place enclosed ant baits along active trails and at entry points, following the label exactly. Give them a week or more to work. A bait that kills too fast never reaches the queen.</li>
</ul>
<p>One last distinction: ants you see only in spring and summer usually nest outdoors and are commuting in for food, so sanitation and exclusion often solve it. Ants present year-round, especially large carpenter ants in winter, usually mean the nest is inside your walls. That is the point to call a professional rather than fighting a colony you cannot reach.</p>`,
}

// applyPostBodyOverride swaps the expanded article body into p when its slug
// carries an override. Call it right after loading the post from the catalog
// and before rendering the blog post page.
func applyPostBodyOverride(p *data.Post) {
	if p == nil {
		return
	}
	if body, ok := postBodyOverrides[p.Slug]; ok {
		p.BodyHTML = body
	}
}
