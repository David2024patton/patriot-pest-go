package view

// Insect pest overrides: region-accurate species and unique copy for the
// Inland Northwest (Spokane WA / Coeur d'Alene ID). Sourced from university
// extension and state agriculture/health publications (see report). The
// audit flagged wrong-region species on these pages; each ScientificName
// below names the species actually doing the damage here.
func init() {
	RegisterPestOverride("ants", PestOverride{
		ScientificName: "Camponotus modoc",
		Description:    `Around Spokane and Coeur d'Alene, the ant doing real structural damage is the western black carpenter ant, Camponotus modoc: big, matte-black workers with reddish legs that carve smooth galleries through damp or decaying wood. They do not eat the wood, they excavate it, and a home infestation is usually a satellite colony with the parent nest up to a hundred yards away in a stump or dead log. The WSU, OSU, and University of Idaho joint publication PNW624 lists them alongside odorous house ants (Tapinoma sessile), which trail into kitchens for sweets, and moisture ants (Lasius pallitarsis), which move into wood already softened by rot. Carpenter ant colonies can reach tens of thousands of workers, and they forage mostly after dark, so a quiet kitchen at noon tells you nothing. Knock on suspect trim and listen: disturbed carpenter ants make a faint rustling inside the void.`,
		Signs: []string{
			"Piles of sawdust-like frass (wood shavings and insect parts) below windowsills, baseboards, or deck posts",
			"Winged swarmers, or discarded wings, indoors in spring, a sign a mature colony is nearby",
			"Faint rustling or crinkling sounds inside walls or hollow doors at night",
			"Steady trails of small dark ants moving between the foundation and kitchen sweets or pet food",
		},
		Treat: `We run recon inside and out to find the parent colony and its satellites, then hit the voids directly with targeted treatments and non-repellent transfer products, and call out any moisture-damaged wood that invited them in.`,
		Prev: []string{
			"Fix roof, plumbing, and crawl space leaks so wood stays dry and uninviting",
			"Trim tree limbs and shrubs so nothing bridges from vegetation to the siding",
			"Store firewood and lumber off the ground and well away from the foundation",
			"Wipe up crumbs, seal sweets, and keep pet food in hard containers with lids",
		},
	})

	RegisterPestOverride("bed-bugs", PestOverride{
		ScientificName: "Cimex lectularius",
		Description:    `The common bed bug, Cimex lectularius, is a hitchhiker, not a dirtiness problem. It rides luggage, used furniture, and moving boxes into clean homes, apartments, dorms, and rentals all along the Spokane to Coeur d'Alene corridor, then feeds on sleeping people at night and hides by day in mattress seams, box springs, headboards, and baseboard cracks. Adults are about the size of an apple seed, flat, and reddish brown, and they can go months between meals, which is why an empty guest room is no guarantee of a clear room. Bites show up as itchy welts, often in lines or clusters, though reactions vary widely from person to person. Because they spread unit to unit through walls and shared laundry, early detection is the whole battle: a small, contained group is a straightforward job, while a colony left to breed for a year turns into a full building operation.`,
		Signs: []string{
			"Rust-colored fecal spots and blood smears along mattress seams, tags, and piping",
			"Translucent shed skins and tiny white eggs tucked into headboard joints and box spring corners",
			"Live flat reddish-brown bugs, apple-seed sized, hiding in seams when the lights come on",
			"Itchy welts in lines or clusters on skin exposed during sleep",
		},
		Treat: `We strip the bed down to the frame and inspect every seam and joint, then combine heat or steam with crack-and-crevice residuals, fit mattress encasements, and return for a verification visit to confirm the colony is actually gone.`,
		Prev: []string{
			"Inspect hotel mattress seams and headboards before unpacking on any trip",
			"Never drag curbside furniture or mattresses inside without a full inspection first",
			"Cut bedroom clutter so there are fewer harborages near the bed",
			"Launder travel clothing on hot and run the dryer on high heat after every trip",
		},
	})

	RegisterPestOverride("bees", PestOverride{
		ScientificName: "Apis mellifera",
		Description:    `The western honey bee, Apis mellifera, is the bee most Inland Northwest homeowners meet: spring swarms clustering on a branch or fence post while scouts house-hunt, and established colonies tucked into wall voids, soffits, or chimneys. Bumble bees (Bombus species) nest in the ground or in old rodent burrows and are generally docile unless stepped on. The wood-borer in this region is the western carpenter bee (Xylocopa tabaniformis), which drills neat half-inch round holes into fascia boards, decks, and eaves and raises young in the tunnels. Honey bees are pollinators worth protecting, so removal prioritizes live relocation through beekeeper partners whenever the colony is accessible. One regional correction: Africanized honey bees are an established hazard in Arizona desert country, but they are not established in the Inland Northwest, so defensive swarms here are almost always ordinary European honey bees.`,
		Signs: []string{
			"A buzzing cluster of bees hanging from a branch, fence, or eave in spring swarm season",
			"Steady bee traffic in and out of a gap in siding, soffits, or around a chimney",
			"Perfectly round half-inch holes in fascia or deck wood with sawdust piles below, the carpenter bee signature",
			"Loud humming inside a wall void on warm afternoons",
		},
		Treat: `Accessible honey bee swarms and colonies are captured for relocation with beekeeper partners whenever possible; carpenter bee galleries get treated after activity ends and the holes sealed, and structural voids are closed so the next swarm keeps moving.`,
		Prev: []string{
			"Seal gaps in eaves, soffits, siding, and chimney caps before spring swarm season",
			"Paint or stain exposed softwood on decks and fascia, carpenter bees prefer bare wood",
			"Keep sugary drinks and hummingbird feeders away from doorways in late spring",
			"Call before spraying any swarm so live relocation stays on the table",
		},
	})

	RegisterPestOverride("box-elder-bugs", PestOverride{
		ScientificName: "Boisea trivittata",
		Description:    `The boxelder bug, Boisea trivittata, is a fall specialist. All summer it feeds on the seeds of boxelder, maple, and ash trees, and when temperatures drop it masses by the hundreds on sun-warmed south and west walls, then squeezes through tiny gaps to overwinter inside wall voids. They do not bite people, do not eat wood or fabric, and do not reproduce indoors, which makes them a nuisance invader rather than a structural threat. The real damage is cosmetic: crushed bugs leave orange-red stains on curtains, walls, and siding. Homes near female (seed-bearing) boxelder trees take the heaviest pressure, and a warm October afternoon can turn a quiet siding wall into a crawling red-and-black carpet. Because the invasion is seasonal and predictable, the winning move is exclusion and an exterior barrier laid down before the fall push, not chasing individual bugs with a vacuum all winter.`,
		Signs: []string{
			"Dense clusters of red-and-black bugs basking on sunny siding in September and October",
			"Bugs slipping indoors around windows, doors, and utility penetrations as nights cool",
			"Orange-red staining on curtains, walls, or siding where bugs were crushed",
			"Dead bugs collecting in window tracks and along baseboards through winter",
		},
		Treat: `We lay an exterior barrier on the sun-struck congregation faces in late summer and fall, vacuum out the indoor stragglers, and seal the entry points they used so next autumn's migration bounces off.`,
		Prev: []string{
			"Caulk gaps around windows, doors, and where utilities enter before fall",
			"Fit tight screens and repair torn ones on windows and vents",
			"Consider removing female boxelder trees planted close to the house",
			"Rake up fallen boxelder and maple seeds near the foundation",
		},
	})

	RegisterPestOverride("cockroaches", PestOverride{
		ScientificName: "Blattella germanica",
		Description:    `In Spokane and North Idaho kitchens, the cockroach that matters is the German cockroach, Blattella germanica: small, tan, with two dark stripes behind the head, and the most common indoor cockroach in the country according to the EPA. German roaches live for warmth, moisture, and tight cracks, which is why infestations center on kitchens, bathrooms, and appliance motors, and a single egg case holds 30 to 40 nymphs, so numbers explode fast. The big American cockroach (Periplaneta americana) that some sites picture is a sewer and basement roach of warmer climates, not the insect behind a Spokane apartment infestation. German roach droppings and shed skins trigger allergies and asthma, especially in kids, and heavy infestations contaminate food and leave a musty, oily odor. Arizona note: the desert flips the script, and in Phoenix the American cockroach (Periplaneta americana) is the classic sewer-dwelling roach coming up through drains and cleanouts.`,
		Signs: []string{
			"Pepper-like black droppings clustered in cabinet corners, drawer tracks, and under appliances",
			"A musty, oily odor in the kitchen that gets stronger at night",
			"Tan purse-shaped egg cases (oothecae) glued under refrigerators, dishwashers, and sinks",
			"Roaches scattering when the kitchen light flips on after dark",
		},
		Treat: `We place gel baits deep in cracks and voids where roaches actually travel and add insect growth regulators that sterilize the next generation, then monitor with traps and re-bait until activity reads zero.`,
		Prev: []string{
			"Fix drips and dry out the cabinet under the sink, moisture is their anchor",
			"Store cereal, flour, and pet food in sealed hard containers",
			"Take trash out nightly and rinse recyclables before binning them",
			"Break down and discard cardboard clutter behind appliances and under sinks",
		},
	})

	RegisterPestOverride("crickets", PestOverride{
		ScientificName: "Acheta domesticus",
		Description:    `The house cricket, Acheta domesticus, is the noisy one: pale brown, about an inch long, with males that chirp all night by rubbing their wings together to call mates. Around here they share the night shift with field crickets (Gryllus species) that wander in from lawns in late summer and native camel crickets (Ceuthophilus species), the humpbacked, long-legged jumpers that haunt damp basements and crawl spaces without ever chirping. Crickets do not bite people, but they chew: wool, silk, cotton, paper, houseplant leaves, and pet food are all on the menu, and a big indoor population stains fabrics with droppings. They are moisture seekers, so a cricket problem in the basement is usually a humidity problem wearing a disguise. Cold weather drives them toward foundation cracks and door gaps, and once inside they settle into wall voids and storage boxes for the winter.`,
		Signs: []string{
			"Rhythmic chirping at night coming from basements, wall voids, or behind appliances",
			"Irregular chewed holes in stored clothing, curtains, or houseplant leaves",
			"Long-legged humpbacked camel crickets jumping in damp basements or crawl spaces",
			"Crickets turning up in glue traps, window wells, and storage boxes",
		},
		Treat: `We dry out the harborage with moisture and clutter reduction, treat cracks, crevices, and entry points where they cross the perimeter, and lay an exterior barrier to cut off the late-summer migration indoors.`,
		Prev: []string{
			"Seal foundation cracks and fit door sweeps so fall migrants cannot walk in",
			"Pull mulch, ground cover, and stacked items back from the foundation",
			"Run a dehumidifier in damp basements and crawl spaces",
			"Store off-season fabrics in sealed bins instead of cardboard boxes",
		},
	})

	RegisterPestOverride("fleas-ticks", PestOverride{
		ScientificName: "Ctenocephalides felis, Dermacentor andersoni",
		Description:    `Fleas and ticks hit from two directions. The cat flea, Ctenocephalides felis, is the flea on Inland Northwest dogs and cats, and its pupae are the reason flea jobs fail: pupae can sit dormant in carpet for months, then hatch in waves when vibration and warmth signal a host is near. Adult fleas are only about five percent of the population, so killing what you see barely dents the operation. Ticks are the outdoor half of this page. Washington Department of Health tick data and an Eastern Washington University study of Spokane County name the Rocky Mountain wood tick (Dermacentor andersoni) and the American dog tick (Dermacentor variabilis) as the two common hard ticks here, questing in grass and brush from spring into early summer and latching onto dogs, hikers, and kids. Both can carry Rocky Mountain spotted fever and tularemia, so a tick found attached is worth identifying, not just flicking away.`,
		Signs: []string{
			"Dogs or cats scratching and chewing, with black flea dirt in bedding that bleeds red on damp paper",
			"Itchy bites clustered around ankles and lower legs after walking on carpet",
			"Tiny white flea eggs and larvae deep in carpet fibers, pet bedding, and furniture seams",
			"Ticks attached to dogs or people after hikes, or crawling on clothing coming indoors",
		},
		Treat: `We treat pet resting areas, carpet edges, and yard harborage with products plus insect growth regulators that break the flea life cycle, cut back tick habitat along fence lines and trails, and coordinate timing with your vet's on-animal flea prevention.`,
		Prev: []string{
			"Wash pet bedding in hot water weekly during flea season",
			"Vacuum floors, furniture seams, and under cushions often, then empty the canister outside",
			"Keep grass mowed and clear leaf litter where ticks quest",
			"Check dogs, kids, and yourself for ticks after time in grass or brush",
		},
	})

	RegisterPestOverride("fruit-flies", PestOverride{
		ScientificName: "Drosophila melanogaster",
		Description:    `Small indoor flies are a lineup, not one insect, and each needs different recon. The classic fruit fly is Drosophila melanogaster: tiny, tan, red-eyed, breeding in fermenting fruit, spilled juice, beer, and the slime under a garbage can lid, with a full generation in about a week. Drain flies (Psychoda species, the fuzzy moth-like ones) breed in the bacterial slime coating seldom-used drains and overflow channels. The house fly (Musca domestica), described in WSU Puyallup's filth-fly bulletin, is the big one on windows and food, breeding in garbage and pet waste and capable of moving bacteria onto food-contact surfaces. None of these are blow flies (Calliphora species), which breed on carrion outdoors. Because every small fly traces back to a breeding source, fogging the air is theater: find the rotten potato, the slimy drain, or the mop bucket, remove it, and the cloud collapses.`,
		Signs: []string{
			"Tiny tan flies with red eyes hovering over fruit bowls, trash cans, or recycling",
			"Fuzzy gray moth-like flies that appear when a seldom-used drain runs",
			"House flies clustering on sunny windows and landing on food",
			"Maggots in garbage disposal splash guards, mop buckets, or under trash can lids",
		},
		Treat: `We identify the exact fly and hunt its breeding source first, then clean drains mechanically, remove the food source, and use traps only as monitors to confirm the population is actually collapsing.`,
		Prev: []string{
			"Refrigerate ripe fruit and vegetables instead of leaving them on the counter",
			"Scrub drains and garbage disposals weekly to remove the slime flies breed in",
			"Take trash and recycling out frequently and keep bins lidded and rinsed",
			"Fix slow or leaky drains so standing organic sludge never builds up",
		},
	})

	RegisterPestOverride("hornets", PestOverride{
		ScientificName: "Dolichovespula maculata",
		Description:    `Here is the correction that matters: the insect Inland Northwest homeowners call a hornet is the bald-faced hornet, Dolichovespula maculata, which WSU entomologists note is technically an aerial yellowjacket, not a true hornet at all. True hornets in the genus Vespa, including the European hornet (Vespa crabro), are an eastern US insect and are not established in Washington. Bald-faced hornets are big, black-and-ivory wasps that build unmistakable gray paper footballs hanging from tree limbs, eaves, and shrubs, housing colonies of several hundred workers by late summer. They are beneficial predators of flies and caterpillars most of the season, but they defend the nest with coordinated, persistent stings, and late-summer colonies are the ones that send people to urgent care. The Asian giant hornet detections in Whatcom County were a targeted eradication campaign, not an established Spokane-area pest.`,
		Signs: []string{
			"A large gray papery football-shaped nest hanging from a tree limb, shrub, or eave",
			"Heavy two-way traffic of big black-and-white wasps around a single point",
			"Wasps stripping gray wood fiber from fences, decks, or firewood to build nest paper",
			"Wasps dive-bombing or chasing people and pets that wander near the nest",
		},
		Treat: `Suited technicians treat the nest in the evening when the foragers are home, then remove the nest and sweep the property for satellite starts before they reach football size.`,
		Prev: []string{
			"Walk the eaves, trees, and shrubs in spring to catch golf-ball-sized starter nests early",
			"Keep trash cans sealed so late-summer workers are not drawn to the yard",
			"Never swat at or spray a nest with a garden hose, it triggers the defense response",
			"Call before the nest grows past softball size, small nests are far simpler jobs",
		},
	})

	RegisterPestOverride("mosquitoes", PestOverride{
		ScientificName: "Culex pipiens, Culex tarsalis",
		Description:    `Washington's mosquitoes are not the tropical kind. The Washington Department of Health names Culex pipiens and Culex tarsalis as the state's primary West Nile virus vectors, with the aggressive floodwater mosquito Aedes vexans surging after spring floods and irrigation season. The yellow fever mosquito, Aedes aegypti, is not established in Washington, so any page naming it as the local threat is working from the wrong map. Around Spokane and Coeur d'Alene, mosquitoes breed in the small stuff: clogged gutters, birdbaths, kiddie pools, tire swings, and the saucers under planters, needing only a bottle cap of stagnant water. They rest by day in cool, shaded vegetation and hunt at dusk and dawn. Arizona note: the desert is different, and Arizona health officials report Aedes aegypti well established across much of the state, including abundant populations in Maricopa County, where it has been tied to locally acquired dengue cases.`,
		Signs: []string{
			"Bites concentrated at dusk and dawn, often around ankles and exposed arms",
			"Tiny wrigglers (larvae) hanging just under the surface of standing water in gutters or birdbaths",
			"Adults flushing out of shaded shrubs and tall grass when disturbed during the day",
			"Swarms rising near ponds, ditches, or irrigated fields after floods or heavy watering",
		},
		Treat: `We drain or treat larval water sources first, then apply barrier treatments to the shaded vegetation where adults rest, timed to the local season so the yard stays usable through the worst weeks.`,
		Prev: []string{
			"Dump and scrub anything holding standing water at least once a week",
			"Clean gutters so water never ponds in the troughs",
			"Use EPA-registered repellent for evening yard time in peak season",
			"Keep window and door screens tight and repair tears promptly",
		},
	})

	RegisterPestOverride("pantry-pests", PestOverride{
		ScientificName: "Plodia interpunctella",
		Description:    `The Indian meal moth, Plodia interpunctella, is the pantry's most common invader: a small moth with coppery bronze wingtips that arrives inside grocery packaging as eggs or larvae and then breeds in flour, cereal, rice, pet food, birdseed, and dried fruit. The larvae do the damage, spinning fine silken webbing through dry goods and leaving frass that ruins the whole box. They share the shelves with sawtoothed grain beetles, flour beetles, and weevils, all of which hitchhike home from the store rather than invading from outdoors. One infested bag of bulk-bin granola can seed every open package in the pantry within weeks, and the moths themselves flutter at dusk, which is usually when homeowners first notice the operation is underway. Pheromone traps confirm the species, but traps alone never end it: the fix is finding every infested product, because a single overlooked bag restarts the whole cycle.`,
		Signs: []string{
			"Small moths with bronze wingtips fluttering in the pantry at dusk",
			"Fine silken webbing inside cereal boxes, flour bags, or pet food",
			"Tiny beetles or weevils crawling in flour, rice, or dry pet food",
			"Clumped grains and gritty frass in otherwise sealed dry goods",
		},
		Treat: `We help you hunt down and discard every infested product, deep-clean the shelf cracks and corners where larvae pupate, set pheromone monitors, and treat voids only if the infestation has spread beyond the food.`,
		Prev: []string{
			"Inspect bulk-bin and bagged dry goods for webbing before shelving them",
			"Move flour, cereal, rice, and pet food into sealed hard containers",
			"Practice first in, first out so nothing sits open for months",
			"Vacuum pantry shelves and shelf-pin holes a few times a year",
		},
	})

	RegisterPestOverride("scorpions", PestOverride{
		ScientificName: "Paruroctonus boreus",
		Description:    `Washington's scorpion story is nothing like the desert's. The northern scorpion, Paruroctonus boreus, is the species WSU researchers documented across the arid and semi-arid regions of central and eastern Washington, and it is a shy, nocturnal hunter of insects whose sting is not medically significant. Bark scorpions in the genus Centruroides do not occur in the Inland Northwest; they are a southern Arizona problem. Northern scorpions hide by day under rocks, bark, and debris and are rare inside homes, with most encounters happening around woodpiles, rock landscaping, and garages in the drier country south and west of Spokane. Under a UV flashlight they fluoresce a ghostly blue-green, which makes nighttime surveys surprisingly effective. Arizona note: the Phoenix area is genuinely different, home to the Arizona bark scorpion (Centruroides sculpturatus), whose sting is medically significant, and the giant desert hairy scorpion (Hadrurus arizonensis).`,
		Signs: []string{
			"Scorpions found under rocks, firewood, or debris near the foundation",
			"Blue-green fluorescence under a UV flashlight during nighttime yard checks",
			"Scorpions wandering into garages or mudrooms near ground-level door gaps",
			"Sightings clustered around rock landscaping, woodpiles, and south-facing slopes",
		},
		Treat: `We survey with UV light at night to map activity, cut back the rock and wood harborage against the house, seal ground-level entry gaps, and apply targeted perimeter treatments where scorpions travel.`,
		Prev: []string{
			"Move firewood, lumber, and rock piles away from the foundation",
			"Seal gaps under doors and around ground-level utility entries",
			"Wear gloves and shake out shoes when working around woodpiles",
			"In desert country, run UV flashlight checks of patios and pool areas at night",
		},
	})

	RegisterPestOverride("silverfish", PestOverride{
		ScientificName: "Lepisma saccharina",
		Description:    `The common silverfish, Lepisma saccharina, is a moisture meter with legs: a silvery, fish-shaped, wingless insect about half an inch long that thrives in damp basements, bathrooms, attics, and laundry rooms. Silverfish eat starches and sugars, which puts book bindings, wallpaper paste, envelopes, photographs, and starched clothing on the menu, leaving irregular feeding marks, yellowish stains, and surface etching. They are nocturnal and fast, famous for the sudden sprint across a bathtub or sink when a light flips on, and they can live a year or more while breeding slowly in undisturbed clutter. A few silverfish are a humidity report; a lot of them, or damage to stored papers and books, means the moisture problem has been running long enough to need both drying and treatment. They do not bite, do not spread disease, and mostly signal that a space needs airflow and dehumidification.`,
		Signs: []string{
			"Silvery fish-shaped insects sprinting across sinks or bathtubs when lights come on",
			"Irregular notches and surface etching on book bindings, wallpaper, or envelopes",
			"Yellowish stains on stored papers, photos, or clothing",
			"Tiny shed skins and pepper-like droppings in undisturbed boxes and files",
		},
		Treat: `We dry out the habitat first with dehumidification and leak repair, then apply crack-and-crevice treatments to the harborages and set monitors to confirm the population is actually declining.`,
		Prev: []string{
			"Run dehumidifiers in damp basements, crawl spaces, and laundry rooms",
			"Store books, papers, and photos in sealed bins instead of cardboard",
			"Fix plumbing leaks and ventilate bathrooms after showers",
			"Caulk baseboards and seal around pipe chases where they travel",
		},
	})

	RegisterPestOverride("spiders", PestOverride{
		ScientificName: "Latrodectus hesperus",
		Description:    `The spider that deserves respect in the Inland Northwest is the western black widow, Latrodectus hesperus: glossy black females with the red hourglass, building tangled webs in dry, undisturbed spots like woodpiles, crawl spaces, and sheds. WSU entomologists call the female black widow the only medically dangerous native spider in Washington, while noting she is timid and bites only when pressed against skin. The hobo spider (Eratigena agrestis), the funnel-web spider common in Northwest homes and gardens, does not deserve its old fearsome reputation: WSU scientists could never replicate the early claims of dangerous venom, and state health officials now list it as not dangerous. The brown recluse does not occur naturally in the Northwest at all. Most of what homeowners see, giant house spiders, orb weavers, and cellar spiders, are harmless pest-eaters, but webs in living spaces and venomous species near doors and play areas still call for action.`,
		Signs: []string{
			"Tangled, messy cobwebs in woodpiles, crawl spaces, and behind stored items",
			"Funnel-shaped webs with a fast brown spider at ground or foundation level",
			"A glossy black spider with a red hourglass marking in a dry, undisturbed spot",
			"Round papery egg sacs tucked into clutter, sheds, or under eaves",
		},
		Treat: `We remove webs and egg sacs, apply crack-and-crevice treatments to harborages and entry points, and cut the insect prey base so the property stops looking like good hunting ground.`,
		Prev: []string{
			"Declutter garages, sheds, and storage so spiders lose their harborages",
			"Wear leather gloves when moving firewood, rocks, or stored boxes",
			"Caulk foundation cracks and keep screens and door sweeps tight",
			"Trim vegetation back from siding and outdoor lighting fixtures",
		},
	})

	RegisterPestOverride("stink-bugs", PestOverride{
		ScientificName: "Halyomorpha halys",
		Description:    `The brown marmorated stink bug, Halyomorpha halys, is an invasive shield-shaped insect that spends summer feeding on fruits, vegetables, and ornamentals, then marches indoors by the dozens when fall cools down. Like boxelder bugs, they overwinter inside wall voids and attics without feeding or reproducing, then wake up confused on warm winter days and blunder into living rooms. Their signature defense is the smell: disturb or crush one and it releases a pungent, cilantro-like odor that lingers on skin and fabric. They do not bite people or damage structures, but they stain curtains and walls, and a house that hosted them one winter will host them again the next unless entry points get sealed. Because they aggregate on sun-warmed walls before pushing inside, the treatment window is late summer through early fall, aimed at the exterior before the migration, not at the ones already sleeping in your walls.`,
		Signs: []string{
			"Shield-shaped mottled-brown bugs gathering on sunny walls and windows in fall",
			"Bugs emerging sluggishly from wall voids and attics on warm winter days",
			"A sharp cilantro-like odor when a bug is disturbed or crushed",
			"Staining on curtains and walls where bugs were crushed",
		},
		Treat: `We apply an exterior barrier on congregation faces before the fall migration, vacuum up indoor bugs without crushing them, and seal the gaps around windows, vents, and siding they used to get in.`,
		Prev: []string{
			"Seal gaps around windows, doors, vents, and where siding meets trim",
			"Cover attic and crawl space vents with fine mesh",
			"Repair torn window screens before the fall invasion",
			"Vacuum indoor stink bugs instead of crushing them",
		},
	})

	RegisterPestOverride("termites", PestOverride{
		ScientificName: "Reticulitermes hesperus",
		Description:    `The termite to watch in the Inland Northwest is the western subterranean termite, Reticulitermes hesperus. A WSU extension bulletin names it as one of Washington's two established termite species, alongside the Pacific dampwood termite (Zootermopsis angusticollis), which needs very wet wood and is mostly a west-side and log-home concern. Western subterranean termites live in soil colonies that can number in the hundreds of thousands, sending workers through mud tubes up foundations to feed on structural wood, paper, and anything cellulose. They swarm on warm days in spring and fall, especially after rain, when winged reproductives pour out to start new colonies. Damage stays hidden for years because they eat wood from the inside out, hollowing studs and sills while the paint looks fine. The telltales are unmistakable once you know them: earthen mud tubes on foundations, discarded wings on windowsills, and wood that sounds hollow when tapped.`,
		Signs: []string{
			"Pencil-width mud tubes climbing foundation walls, piers, or crawl space posts",
			"Piles of discarded translucent wings on windowsills after a warm rainy day",
			"Wood that sounds hollow or papery when tapped, especially sills and framing",
			"Swarming winged insects pouring from soil or wood in spring or fall",
		},
		Treat: `We run a full structural inspection to map activity, then treat the soil perimeter or install baiting stations with ongoing monitoring, and flag the moisture conditions that let the colony thrive.`,
		Prev: []string{
			"Eliminate wood-to-soil contact: no siding, steps, or posts buried in grade",
			"Grade soil and gutters so water drains away from the foundation",
			"Fix plumbing and roof leaks that keep crawl spaces and sills damp",
			"Store firewood and scrap lumber off the ground and away from the house",
		},
	})

	RegisterPestOverride("wasps", PestOverride{
		ScientificName: "Polistes dominula",
		Description:    `The paper wasp most Inland Northwest homeowners meet is the European paper wasp, Polistes dominula, an introduced species that WSU's extension fact sheet describes as well established across the Pacific Northwest. It builds the classic open-comb umbrella nest, a single layer of papery cells hanging from eaves, deck rails, play structures, and grill covers, starting in spring from a lone overwintered queen and growing to a few dozen workers by late summer. These wasps are yellow-and-black, slender, with long legs dangling in flight, and they are beneficial predators that hunt caterpillars in gardens all season. They are also the wasp most likely to sting a homeowner, because their nests hang exactly where hands and heads go: under eaves, around door frames, and inside sheds. Unlike yellowjackets they are not aggressive away from the nest, so a wasp working your flower beds can be left alone, but a nest over the back door cannot.`,
		Signs: []string{
			"Open umbrella-shaped paper combs hanging under eaves, deck rails, or play equipment",
			"Slender yellow-and-black wasps with dangling legs scraping wood fiber from fences and decks",
			"Wasps hovering persistently around door frames, sheds, and grill covers",
			"A single large queen working a small starter nest in early spring",
		},
		Treat: `We treat the nest in the evening when workers have returned, remove the comb, and close off the eave voids and cavities so the next queen cannot rebuild in the same spot.`,
		Prev: []string{
			"Walk eaves and structures in spring and knock down golf-ball-sized starter nests early",
			"Seal gaps in eaves, soffits, and around door and window frames",
			"Keep trash cans sealed and food covered during late-summer cookouts",
			"Leave foraging wasps in the garden alone and act only on nests near people",
		},
	})

	RegisterPestOverride("yellow-jackets", PestOverride{
		ScientificName: "Vespula pensylvanica",
		Description:    `The yellowjacket running the Inland Northwest is the western yellowjacket, Vespula pensylvanica: a native, ground-nesting species that builds enclosed paper nests in old rodent burrows, wall voids, and landscape timbers, growing to several thousand workers by August and September. WSU's yellowjacket bulletin describes the pattern exactly: queens start small nests in spring, workers take over through summer, and the colony turns aggressive in late summer as natural food dwindles and workers switch to scavenging meat, soda, and garbage. That is why a calm June nest becomes a picnic-ruining menace by Labor Day, and why most stings happen when someone mows over, steps on, or leans against a hidden entrance. The German yellowjacket (Vespula germanica) is an introduced species present in parts of the urban Northwest, but the native western yellowjacket is the dominant ground nester around Spokane. Colonies die with the first hard frosts; only mated queens overwinter.`,
		Signs: []string{
			"A steady stream of yellow-and-black wasps flying in and out of a hole in the ground or a wall void",
			"Yellowjackets swarming meat, soda cans, and garbage at late-summer cookouts",
			"Wasp traffic around landscape timbers, sheds, or abandoned rodent burrows",
			"Sudden aggressive stinging when mowing, digging, or brushing against a hidden nest entrance",
		},
		Treat: `We locate the nest entrance during peak traffic, treat it in the evening when the colony is home, and open and treat wall voids where nests are hidden inside structures.`,
		Prev: []string{
			"Keep outdoor food and drinks covered and trash cans tightly lidded in late summer",
			"Fill abandoned rodent burrows and holes near patios, play areas, and walkways",
			"Seal gaps in siding and foundations where queens scout for void nest sites in spring",
			"Watch for heavy wasp traffic in spring and act while the nest is still small",
		},
	})
}
