<script setup lang="ts">
import CreateRecipe from "~/layouts/create-recipe.vue";
import MinusCircle from "~/components/icons/MinusCircle.vue";
import Button from "~/components/core/Button.vue";
import UpCircle from "~/components/icons/UpCircle.vue";
import DownCircle from "~/components/icons/DownCircle.vue";
import IconButton from "~/components/core/IconButton.vue";

useHead({
  title: "Create a new recipe",
});

definePageMeta({
  layout: "authenticated",
  middleware: ["auth"],
});

type FormData = {
  title: string;
  description: string;
  ingredients: {
    amount: string;
    unit: string;
    ingredient: string;
  }[];
  instructions: {
    step: number;
    instruction: string;
  }[];
  cookTime: string;
  servings: number;
};

const formData = reactive<FormData>({
  title: "",
  description: "",
  ingredients: [
    {
      amount: "",
      unit: "",
      ingredient: "",
    },
  ],
  instructions: [
    {
      step: 1,
      instruction: "",
    },
    {
      step: 2,
      instruction: "",
    },
    {
      step: 3,
      instruction: "",
    },
  ],
  cookTime: "",
  servings: 0,
});

watch(formData.instructions, (value) => {
  console.log(value);
});

const addIngredient = () => {
  formData.ingredients.push({
    amount: "",
    unit: "",
    ingredient: "",
  });
};

const removeIngredient = (i: number) => {
  formData.ingredients.splice(i, 1);
};

const addStep = () => {
  formData.instructions.push({
    step: formData.instructions.length + 1,
    instruction: "",
  });
};

const removeStep = (i: number) => {
  formData.instructions.splice(i, 1);
  resetInstructionNumbers();
};

const moveStep = (i: number, direction: "up" | "down") => {
  const step = formData.instructions[i];
  const swapStepIndex = direction === "up" ? i - 1 : i + 1;

  formData.instructions[i] = formData.instructions[swapStepIndex];
  formData.instructions[swapStepIndex] = step;
  // Update the step numbers
  resetInstructionNumbers();
};

const resetInstructionNumbers = () => {
  formData.instructions = formData.instructions.map((step, index) => ({
    ...step,
    step: index + 1,
  }));
};

const submitForm = () => {
  console.log(formData);
};
</script>

<template>
  <main class="mx-auto max-w-screen-md pt-8">
    <CreateRecipe title="Create a new recipe">
      <div class="mx-auto max-w-2xl py-12 space-y-6">
        <div class="text-center space-y-2">
          <p class="text-gray-500 dark:text-gray-400">
            Share your culinary creations with the community.
          </p>
          <div class="text-gray-500 dark:text-gray-400">
            Describe your recipe in detail so others can recreate it.
          </div>
        </div>
        <form @submit.prevent="submitForm" class="space-y-6">
          <div class="space-y-2">
            <label
              class="peer-disabled:cursor-not-allowed text-sm font-medium leading-none peer-disabled:opacity-70"
              for="title"
            >
              Recipe Title
            </label>
            <input
              class="flex h-10 w-full rounded-md file:border-0 border file:bg-transparent px-3 py-2 text-sm file:text-sm file:font-medium border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              id="title"
              placeholder="Enter recipe title"
              required
              v-model="formData.title"
            />
          </div>
          <div class="space-y-2">
            <label
              class="peer-disabled:cursor-not-allowed text-sm font-medium leading-none peer-disabled:opacity-70"
              for="description"
            >
              Description
            </label>
            <textarea
              class="flex w-full rounded-md border px-3 py-2 text-sm min-h-[80px] border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
              id="description"
              placeholder="Describe your recipe"
              v-model="formData.description"
            ></textarea>
          </div>
          <div class="space-y-2">
            <label
              class="peer-disabled:cursor-not-allowed text-sm font-medium leading-none peer-disabled:opacity-70"
              for="ingredients"
            >
              Ingredients
            </label>
            <div
              class="grid grid-cols-[1fr_1fr_1fr_auto] gap-4 items-center"
              v-for="(ingredient, i) in formData.ingredients"
            >
              <div class="mb-2">
                <label class="text-sm" for="amount">Amount</label>
                <input
                  id="amount"
                  class="flex h-10 w-full rounded-md file:border-0 border file:bg-transparent px-3 py-2 text-sm file:text-sm file:font-medium border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  placeholder="Amount"
                  v-model="ingredient.amount"
                />
              </div>
              <div class="mb-2">
                <label class="text-sm" for="unit">Unit</label>
                <input
                  id="unit"
                  class="flex h-10 w-full rounded-md file:border-0 border file:bg-transparent px-3 py-2 text-sm file:text-sm file:font-medium border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  placeholder="Unit"
                  v-model="ingredient.unit"
                />
              </div>
              <div class="mb-2">
                <label class="text-sm" for="ingredient">Ingredient</label>
                <input
                  id="ingredient"
                  class="flex h-10 w-full rounded-md file:border-0 border file:bg-transparent px-3 py-2 text-sm file:text-sm file:font-medium border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  placeholder="Ingredient"
                  v-model="ingredient.ingredient"
                />
              </div>
              <div class="translate-y-2 flex items-center">
                <IconButton
                  button-style="danger"
                  alt-text="Remove ingredient"
                  @click="() => removeIngredient(i)"
                >
                  <MinusCircle />
                </IconButton>
              </div>
            </div>
            <Button @click="addIngredient" button-style="outline"
              >Add Ingredient
            </Button>
          </div>
          <div class="space-y-2">
            <label
              class="peer-disabled:cursor-not-allowed text-sm font-medium leading-none peer-disabled:opacity-70"
              for="instructions"
            >
              Instructions
            </label>
            <div class="space-y-2">
              <div
                class="grid items-center gap-4 grid-cols-[5ch_1fr_auto]"
                v-for="(step, i) in formData.instructions"
              >
                <div class="font-medium">Step {{ step.step }}</div>
                <textarea
                  v-model="step.instruction"
                  class="flex w-full rounded-md border px-3 py-2 text-sm border-input bg-background ring-offset-background placeholder:text-muted-foreground min-h-[48px] focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                  placeholder="Provide step-by-step instructions"
                  required
                ></textarea>
                <div class="flex items-center">
                  <div class="flex flex-col">
                    <IconButton
                      v-if="i !== 0"
                      :alt-text="`Move step  ${step.step} to ${step.step - 1}`"
                      @click="() => moveStep(i, 'up')"
                    >
                      <UpCircle />
                    </IconButton>
                    <IconButton
                      v-if="i !== formData.instructions.length - 1"
                      :alt-text="`Move step  ${step.step} to ${step.step + 1}`"
                      @click="() => moveStep(i, 'down')"
                    >
                      <DownCircle />
                    </IconButton>
                  </div>
                  <button
                    class="text-red-600 hover:text-red-600/80 transition-colors size-[40px] grid place-items-center focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2"
                    type="button"
                    @click="() => removeStep(i)"
                  >
                    <span class="sr-only">Remove step</span>
                    <MinusCircle />
                  </button>
                </div>
              </div>
            </div>
            <Button @click="addStep" button-style="outline">Add Step</Button>
          </div>
          <div class="grid grid-cols-2 gap-4">
            <div class="space-y-2">
              <label
                class="peer-disabled:cursor-not-allowed text-sm font-medium leading-none peer-disabled:opacity-70"
                for="cook-time"
              >
                Cook Time
              </label>
              <input
                v-model="formData.cookTime"
                class="flex h-10 w-full rounded-md file:border-0 border file:bg-transparent px-3 py-2 text-sm file:text-sm file:font-medium border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                id="cook-time"
                placeholder="e.g. 30 minutes"
                required
                type="text"
              />
            </div>
            <div class="space-y-2">
              <label
                class="peer-disabled:cursor-not-allowed text-sm font-medium leading-none peer-disabled:opacity-70"
                for="servings"
              >
                Servings
              </label>
              <input
                v-model="formData.servings"
                class="flex h-10 w-full rounded-md file:border-0 border file:bg-transparent px-3 py-2 text-sm file:text-sm file:font-medium border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"
                id="servings"
                placeholder="Number of servings"
                required
                type="number"
              />
            </div>
          </div>
          <!--          <div class="space-y-2">-->
          <!--            <label-->
          <!--              class="peer-disabled:cursor-not-allowed text-sm font-medium leading-none peer-disabled:opacity-70"-->
          <!--              for="image"-->
          <!--            >-->
          <!--              Recipe Image (optional)-->
          <!--            </label>-->
          <!--            <input-->
          <!--              class="flex h-10 w-full rounded-md file:border-0 border file:bg-transparent px-3 py-2 text-sm file:text-sm file:font-medium border-input bg-background ring-offset-background placeholder:text-muted-foreground focus-visible:ring-ring focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:cursor-not-allowed disabled:opacity-50"-->
          <!--              id="image"-->
          <!--              type="file"-->
          <!--            />-->
          <!--          </div>-->
          <Button button-style="primary" type="submit"
            >Submit your recipe!
          </Button>
        </form>
      </div>
    </CreateRecipe>
  </main>
</template>

<style scoped></style>
