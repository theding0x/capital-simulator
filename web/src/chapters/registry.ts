export type ChapterStatus = "done" | "pending"

export type ChapterDef = {
  id: string
  number: number
  title: string
  part: string
  status: ChapterStatus
}

export const CHAPTERS: ChapterDef[] = [
  // Part I — Commodities & Money
  { id: "ch01", number: 1,  title: "The Commodity",                                        part: "Part I — Commodities & Money",                         status: "done"    },
  { id: "ch02", number: 2,  title: "Exchange",                                              part: "Part I — Commodities & Money",                         status: "done"    },
  { id: "ch03", number: 3,  title: "Money, or the Circulation of Commodities",              part: "Part I — Commodities & Money",                         status: "done"    },
  // Part II — The Transformation of Money into Capital
  { id: "ch04", number: 4,  title: "The General Formula for Capital",                       part: "Part II — The Transformation of Money into Capital",   status: "done"    },
  { id: "ch05", number: 5,  title: "Contradictions in the General Formula",                 part: "Part II — The Transformation of Money into Capital",   status: "done"    },
  { id: "ch06", number: 6,  title: "The Sale and Purchase of Labour-Power",                 part: "Part II — The Transformation of Money into Capital",   status: "done"    },
  // Part III — The Production of Absolute Surplus-Value
  { id: "ch07", number: 7,  title: "The Labour-Process and the Valorization Process",       part: "Part III — The Production of Absolute Surplus-Value",  status: "done"    },
  { id: "ch08", number: 8,  title: "Constant Capital and Variable Capital",                 part: "Part III — The Production of Absolute Surplus-Value",  status: "pending" },
  { id: "ch09", number: 9,  title: "The Rate of Surplus-Value",                             part: "Part III — The Production of Absolute Surplus-Value",  status: "pending" },
  { id: "ch10", number: 10, title: "The Working Day",                                       part: "Part III — The Production of Absolute Surplus-Value",  status: "pending" },
  { id: "ch11", number: 11, title: "The Rate and Mass of Surplus-Value",                    part: "Part III — The Production of Absolute Surplus-Value",  status: "pending" },
  // Part IV — The Production of Relative Surplus-Value
  { id: "ch12", number: 12, title: "The Concept of Relative Surplus-Value",                 part: "Part IV — The Production of Relative Surplus-Value",   status: "pending" },
  { id: "ch13", number: 13, title: "Co-operation",                                          part: "Part IV — The Production of Relative Surplus-Value",   status: "pending" },
  { id: "ch14", number: 14, title: "Division of Labour and Manufacture",                    part: "Part IV — The Production of Relative Surplus-Value",   status: "pending" },
  { id: "ch15", number: 15, title: "Machinery and Modern Industry",                         part: "Part IV — The Production of Relative Surplus-Value",   status: "pending" },
  // Part V — Absolute and Relative Surplus-Value
  { id: "ch16", number: 16, title: "Absolute and Relative Surplus-Value",                   part: "Part V — Absolute and Relative Surplus-Value",         status: "pending" },
  { id: "ch17", number: 17, title: "Changes of Magnitude in the Price of Labour-Power",     part: "Part V — Absolute and Relative Surplus-Value",         status: "pending" },
  { id: "ch18", number: 18, title: "Different Formulae for the Rate of Surplus-Value",      part: "Part V — Absolute and Relative Surplus-Value",         status: "pending" },
  // Part VI — Wages
  { id: "ch19", number: 19, title: "The Transformation of the Value of Labour-Power into Wages", part: "Part VI — Wages",                                status: "pending" },
  { id: "ch20", number: 20, title: "Time-Wages",                                            part: "Part VI — Wages",                                      status: "pending" },
  { id: "ch21", number: 21, title: "Piece-Wages",                                           part: "Part VI — Wages",                                      status: "pending" },
  { id: "ch22", number: 22, title: "National Differences in Wages",                         part: "Part VI — Wages",                                      status: "pending" },
  // Part VII — The Accumulation of Capital
  { id: "ch23", number: 23, title: "Simple Reproduction",                                   part: "Part VII — The Accumulation of Capital",               status: "pending" },
  { id: "ch24", number: 24, title: "The Transformation of Surplus-Value into Capital",      part: "Part VII — The Accumulation of Capital",               status: "pending" },
  { id: "ch25", number: 25, title: "The General Law of Capitalist Accumulation",            part: "Part VII — The Accumulation of Capital",               status: "pending" },
  // Part VIII — So-Called Primitive Accumulation
  { id: "ch26", number: 26, title: "The Secret of Primitive Accumulation",                  part: "Part VIII — So-Called Primitive Accumulation",         status: "pending" },
  { id: "ch27", number: 27, title: "Expropriation of the Agricultural Population",          part: "Part VIII — So-Called Primitive Accumulation",         status: "pending" },
  { id: "ch28", number: 28, title: "Bloody Legislation against the Expropriated",           part: "Part VIII — So-Called Primitive Accumulation",         status: "pending" },
  { id: "ch29", number: 29, title: "Genesis of the Capitalist Farmer",                      part: "Part VIII — So-Called Primitive Accumulation",         status: "pending" },
  { id: "ch30", number: 30, title: "Reaction of the Agricultural Revolution on Industry",   part: "Part VIII — So-Called Primitive Accumulation",         status: "pending" },
  { id: "ch31", number: 31, title: "Genesis of the Industrial Capitalist",                  part: "Part VIII — So-Called Primitive Accumulation",         status: "pending" },
  { id: "ch32", number: 32, title: "The Historical Tendency of Capitalist Accumulation",    part: "Part VIII — So-Called Primitive Accumulation",         status: "pending" },
  { id: "ch33", number: 33, title: "The Modern Theory of Colonisation",                     part: "Part VIII — So-Called Primitive Accumulation",         status: "pending" },
]
