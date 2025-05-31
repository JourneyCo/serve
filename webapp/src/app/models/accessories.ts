export type Accessory = {
    id: number,
    value: string
}

export const Types: Record<number, string> = {
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
    16: "Yard Work",
    17: "Kid Friendly",
    18: "Collection/Drop Off",
    19: "Food Prep/Distribution",

} as const;

export const Ages = [
    'All ages',
    "Families with Small Children ONLY",
    "Families with Young Children",
    "1 Year and Older",
    "2 Years and Older",
    "3 Years and Older",
    "4 Years and Older",
    "5 Years and Older",
    "6 Years and Older",
    "7 Years and Older",
    "8 Years and Older",
    "9 Years and Older",
    "10 Years and Older",
    "11 Years and Older",
    "12 Years and Older",
    "13 Years and Older",
    "14 Years and Older",
    "15 Years and Older",
    "16 Years and Older",
    "17 Years and Older",
    "18 Years and Older",
    "19 Years and Older",
    "20 Years and Older",
    "21 Years and Older",
] as const;
