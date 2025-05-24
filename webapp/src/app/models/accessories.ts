export type Accessory = {
    id: number,
    value: string
}

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
