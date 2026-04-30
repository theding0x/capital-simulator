// TypeScript shapes that mirror the Go structs in commodity-service.
// Keep these in sync with services/commodity-service/internal/commodity/*.go.

export interface ConcreteLabour {
  kind: string;
  description?: string;
}

export interface UseValue {
  description: string;
  unit: string;
}

export interface Commodity {
  id: string;
  name: string;
  use_value: UseValue;
  concrete_labour: ConcreteLabour;
  snlt_per_unit: number; // labour-minutes per unit
  created_at: string;
  updated_at: string;
}

export interface ValueResponse {
  commodity: Commodity;
  quantity: number;
  value: number; // labour-minutes
  value_hours: number;
}

export interface ExchangeRatio {
  base: Commodity;
  quote: Commodity;
  base_qty: number;
  quote_qty: number;
  common_value: number;
}

export interface SimpleForm {
  relative: Commodity;
  equivalent: Commodity;
  relative_qty: number;
  equivalent_qty: number;
  common_value: number;
}

export interface LabourRelation {
  counterpart: Commodity;
  counterpart_qty: number;
  labour_time: number;
}

export interface SocialRelations {
  subject: Commodity;
  labour_per_unit: number;
  concrete_labour: ConcreteLabour;
  labour_relations: LabourRelation[];
  note: string;
}

export interface CreateCommodityInput {
  name: string;
  use_value: UseValue;
  concrete_labour: ConcreteLabour;
  snlt_per_unit: number;
}

export interface UpdateCommodityInput {
  name?: string;
  use_value_description?: string;
  use_value_unit?: string;
  concrete_labour_kind?: string;
  concrete_labour_description?: string;
  snlt_per_unit?: number;
}
