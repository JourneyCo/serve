export interface ProjectAccessory {
    id: number;
    name: string;
}

export const Tools: Record<number, string> = {
    1: "Chainsaw",
    2: "Drill",
    3: "Gloves",
    4: "Hacksaw",
    5: "Hammer",
    6: "Hand saw",
    7: "Hedge trimmers",
    8: "Hoe",
    9: "Ladder",
    10: "Lawn mower",
    11: "Leaf blower",
    12: "Level",
    13: "Miter saw",
    14: "Nail gun",
    15: "Paint gun",
    16: "Pickaxe",
    17: "Pitch fork",
    18: "Pliers",
    19: "Pressure washer",
    20: "Pruners",
    21: "Putty knife",
    22: "Rake",
    23: "Sander",
    24: "Sawhorse",
    25: "Screwdriver",
    26: "Shovel",
    27: "Sledgehammer",
    28: "Socket set",
    29: "Spanner wrench",
    30: "Tape measure",
    31: "Utility knife",
    32: "Weedeater",
    33: "Wheelbarrow",
    34: "Wrench",
} as const;

export const Categories: Record<number, string> = {
    1: "Arts & Crafts",
    2: "Community Outreach",
    3: "Español",
    4: "Family/Kid Friendly",
    5: "Food Prep & Distribution",
    6: "Indoors",
    7: "Landscaping",
    8: "Minor Home Repairs",
    9: "Outdoor",
    10: "Painting",
    11: "Prayer & Visitations",
    12: "Skilled Construction",
    13: "Sorting/Assembly",
    14: "Block Party",
    15: "Kids Ministry",
} as const;

export const Skills: Record<number, string> = {
    1: "Carpentry",
    2: "Communication",
    3: "Construction",
    4: "Cooking",
    5: "Hospitality",
    6: "Landscaping",
    7: "Musical",
    8: "Organizational",
    9: "Painting",
    10: "Photography",
} as const;

export const Supplies: Record<number, string> = {
    1: "Bleach",
    2: "Car Wash Supplies",
    3: "Cleaning supplies",
    4: "Craft supplies",
    5: "Duct tape",
    6: "Grilling Supplies",
    7: "Landscape supplies",
    8: "Nails",
    9: "Paint supplies",
    10: "Screws",
} as const;

export const Ages: Record<number, string> = {
    1: "All ages",
    2: "11-14",
    3: "15-18",
    4: "19-22",
    5: "20s",
    6: "30s",
    7: "40s",
    8: "50s",
    9: "60s+",
} as const;