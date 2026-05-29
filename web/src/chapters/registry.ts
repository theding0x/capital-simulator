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

  // Volume II — The Process of Circulation of Capital
  // Part I — The Metamorphoses of Capital and Their Circuits
  { id: "v2-ch01", volume: 2, number: 1,  title: "The Circuit of Money-Capital",                              part: "Part I — The Metamorphoses of Capital and Their Circuits",          circuitNode: ["whole", "M", "M-prime"],          status: "done" },
  { id: "v2-ch02", volume: 2, number: 2,  title: "The Circuit of Productive Capital",                         part: "Part I — The Metamorphoses of Capital and Their Circuits",          circuitNode: ["whole", "P"],                     status: "done" },
  { id: "v2-ch03", volume: 2, number: 3,  title: "The Circuit of Commodity-Capital",                          part: "Part I — The Metamorphoses of Capital and Their Circuits",          circuitNode: ["whole", "C-prime", "C-M-prime"],  status: "done" },
  { id: "v2-ch04", volume: 2, number: 4,  title: "The Three Formulas of the Circuit",                         part: "Part I — The Metamorphoses of Capital and Their Circuits",          circuitNode: ["whole"],                          status: "done"    },
  { id: "v2-ch05", volume: 2, number: 5,  title: "The Time of Circulation",                                   part: "Part I — The Metamorphoses of Capital and Their Circuits",          circuitNode: ["M-C", "C-M-prime"],               status: "done" },
  { id: "v2-ch06", volume: 2, number: 6,  title: "The Costs of Circulation",                                  part: "Part I — The Metamorphoses of Capital and Their Circuits",          circuitNode: ["M-C", "C-M-prime"],               status: "done" },
  // Part II — The Turnover of Capital
  { id: "v2-ch07", volume: 2, number: 7,  title: "The Turnover Time and the Number of Turnovers",             part: "Part II — The Turnover of Capital",                                 circuitNode: ["whole"],                          status: "done"    },
  { id: "v2-ch08", volume: 2, number: 8,  title: "Fixed Capital and Circulating Capital",                     part: "Part II — The Turnover of Capital",                                 circuitNode: ["whole"],                          status: "done"    },
  { id: "v2-ch09", volume: 2, number: 9,  title: "The Aggregate Turnover of Advanced Capital",                part: "Part II — The Turnover of Capital",                                 circuitNode: ["whole"],                          status: "done"    },
  { id: "v2-ch10", volume: 2, number: 10, title: "Theories of Fixed and Circulating Capital — Physiocrats and Adam Smith", part: "Part II — The Turnover of Capital",                       circuitNode: ["whole"],                          status: "done"    },
  { id: "v2-ch11", volume: 2, number: 11, title: "Theories of Fixed and Circulating Capital — Ricardo",       part: "Part II — The Turnover of Capital",                                 circuitNode: ["whole"],                          status: "done"    },
  { id: "v2-ch12", volume: 2, number: 12, title: "The Working Period",                                        part: "Part II — The Turnover of Capital",                                 circuitNode: ["P"],                              status: "done"    },
  { id: "v2-ch13", volume: 2, number: 13, title: "The Time of Production",                                    part: "Part II — The Turnover of Capital",                                 circuitNode: ["P"],                              status: "done"    },
  { id: "v2-ch14", volume: 2, number: 14, title: "The Time of Circulation",                                   part: "Part II — The Turnover of Capital",                                 circuitNode: ["M-C", "C-M-prime"],               status: "done"    },
  { id: "v2-ch15", volume: 2, number: 15, title: "The Effects of a Change of Prices",                         part: "Part II — The Turnover of Capital",                                 circuitNode: ["whole"],                          status: "done"    },
  { id: "v2-ch16", volume: 2, number: 16, title: "The Turnover of Variable Capital",                          part: "Part II — The Turnover of Capital",                                 circuitNode: ["P", "delta-M"],                   status: "done"    },
  { id: "v2-ch17", volume: 2, number: 17, title: "The Circulation of Surplus-Value",                          part: "Part II — The Turnover of Capital",                                 circuitNode: ["delta-M", "whole"],               status: "done"    },
  // Part III — The Reproduction and Circulation of the Total Social Capital
  { id: "v2-ch18", volume: 2, number: 18, title: "The Role of Money-Capital in Reproduction",                 part: "Part III — The Reproduction and Circulation of the Total Social Capital", circuitNode: ["whole"],                    status: "done"    },
  { id: "v2-ch19", volume: 2, number: 19, title: "Former Presentations of the Subject",                       part: "Part III — The Reproduction and Circulation of the Total Social Capital", circuitNode: ["whole", "historical"],      status: "done"    },
  { id: "v2-ch20", volume: 2, number: 20, title: "Simple Reproduction",                                       part: "Part III — The Reproduction and Circulation of the Total Social Capital", circuitNode: ["whole"],                    status: "done"    },
  { id: "v2-ch21", volume: 2, number: 21, title: "Accumulation and Reproduction on an Extended Scale",        part: "Part III — The Reproduction and Circulation of the Total Social Capital", circuitNode: ["whole", "delta-M"],         status: "done"    },

  // Volume III — The Process of Capitalist Production as a Whole
  // Part I — The Transformation of Surplus-Value into Profit
  { id: "v3-ch01", volume: 3, number: 1,  title: "Cost-Price and Profit",                                     part: "Part I — The Transformation of Surplus-Value into Profit",          circuitNode: ["delta-M"],                        status: "done"    },
  { id: "v3-ch02", volume: 3, number: 2,  title: "The Rate of Profit",                                        part: "Part I — The Transformation of Surplus-Value into Profit",          circuitNode: ["delta-M"],                        status: "done"    },
  { id: "v3-ch03", volume: 3, number: 3,  title: "The Relation of the Rate of Profit to the Rate of Surplus-Value", part: "Part I — The Transformation of Surplus-Value into Profit",    circuitNode: ["delta-M"],                        status: "done"    },
  { id: "v3-ch04", volume: 3, number: 4,  title: "The Effect of the Turnover on the Rate of Profit",          part: "Part I — The Transformation of Surplus-Value into Profit",          circuitNode: ["M-C", "C-M-prime", "delta-M"],    status: "done"    },
  { id: "v3-ch05", volume: 3, number: 5,  title: "Economy in the Employment of Constant Capital",             part: "Part I — The Transformation of Surplus-Value into Profit",          circuitNode: ["P", "delta-M"],                   status: "pending" },
  { id: "v3-ch06", volume: 3, number: 6,  title: "The Effect of Price Fluctuation",                           part: "Part I — The Transformation of Surplus-Value into Profit",          circuitNode: ["C-M-prime", "delta-M"],           status: "pending" },
  { id: "v3-ch07", volume: 3, number: 7,  title: "Supplementary Remarks",                                     part: "Part I — The Transformation of Surplus-Value into Profit",          circuitNode: ["delta-M"],                        status: "pending" },
  // Part II — Conversion of Profit into Average Profit
  { id: "v3-ch08", volume: 3, number: 8,  title: "Different Compositions of Capitals in Different Branches",  part: "Part II — Conversion of Profit into Average Profit",                circuitNode: ["delta-M", "whole"],               status: "pending" },
  { id: "v3-ch09", volume: 3, number: 9,  title: "Formation of a General Rate of Profit",                     part: "Part II — Conversion of Profit into Average Profit",                circuitNode: ["delta-M", "whole"],               status: "pending" },
  { id: "v3-ch10", volume: 3, number: 10, title: "Equalisation of the General Rate of Profit",                part: "Part II — Conversion of Profit into Average Profit",                circuitNode: ["C-M-prime", "delta-M"],           status: "pending" },
  { id: "v3-ch11", volume: 3, number: 11, title: "Effects of General Wage Fluctuations on Prices of Production", part: "Part II — Conversion of Profit into Average Profit",             circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch12", volume: 3, number: 12, title: "Supplementary Remarks",                                     part: "Part II — Conversion of Profit into Average Profit",                circuitNode: ["delta-M"],                        status: "pending" },
  // Part III — The Law of the Tendential Fall in the Rate of Profit
  { id: "v3-ch13", volume: 3, number: 13, title: "The Law as Such",                                           part: "Part III — The Law of the Tendential Fall in the Rate of Profit",   circuitNode: ["delta-M", "whole"],               status: "pending" },
  { id: "v3-ch14", volume: 3, number: 14, title: "Counteracting Influences",                                  part: "Part III — The Law of the Tendential Fall in the Rate of Profit",   circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch15", volume: 3, number: 15, title: "Internal Contradictions of the Law",                        part: "Part III — The Law of the Tendential Fall in the Rate of Profit",   circuitNode: ["delta-M", "whole"],               status: "pending" },
  // Part IV — Conversion of Commodity-Capital and Money-Capital into Commercial Capital
  { id: "v3-ch16", volume: 3, number: 16, title: "Commercial Capital",                                        part: "Part IV — Conversion of Commodity-Capital and Money-Capital into Commercial Capital", circuitNode: ["C-M-prime", "delta-M"],   status: "pending" },
  { id: "v3-ch17", volume: 3, number: 17, title: "Commercial Profit",                                         part: "Part IV — Conversion of Commodity-Capital and Money-Capital into Commercial Capital", circuitNode: ["C-M-prime", "delta-M"],   status: "pending" },
  { id: "v3-ch18", volume: 3, number: 18, title: "The Turnover of Merchant's Capital",                        part: "Part IV — Conversion of Commodity-Capital and Money-Capital into Commercial Capital", circuitNode: ["C-M-prime", "delta-M"],   status: "pending" },
  { id: "v3-ch19", volume: 3, number: 19, title: "Money-Dealing Capital",                                     part: "Part IV — Conversion of Commodity-Capital and Money-Capital into Commercial Capital", circuitNode: ["M", "delta-M"],           status: "pending" },
  { id: "v3-ch20", volume: 3, number: 20, title: "Historical Facts about Merchant's Capital",                 part: "Part IV — Conversion of Commodity-Capital and Money-Capital into Commercial Capital", circuitNode: ["historical", "delta-M"],  status: "pending" },
  // Part V — Division of Profit into Interest and Profit of Enterprise
  { id: "v3-ch21", volume: 3, number: 21, title: "Interest-Bearing Capital",                                  part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["delta-M"],                       status: "pending" },
  { id: "v3-ch22", volume: 3, number: 22, title: "Division of Profit; Rate of Interest",                      part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["delta-M"],                       status: "pending" },
  { id: "v3-ch23", volume: 3, number: 23, title: "Interest and Profit of Enterprise",                         part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["delta-M"],                       status: "pending" },
  { id: "v3-ch24", volume: 3, number: 24, title: "Externalisation of the Relations of Capital",               part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["delta-M", "whole"],              status: "pending" },
  { id: "v3-ch25", volume: 3, number: 25, title: "Credit and Fictitious Capital",                             part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["delta-M"],                       status: "pending" },
  { id: "v3-ch26", volume: 3, number: 26, title: "Accumulation of Money-Capital",                             part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["M", "M-prime", "delta-M"],       status: "pending" },
  { id: "v3-ch27", volume: 3, number: 27, title: "The Role of Credit in Capitalist Production",               part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["delta-M"],                       status: "pending" },
  { id: "v3-ch28", volume: 3, number: 28, title: "Medium of Circulation and Capital — Tooke and Fullarton",   part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["M", "delta-M"],                  status: "pending" },
  { id: "v3-ch29", volume: 3, number: 29, title: "Component Parts of Bank Capital",                           part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["delta-M"],                       status: "pending" },
  { id: "v3-ch30", volume: 3, number: 30, title: "Money-Capital and Real Capital, I",                         part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["M", "delta-M"],                  status: "pending" },
  { id: "v3-ch31", volume: 3, number: 31, title: "Money-Capital and Real Capital, II",                        part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["M", "delta-M"],                  status: "pending" },
  { id: "v3-ch32", volume: 3, number: 32, title: "Money-Capital and Real Capital, III",                       part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["M", "delta-M"],                  status: "pending" },
  { id: "v3-ch33", volume: 3, number: 33, title: "The Medium of Circulation in the Credit System",            part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["M", "delta-M"],                  status: "pending" },
  { id: "v3-ch34", volume: 3, number: 34, title: "The Currency Principle and the English Bank Legislation",   part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["M", "delta-M", "historical"],    status: "pending" },
  { id: "v3-ch35", volume: 3, number: 35, title: "Precious Metal and Rate of Exchange",                       part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["M"],                             status: "pending" },
  { id: "v3-ch36", volume: 3, number: 36, title: "Pre-Capitalist Relationships",                              part: "Part V — Division of Profit into Interest and Profit of Enterprise", circuitNode: ["historical"],                    status: "pending" },
  // Part VI — Transformation of Surplus-Profit into Ground-Rent
  { id: "v3-ch37", volume: 3, number: 37, title: "Introduction",                                              part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch38", volume: 3, number: 38, title: "General Remarks on Differential Rent",                      part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch39", volume: 3, number: 39, title: "First Form of Differential Rent — Differential Rent I",     part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch40", volume: 3, number: 40, title: "Second Form of Differential Rent — Differential Rent II",   part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch41", volume: 3, number: 41, title: "Differential Rent II — Constant Price of Production",       part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch42", volume: 3, number: 42, title: "Differential Rent II — Falling Price of Production",        part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch43", volume: 3, number: 43, title: "Differential Rent II — Rising Price of Production",         part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch44", volume: 3, number: 44, title: "Differential Rent Also on the Worst Cultivated Soil",       part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch45", volume: 3, number: 45, title: "Absolute Ground-Rent",                                      part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch46", volume: 3, number: 46, title: "Building Site Rent; Rent in Mining",                        part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["delta-M"],                        status: "pending" },
  { id: "v3-ch47", volume: 3, number: 47, title: "Genesis of Capitalist Ground-Rent",                         part: "Part VI — Transformation of Surplus-Profit into Ground-Rent",       circuitNode: ["historical", "delta-M"],          status: "pending" },
  // Part VII — The Revenues and Their Sources
  { id: "v3-ch48", volume: 3, number: 48, title: "The Trinity Formula",                                       part: "Part VII — The Revenues and Their Sources",                         circuitNode: ["whole", "delta-M"],               status: "pending" },
  { id: "v3-ch49", volume: 3, number: 49, title: "Concerning the Analysis of the Process of Production",      part: "Part VII — The Revenues and Their Sources",                         circuitNode: ["whole", "P", "delta-M"],          status: "pending" },
  { id: "v3-ch50", volume: 3, number: 50, title: "Illusions Created by Competition",                          part: "Part VII — The Revenues and Their Sources",                         circuitNode: ["whole", "delta-M"],               status: "pending" },
  { id: "v3-ch51", volume: 3, number: 51, title: "Distribution Relations and Production Relations",           part: "Part VII — The Revenues and Their Sources",                         circuitNode: ["whole"],                          status: "pending" },
  { id: "v3-ch52", volume: 3, number: 52, title: "Classes",                                                   part: "Part VII — The Revenues and Their Sources",                         circuitNode: ["historical", "whole"],            status: "pending" },
]
