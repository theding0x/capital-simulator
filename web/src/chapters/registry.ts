export type ChapterStatus = "done" | "pending"

/**
 * Moments of the circuit M—C(Lp+Mp)…P…C'—M'.
 *
 *   M           money capital advanced
 *   M-C         purchase phase: buying labour-power and means of production
 *   P           production / valorization
 *   C-prime     finished commodity carrying surplus-value
 *   C-M-prime   sale phase: realising value in money
 *   M-prime     money capital returned (= M + ΔM)
 *   delta-M     surplus-value as such (its distribution lives mainly in Vol. III)
 *   whole       the totality of the circuit — used for chapters that read the
 *               whole motion rather than a single moment
 *   historical  the conditions that produced the circuit in the first place
 *               (primitive accumulation; secular tendencies)
 */
export type CircuitNode =
  | "M"
  | "M-C"
  | "P"
  | "C-prime"
  | "C-M-prime"
  | "M-prime"
  | "delta-M"
  | "whole"
  | "historical"

export type Volume = 1 | 2 | 3

export type ChapterDef = {
  id: string
  volume: Volume
  number: number
  title: string
  part: string
  circuitNode: CircuitNode[]
  status: ChapterStatus
}

export const CHAPTERS: ChapterDef[] = [
  // Volume I — The Process of Production of Capital
  // Part I — Commodities & Money
  { id: "v1-ch01", volume: 1, number: 1,  title: "The Commodity",                                            part: "Part I — Commodities & Money",                         circuitNode: ["whole"],                          status: "done"    },
  { id: "v1-ch02", volume: 1, number: 2,  title: "Exchange",                                                  part: "Part I — Commodities & Money",                         circuitNode: ["M-C", "C-M-prime"],               status: "done"    },
  { id: "v1-ch03", volume: 1, number: 3,  title: "Money, or the Circulation of Commodities",                  part: "Part I — Commodities & Money",                         circuitNode: ["M", "M-prime"],                   status: "done"    },
  // Part II — The Transformation of Money into Capital
  { id: "v1-ch04", volume: 1, number: 4,  title: "The General Formula for Capital",                           part: "Part II — The Transformation of Money into Capital",   circuitNode: ["whole"],                          status: "done"    },
  { id: "v1-ch05", volume: 1, number: 5,  title: "Contradictions in the General Formula",                     part: "Part II — The Transformation of Money into Capital",   circuitNode: ["whole"],                          status: "done"    },
  { id: "v1-ch06", volume: 1, number: 6,  title: "The Sale and Purchase of Labour-Power",                     part: "Part II — The Transformation of Money into Capital",   circuitNode: ["M-C"],                            status: "done"    },
  // Part III — The Production of Absolute Surplus-Value
  { id: "v1-ch07", volume: 1, number: 7,  title: "The Labour-Process and the Valorization Process",           part: "Part III — The Production of Absolute Surplus-Value",  circuitNode: ["P"],                              status: "done"    },
  { id: "v1-ch08", volume: 1, number: 8,  title: "Constant Capital and Variable Capital",                     part: "Part III — The Production of Absolute Surplus-Value",  circuitNode: ["P"],                              status: "done"    },
  { id: "v1-ch09", volume: 1, number: 9,  title: "The Rate of Surplus-Value",                                 part: "Part III — The Production of Absolute Surplus-Value",  circuitNode: ["P", "delta-M"],                   status: "done"    },
  { id: "v1-ch10", volume: 1, number: 10, title: "The Working Day",                                           part: "Part III — The Production of Absolute Surplus-Value",  circuitNode: ["P"],                              status: "done"    },
  { id: "v1-ch11", volume: 1, number: 11, title: "The Rate and Mass of Surplus-Value",                        part: "Part III — The Production of Absolute Surplus-Value",  circuitNode: ["P", "delta-M"],                   status: "done"    },
  // Part IV — The Production of Relative Surplus-Value
  { id: "v1-ch12", volume: 1, number: 12, title: "The Concept of Relative Surplus-Value",                     part: "Part IV — The Production of Relative Surplus-Value",   circuitNode: ["P"],                              status: "done"    },
  { id: "v1-ch13", volume: 1, number: 13, title: "Co-operation",                                              part: "Part IV — The Production of Relative Surplus-Value",   circuitNode: ["P"],                              status: "done"    },
  { id: "v1-ch14", volume: 1, number: 14, title: "Division of Labour and Manufacture",                        part: "Part IV — The Production of Relative Surplus-Value",   circuitNode: ["P"],                              status: "done"    },
  { id: "v1-ch15", volume: 1, number: 15, title: "Machinery and Modern Industry",                             part: "Part IV — The Production of Relative Surplus-Value",   circuitNode: ["P"],                              status: "done"    },
  // Part V — Absolute and Relative Surplus-Value
  { id: "v1-ch16", volume: 1, number: 16, title: "Absolute and Relative Surplus-Value",                       part: "Part V — Absolute and Relative Surplus-Value",         circuitNode: ["P", "delta-M"],                   status: "done"    },
  { id: "v1-ch17", volume: 1, number: 17, title: "Changes of Magnitude in the Price of Labour-Power",         part: "Part V — Absolute and Relative Surplus-Value",         circuitNode: ["P", "delta-M"],                   status: "done"    },
  { id: "v1-ch18", volume: 1, number: 18, title: "Different Formulae for the Rate of Surplus-Value",          part: "Part V — Absolute and Relative Surplus-Value",         circuitNode: ["delta-M"],                        status: "done"    },
  // Part VI — Wages
  { id: "v1-ch19", volume: 1, number: 19, title: "The Transformation of the Value of Labour-Power into Wages", part: "Part VI — Wages",                                     circuitNode: ["M-C", "delta-M"],                 status: "done"    },
  { id: "v1-ch20", volume: 1, number: 20, title: "Time-Wages",                                                part: "Part VI — Wages",                                      circuitNode: ["M-C"],                            status: "done"    },
  { id: "v1-ch21", volume: 1, number: 21, title: "Piece-Wages",                                               part: "Part VI — Wages",                                      circuitNode: ["M-C"],                            status: "done"    },
  { id: "v1-ch22", volume: 1, number: 22, title: "National Differences in Wages",                             part: "Part VI — Wages",                                      circuitNode: ["M-C"],                            status: "done"    },
  // Part VII — The Accumulation of Capital
  { id: "v1-ch23", volume: 1, number: 23, title: "Simple Reproduction",                                       part: "Part VII — The Accumulation of Capital",               circuitNode: ["whole"],                          status: "done"    },
  { id: "v1-ch24", volume: 1, number: 24, title: "The Transformation of Surplus-Value into Capital",          part: "Part VII — The Accumulation of Capital",               circuitNode: ["whole", "delta-M"],               status: "done"    },
  { id: "v1-ch25", volume: 1, number: 25, title: "The General Law of Capitalist Accumulation",                part: "Part VII — The Accumulation of Capital",               circuitNode: ["whole", "delta-M", "historical"], status: "done"    },
  // Part VIII — So-Called Primitive Accumulation
  { id: "v1-ch26", volume: 1, number: 26, title: "The Secret of Primitive Accumulation",                      part: "Part VIII — So-Called Primitive Accumulation",         circuitNode: ["historical"],                     status: "done"    },
  { id: "v1-ch27", volume: 1, number: 27, title: "Expropriation of the Agricultural Population",              part: "Part VIII — So-Called Primitive Accumulation",         circuitNode: ["historical"],                     status: "done"    },
  { id: "v1-ch28", volume: 1, number: 28, title: "Bloody Legislation against the Expropriated",               part: "Part VIII — So-Called Primitive Accumulation",         circuitNode: ["historical"],                     status: "done"    },
  { id: "v1-ch29", volume: 1, number: 29, title: "Genesis of the Capitalist Farmer",                          part: "Part VIII — So-Called Primitive Accumulation",         circuitNode: ["historical"],                     status: "done"    },
  { id: "v1-ch30", volume: 1, number: 30, title: "Reaction of the Agricultural Revolution on Industry",       part: "Part VIII — So-Called Primitive Accumulation",         circuitNode: ["historical"],                     status: "done"    },
  { id: "v1-ch31", volume: 1, number: 31, title: "Genesis of the Industrial Capitalist",                      part: "Part VIII — So-Called Primitive Accumulation",         circuitNode: ["historical"],                     status: "done"    },
  { id: "v1-ch32", volume: 1, number: 32, title: "The Historical Tendency of Capitalist Accumulation",        part: "Part VIII — So-Called Primitive Accumulation",         circuitNode: ["historical"],                     status: "done"    },
  { id: "v1-ch33", volume: 1, number: 33, title: "The Modern Theory of Colonisation",                         part: "Part VIII — So-Called Primitive Accumulation",         circuitNode: ["historical"],                     status: "done"    },
]
