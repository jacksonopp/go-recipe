export type Ingredient = {
  id: number;
  name: string;
  quantity: string;
  unit: string;
};

export type Instruction = {
  id: number;
  step: number;
  contents: string;
};

export type RecipeType = {
  id: number;
  description: string;
  created_at: string;
  name: string;
  ingredients: Ingredient[];
  instructions: Instruction[];
};
