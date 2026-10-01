package view

// Region-accurate copy for the rodent and wildlife pest pages. Every entry
// replaces the duplicated category-generic Signs/Treat/Prev blocks and the
// wrong-region scientific names flagged by the Keystone audit (moles named
// Scapanus latimanus, a California species; rodents named only Rattus
// norvegicus; pack rats described as desert woodrats; squirrels named as the
// introduced eastern gray). Species below were sourced from WDFW, USFWS,
// Animal Diversity Web, Idaho Fish and Game, and WSU Extension.

func init() {
	RegisterPestOverride("bats", PestOverride{
		ScientificName: "Myotis lucifugus",
		Description:    `Washington's most common bat, the little brown bat, owns the night shift over the Inland Northwest. Summer maternity colonies settle into the warm attics, barns, and bat houses of the Spokane and Coeur d'Alene area, each female raising a single pup through June and July, which is exactly when bat calls spike. At dusk they pour out of entry gaps to hunt insects over rivers, lakes, and wetlands, returning before dawn to the same roost. You will rarely see the bats themselves; you will see their work: crumbly guano piling up under roosts, dark oily rub stains on siding below entry points, and a sharp ammonia smell in a closed-up summer attic. Bats are protected wildlife in Washington, so there is no lethal option and no exclusion during maternity season, when pups cannot fly. The correct play is one-way doors after the pups fledge: the colony leaves for the night feed and finds its billet sealed behind it.`,
		Signs: []string{
			"High-pitched chirping or chittering from the attic or walls at dusk and dawn",
			"Dark oily rub stains around entry gaps under eaves, soffits, or gable vents",
			"Piles of dry crumbly guano beneath roosts or on siding (crumbles to dust, unlike mouse droppings)",
			"Sharp ammonia odor in a closed-up attic during summer months",
		},
		Treat: `Bats are protected wildlife in Washington, so control is exclusion only, never lethal. We run a dusk emergence count, install one-way doors on every entry after the June and July maternity season so no pups are sealed inside, then finish with a full cleanout and sanitation of accumulated guano.`,
		Prev: []string{
			"Cover gable vents, louvered vents, and chimney openings with 1/4 inch hardware cloth and wildlife-proof caps",
			"Seal gaps where siding meets the roofline, and repair damaged soffits and fascia",
			"Mount a bat house on a sun-facing wall away from the home to give them a legal billet",
			"Keep attic scuttle holes and access panels latched tight year-round",
		},
	})

	RegisterPestOverride("gophers", PestOverride{
		ScientificName: "Thomomys talpoides",
		Description:    `The northern pocket gopher is the smallest and most widespread pocket gopher in eastern Washington, and Spokane County sits squarely in its territory. You will almost never see one; they live their whole lives underground and surface only to shove dirt out of the burrow. You see the evidence instead: fresh crescent-shaped soil mounds with a plugged hole off to one side, often appearing overnight after rain or irrigation softens the ground. One gopher works an entire burrow system alone, eating roots, bulbs, and garden plants from below, which is why a row of vegetables can die one plant at a time with no visible culprit above ground. Their tunnels also chew through drip lines, gnaw underground cable, and trip up mowers. A single tunnel rat can hold a whole garden against your landscaping budget.`,
		Signs: []string{
			"Crescent-shaped mounds of fresh loose soil with the entrance plugged off to one side (moles make volcano mounds, gophers make crescents)",
			"Plants wilting or dying in a row, pulled down from the roots with no above-ground culprit",
			"Gnawed drip lines, irrigation tubing, or chewed underground cable",
			"New mounds appearing overnight after watering or rain, often along fence lines and garden beds",
		},
		Treat: `We set targeted traps in the active tunnels identified by probing for the main runway, then fill and compact the runs to deny re-entry. Pocket gophers are solitary and territorial, so removing the resident usually ends the campaign; repellents and sonic stakes do not work.`,
		Prev: []string{
			"Bury 1/2 inch hardware cloth 2 feet deep around garden beds, with 1 foot left above ground",
			"Plant bulbs and young trees in gopher-proof wire baskets",
			"Keep grass mowed short and clear weedy cover along fences where they stage",
			"Trap newcomers early in spring, before the single annual breeding season starts",
		},
	})

	RegisterPestOverride("moles", PestOverride{
		ScientificName: "Scapanus orarius",
		Description:    `The coast mole, Scheffer's subspecies, is the mole of the Inland Northwest, working the lawns, pastures, and garden beds around Spokane and Coeur d'Alene. It replaces the broad-footed mole of California that the old page named, a species that does not live here. Coast moles hunt earthworms and grubs in tunnels seven to ninety centimeters down, surfacing only to shove up the volcano-shaped molehills that mark their territory. Moles do not eat your plants; the grubs they chase do that. But their shallow surface runs lift and dry out turf roots, leaving brown ridges across the lawn, and one mole can work a quarter acre in a season. They stay active year-round and breed once a year, so one resident in worm-rich, moist soil can become the whole neighborhood's problem. They are not after your grass; they are after the supply lines underneath it.`,
		Signs: []string{
			"Volcano-shaped molehills of fine soil pushed up from below (gopher mounds are crescent-shaped and plugged)",
			"Raised surface ridges running across the lawn where shallow tunnels lifted the sod",
			"Spongy, lifted turf that dries out and browns along the tunnel runs",
			"Most activity in spring and fall, when moist soil brings worms near the surface",
		},
		Treat: `We place scissor or harpoon traps directly in the active main runways, verified by collapsing a section and confirming it is reopened within 24 hours. Poison baits and home remedies fail against an animal that eats live worms, not grain; trapping the solitary resident usually ends the damage.`,
		Prev: []string{
			"Reduce grub and earthworm pressure with targeted lawn treatments so the soil holds less prey",
			"Tamp down and repair surface runs quickly so the lawn re-roots before it dries out",
			"Avoid overwatering, which draws worms toward the surface and invites moles",
			"Install an L-shaped underground wire barrier around high-value garden beds",
		},
	})

	RegisterPestOverride("pack-rats", PestOverride{
		ScientificName: "Neotoma cinerea",
		Description:    `The bushy-tailed woodrat is the Inland Northwest's pack rat, and it is not the desert woodrat the old page described. This species ranges all across Idaho and eastern Washington, building its midden, a sprawling fort of sticks, twigs, and stolen treasure, in rockslides, hollow logs, sheds, cabins, and anywhere else it finds cover. Given a gap in the wall, it moves the operation into attics, garages, and vehicle engine bays. Pack rats are compulsive collectors of anything shiny or interesting, from tools to foil insulation, and the real damage is the gnawing: chewed vehicle wiring, RV wiring, and attic electrical. They are nocturnal, loud, and bold, urinating constantly on the midden, and one animal's nest can grow for years if nobody evicts it. Part quartermaster, part saboteur: it steals your gear and chews your comms.`,
		Signs: []string{
			"A messy stick-and-debris midden in a shed corner, under a deck, or behind stored boxes",
			"Chewed vehicle or RV wiring, spark plug leads, or attic electrical with no other obvious culprit",
			"Large dry pellet droppings, bigger and blunter than mouse droppings, near the midden",
			"Loud nocturnal thumping, rustling, and vocal chatter from the attic or walls",
		},
		Treat: `We live-trap or snap-trap the resident woodrats, remove the entire midden and contaminated nesting material, and seal every entry larger than a half inch. The midden's scent trail keeps drawing newcomers back to the same spot until it is gone.`,
		Prev: []string{
			"Declutter sheds, garages, and crawlspaces so there is nowhere to build a midden",
			"Keep vehicle hoods closed and engine bays sealed; store RVs away from brush and rock piles",
			"Remove rock piles, brush piles, and firewood stacked against the house",
			"Seal attic vents and eave gaps with 1/4 inch hardware cloth",
		},
	})

	RegisterPestOverride("raccoons", PestOverride{
		ScientificName: "Procyon lotor",
		Description:    `The raccoon needs no introduction to anyone with a garbage can in Spokane or Post Falls. Raccoons live statewide in Washington, and the suburbs hand them the perfect package: storm drains for travel corridors, uncapped chimneys and attics for dens, and a steady supply of pet food, grills, and trash. Breeding runs January through March, and by April females are hunting for a dark, warm den to raise three to five kits, which is exactly when chimney and attic calls spike. A denning mother shreds insulation for nesting material, and the shared latrines they leave on roofs and decks carry roundworm eggs that survive in soil for years, making droppings a genuine sanitation problem. They do not hibernate, so winter calls still happen, just on the warmer nights. Smart, strong, nocturnal, and already inside your perimeter.`,
		Signs: []string{
			"Heavy thumping and chattering in the attic or chimney at night, especially April through June",
			"Torn soffits, pried roof vents, or a chimney cap knocked loose",
			"Overturned garbage cans with lids removed and contents sorted with surgical precision",
			"Large blunt droppings in a shared latrine spot on the roof, deck, or at the base of a tree",
		},
		Treat: `We use humane live traps or one-way doors on the den entry, with a full check for dependent kits before any exclusion so no young are sealed inside. Then we repair the entry with sheet metal and hardware cloth and sanitize any latrine areas.`,
		Prev: []string{
			"Secure garbage cans with locking lids and bring pet food in at night",
			"Cap the chimney with a wildlife-proof stainless steel cap",
			"Trim tree limbs back at least 6 feet from the roofline to cut off roof access",
			"Close denning sites under porches and sheds with buried 1/4 inch hardware cloth",
		},
	})

	RegisterPestOverride("rodents", PestOverride{
		ScientificName: "Mus musculus, Rattus norvegicus",
		Description:    `The old page named only the Norway rat, but Spokane homes fight a two-front war. The house mouse slips through a gap the width of a pencil and can live its entire life inside your walls, breeding year-round with five to ten litters a year, so a couple of mice becomes an infantry company in months. The Norway rat is the heavy burrower: it tunnels under foundations, favors crawlspaces and ground-level runs, gnaws wiring and framing, and contaminates stored food. Fall is the invasion season for both, as cold weather pushes field populations indoors. Roof rats show up in the region too, but the house mouse and the Norway rat are the two doing the real damage inside Inland Northwest homes. Droppings tell them apart at a glance: rice-grain mouse pellets versus the blunt capsules of a rat. Two different enemies, two different doctrines: identify yours before you pick your weapon.`,
		Signs: []string{
			"Mouse droppings like dark grains of rice along baseboards, in drawers, and behind appliances (rat droppings are larger blunt capsules)",
			"Gnaw marks on food packaging, wiring, and wood, plus greasy rub marks along baseboards from rat fur",
			"Scratching and scurrying in walls or ceilings; mice at any hour, rats mostly at night",
			"Burrow holes about 2 inches wide at the foundation or under decks, with packed smooth edges (rat sign)",
		},
		Treat: `We run a snap-trap and exclusion program tuned to the species: tight trap grids along mouse runways inside, burrow trapping outside for rats, then seal every entry with steel wool, copper mesh, and caulk. Poison alone just creates new vacancies in a food-rich territory.`,
		Prev: []string{
			"Seal gaps around pipes, vents, and foundations with steel wool and caulk; nothing larger than a dime for mice",
			"Store food, pet food, and birdseed in sealed hard containers",
			"Keep firewood and debris piles at least 20 feet from the house",
			"Cut back vegetation touching the house and keep gutters clean and draining",
		},
	})

	RegisterPestOverride("squirrels", PestOverride{
		ScientificName: "Tamiasciurus hudsonicus",
		Description:    `The native tree squirrel of the Inland Northwest is the American red squirrel, a rust-red, white-bellied pine squirrel that Washington's own wildlife agency places squarely in the conifer forests and semi-open woods of the northeast corner. The eastern gray squirrel named on the old page is an introduced species here; the red squirrel is the regional native doing the attic damage. Red squirrels are extremely territorial and loud, chattering and tail-flicking at anything that crosses their turf. In late summer and fall they cut green pinecones and cache them in middens, sometimes in gutters and downspouts, which clogs drainage. They reach the attic through any gap the size of a golf ball, nesting in insulation and gnawing wood and wiring, and spring litters from February through May make that the prime attic-invasion season. Tiny, loud, and extremely territorial: the neighborhood watch with teeth.`,
		Signs: []string{
			"Loud daytime chattering and rapid scurrying in the attic, loudest in the morning (squirrels are diurnal; night noise points to rats)",
			"Chewed entry holes at eaves, gable vents, or fascia, often with gnawed wood around a golf-ball-sized opening",
			"Pinecone scales and debris piles, the squirrel's midden, in gutters or near the den",
			"Nests of shredded insulation and leaves in the attic, with gnawed wiring nearby",
		},
		Treat: `We install one-way doors on attic entries after confirming no dependent young are inside, then seal the chewed entries with metal flashing and hardware cloth. Squirrels chew straight through wood and plastic patches but give up on steel.`,
		Prev: []string{
			"Trim branches at least 8 to 10 feet from the roof so squirrels cannot bridge the gap",
			"Cover gable and roof vents with 1/2 inch hardware cloth",
			"Repair fascia and soffit damage promptly before it becomes a doorway",
			"Keep bird feeders away from the house or use squirrel-proof feeders to stop luring them to the eaves",
		},
	})
}
